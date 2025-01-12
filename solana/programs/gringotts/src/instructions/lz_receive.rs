use crate::constants::{CHAIN_TRANSFER_DECIMALS, NATIVE_MINT};
use crate::msg_codec::{ChainTransfer, Message, SolanaTransfer, CHAIN_COMPLETION_TYPE, CHAIN_REGISTER_ORDER_TYPE, CHAIN_TRANSFER_TYPE};
use crate::state::{Gringotts, LzReceiveTypesAccounts, Order, OrderState};
use crate::utils::change_decimals;
use crate::*;
use anchor_lang::system_program;
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

    #[account(mut, seeds = [VAULT_SEED], bump)]
    pub vault: SystemAccount<'info>,

    #[account(mut,seeds = [LZ_RECEIVE_TYPES_SEED, &gringotts.key().to_bytes()], bump)]
    pub lz_receive_types_accounts: Account<'info, LzReceiveTypesAccounts>,

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

        let message = Message::decode(params.message.as_slice());

        let vault = &ctx.accounts.vault;
        let vault_bump = ctx.bumps.vault;

        let lz_receive_types_accounts = &mut ctx.accounts.lz_receive_types_accounts;
        let system_program = &ctx.accounts.system_program;
        let mut clear_lz_accounts = true;

        if message.header == CHAIN_REGISTER_ORDER_TYPE {
            let order_account = &ctx.remaining_accounts[0];

            // message.payload is Message
            create_order(
                order_account,
                vault,
                vault_bump,
                system_program,
                params.message.clone(),
                &params.guid,
                params.nonce,
            )?;

            lz_receive_types_accounts.add_fc_item(order_account.key());
        } else if message.header == CHAIN_TRANSFER_TYPE {
            let order_item = &mut Account::<Order>::try_from(&ctx.remaining_accounts[0])?;
            require!(lz_receive_types_accounts.contains_fc_item(order_item.key()), LzReceiveErrorCode::InvalidParams);
            order_item.state = OrderState::Processed;

            clear_lz_accounts = false;

            let associated_token_program = &Program::<AssociatedToken>::try_from(&ctx.remaining_accounts[1])?;
            let token_program = &Program::<Token>::try_from(&ctx.remaining_accounts[2])?;
            let jupiter = &ctx.remaining_accounts[3];

            let stable_coin_mint = &Account::<Mint>::try_from(&ctx.remaining_accounts[4])?;
            let gringotts_stable_coin_account = &mut Account::<TokenAccount>::try_from(&ctx.remaining_accounts[5])?;

            // It should be as same as `addressMap` in oracle
            let virtual_accounts: &Vec<AccountInfo<'info>> = &[
                vec![
                    gringotts.to_account_info(),
                    associated_token_program.to_account_info(),
                    token_program.to_account_info(),
                    system_program.to_account_info(),
                    jupiter.to_account_info(),
                    stable_coin_mint.to_account_info(),
                    gringotts_stable_coin_account.to_account_info(),
                ],
                ctx.remaining_accounts[6..].to_vec(),
            ].concat();

            let original_message = Message::decode(order_item.message.as_slice());
            let chain_transfer = ChainTransfer::decode(original_message.message);
            let solana_transfer = SolanaTransfer::decode(chain_transfer.message, chain_transfer.items.len());

            for i in 0..chain_transfer.items.len() {
                let chain_item = &chain_transfer.items[i];
                let solana_item = &solana_transfer.items[i];

                let desired_token_mint_account = get_account(virtual_accounts, solana_item.mint_index);
                let desired_token_mint_info = Mint::try_deserialize(&mut &desired_token_mint_account.try_borrow_data()?[..])?;

                let recipient = get_account(virtual_accounts, solana_item.recipient_index);
                /* if user wants native SOL, this will belong to gringotts - WSOL */
                let recipient_or_gringotts_desired_token_account = get_account(virtual_accounts, solana_item.recipient_token_index);
                let authority = if desired_token_mint_account.key() != NATIVE_MINT { recipient.to_account_info() } else { gringotts.to_account_info() };

                init_token_account_if_needed(
                    &recipient_or_gringotts_desired_token_account.to_account_info(),
                    &authority,
                    desired_token_mint_account,
                    vault,
                    vault_bump,
                    associated_token_program,
                    token_program,
                    system_program,
                )?;

                if gringotts.has_stable_coin(desired_token_mint_account.key().to_bytes()) {
                    let cpi_ctx = CpiContext::new_with_signer(
                        token_program.to_account_info(),
                        Transfer {
                            from: gringotts_stable_coin_account.to_account_info(),
                            to: recipient_or_gringotts_desired_token_account.to_account_info(),
                            authority: gringotts.to_account_info(),
                        },
                        signer_seeds,
                    );

                    token::transfer(
                        cpi_ctx,
                        change_decimals(chain_item.amount_usdx, CHAIN_TRANSFER_DECIMALS, desired_token_mint_info.decimals),
                    )?;
                } else {
                    let swap_accounts = get_accounts(virtual_accounts, solana_item.swap_accounts);
                    let before_swap_stable_coin_amount = gringotts_stable_coin_account.amount;

                    let _ = swap_on_jupiter(
                        gringotts,
                        jupiter,
                        swap_accounts.as_slice(),
                        solana_item.swap_command.to_vec(),
                    );

                    gringotts_stable_coin_account.reload()?;
                    let swap_stable_coin_use = before_swap_stable_coin_amount - gringotts_stable_coin_account.amount;

                    require!(
                            change_decimals(swap_stable_coin_use, stable_coin_mint.decimals, CHAIN_TRANSFER_DECIMALS) <= chain_item.amount_usdx,
                            LzReceiveErrorCode::InvalidParams
                        );

                    if desired_token_mint_account.key() == NATIVE_MINT {
                        close_wsol_token(
                            gringotts,
                            recipient_or_gringotts_desired_token_account,
                            recipient,
                            token_program,
                        )?;
                    }
                }
            }
        } else if message.header == CHAIN_COMPLETION_TYPE {
            let order = &Account::<Order>::try_from(&ctx.remaining_accounts[0])?;
            require!(lz_receive_types_accounts.remove_fc_item(order.key()), LzReceiveErrorCode::InvalidParams);

            // The swap/transfer has been failed, need to do fallback
            if order.state == OrderState::None {
                let original_message = Message::decode(order.message.as_slice());
                let chain_transfer = ChainTransfer::decode(original_message.message);

                let associated_token_program = &Program::<AssociatedToken>::try_from(&ctx.remaining_accounts[1])?;
                let token_program = &Program::<Token>::try_from(&ctx.remaining_accounts[2])?;

                let stable_coin_mint = &Account::<Mint>::try_from(&ctx.remaining_accounts[3])?;
                let gringotts_stable_coin_account = &Account::<TokenAccount>::try_from(&ctx.remaining_accounts[4])?;

                let mut r_index = 5;

                for chain_item in chain_transfer.items.iter() {
                    fallback(
                        chain_item.amount_usdx,
                        gringotts,
                        &gringotts_stable_coin_account.to_account_info(),
                        &ctx.remaining_accounts[r_index],
                        &ctx.remaining_accounts[r_index + 1],
                        stable_coin_mint,
                        vault,
                        vault_bump,
                        associated_token_program,
                        token_program,
                        system_program,
                    )?;

                    r_index += 2;
                }

                let accounts_for_clear = &ctx.remaining_accounts[ctx.remaining_accounts.len() - 2 * Clear::MIN_ACCOUNTS_LEN..];
                let _ = oapp::endpoint_cpi::clear(
                    gringotts.lz_endpoint_program,
                    gringotts.key(),
                    accounts_for_clear,
                    seeds,
                    ClearParams {
                        receiver: gringotts.key(),
                        src_eid: params.src_eid,
                        sender: params.sender,
                        nonce: order.nonce,
                        guid: order.message_id,
                        message: order.message.to_vec(),
                    },
                )?;
            }

            close_order_account(&order.to_account_info(), vault)?;
        } else {
            require!(false, LzReceiveErrorCode::InvalidParams);
        }

        if clear_lz_accounts {
            let accounts_for_clear = &ctx.remaining_accounts[ctx.remaining_accounts.len() - Clear::MIN_ACCOUNTS_LEN..];
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
        }

        Ok(())
    }
}

fn fallback<'info>(
    amount_usdx: u64,
    gringotts: &Account<'info, Gringotts>,
    gringotts_stable_coin_account: &AccountInfo<'info>,
    recipient: &AccountInfo<'info>,
    recipient_stable_coin_account: &AccountInfo<'info>,
    stable_coin_mint: &Account<'info, Mint>,
    vault: &AccountInfo<'info>,
    vault_bump: u8,
    associated_token_program: &Program<'info, AssociatedToken>,
    token_program: &Program<'info, Token>,
    system_program: &Program<'info, System>,
) -> Result<()> {
    let gringotts_seeds: &[&[u8]] = &[GRINGOTTS_SEED, &[gringotts.bump]];
    let signer_seeds = &[gringotts_seeds];

    init_token_account_if_needed(
        recipient_stable_coin_account,
        &recipient.to_account_info(),
        &stable_coin_mint.to_account_info(),
        vault,
        vault_bump,
        associated_token_program,
        token_program,
        system_program,
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
        change_decimals(amount_usdx, CHAIN_TRANSFER_DECIMALS, stable_coin_mint.decimals),
    )?;

    Ok(())
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
    token_mint: &AccountInfo<'info>,
    vault: &AccountInfo<'info>,
    vault_bump: u8,
    associated_token_program: &Program<'info, AssociatedToken>,
    token_program: &Program<'info, Token>,
    system_program: &Program<'info, System>,
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

pub fn create_order<'info>(
    order_account: &AccountInfo<'info>,
    vault_account: &AccountInfo<'info>,
    vault_bump: u8,
    system_program: &AccountInfo<'info>,
    message: Vec<u8>,
    message_id: &[u8; 32],
    nonce: u64,
) -> Result<()> {
    let (_, order_bump) = Pubkey::find_program_address(
        &[ORDER_SEED, message_id],
        &id(),
    );

    let space = 8 + Order::INIT_SPACE;
    let rent = Rent::get()?.minimum_balance(space);

    let vault_seeds: &[&[u8]] = &[VAULT_SEED, &[vault_bump]];
    let order_seeds: &[&[u8]] = &[ORDER_SEED, message_id, &[order_bump]];
    let signer_seeds = &[vault_seeds, order_seeds];

    let cpi_context = CpiContext::new_with_signer(
        system_program.to_account_info(),
        system_program::CreateAccount {
            from: vault_account.to_account_info(),
            to: order_account.to_account_info(),
        },
        signer_seeds,
    );

    system_program::create_account(cpi_context, rent, space as u64, &id())?;

    let order = Order {
        message: message,
        message_id: message_id.clone(),
        nonce: nonce,
        timestamp: Clock::get()?.unix_timestamp as u64,
        state: OrderState::None,
    };

    let mut data = order_account.data.borrow_mut();
    order.try_serialize(&mut *data)?;

    Ok(())
}

pub fn close_order_account<'info>(
    order_account: &AccountInfo<'info>,
    vault_account: &AccountInfo<'info>,
) -> Result<()> {
    let lamports = order_account.to_account_info().lamports();
    **order_account.to_account_info().lamports.borrow_mut() = 0;
    **vault_account.to_account_info().lamports.borrow_mut() += lamports;

    Ok(())
}

fn get_account<'b, 'info>(
    remaining_accounts: &'b [AccountInfo<'info>],
    index: usize,
) -> &'b AccountInfo<'info> {
    &remaining_accounts[index]
}

fn get_accounts<'a, 'info>(
    remaining_accounts: &'a [AccountInfo<'info>],
    indexes: &'a [u8],
) -> Vec<AccountInfo<'info>> {
    indexes.iter().map(|i| get_account(remaining_accounts, *i as usize).clone()).collect()
}

#[error_code]
pub enum LzReceiveErrorCode {
    #[msg("Invalid params.")]
    InvalidParams,
}
