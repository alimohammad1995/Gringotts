use anchor_lang::prelude::*;

#[event]
pub struct SendChainTransferEvent {
    pub id: u64,
    pub sender: Pubkey,
    pub chain_id: u8,
    pub message_id: [u8; 32],
    pub amount_usdx: u64,
}

#[event]
pub struct ReceiveChainTransferEvent {
    pub chain_id_lz: u32,
    pub message_id: [u8; 32],
    pub asset: Pubkey,
    pub recipient: Pubkey,
    pub amount_usdx: u64,
}