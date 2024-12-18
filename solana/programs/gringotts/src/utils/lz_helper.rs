const TYPE_3: u16 = 3;
const OPTION_TYPE_LZRECEIVE: u8 = 1;
const WORKER_ID: u8 = 1;

#[derive(Debug, Clone)]
pub struct OptionsBuilder {
    options: Vec<u8>,
}

impl OptionsBuilder {
    pub(crate) fn new() -> Self {
        OptionsBuilder {
            options: TYPE_3.to_be_bytes().to_vec(),
        }
    }

    pub(crate) fn add_executor_lz_receive_option(&mut self, gas: u64, value: u64) {
        self.add_executor_option(
            OPTION_TYPE_LZRECEIVE,
            Self::encode_lz_receive_option(gas, value),
        );
    }

    fn add_executor_option(&mut self, option_type: u8, option: Vec<u8>) {
        let option_size = option.len() as u16;

        self.options.extend(WORKER_ID.to_be_bytes());
        self.options.extend((option_size + 1).to_be_bytes());
        self.options.extend(option_type.to_be_bytes());
        self.options.extend(option);
    }

    fn encode_lz_receive_option(gas: u64, value: u64) -> Vec<u8> {
        let mut res = vec![0u8; 8];
        res.extend_from_slice(&gas.to_be_bytes());

        if value > 0 {
            let mut value_bytes = vec![0u8; 8];
            value_bytes.extend_from_slice(&value.to_be_bytes());
            res.extend_from_slice(&value_bytes);
        }
        res
    }

    pub(crate) fn options(self) -> Vec<u8> {
        self.options
    }
}
