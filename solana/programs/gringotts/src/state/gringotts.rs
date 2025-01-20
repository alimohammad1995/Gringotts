use crate::constants::{MAX_QUEUE_SIZE, MAX_STABLE_COINS};
use anchor_lang::prelude::*;

#[account]
#[derive(InitSpace)]
pub struct Gringotts {
    pub owner: Pubkey,
    pub bump: u8,
    pub chain_id: u8,

    pub lz_endpoint_program: Pubkey,

    pub pyth_price_feed_id: [u8; 32],

    #[max_len(MAX_STABLE_COINS)]
    pub stable_coins: Vec<[u8; 32]>,

    pub commission_micro_bps: u32,
    pub commission_discount_bps: u32,
    pub gas_discount_bps: u32,

    pub tx_count: u64,
}

impl Gringotts {
    pub fn has_stable_coin(&self, stable_coin: [u8; 32]) -> bool {
        self.stable_coins.iter().any(|&coin| coin == stable_coin)
    }

    pub fn generate_id(&self) -> u64 {
        ((self.chain_id as u64) << 56) | self.tx_count
    }
}

#[account]
#[derive(InitSpace)]
pub struct LzReceiveTypesAccounts {
    pub gringotts: Pubkey,
    pub orders: [Pubkey; MAX_QUEUE_SIZE],
}

impl LzReceiveTypesAccounts {
    pub fn add_fc_item(&mut self, address: Pubkey) -> bool {
        let mut free_index = -1;

        for i in 0..self.orders.len() {
            if self.orders[i] == address {
                return false;
            } else if self.orders[i] == Pubkey::default() {
                free_index = i as i32;
            }
        }

        if free_index <= 0 {
            return false;
        }

        self.orders[free_index as usize] = address;
        true
    }

    pub fn remove_fc_item(&mut self, address: Pubkey) -> bool {
        for i in 0..self.orders.len() {
            if self.orders[i] == address {
                self.orders[i] = Pubkey::default();
                return true;
            }
        }

        false
    }

    pub fn contains_fc_item(&self, address: Pubkey) -> bool {
        for fc in self.orders.iter() {
            if fc.key() == address {
                return true;
            }
        }

        false
    }
}
