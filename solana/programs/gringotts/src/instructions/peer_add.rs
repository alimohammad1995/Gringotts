use crate::state::{Gringotts, Peer};
use crate::*;

#[derive(Accounts)]
#[instruction(params: PeerAddParams)]
pub struct PeerAdd<'info> {
    #[account(
        init,
        payer = owner,
        space = 8 + Peer::INIT_SPACE,
        seeds = [PEER_SEED, &params.lz_eid.to_le_bytes()],
        bump,
    )]
    pub peer: Account<'info, Peer>,
    #[account(has_one = owner, seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,
    #[account(mut)]
    pub owner: Signer<'info>,
    pub system_program: Program<'info, System>,
}

impl PeerAdd<'_> {
    pub fn apply(ctx: Context<PeerAdd>, params: &PeerAddParams) -> Result<()> {
        let peer = &mut ctx.accounts.peer;
        peer.bump = ctx.bumps.peer;

        peer.chain_id = params.chain_id;
        peer.lz_eid = params.lz_eid;
        peer.address = params.address;

        peer.stable_coins = params.stable_coins.clone();
        peer.base_gas_estimate = params.base_gas_estimate;

        Ok(())
    }
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct PeerAddParams {
    pub chain_id: u8,
    pub lz_eid: u32,
    pub address: [u8; 32],

    pub stable_coins: Vec<[u8; 32]>,

    pub base_gas_estimate: u64,
}
