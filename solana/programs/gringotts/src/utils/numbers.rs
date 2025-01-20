pub const BPS: u64 = 10_000;
pub const MICRO_BPS: u64 = BPS * 1000;

pub(crate) fn micro_bps(value: u64, micro_bps: u32) -> u64 {
    (value as f64 * micro_bps as f64 / MICRO_BPS as f64) as u64
}

pub(crate) fn bps(value: u64, bps: u32) -> u64 {
    micro_bps(value, bps * 1000)
}

pub(crate) fn change_decimals(value: u64, current_decimals: u8, new_decimals: u8) -> u64 {
    value * 10u64.pow(new_decimals as u32) / 10u64.pow(current_decimals as u32)
}