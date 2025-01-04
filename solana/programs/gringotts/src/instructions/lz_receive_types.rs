use crate::msg_codec::{ChainTransfer, Message, SolanaTransfer, CHAIN_TRANSFER_TYPE};
use crate::state::Gringotts;
use crate::*;
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

        let (vault, _) = Pubkey::find_program_address(
            &[VAULT_SEED],
            ctx.program_id,
        );

        let (self_peer, _) = Pubkey::find_program_address(
            &[PEER_SEED, &gringotts.lz_eid.to_be_bytes()],
            ctx.program_id,
        );
        let (peer, _) = Pubkey::find_program_address(
            &[PEER_SEED, &params.src_eid.to_be_bytes()],
            ctx.program_id,
        );

        let mut accounts = vec![
            LzAccount {
                pubkey: gringotts.key(),
                is_signer: false,
                is_writable: false,
            },
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
            LzAccount {
                pubkey: vault,
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
                pubkey: System::id(),
                is_signer: false,
                is_writable: false,
            },
        ];

        let message = Message::decode(params.message.as_slice());

        if message.header == CHAIN_TRANSFER_TYPE {
            let chain_transfer = ChainTransfer::decode(message.payload);
            let solana_transfer = SolanaTransfer::decode(chain_transfer.message, chain_transfer.items.len());

            for i in 0..solana_transfer.accounts_address.len() {
                let is_writeable_flag_index = i / 8;
                let is_writeable_flag_mask = 7 - (i % 8);

                accounts.push(LzAccount {
                    pubkey: Pubkey::new_from_array(*solana_transfer.accounts_address[i]),
                    is_writable: solana_transfer.accounts_flags[is_writeable_flag_index] & (1 << is_writeable_flag_mask) != 0,
                    is_signer: false,
                });
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
