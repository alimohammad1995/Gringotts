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
use std::fmt::Pointer;

#[derive(Accounts)]
#[instruction(params: LzReceiveParams)]
pub struct LzReceive<'info> {
    #[account(
        seeds = [PEER_SEED, &gringotts.lz_eid.to_le_bytes()], bump = self_peer.bump,
    )]
    pub self_peer: Account<'info, Peer>,
    #[account(
        seeds = [PEER_SEED, &params.src_eid.to_le_bytes()], bump = peer.bump,
        constraint = params.sender == peer.address
    )]
    pub peer: Account<'info, Peer>,

    #[account(seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,
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
        let seeds: &[&[u8]] = &[GRINGOTTS_SEED, &[ctx.accounts.gringotts.bump]];
        let signer_seeds = &[seeds];

        let gringotts = &ctx.accounts.gringotts;

        let self_peer = &ctx.accounts.self_peer;
        let peer = &ctx.accounts.peer;

        let vault = &ctx.accounts.vault;
        let vault_bump = ctx.bumps.vault;

        let remaining_accounts = ctx.remaining_accounts;

        let system_program = &ctx.accounts.system_program;
        let token_program = &ctx.accounts.token_program;
        let associated_token_program = &ctx.accounts.associated_token_program;

        let mut r = 0;

        let message = Message::decode(params.message.as_slice());

        if message.header == CHAIN_TRANSFER_TYPE {
            let chain_transfer = ChainTransfer::decode(message.payload);

            for item in &chain_transfer.items {
                let recipient = &remaining_accounts[r];
                let mint: &Account<Mint> = &Account::try_from(&remaining_accounts[r + 1])?;
                let recipient_or_gringotts_token_account = &remaining_accounts[r + 2];

                require!(
                    *item.asset == mint.key().to_bytes(),
                    LzReceiveErrorCode::InvalidParams
                );

                init_token_account_if_needed(
                    recipient_or_gringotts_token_account,
                    recipient,
                    mint,
                    vault,
                    vault_bump,
                    system_program,
                    token_program,
                    associated_token_program,
                )?;

                r += 3;

                if self_peer.stable_coins.contains(&item.asset) {
                    let gringotts_stable_coin = &remaining_accounts[r];

                    let cpi_ctx = CpiContext::new_with_signer(
                        token_program.to_account_info(),
                        Transfer {
                            from: gringotts_stable_coin.to_account_info(),
                            to: recipient_or_gringotts_token_account.to_account_info(),
                            authority: gringotts.to_account_info(),
                        },
                        signer_seeds,
                    );

                    token::transfer(
                        cpi_ctx,
                        change_decimals(item.amount_usdx, CHAIN_TRANSFER_DECIMALS, mint.decimals),
                    )?;

                    r += 1;
                } else {
                    let swap_program = &remaining_accounts[r];
                    let stable_coin_mint: Account<Mint> = Account::try_from(&remaining_accounts[r + 1])?;
                    let mut gringotts_stable_coin: Account<TokenAccount> = Account::try_from(&remaining_accounts[r + 2])?;
                    let recipient_stable_coin = &remaining_accounts[r + 3];
                    r = r + 4;

                    let swap_account_counts = item.metadata[0] as usize;
                    let before_swap_stable_coin_amount = gringotts_stable_coin.amount;

                    let res = swap_on_jupiter(
                        gringotts,
                        swap_program,
                        &remaining_accounts[r..r + swap_account_counts],
                        item.command.to_vec(),
                    );

                    r = r + swap_account_counts;

                    if res.is_ok() {
                        gringotts_stable_coin.reload()?;
                        let swap_stable_coin_use = gringotts_stable_coin.amount - before_swap_stable_coin_amount;

                        require!(
                            change_decimals(swap_stable_coin_use, stable_coin_mint.decimals, CHAIN_TRANSFER_DECIMALS) >= item.amount_usdx,
                            LzReceiveErrorCode::InvalidParams
                        );

                        if *item.asset == [0; 32] {
                            close_wsol_token(
                                gringotts,
                                recipient_or_gringotts_token_account,
                                recipient,
                                token_program,
                            )?;
                        }
                    } else {
                        let cpi_ctx = CpiContext::new_with_signer(
                            token_program.to_account_info(),
                            Transfer {
                                from: gringotts_stable_coin.to_account_info(),
                                to: recipient_stable_coin.to_account_info(),
                                authority: gringotts.to_account_info(),
                            },
                            signer_seeds,
                        );

                        token::transfer(
                            cpi_ctx,
                            change_decimals(
                                item.amount_usdx,
                                CHAIN_TRANSFER_DECIMALS,
                                stable_coin_mint.decimals,
                            ),
                        )?;
                    }
                }
            }
        }

        let accounts_for_clear = &ctx.remaining_accounts[r..r + Clear::MIN_ACCOUNTS_LEN];
        let _ = oapp::endpoint_cpi::clear(
            ctx.accounts.gringotts.lz_endpoint_program,
            ctx.accounts.gringotts.key(),
            accounts_for_clear,
            seeds,
            ClearParams {
                receiver: ctx.accounts.gringotts.key(),
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
    token_program: &Program<'info, Token>,
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
    token_mint: &Account<'info, Mint>,
    vault: &AccountInfo<'info>,
    vault_bump: u8,
    system_program: &Program<'info, System>,
    token_program: &Program<'info, Token>,
    associated_token_program: &Program<'info, AssociatedToken>,
) -> Result<()> {
    let ata = associated_token::get_associated_token_address(&authority.key(), &token_mint.key());

    require!(
        ata.key() == token_account.key(),
        BridgeErrorCode::InvalidParams
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

#[error_code]
pub enum LzReceiveErrorCode {
    #[msg("Invalid params.")]
    InvalidParams,
}
