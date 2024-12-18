pub const CHAIN_TRANSFER_TYPE: u8 = 1;

#[derive(Debug, Clone)]
pub struct Message {
    pub header: u8,
    pub payload: Vec<u8>,
}

impl Message {
    pub fn new(header: u8, payload: Vec<u8>) -> Message {
        Message { header, payload }
    }

    pub fn encode(&self) -> Vec<u8> {
        let mut encoded: Vec<u8> = Vec::new();
        encoded.push(self.header);
        encoded.extend_from_slice(&self.payload);
        encoded
    }

    pub fn decode(data: &[u8]) -> Message {
        Self::new(data[0], data[1..].to_vec())
    }
}

#[derive(Debug, Clone)]
pub struct ChainTransferItem {
    pub amount_usdx: u64,
    pub asset: [u8; 32],
    pub recipient: [u8; 32],
    pub executor: [u8; 32],
    pub stable_token: [u8; 32],
    pub command: Vec<u8>,
    pub metadata: Vec<u8>,
}

#[derive(Debug, Clone)]
pub struct ChainTransfer {
    pub items: Vec<ChainTransferItem>,
}

impl ChainTransfer {
    pub fn encode(&self) -> Vec<u8> {
        let mut encoded: Vec<u8> = Vec::new();

        let items_count = self.items.len() as u8;
        encoded.extend_from_slice(&items_count.to_be_bytes());

        for item in &self.items {
            encoded.extend_from_slice(&item.amount_usdx.to_be_bytes());

            encoded.extend_from_slice(&item.asset);
            encoded.extend_from_slice(&item.recipient);
            encoded.extend_from_slice(&item.executor);
            encoded.extend_from_slice(&item.stable_token);

            encoded.extend_from_slice(&(item.command.len() as u16).to_be_bytes());
            encoded.extend_from_slice(&item.command);

            encoded.extend_from_slice(&(item.metadata.len() as u16).to_be_bytes());
            encoded.extend_from_slice(&item.metadata);
        }

        encoded
    }

    pub fn decode(data: &[u8]) -> ChainTransfer {
        let mut offset = 0;

        // Read itemsCount (uint8, 1 byte)
        let items_count_bytes = &data[offset..offset + 1];
        offset += 1;
        let items_count = u8::from_be_bytes(items_count_bytes.try_into().unwrap());

        let mut items = Vec::new();

        for _ in 0..items_count {
            // Read amountUSDX (uint64, 8 bytes)
            let amount_usdx_bytes = &data[offset..offset + 8];
            offset += 8;
            let amount_usdx = u64::from_be_bytes(amount_usdx_bytes.try_into().unwrap());

            // Read asset (32 bytes)
            let asset = data[offset..offset + 32].try_into().unwrap();
            offset += 32;

            // Read recipient (32 bytes)
            let recipient = data[offset..offset + 32].try_into().unwrap();
            offset += 32;

            // Read executor (32 bytes)
            let executor = data[offset..offset + 32].try_into().unwrap();
            offset += 32;

            // Read stableToken (32 bytes)
            let stable_token = data[offset..offset + 32].try_into().unwrap();
            offset += 32;

            // Read command (uint16, 2 bytes)
            let command_length =
                u16::from_be_bytes((&data[offset..offset + 2]).try_into().unwrap()) as usize;
            offset += 2;
            let command = data[offset..offset + command_length].to_vec();
            offset += command_length;

            // Read metadata (uint16, 2 bytes)
            let metadata_length =
                u16::from_be_bytes((&data[offset..offset + 2]).try_into().unwrap()) as usize;
            offset += 2;
            let metadata = data[offset..offset + metadata_length].to_vec();
            offset += metadata_length;

            let item = ChainTransferItem {
                amount_usdx,
                asset,
                recipient,
                executor,
                stable_token,
                command,
                metadata,
            };

            items.push(item);
        }

        ChainTransfer { items }
    }
}
