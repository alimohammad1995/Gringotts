use crate::constants::MAX_STABLE_COINS;
use anchor_lang::prelude::*;

#[account]
#[derive(InitSpace)]
pub struct Peer {
    pub bump: u8,
    pub chain_id: u8,
    pub lz_eid: u32,
    pub address: [u8; 32],

    #[max_len(MAX_STABLE_COINS)]
    pub stable_coins: Vec<[u8; 32]>,

    pub base_gas_estimate: u64,
}

impl Peer {
    pub fn has_stable_coin(&self, stable_coin: [u8; 32]) -> bool {
        self.stable_coins.iter().any(|&coin| coin == stable_coin)
    }
}
