use crate::instructions::estimate_impl::estimate_marketplace;
use crate::state::{Gringotts, Peer};
use crate::*;
use pyth_solana_receiver_sdk::price_update::PriceUpdateV2;

#[derive(Accounts)]
#[instruction(params: EstimateRequest)]
pub struct Estimate<'info> {
    pub price_feed: Account<'info, PriceUpdateV2>,
    #[account(seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,
}

impl Estimate<'_> {
    pub fn apply(ctx: Context<Estimate>, params: &EstimateRequest) -> Result<EstimateResponse> {
        let gringotts = &ctx.accounts.gringotts;
        let price_feed = &ctx.accounts.price_feed;

        let peers = params.process(&ctx.remaining_accounts)?;

        let estimate_result = estimate_marketplace(
            params,
            gringotts,
            peers.as_slice(),
            price_feed,
            &ctx.remaining_accounts[peers.len()..],
            false,
        )?;

        require!(
            params.inbound.amount_usdx
                >= estimate_result.commission_usdx + estimate_result.transfer_gas_usdx,
            EstimateErrorCode::InvalidParams
        );

        Ok(estimate_result)
    }
}

impl EstimateRequest {
    pub fn process(&self, remaining_accounts: &[AccountInfo]) -> Result<Vec<Peer>> {
        let raw_peers: &[AccountInfo] = &remaining_accounts[0..self.outbounds.len()];
        let mut peers = Vec::with_capacity(self.outbounds.len());

        for i in 0usize..raw_peers.len() {
            let mut data = &raw_peers[i].try_borrow_data()?[..];
            let peer = Peer::try_deserialize(&mut data)?;
            peers.push(peer);
        }

        Ok(peers)
    }
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct EstimateInboundTransfer {
    pub amount_usdx: u64,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct EstimateOutboundTransferItem {
    pub asset: [u8; 32],
    pub execution_gas: u64,
    pub command_length: u16,
    pub metadata_length: u16,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct EstimateOutboundTransfer {
    pub chain_id: u8,
    pub metadata_length: u16,
    pub items: Vec<EstimateOutboundTransferItem>,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct EstimateRequest {
    pub inbound: EstimateInboundTransfer,
    pub outbounds: Vec<EstimateOutboundTransfer>,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct EstimateOutboundDetails {
    pub chain_id: u8,
    pub execution_gas: u64,
    pub execution_gas_usdx: u64,
    pub transfer_gas: u64,
    pub transfer_gas_usdx: u64,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct EstimateResponse {
    pub commission_usdx: u64,
    pub commission_discount_usdx: u64,
    pub transfer_gas_usdx: u64,
    pub transfer_gas_discount_usdx: u64,
    pub outbound_details: Vec<EstimateOutboundDetails>,
}

#[error_code]
pub enum EstimateErrorCode {
    #[msg("Invalid params.")]
    InvalidParams,
}
