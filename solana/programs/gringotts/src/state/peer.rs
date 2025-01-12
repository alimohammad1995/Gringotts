use anchor_lang::prelude::*;

#[account]
#[derive(InitSpace)]
pub struct Peer {
    pub bump: u8,
    pub chain_id: u8,
    pub lz_eid: u32,
    pub address: [u8; 32],
    pub multi_send: bool,
    pub base_gas_estimate: u64,
}
