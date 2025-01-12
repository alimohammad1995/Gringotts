pub const CHAIN_TRANSFER_TYPE: u8 = 1;
pub const CHAIN_REGISTER_ORDER_TYPE: u8 = 2;
pub const CHAIN_COMPLETION_TYPE: u8 = 3;

#[derive(Debug, Clone)]
pub struct Message<'a> {
    pub header: u8,
    pub message: &'a [u8],
}

impl<'a> Message<'a> {
    pub fn new(header: u8, payload: &'a [u8]) -> Message<'a> {
        Message { header, message: payload }
    }

    pub fn encode(&self) -> Vec<u8> {
        let mut encoded: Vec<u8> = Vec::new();
        encoded.push(self.header);
        encoded.extend_from_slice(&self.message);
        encoded
    }

    pub fn decode(data: &'a [u8]) -> Message<'a> {
        Self::new(data[0], &data[1..])
    }
}

#[derive(Debug, Clone)]
pub struct ChainTransferItem {
    pub amount_usdx: u64,
}

#[derive(Debug, Clone)]
pub struct ChainTransfer<'info> {
    pub items: Vec<ChainTransferItem>,
    pub message: &'info [u8],
}

impl<'info> ChainTransfer<'info> {
    pub fn encode(&self) -> Vec<u8> {
        let mut encoded: Vec<u8> = Vec::new();

        let items_count = self.items.len() as u8;
        encoded.extend_from_slice(&items_count.to_be_bytes());

        for item in &self.items {
            encoded.extend_from_slice(&item.amount_usdx.to_be_bytes());
        }

        encoded.extend_from_slice(&self.message);
        encoded
    }

    pub fn decode(data: &'info [u8]) -> ChainTransfer<'info> {
        let mut offset = 0;

        // Read itemsCount (uint8, 1 byte)
        let items_count = u8::from_be_bytes((&data[offset..offset + 1]).try_into().unwrap());
        offset += 1;

        let mut items = Vec::with_capacity(items_count as usize);

        for _ in 0..items_count {
            // Read amountUSDX (uint64, 8 bytes)
            let amount_usdx = u64::from_be_bytes((&data[offset..offset + 8]).try_into().unwrap());
            offset += 8;

            items.push(ChainTransferItem {
                amount_usdx
            });
        }

        let message = &data[offset..];

        ChainTransfer { items, message }
    }
}

#[derive(Debug, Clone)]
pub struct SolanaTransferItem<'info> {
    pub mint_index: usize,
    pub recipient_index: usize,
    pub recipient_token_index: usize,
    pub swap_accounts: &'info [u8],
    pub swap_command: &'info [u8],
}

#[derive(Debug, Clone)]
pub struct SolanaTransfer<'info> {
    pub stable_coin_index: usize,
    pub accounts_address: Vec<&'info [u8; 32]>,
    pub accounts_flags: &'info [u8],
    pub items: Vec<SolanaTransferItem<'info>>,
}

const STABLE_COIN_MINT_INDEX: usize = 5;

impl<'info> SolanaTransfer<'info> {
    pub fn decode(data: &'info [u8], chain_transfer_count: usize) -> SolanaTransfer<'info> {
        let mut offset = 0;

        // Read stable_coin_index (uint8, 1 byte)
        let stable_coin_index = u8::from_be_bytes((&data[offset..offset + 1]).try_into().unwrap()) as usize;
        offset += 1;

        // Read accounts_count (uint8, 1 byte)
        let accounts_count = u8::from_be_bytes((&data[offset..offset + 1]).try_into().unwrap());
        offset += 1;

        // Read accounts
        let mut accounts_address = Vec::<&[u8; 32]>::with_capacity(accounts_count as usize);
        for _ in 0..accounts_count {
            accounts_address.push(data[offset..offset + 32].try_into().unwrap());
            offset += 32;
        }

        // Read accounts IsWritable
        let account_flags_count = ((accounts_count + 7) / 8) as usize;
        let accounts_flags = &data[offset..offset + account_flags_count];
        offset += account_flags_count;

        let mut items = Vec::with_capacity(chain_transfer_count);
        let message = &data[offset..];
        let mut message_offset = 0;

        for _ in 0..chain_transfer_count {
            let recipient_index = message[message_offset] as usize;
            let mint_index = message[message_offset + 1] as usize;
            let recipient_token_index = message[message_offset + 2] as usize;

            message_offset += 3;

            let mut swap_accounts: &[u8] = &[];
            let mut swap_command: &[u8] = &[];

            if mint_index != STABLE_COIN_MINT_INDEX {
                let swap_accounts_count = message[message_offset] as usize;
                message_offset += 1;

                swap_accounts = &message[message_offset..message_offset + swap_accounts_count];
                message_offset += swap_accounts_count;

                let command_length = u16::from_be_bytes((&message[message_offset..message_offset + 2]).try_into().unwrap()) as usize;
                message_offset += 2;

                swap_command = &message[message_offset..message_offset + command_length];
                message_offset += command_length;
            }

            items.push(SolanaTransferItem {
                mint_index,
                recipient_index,
                recipient_token_index,
                swap_accounts,
                swap_command,
            });
        }

        SolanaTransfer {
            stable_coin_index,
            accounts_address,
            accounts_flags,
            items,
        }
    }
}

#[derive(Debug, Clone)]
pub struct ChainExecution<'info> {
    pub message_id: &'info [u8; 32],
}
impl<'info> ChainExecution<'info> {
    pub fn decode(data: &'info [u8]) -> ChainExecution<'info> {
        ChainExecution { message_id: data.try_into().unwrap() }
    }
}

#[derive(Debug, Clone)]
pub struct ChainCompletion<'info> {
    pub message_id: &'info [u8; 32],
}
impl<'info> ChainCompletion<'info> {
    pub fn decode(data: &'info [u8]) -> ChainCompletion<'info> {
        ChainCompletion { message_id: data.try_into().unwrap() }
    }
}
