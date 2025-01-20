use crate::state::{Gringotts, Peer};
use crate::*;

#[derive(Accounts)]
#[instruction(params: PeerUpdateParams)]
pub struct PeerUpdate<'info> {
    #[account(mut, seeds = [PEER_SEED, &params.lz_eid.to_be_bytes()], bump = peer.bump)]
    pub peer: Account<'info, Peer>,
    #[account(has_one = owner, seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,
    #[account(mut)]
    pub owner: Signer<'info>,
    pub system_program: Program<'info, System>,
}

impl PeerUpdate<'_> {
    pub fn apply(ctx: Context<PeerUpdate>, params: &PeerUpdateParams) -> Result<()> {
        let peer = &mut ctx.accounts.peer;

        if let Some(chain_id) = params.chain_id {
            peer.chain_id = chain_id;
        }
        if let Some(address) = params.address {
            peer.address = address;
        }
        if let Some(base_gas_estimate) = params.base_gas_estimate {
            peer.base_gas_estimate = base_gas_estimate;
        }
        if let Some(multi_send) = params.multi_send {
            peer.multi_send = multi_send;
        }

        Ok(())
    }
}

#[derive(Debug, Clone, AnchorSerialize, AnchorDeserialize)]
pub struct PeerUpdateParams {
    pub lz_eid: u32,
    pub chain_id: Option<u8>,
    pub address: Option<[u8; 32]>,
    pub multi_send: Option<bool>,
    pub base_gas_estimate: Option<u64>,
}
