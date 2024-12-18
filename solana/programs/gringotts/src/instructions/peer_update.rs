use crate::state::{Gringotts, Peer};
use crate::*;

#[derive(Accounts)]
#[instruction(params: PeerUpdateParams)]
pub struct PeerUpdate<'info> {
    #[account(mut, seeds = [PEER_SEED, &params.lz_eid.to_le_bytes()], bump = peer.bump)]
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

        if let Some(address) = params.address {
            peer.address = address;
        }
        if let Some(stable_coins) = &params.stable_coins {
            peer.stable_coins = stable_coins.clone();
        }
        if let Some(base_gas_estimate) = params.base_gas_estimate {
            peer.base_gas_estimate = base_gas_estimate;
        }

        Ok(())
    }
}

#[derive(Debug, Clone, AnchorSerialize, AnchorDeserialize)]
pub struct PeerUpdateParams {
    pub lz_eid: u32,
    pub address: Option<[u8; 32]>,
    pub stable_coins: Option<Vec<[u8; 32]>>,
    pub base_gas_estimate: Option<u64>,
}
