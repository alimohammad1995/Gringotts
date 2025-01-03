use anchor_lang::prelude::*;

#[account]
#[derive(InitSpace)]
pub struct Gringotts {
    pub owner: Pubkey,
    pub bump: u8,
    pub vault_bump: u8,
    pub chain_id: u8,

    pub lz_eid: u32,
    pub lz_endpoint_program: Pubkey,

    pub pyth_price_feed_id: [u8; 32],

    pub commission_micro_bps: u32,

    pub commission_discount_bps: u32,
    pub gas_discount_bps: u32,
}

#[account]
#[derive(InitSpace)]
pub struct LzReceiveTypesAccounts {
    pub gringotts: Pubkey,
}
