use crate::msg_codec::{ChainTransfer, Message, CHAIN_TRANSFER_TYPE};
use crate::state::Gringotts;
use crate::utils::NATIVE_MINT;
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

impl LzReceiveTypes<'_> {
    pub fn apply(
        ctx: &Context<LzReceiveTypes>,
        params: &LzReceiveParams,
    ) -> Result<Vec<LzAccount>> {
        let gringotts = &ctx.accounts.gringotts;

        let (self_peer, _) = Pubkey::find_program_address(
            &[PEER_SEED, &gringotts.lz_eid.to_le_bytes()],
            ctx.program_id,
        );
        let (peer, _) = Pubkey::find_program_address(
            &[PEER_SEED, &params.src_eid.to_le_bytes()],
            ctx.program_id,
        );

        let mut accounts = vec![
            LzAccount {
                pubkey: self_peer,
                is_signer: false,
                is_writable: false,
            },
            LzAccount {
                pubkey: peer,
                is_signer: false,
                is_writable: false,
            },
        ];

        let message = Message::decode(params.message.as_slice());

        if message.header == CHAIN_TRANSFER_TYPE {
            let chain_transfer = ChainTransfer::decode(message.payload);

            for item in &chain_transfer.items {
                let recipient = Pubkey::new_from_array(*item.recipient);

                accounts.push(LzAccount {
                    pubkey: recipient,
                    is_signer: false,
                    is_writable: true,
                }); //XXXX

                if *item.executor == [0; 32] {
                    accounts.push(LzAccount {
                        pubkey: Pubkey::new_from_array(*item.asset),
                        is_signer: false,
                        is_writable: false,
                    });
                    accounts.push(LzAccount {
                        pubkey: associated_token::get_associated_token_address(
                            &recipient,
                            &Pubkey::new_from_array(*item.asset),
                        ),
                        is_signer: false,
                        is_writable: true,
                    });
                    accounts.push(LzAccount {
                        pubkey: associated_token::get_associated_token_address(
                            &gringotts.key(),
                            &Pubkey::new_from_array(*item.asset),
                        ),
                        is_signer: false,
                        is_writable: true,
                    });
                } else {
                    if *item.asset != [0; 32] {
                        accounts.push(LzAccount {
                            pubkey: Pubkey::new_from_array(*item.asset),
                            is_signer: false,
                            is_writable: false,
                        });
                        accounts.push(LzAccount {
                            pubkey: associated_token::get_associated_token_address(
                                &recipient,
                                &Pubkey::new_from_array(*item.asset),
                            ),
                            is_signer: false,
                            is_writable: true,
                        });
                    } else {
                        accounts.push(LzAccount {
                            pubkey: NATIVE_MINT,
                            is_signer: false,
                            is_writable: false,
                        }); //XXXX
                        accounts.push(LzAccount {
                            pubkey: associated_token::get_associated_token_address(
                                &gringotts.key(),
                                &NATIVE_MINT,
                            ),
                            is_signer: false,
                            is_writable: true,
                        }); //XXXX
                    }

                    accounts.push(LzAccount {
                        pubkey: Pubkey::new_from_array(*item.executor),
                        is_signer: false,
                        is_writable: false,
                    }); //XXXX
                    accounts.push(LzAccount {
                        pubkey: Pubkey::new_from_array(*item.stable_token),
                        is_signer: false,
                        is_writable: false,
                    }); //XXXX
                    accounts.push(LzAccount {
                        pubkey: associated_token::get_associated_token_address(
                            &gringotts.key(),
                            &Pubkey::new_from_array(*item.stable_token),
                        ),
                        is_signer: false,
                        is_writable: true,
                    }); //XXXX
                    accounts.push(LzAccount {
                        pubkey: associated_token::get_associated_token_address(
                            &recipient.key(),
                            &Pubkey::new_from_array(*item.stable_token),
                        ),
                        is_signer: false,
                        is_writable: true,
                    }); //XXXX

                    let swap_accounts_len = item.metadata[0] as usize;
                    let mut m_index = 1;
                    for _ in 0..swap_accounts_len {
                        accounts.push(LzAccount {
                            pubkey: Pubkey::new_from_array(
                                item.metadata[m_index..m_index + 32]
                                    .try_into()
                                    .unwrap(),
                            ),
                            is_signer: false,
                            is_writable: item.metadata[m_index + 32] == 0,
                        });
                        m_index += 33
                    }
                }
            }
        }

        let accounts_for_clear = get_accounts_for_clear(
            gringotts.lz_endpoint_program,
            &gringotts.key(),
            params.src_eid,
            &params.sender,
            params.nonce,
        );

        accounts.extend(accounts_for_clear);

        Ok(accounts)
    }
}
