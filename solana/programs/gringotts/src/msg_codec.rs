pub const CHAIN_TRANSFER_TYPE: u8 = 1;

#[derive(Debug, Clone)]
pub struct Message<'a> {
    pub header: u8,
    pub payload: &'a [u8],
}

impl<'a> Message<'a> {
    pub fn new(header: u8, payload: &'a [u8]) -> Message<'a> {
        Message { header, payload }
    }

    pub fn encode(&self) -> Vec<u8> {
        let mut encoded: Vec<u8> = Vec::new();
        encoded.push(self.header);
        encoded.extend_from_slice(&self.payload);
        encoded
    }

    pub fn decode(data: &'a [u8]) -> Message {
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
    pub swap_account_count: usize,
    pub swap_command: &'info [u8],
}

#[derive(Debug, Clone)]
pub struct SolanaTransfer<'info> {
    pub accounts_address: Vec<&'info [u8; 32]>,
    pub accounts_flags: &'info [u8],
    pub accounts_mapping: &'info [u8],
    pub items: Vec<Option<SolanaTransferItem<'info>>>,
}

impl<'info> SolanaTransfer<'info> {
    pub fn decode(data: &'info [u8], chain_transfer_count: usize) -> SolanaTransfer<'info> {
        let mut offset = 0;

        // Read accounts_count (uint8, 1 byte)
        let accounts_count = u8::from_be_bytes((&data[offset..offset + 1]).try_into().unwrap());
        offset += 1;

        let mut accounts = Vec::<&[u8; 32]>::with_capacity(accounts_count as usize);
        for _ in 0..accounts_count {
            let slice = &data[offset..offset + 32];
            accounts.push(slice.try_into().unwrap());
            offset += 32;
        }

        let account_flags_count = ((accounts_count + 7) / 8) as usize;
        let accounts_flags = &data[offset..offset + account_flags_count];
        offset += account_flags_count;

        let mapping_counts = u8::from_be_bytes((&data[offset..offset + 1]).try_into().unwrap()) as usize;
        offset += 1;
        let accounts_mapping = &data[offset..offset + mapping_counts];
        offset += mapping_counts;

        let mut items = Vec::with_capacity(chain_transfer_count);
        let items_flags = data[offset..offset + 1][0];
        offset += 1;

        let message = &data[offset..];
        let mut message_offset = 0;

        for i in 0..chain_transfer_count {
            if items_flags & (1 << (7 - i)) == 0 {
                items.push(None)
            } else {
                let swap_account_counts = u8::from_be_bytes((&message[message_offset..message_offset + 1]).try_into().unwrap()) as usize;
                message_offset += 1;

                let command_length = u16::from_be_bytes((&message[message_offset..message_offset + 2]).try_into().unwrap()) as usize;
                message_offset += 2;

                let command = &message[message_offset..message_offset + command_length];
                message_offset += command_length;
                items.push(Some(SolanaTransferItem { swap_command: command, swap_account_count: swap_account_counts }));
            }
        }

        SolanaTransfer {
            accounts_address: accounts,
            accounts_flags: accounts_flags,
            accounts_mapping: accounts_mapping,
            items: items,
        }
    }
}