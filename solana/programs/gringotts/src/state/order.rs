use crate::constants::MAX_MESSAGE_LENGTH;
use anchor_lang::prelude::*;

#[derive(InitSpace, AnchorSerialize, AnchorDeserialize, PartialEq, Clone)]
pub enum OrderState {
    None,
    Processed,
}

#[account]
#[derive(InitSpace)]
pub struct Order {
    #[max_len(MAX_MESSAGE_LENGTH)]
    pub message: Vec<u8>,
    pub message_id: [u8; 32],
    pub nonce: u64,
    pub timestamp: u64,
    pub state: OrderState,
}
