use anchor_lang::prelude::Pubkey;

pub const NETWORK_DECIMALS: u8 = 9;
pub const MAX_PRICE_AGE: u64 = 5 * 60;
pub const CHAIN_TRANSFER_DECIMALS: u8 = 6;
pub const MAX_STABLE_COINS: usize = 16;
pub const MAX_QUEUE_SIZE: usize = 16;
pub const MAX_MESSAGE_LENGTH: usize = 1024;

pub const NATIVE_MINT: Pubkey = Pubkey::new_from_array([
    6, 155, 136, 87, 254, 171, 129, 132, 251, 104, 127, 99, 70, 24, 192, 53, 218, 196, 57, 220, 26, 235, 59, 85, 152, 160, 240, 0, 0, 0, 0, 1
]);

pub const JUPITER: Pubkey = Pubkey::new_from_array([
    4, 121, 213, 91, 242, 49, 192, 110, 238, 116, 197, 110, 206, 104, 21, 7, 253, 177, 178, 222, 163, 244, 142, 81, 2, 177, 205, 162, 86, 188, 19, 143
]);
