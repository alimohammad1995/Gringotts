pub const CHAIN_TRANSFER_TYPE: u8 = 1;

#[derive(Debug, Clone)]
pub struct Message<'info> {
    pub header: u8,
    pub payload: &'info [u8],
}

impl<'info> Message<'info> {
    pub fn new(header: u8, payload: &'info [u8]) -> Message<'info> {
        Message { header, payload }
    }

    pub fn encode(&self) -> Vec<u8> {
        let mut encoded: Vec<u8> = Vec::new();
        encoded.push(self.header);
        encoded.extend_from_slice(&self.payload);
        encoded
    }

    pub fn decode(data: &'info [u8]) -> Message {
        Self::new(data[0], &data[1..])
    }
}

#[derive(Debug, Clone)]
pub struct ChainTransferItem<'info> {
    pub amount_usdx: u64,
    pub asset: &'info [u8; 32],
    pub recipient: &'info [u8; 32],
    pub executor: &'info [u8; 32],
    pub stable_token: &'info [u8; 32],
    pub command: &'info [u8],
    pub metadata: &'info [u8],
}

#[derive(Debug, Clone)]
pub struct ChainTransfer<'info> {
    pub items: Vec<ChainTransferItem<'info>>,
    pub metadata: &'info [u8],
}

impl<'info> ChainTransfer<'info> {
    pub fn encode(&self) -> Vec<u8> {
        let mut encoded: Vec<u8> = Vec::new();

        let items_count = self.items.len() as u8;
        encoded.extend_from_slice(&items_count.to_be_bytes());

        for item in &self.items {
            encoded.extend_from_slice(&item.amount_usdx.to_be_bytes());

            encoded.extend_from_slice(item.asset);
            encoded.extend_from_slice(item.recipient);

            let need_swap = if item.command.len() > 0 { 1u8 } else { 0u8 };
            encoded.extend_from_slice(&need_swap.to_be_bytes());

            if need_swap > 0 {
                encoded.extend_from_slice(item.executor);
                encoded.extend_from_slice(item.stable_token);

                encoded.extend_from_slice(&(item.command.len() as u16).to_be_bytes());
                encoded.extend_from_slice(&item.command);

                encoded.extend_from_slice(&(item.metadata.len() as u16).to_be_bytes());
                encoded.extend_from_slice(&item.metadata);
            }
        }

        encoded.extend_from_slice(&(self.metadata.len() as u16).to_be_bytes());
        encoded.extend_from_slice(&self.metadata);

        encoded
    }

    pub fn decode(data: &'info [u8]) -> ChainTransfer<'info> {
        let mut offset = 0;

        // Read itemsCount (uint8, 1 byte)
        let items_count_bytes = &data[offset..offset + 1];
        offset += 1;
        let items_count = u8::from_be_bytes(items_count_bytes.try_into().unwrap());

        let mut items = Vec::new();

        for _ in 0..items_count {
            // Read amountUSDX (uint64, 8 bytes)
            let amount_usdx = u64::from_be_bytes((&data[offset..offset + 8]).try_into().unwrap());
            offset += 8;

            // Read asset (32 bytes)
            let asset = &data[offset..offset + 32];
            offset += 32;

            // Read recipient (32 bytes)
            let recipient = &data[offset..offset + 32];
            offset += 32;

            // Read needSwap (uint8, 1 byte)
            let need_swap = &data[offset..offset + 1];
            offset += 1;

            let item: ChainTransferItem;

            if need_swap[0] > 0 {
                // Read executor (32 bytes)
                let executor = &data[offset..offset + 32];
                offset += 32;

                // Read stableToken (32 bytes)
                let stable_token = &data[offset..offset + 32];
                offset += 32;

                // Read command (uint16, 2 bytes)
                let command_length = u16::from_be_bytes((&data[offset..offset + 2]).try_into().unwrap()) as usize;
                offset += 2;
                let command = &data[offset..offset + command_length];
                offset += command_length;

                // Read metadata (uint16, 2 bytes)
                let metadata_length = u16::from_be_bytes((&data[offset..offset + 2]).try_into().unwrap()) as usize;
                offset += 2;
                let metadata = &data[offset..offset + metadata_length];
                offset += metadata_length;

                item = ChainTransferItem {
                    amount_usdx,
                    asset: asset.try_into().unwrap(),
                    recipient: recipient.try_into().unwrap(),
                    executor: executor.try_into().unwrap(),
                    stable_token: stable_token.try_into().unwrap(),
                    command,
                    metadata,
                };
            } else {
                item = ChainTransferItem {
                    amount_usdx,
                    asset: asset.try_into().unwrap(),
                    recipient: recipient.try_into().unwrap(),
                    executor: &[0u8; 32],
                    stable_token: &[0u8; 32],
                    command: &[],
                    metadata: &[],
                };
            }

            items.push(item);
        }

        let metadata_length = u16::from_be_bytes((&data[offset..offset + 2]).try_into().unwrap()) as usize;
        offset += 2;
        let metadata = data[offset..offset + metadata_length].as_ref();

        ChainTransfer { items, metadata }
    }
}
