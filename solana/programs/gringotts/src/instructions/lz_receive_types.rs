use crate::constants::{JUPITER, NATIVE_MINT};
use crate::msg_codec::{ChainExecution, ChainCompletion, ChainTransfer, Message, SolanaTransfer, CHAIN_COMPLETION_TYPE, CHAIN_REGISTER_ORDER_TYPE, CHAIN_TRANSFER_TYPE};
use crate::state::{Gringotts, Order, OrderState};
use crate::*;
use anchor_spl::associated_token;
use anchor_spl::associated_token::AssociatedToken;
use anchor_spl::token::Token;
use oapp::endpoint_cpi::{get_accounts_for_clear, LzAccount};

#[derive(Accounts)]
pub struct LzReceiveTypes<'info> {
    #[account(seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,
}

impl<'info> LzReceiveTypes<'info> {
    pub fn apply(
        ctx: &Context<'_, '_, 'info, 'info, LzReceiveTypes<'info>>,
        params: &LzReceiveParams,
    ) -> Result<Vec<LzAccount>> {
        let gringotts = &ctx.accounts.gringotts;

        let (vault, _) = Pubkey::find_program_address(
            &[VAULT_SEED],
            ctx.program_id,
        );
        let (lz_receive_type, _) = Pubkey::find_program_address(
            &[LZ_RECEIVE_TYPES_SEED, &gringotts.key().to_bytes()],
            ctx.program_id,
        );

        let mut accounts = vec![
            LzAccount {
                pubkey: gringotts.key(),
                is_signer: false,
                is_writable: false,
            },
            LzAccount {
                pubkey: vault,
                is_signer: false,
                is_writable: true,
            },
            LzAccount {
                pubkey: lz_receive_type,
                is_signer: false,
                is_writable: true,
            },
            LzAccount {
                pubkey: System::id(),
                is_signer: false,
                is_writable: false,
            },
        ];

        let message = Message::decode(params.message.as_slice());

        if message.header == CHAIN_REGISTER_ORDER_TYPE {
            let (order_account, _) = Pubkey::find_program_address(
                &[ORDER_SEED, &params.guid],
                ctx.program_id,
            );

            accounts.push(
                LzAccount {
                    pubkey: order_account,
                    is_signer: false,
                    is_writable: true,
                },
            );

            let accounts_for_clear = get_accounts_for_clear(
                gringotts.lz_endpoint_program,
                &gringotts.key(),
                params.src_eid,
                &params.sender,
                params.nonce,
            );

            accounts.extend(accounts_for_clear);
        } else if message.header == CHAIN_TRANSFER_TYPE {
            let chain_execution = ChainExecution::decode(message.message);
            let (order_account_pubkey, _) = Pubkey::find_program_address(
                &[ORDER_SEED, chain_execution.message_id],
                ctx.program_id,
            );

            let order = &Account::<Order>::try_from(
                ctx.remaining_accounts.iter().find(|acc| acc.key() == order_account_pubkey).unwrap()
            )?;

            let original_message = Message::decode(order.message.as_slice());
            let chain_transfer = ChainTransfer::decode(original_message.message);
            let solana_transfer = SolanaTransfer::decode(chain_transfer.message, chain_transfer.items.len());

            let stable_coin_mint = Pubkey::from(gringotts.stable_coins[solana_transfer.stable_coin_index]);
            let gringotts_stable_coin_account = associated_token::get_associated_token_address(&gringotts.key(), &stable_coin_mint.key());

            accounts.extend(vec![
                LzAccount {
                    pubkey: order.key(),
                    is_signer: false,
                    is_writable: true,
                },
                LzAccount {
                    pubkey: AssociatedToken::id(),
                    is_signer: false,
                    is_writable: false,
                },
                LzAccount {
                    pubkey: Token::id(),
                    is_signer: false,
                    is_writable: false,
                },
                LzAccount {
                    pubkey: JUPITER,
                    is_signer: false,
                    is_writable: false,
                },
                LzAccount {
                    pubkey: stable_coin_mint,
                    is_signer: false,
                    is_writable: false,
                },
                LzAccount {
                    pubkey: gringotts_stable_coin_account,
                    is_signer: false,
                    is_writable: true,
                },
            ]);

            let gringotts_pubkey = &gringotts.key().to_bytes();
            let default_pubkey = &Pubkey::default().to_bytes();
            let stable_coin_mint_pubkey = &stable_coin_mint.to_bytes();
            let gringotts_stable_coin_account_pubkey = &gringotts_stable_coin_account.to_bytes();

            // It should be as same as `addressMap` in oracle
            let virtual_accounts = [
                vec![
                    gringotts_pubkey,
                    default_pubkey, // associated_token_program
                    default_pubkey, // token_program
                    default_pubkey, // system_program
                    default_pubkey, // jupiter
                    stable_coin_mint_pubkey,
                    gringotts_stable_coin_account_pubkey,
                ],
                solana_transfer.accounts_address.clone(),
            ].concat();

            for i in 0..solana_transfer.accounts_address.len() {
                let is_writeable_flag_index = i / 8;
                let is_writeable_flag_mask = 7 - (i % 8);

                accounts.push(LzAccount {
                    pubkey: Pubkey::new_from_array(*solana_transfer.accounts_address[i]),
                    is_writable: solana_transfer.accounts_flags[is_writeable_flag_index] & (1 << is_writeable_flag_mask) != 0,
                    is_signer: false,
                });
            }

            for item in solana_transfer.items.iter() {
                let mint_pubkey = Pubkey::from(*virtual_accounts[item.mint_index]);
                let recipient_pubkey = Pubkey::from(*virtual_accounts[item.recipient_index]);

                let recipient_or_gringotts_token_account = if mint_pubkey != NATIVE_MINT {
                    associated_token::get_associated_token_address(&recipient_pubkey, &mint_pubkey.key())
                } else {
                    associated_token::get_associated_token_address(&gringotts.key(), &mint_pubkey.key())
                };

                let found_recipient_or_gringotts_token_account = accounts.iter().any(|acc| acc.pubkey == recipient_or_gringotts_token_account);

                if !found_recipient_or_gringotts_token_account {
                    accounts.push(LzAccount {
                        pubkey: recipient_or_gringotts_token_account,
                        is_writable: true,
                        is_signer: false,
                    });
                }
            }
        } else if message.header == CHAIN_COMPLETION_TYPE {
            let chain_completion = ChainCompletion::decode(message.message);
            let (order_account_pubkey, _) = Pubkey::find_program_address(
                &[ORDER_SEED, chain_completion.message_id],
                ctx.program_id,
            );

            let order = &Account::<Order>::try_from(
                ctx.remaining_accounts.iter().find(|acc| acc.key() == order_account_pubkey).unwrap()
            )?;

            let original_message = Message::decode(order.message.as_slice());
            let chain_transfer = ChainTransfer::decode(original_message.message);
            let solana_transfer = SolanaTransfer::decode(chain_transfer.message, chain_transfer.items.len());

            let stable_coin_mint = Pubkey::from(gringotts.stable_coins[solana_transfer.stable_coin_index]);
            let gringotts_stable_coin_account = associated_token::get_associated_token_address(&gringotts.key(), &stable_coin_mint.key());

            accounts.extend(vec![
                LzAccount {
                    pubkey: order.key(),
                    is_signer: false,
                    is_writable: true,
                },
                LzAccount {
                    pubkey: AssociatedToken::id(),
                    is_signer: false,
                    is_writable: false,
                },
                LzAccount {
                    pubkey: Token::id(),
                    is_signer: false,
                    is_writable: false,
                },
                LzAccount {
                    pubkey: stable_coin_mint,
                    is_signer: false,
                    is_writable: false,
                },
                LzAccount {
                    pubkey: gringotts_stable_coin_account,
                    is_signer: false,
                    is_writable: true,
                },
            ]);

            if order.state == OrderState::None {
                let original_message = Message::decode(order.message.as_slice());
                let chain_transfer = ChainTransfer::decode(original_message.message);
                let solana_transfer = SolanaTransfer::decode(chain_transfer.message, chain_transfer.items.len());

                let default_pubkey = &Pubkey::default().to_bytes();

                // It should be as same as `addressMap` in oracle
                let virtual_accounts = [
                    vec![
                        default_pubkey, // gringotts
                        default_pubkey, // associated_token_program
                        default_pubkey, // token_program
                        default_pubkey, // system_program
                        default_pubkey, // jupiter
                        default_pubkey, // stable_coin_mint_pubkey
                        default_pubkey, // gringotts_stable_coin_account_pubkey
                    ],
                    solana_transfer.accounts_address.clone(),
                ].concat();

                for item in solana_transfer.items.iter() {
                    let recipient_pubkey = Pubkey::from(*virtual_accounts[item.recipient_index]);
                    let recipient_stable_coin_account = associated_token::get_associated_token_address(&recipient_pubkey, &stable_coin_mint.key());

                    accounts.extend(vec![
                        LzAccount {
                            pubkey: recipient_pubkey,
                            is_writable: false,
                            is_signer: false,
                        },
                        LzAccount {
                            pubkey: recipient_stable_coin_account,
                            is_writable: true,
                            is_signer: false,
                        },
                    ]);
                }
            }

            let accounts_for_clear_original = get_accounts_for_clear(
                gringotts.lz_endpoint_program,
                &gringotts.key(),
                params.src_eid,
                &params.sender,
                order.nonce,
            );
            accounts.extend(accounts_for_clear_original);

            let accounts_for_clear = get_accounts_for_clear(
                gringotts.lz_endpoint_program,
                &gringotts.key(),
                params.src_eid,
                &params.sender,
                params.nonce,
            );
            accounts.extend(accounts_for_clear);
        } else {
            require!(false, LzReceiveTypesErrorCode::InvalidParams);
        }

        Ok(accounts)
    }
}

#[error_code]
pub enum LzReceiveTypesErrorCode {
    #[msg("Invalid params.")]
    InvalidParams,
}
