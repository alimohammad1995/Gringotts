use crate::msg_codec::{ChainTransfer, Message, CHAIN_TRANSFER_TYPE};
use crate::state::{Gringotts, Peer};
use crate::utils::change_decimals;
use crate::*;
use anchor_spl::associated_token::AssociatedToken;
use anchor_spl::token::{CloseAccount, Mint, Token, TokenAccount, Transfer};
use anchor_spl::{associated_token, token};
use oapp::endpoint::cpi::accounts::Clear;
use oapp::endpoint::instructions::ClearParams;
use oapp::endpoint::ConstructCPIContext;

#[derive(Accounts)]
#[instruction(params: LzReceiveParams)]
pub struct LzReceive<'info> {
    #[account(seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,
    #[account(
        seeds = [PEER_SEED, &gringotts.lz_eid.to_le_bytes()],
        bump = self_peer.bump
    )]
    pub self_peer: Account<'info, Peer>,
    #[account(
        seeds = [PEER_SEED, &params.src_eid.to_le_bytes()],
        bump = peer.bump,
        constraint = params.sender == peer.address
    )]
    pub peer: Account<'info, Peer>,
    #[account(mut, seeds = [VAULT_SEED], bump)]
    pub vault: SystemAccount<'info>,

    pub associated_token_program: Program<'info, AssociatedToken>,
    pub token_program: Program<'info, Token>,
    pub system_program: Program<'info, System>,
}

impl<'info> LzReceive<'info> {
    pub fn apply(
        ctx: Context<'_, '_, 'info, 'info, LzReceive<'info>>,
        params: &LzReceiveParams,
    ) -> Result<()> {
        let gringotts = &ctx.accounts.gringotts;
        let seeds: &[&[u8]] = &[GRINGOTTS_SEED, &[gringotts.bump]];
        let signer_seeds = &[seeds];

        let self_peer = &ctx.accounts.self_peer;
        let peer = &ctx.accounts.peer;

        let all_accounts: Vec<AccountInfo<'info>> = ctx.remaining_accounts
            .iter()
            .map(|account| account.to_account_info())
            .chain([
                ctx.accounts.gringotts.to_account_info(),
                ctx.accounts.vault.to_account_info(),
                ctx.accounts.associated_token_program.to_account_info(),
                ctx.accounts.token_program.to_account_info(),
                ctx.accounts.system_program.to_account_info(),
            ])
            .collect();
        let remaining_accounts = &all_accounts;

        let mut r = 0;

        let message = Message::decode(params.message.as_slice());

        if message.header == CHAIN_TRANSFER_TYPE {
            let chain_transfer = ChainTransfer::decode(message.payload);

            let accounts_count = chain_transfer.metadata[0] as usize;
            let mappings = &chain_transfer.metadata[accounts_count * 33 + 1..];

            let vault = get_account(remaining_accounts, r, mappings);
            let (_, vault_bump) = Pubkey::find_program_address(&[VAULT_SEED], ctx.program_id);

            let system_program = get_account(remaining_accounts, r + 1, mappings);
            let token_program = get_account(remaining_accounts, r + 2, mappings);
            let associated_token_program = get_account(remaining_accounts, r + 3, mappings);

            r += 4;

            for item in &chain_transfer.items {
                let recipient = &remaining_accounts[r];
                r += 1;

                if self_peer.has_stable_coin(*item.asset) {
                    let stable_coin_mint = get_account(remaining_accounts, r, mappings);
                    let mut data = &stable_coin_mint.try_borrow_data()?[..];
                    let stable_coin_info = Mint::try_deserialize(&mut data)?;

                    let gringotts_stable_coin_account = get_account(remaining_accounts, r + 1, mappings);
                    let recipient_stable_coin_account = get_account(remaining_accounts, r + 2, mappings);

                    init_token_account_if_needed(
                        &recipient_stable_coin_account.to_account_info(),
                        &recipient.to_account_info(),
                        stable_coin_mint,
                        vault,
                        vault_bump,
                        system_program,
                        token_program,
                        associated_token_program,
                    )?;

                    let cpi_ctx = CpiContext::new_with_signer(
                        token_program.to_account_info(),
                        Transfer {
                            from: gringotts_stable_coin_account.to_account_info(),
                            to: recipient_stable_coin_account.to_account_info(),
                            authority: gringotts.to_account_info(),
                        },
                        signer_seeds,
                    );

                    token::transfer(
                        cpi_ctx,
                        change_decimals(item.amount_usdx, CHAIN_TRANSFER_DECIMALS, stable_coin_info.decimals),
                    )?;

                    r += 3;
                } else if *item.asset != [0; 32] {
                    let stable_coin_mint = get_account(remaining_accounts, r, mappings);
                    let mut data = &stable_coin_mint.try_borrow_data()?[..];
                    let stable_coin_mint_info = Mint::try_deserialize(&mut data)?;

                    let gringotts_stable_coin_account = get_account(remaining_accounts, r + 1, mappings);
                    let mut data = &gringotts_stable_coin_account.try_borrow_data()?[..];
                    let stable_coin_info = TokenAccount::try_deserialize(&mut data)?;

                    let desired_token_mint = get_account(remaining_accounts, r + 2, mappings);
                    let recipient_desired_token_account = get_account(remaining_accounts, r + 3, mappings);
                    let swap_program = get_account(remaining_accounts, r + 4, mappings);
                    let recipient_stable_coin_account = get_account(remaining_accounts, r + 5, mappings);

                    r = r + 6;

                    init_token_account_if_needed(
                        &recipient_desired_token_account.to_account_info(),
                        &recipient.to_account_info(),
                        desired_token_mint,
                        vault,
                        vault_bump,
                        system_program,
                        token_program,
                        associated_token_program,
                    )?;

                    let swap_account_counts = item.metadata[0] as usize;
                    let swap_accounts = get_accounts(remaining_accounts, r, r + swap_account_counts, mappings);

                    let before_swap_stable_coin_amount = stable_coin_info.amount;

                    let res = swap_on_jupiter(
                        gringotts,
                        swap_program,
                        swap_accounts.as_slice(),
                        item.command.to_vec(),
                    );

                    r = r + swap_account_counts;

                    if res.is_ok() {
                        let mut data = &gringotts_stable_coin_account.try_borrow_data()?[..];
                        let stable_coin_info = TokenAccount::try_deserialize(&mut data)?;
                        let swap_stable_coin_use = before_swap_stable_coin_amount - stable_coin_info.amount;

                        require!(
                            change_decimals(swap_stable_coin_use, stable_coin_mint_info.decimals, CHAIN_TRANSFER_DECIMALS) <= item.amount_usdx,
                            LzReceiveErrorCode::InvalidParams
                        );
                    } else {
                        init_token_account_if_needed(
                            recipient_stable_coin_account,
                            &recipient.to_account_info(),
                            stable_coin_mint,
                            vault,
                            vault_bump,
                            system_program,
                            token_program,
                            associated_token_program,
                        )?;

                        let cpi_ctx = CpiContext::new_with_signer(
                            token_program.to_account_info(),
                            Transfer {
                                from: gringotts_stable_coin_account.to_account_info(),
                                to: recipient_stable_coin_account.to_account_info(),
                                authority: gringotts.to_account_info(),
                            },
                            signer_seeds,
                        );

                        token::transfer(
                            cpi_ctx,
                            change_decimals(item.amount_usdx, CHAIN_TRANSFER_DECIMALS, stable_coin_mint_info.decimals),
                        )?;
                    }
                } else {
                    let stable_coin_mint = get_account(remaining_accounts, r, mappings);
                    let mut data = &stable_coin_mint.try_borrow_data()?[..];
                    let stable_coin_mint_info = Mint::try_deserialize(&mut data)?;

                    let gringotts_stable_coin_account = get_account(remaining_accounts, r + 1, mappings);
                    let mut data = &gringotts_stable_coin_account.try_borrow_data()?[..];
                    let stable_coin_info = TokenAccount::try_deserialize(&mut data)?;

                    let native_mint = get_account(remaining_accounts, r + 2, mappings);
                    let gringotts_native_account = get_account(remaining_accounts, r + 3, mappings);
                    let swap_program = get_account(remaining_accounts, r + 4, mappings);
                    let recipient_stable_coin_account = get_account(remaining_accounts, r + 5, mappings);

                    r = r + 6;

                    init_token_account_if_needed(
                        &gringotts_native_account.to_account_info(),
                        &gringotts.to_account_info(),
                        native_mint,
                        vault,
                        vault_bump,
                        system_program,
                        token_program,
                        associated_token_program,
                    )?;

                    let swap_account_counts = item.metadata[0] as usize;
                    let swap_accounts = get_accounts(remaining_accounts, r, r + swap_account_counts, mappings);

                    let before_swap_stable_coin_amount = stable_coin_info.amount;

                    let res = swap_on_jupiter(
                        gringotts,
                        swap_program,
                        swap_accounts.as_slice(),
                        item.command.to_vec(),
                    );

                    r = r + swap_account_counts;

                    if res.is_ok() {
                        let mut data = &gringotts_stable_coin_account.try_borrow_data()?[..];
                        let stable_coin_info = TokenAccount::try_deserialize(&mut data)?;
                        let swap_stable_coin_use = before_swap_stable_coin_amount - stable_coin_info.amount;

                        require!(
                            change_decimals(swap_stable_coin_use, stable_coin_mint_info.decimals, CHAIN_TRANSFER_DECIMALS) <= item.amount_usdx,
                            LzReceiveErrorCode::InvalidParams
                        );

                        close_wsol_token(
                            gringotts,
                            &gringotts_native_account.to_account_info(),
                            recipient,
                            token_program,
                        )?;
                    } else {
                        init_token_account_if_needed(
                            recipient_stable_coin_account,
                            &recipient.to_account_info(),
                            stable_coin_mint,
                            vault,
                            vault_bump,
                            system_program,
                            token_program,
                            associated_token_program,
                        )?;

                        let cpi_ctx = CpiContext::new_with_signer(
                            token_program.to_account_info(),
                            Transfer {
                                from: gringotts_stable_coin_account.to_account_info(),
                                to: recipient_stable_coin_account.to_account_info(),
                                authority: gringotts.to_account_info(),
                            },
                            signer_seeds,
                        );

                        token::transfer(
                            cpi_ctx,
                            change_decimals(item.amount_usdx, CHAIN_TRANSFER_DECIMALS, stable_coin_mint_info.decimals),
                        )?;
                    }
                }
            }
        }

        let accounts_for_clear = &ctx.remaining_accounts[r..r + Clear::MIN_ACCOUNTS_LEN];
        let _ = oapp::endpoint_cpi::clear(
            gringotts.lz_endpoint_program,
            gringotts.key(),
            accounts_for_clear,
            seeds,
            ClearParams {
                receiver: gringotts.key(),
                src_eid: params.src_eid,
                sender: params.sender,
                nonce: params.nonce,
                guid: params.guid,
                message: params.message.clone(),
            },
        )?;

        Ok(())
    }
}

fn close_wsol_token<'info>(
    gringotts: &Account<'info, Gringotts>,
    token_account: &AccountInfo<'info>,
    recipient: &AccountInfo<'info>,
    token_program: &AccountInfo<'info>,
) -> Result<()> {
    let gringotts_seeds: &[&[u8]] = &[GRINGOTTS_SEED, &[gringotts.bump]];
    let signer_seeds = &[gringotts_seeds];

    token::sync_native(CpiContext::new_with_signer(
        token_program.to_account_info(),
        token::SyncNative {
            account: token_account.to_account_info(),
        },
        signer_seeds,
    ))?;

    token::close_account(CpiContext::new_with_signer(
        token_program.to_account_info(),
        CloseAccount {
            account: token_account.to_account_info(),
            destination: recipient.to_account_info(),
            authority: gringotts.to_account_info(),
        },
        signer_seeds,
    ))?;

    Ok(())
}

fn init_token_account_if_needed<'info>(
    token_account: &AccountInfo<'info>,
    authority: &AccountInfo<'info>,
    token_mint: &AccountInfo<'info>,
    vault: &AccountInfo<'info>,
    vault_bump: u8,
    system_program: &AccountInfo<'info>,
    token_program: &AccountInfo<'info>,
    associated_token_program: &AccountInfo<'info>,
) -> Result<()> {
    let ata = associated_token::get_associated_token_address(&authority.key(), &token_mint.key());

    require!(
        ata.key() == token_account.key(), BridgeErrorCode::InvalidParams
    );

    if token_account.data_is_empty() {
        let cpi_accounts = associated_token::Create {
            payer: vault.to_account_info(),
            associated_token: token_account.clone(),
            authority: authority.to_account_info(),
            mint: token_mint.to_account_info(),
            system_program: system_program.to_account_info(),
            token_program: token_program.to_account_info(),
        };

        let vault_seeds: &[&[u8]] = &[VAULT_SEED, &[vault_bump]];
        let signer_seeds = &[vault_seeds];

        let cpi_context = CpiContext::new_with_signer(
            associated_token_program.to_account_info(),
            cpi_accounts,
            signer_seeds,
        );
        associated_token::create(cpi_context)?;
    }

    Ok(())
}

fn get_account<'b, 'info>(
    remaining_accounts: &'b [AccountInfo<'info>],
    index: usize,
    mapping: &[u8],
) -> &'b AccountInfo<'info> {
    &remaining_accounts[mapping[index] as usize]
}

fn get_accounts<'a, 'info>(
    remaining_accounts: &'a [AccountInfo<'info>],
    start_index: usize,
    end_index: usize,
    mapping: &[u8],
) -> Vec<AccountInfo<'info>> {
    (start_index..end_index).map(|i| get_account(remaining_accounts, i, mapping).clone()).collect()
}

#[error_code]
pub enum LzReceiveErrorCode {
    #[msg("Invalid params.")]
    InvalidParams,
}
