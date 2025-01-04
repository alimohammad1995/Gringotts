use crate::msg_codec::{ChainTransfer, Message, SolanaTransfer, CHAIN_TRANSFER_TYPE};
use crate::state::{Gringotts, Peer};
use crate::utils::{change_decimals, NATIVE_MINT};
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

    #[account(seeds = [PEER_SEED, &gringotts.lz_eid.to_be_bytes()], bump = self_peer.bump)]
    pub self_peer: Account<'info, Peer>,

    #[account( seeds = [PEER_SEED, &params.src_eid.to_be_bytes()], bump = peer.bump, constraint = params.sender == peer.address)]
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

        let remaining_accounts: &Vec<AccountInfo<'info>> = &[
            vec![
                ctx.accounts.gringotts.to_account_info(),
                ctx.accounts.vault.to_account_info(),
                ctx.accounts.associated_token_program.to_account_info(),
                ctx.accounts.token_program.to_account_info(),
                ctx.accounts.system_program.to_account_info(),
            ],
            ctx.remaining_accounts.to_vec(),
        ].concat();

        let mut r = 0;

        let message = Message::decode(params.message.as_slice());

        if message.header == CHAIN_TRANSFER_TYPE {
            let chain_transfer = ChainTransfer::decode(message.payload);
            let solana_transfer = SolanaTransfer::decode(chain_transfer.message, chain_transfer.items.len());

            let mappings = solana_transfer.accounts_mapping;

            let vault = &ctx.accounts.vault;
            let vault_bump = ctx.bumps.vault;

            let system_program = &ctx.accounts.system_program.to_account_info();
            let token_program = &ctx.accounts.token_program.to_account_info();
            let associated_token_program = &ctx.accounts.associated_token_program.to_account_info();

            for i in 0..chain_transfer.items.len() {
                let chain_item = &chain_transfer.items[i];

                let recipient = get_account(remaining_accounts, r, mappings);

                let desired_token_mint_account = get_account(remaining_accounts, r + 1, mappings);
                let desired_token_mint_info = Mint::try_deserialize(&mut &desired_token_mint_account.try_borrow_data()?[..])?;

                r += 2;

                if self_peer.has_stable_coin(desired_token_mint_account.key().to_bytes()) {
                    let gringotts_stable_coin_account = get_account(remaining_accounts, r, mappings);
                    let recipient_stable_coin_account = get_account(remaining_accounts, r + 1, mappings);

                    init_token_account_if_needed(
                        &recipient_stable_coin_account.to_account_info(),
                        &recipient.to_account_info(),
                        desired_token_mint_account,
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
                        change_decimals(chain_item.amount_usdx, CHAIN_TRANSFER_DECIMALS, desired_token_mint_info.decimals),
                    )?;

                    r += 2;
                } else if desired_token_mint_account.key() != NATIVE_MINT {
                    let solana_item = solana_transfer.items[i].clone().unwrap();

                    let recipient_desired_token_account = get_account(remaining_accounts, r, mappings);
                    let recipient_stable_coin_account = get_account(remaining_accounts, r + 1, mappings);

                    let stable_coin_mint_account = get_account(remaining_accounts, r + 2, mappings);
                    let stable_coin_mint_info = Mint::try_deserialize(&mut &stable_coin_mint_account.try_borrow_data()?[..])?;

                    let gringotts_stable_coin_account = get_account(remaining_accounts, r + 3, mappings);
                    let gringotts_stable_coin_info = TokenAccount::try_deserialize(&mut &gringotts_stable_coin_account.try_borrow_data()?[..])?;

                    let swap_program = get_account(remaining_accounts, r + 4, mappings);

                    r = r + 5;

                    init_token_account_if_needed(
                        &recipient_desired_token_account.to_account_info(),
                        &recipient.to_account_info(),
                        desired_token_mint_account,
                        vault,
                        vault_bump,
                        system_program,
                        token_program,
                        associated_token_program,
                    )?;

                    let swap_account_counts = solana_item.swap_account_count;
                    let swap_accounts = get_accounts(remaining_accounts, r, r + swap_account_counts as usize, mappings);

                    let before_swap_stable_coin_amount = gringotts_stable_coin_info.amount;

                    let res = swap_on_jupiter(
                        gringotts,
                        swap_program,
                        swap_accounts.as_slice(),
                        solana_item.swap_command.to_vec(),
                    );

                    r = r + swap_account_counts;

                    if res.is_ok() {
                        let gringotts_stable_coin_info = TokenAccount::try_deserialize(&mut &gringotts_stable_coin_account.try_borrow_data()?[..])?;
                        let swap_stable_coin_use = before_swap_stable_coin_amount - gringotts_stable_coin_info.amount;

                        require!(
                            change_decimals(swap_stable_coin_use, stable_coin_mint_info.decimals, CHAIN_TRANSFER_DECIMALS) <= chain_item.amount_usdx,
                            LzReceiveErrorCode::InvalidParams
                        );
                    } else {
                        fallback(
                            chain_item.amount_usdx,
                            gringotts,
                            gringotts_stable_coin_account,
                            recipient,
                            recipient_stable_coin_account,
                            stable_coin_mint_account,
                            &stable_coin_mint_info,
                            vault,
                            vault_bump,
                            system_program,
                            token_program,
                            associated_token_program,
                        )?;
                    }
                } else {
                    let solana_item = solana_transfer.items[i].clone().unwrap();

                    let gringotts_native_account = get_account(remaining_accounts, r, mappings);
                    let recipient_stable_coin_account = get_account(remaining_accounts, r + 1, mappings);

                    let stable_coin_mint_account = get_account(remaining_accounts, r + 2, mappings);
                    let stable_coin_mint_info = Mint::try_deserialize(&mut &stable_coin_mint_account.try_borrow_data()?[..])?;

                    let gringotts_stable_coin_account = get_account(remaining_accounts, r + 3, mappings);
                    let gringotts_stable_coin_info = TokenAccount::try_deserialize(&mut &gringotts_stable_coin_account.try_borrow_data()?[..])?;

                    let swap_program = get_account(remaining_accounts, r + 4, mappings);

                    r = r + 6;

                    init_token_account_if_needed(
                        &gringotts_native_account.to_account_info(),
                        &gringotts.to_account_info(),
                        desired_token_mint_account,
                        vault,
                        vault_bump,
                        system_program,
                        token_program,
                        associated_token_program,
                    )?;

                    let swap_account_counts = solana_item.swap_account_count;
                    let swap_accounts = get_accounts(remaining_accounts, r, r + swap_account_counts as usize, mappings);

                    let before_swap_stable_coin_amount = gringotts_stable_coin_info.amount;

                    let res = swap_on_jupiter(
                        gringotts,
                        swap_program,
                        swap_accounts.as_slice(),
                        solana_item.swap_command.to_vec(),
                    );

                    r = r + swap_account_counts as usize;

                    if res.is_ok() {
                        let gringotts_stable_coin_info = TokenAccount::try_deserialize(&mut &gringotts_stable_coin_account.try_borrow_data()?[..])?;
                        let swap_stable_coin_use = before_swap_stable_coin_amount - gringotts_stable_coin_info.amount;

                        require!(
                            change_decimals(swap_stable_coin_use, stable_coin_mint_info.decimals, CHAIN_TRANSFER_DECIMALS) <= chain_item.amount_usdx,
                            LzReceiveErrorCode::InvalidParams
                        );

                        close_wsol_token(
                            gringotts,
                            &gringotts_native_account.to_account_info(),
                            recipient,
                            token_program,
                        )?;
                    } else {
                        fallback(
                            chain_item.amount_usdx,
                            gringotts,
                            gringotts_stable_coin_account,
                            recipient,
                            recipient_stable_coin_account,
                            stable_coin_mint_account,
                            &stable_coin_mint_info,
                            vault,
                            vault_bump,
                            system_program,
                            token_program,
                            associated_token_program,
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

fn fallback<'info>(
    amount_usdx: u64,
    gringotts: &Account<'info, Gringotts>,
    gringotts_stable_coin_account: &AccountInfo<'info>,
    recipient: &AccountInfo<'info>,
    recipient_stable_coin_account: &AccountInfo<'info>,
    stable_coin_mint_account: &AccountInfo<'info>,
    stable_coin_mint_info: &Mint,
    vault: &AccountInfo<'info>,
    vault_bump: u8,
    system_program: &AccountInfo<'info>,
    token_program: &AccountInfo<'info>,
    associated_token_program: &AccountInfo<'info>,
) -> Result<()> {
    let gringotts_seeds: &[&[u8]] = &[GRINGOTTS_SEED, &[gringotts.bump]];
    let signer_seeds = &[gringotts_seeds];

    init_token_account_if_needed(
        recipient_stable_coin_account,
        &recipient.to_account_info(),
        stable_coin_mint_account,
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
        change_decimals(amount_usdx, CHAIN_TRANSFER_DECIMALS, stable_coin_mint_info.decimals),
    )?;

    Ok(())
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
