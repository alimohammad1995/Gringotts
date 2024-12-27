use crate::msg_codec::{ChainTransfer, Message, CHAIN_TRANSFER_TYPE};
use crate::state::Gringotts;
use crate::*;
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
        ];

        let message = Message::decode(params.message.as_slice());

        if message.header == CHAIN_TRANSFER_TYPE {
            let chain_transfer = ChainTransfer::decode(message.payload);

            let account_size = chain_transfer.metadata[0] as usize;
            let mut m_index = 1;

            for _ in 0..account_size {
                accounts.push(LzAccount {
                    pubkey: Pubkey::new_from_array(
                        chain_transfer.metadata[m_index..m_index + 32]
                            .try_into()
                            .unwrap(),
                    ),
                    is_writable: chain_transfer.metadata[m_index + 32] == 1,
                    is_signer: false,
                });
                m_index += 33
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
