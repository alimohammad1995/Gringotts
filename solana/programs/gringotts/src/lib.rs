use anchor_lang::prelude::*;

pub mod constants;
pub mod instructions;
pub mod msg_codec;
pub mod state;
pub mod utils;

use instructions::*;

use oapp::{endpoint_cpi::LzAccount, LzReceiveParams};

declare_id!("D4qgtRF5t7MzA5nASfag3ef5Y68xb5hBcp9ac9Wttfkr");

pub const MAX_TRANSFERS: u8 = 4;
pub const MAX_PRICE_AGE: u64 = 5 * 60;
pub const NETWORK_DECIMALS: u8 = 9;
pub const CHAIN_TRANSFER_DECIMALS: u8 = 6;

pub const GRINGOTTS_SEED: &[u8] = b"Gringotts";
pub const VAULT_SEED: &[u8] = b"Vault";
pub const ASSET_SEED: &[u8] = b"Asset";
pub const PEER_SEED: &[u8] = b"Peer";
pub const LZ_RECEIVE_TYPES_SEED: &[u8] = b"LzReceiveTypes";

#[program]
pub mod gringotts {
    use super::*;

    pub fn initialize(ctx: Context<Initialize>, params: InitializeParams) -> Result<()> {
        Initialize::apply(ctx, &params)
    }

    pub fn destroy(ctx: Context<Destroy>) -> Result<()> {
        instructions::destroy(ctx)
    }

    pub fn peer_add(ctx: Context<PeerAdd>, params: PeerAddParams) -> Result<()> {
        PeerAdd::apply(ctx, &params)
    }

    pub fn peer_update(ctx: Context<PeerUpdate>, params: PeerUpdateParams) -> Result<()> {
        PeerUpdate::apply(ctx, &params)
    }

    pub fn token_withdraw(ctx: Context<TokenWithdraw>, params: TokenWithdrawParams) -> Result<()> {
        TokenWithdraw::apply(ctx, &params)
    }

    pub fn vault_withdraw(ctx: Context<VaultWithdraw>, params: VaultWithdrawParams) -> Result<()> {
        VaultWithdraw::apply(ctx, &params)
    }

    pub fn lz_receive<'a>(
        ctx: Context<'_, '_, 'a, 'a, LzReceive<'a>>,
        params: LzReceiveParams,
    ) -> Result<()> {
        LzReceive::apply(ctx, &params)
    }

    pub fn lz_receive_types(
        ctx: Context<LzReceiveTypes>,
        params: LzReceiveParams,
    ) -> Result<Vec<LzAccount>> {
        LzReceiveTypes::apply(&ctx, &params)
    }

    pub fn estimate(ctx: Context<Estimate>, params: EstimateRequest) -> Result<EstimateResponse> {
        Estimate::apply(ctx, &params)
    }

    pub fn bridge<'a>(
        ctx: Context<'_, '_, 'a, 'a, Bridge<'a>>,
        params: BridgeRequest,
    ) -> Result<BridgeResponse> {
        Bridge::apply(ctx, &params)
    }
}
