use anchor_lang::prelude::*;

pub mod constants;
pub mod instructions;
pub mod msg_codec;
pub mod state;
pub mod utils;

use instructions::*;

use oapp::{endpoint_cpi::LzAccount, LzReceiveParams};

declare_id!("Dh3ak9SbvubmtQTeq8kgDXhZeUKorKa5bt1x4PZ8Vebf");

pub const GRINGOTTS_SEED: &[u8] = b"Gringotts";
pub const ORDER_SEED: &[u8] = b"OrderItem";
pub const VAULT_SEED: &[u8] = b"Vault";
pub const PEER_SEED: &[u8] = b"Peer";
pub const LZ_RECEIVE_TYPES_SEED: &[u8] = b"LzReceiveTypes";

#[program]
pub mod gringotts {
    use super::*;

    pub fn gringotts_initialize(ctx: Context<GringottsInitialize>, params: GringottsInitializeParams) -> Result<()> {
        GringottsInitialize::apply(ctx, &params)
    }

    pub fn gringotts_update(ctx: Context<GringottsUpdate>, params: GringottsUpdateParams) -> Result<()> {
        GringottsUpdate::apply(ctx, &params)
    }

    pub fn account_destroy(ctx: Context<AccountDestroy>, params: AccountDestroyParams) -> Result<()> {
        AccountDestroy::apply(&ctx, &params)
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

    pub fn token_fund(ctx: Context<TokenFund>, params: TokenFundParams) -> Result<()> {
        TokenFund::apply(ctx, &params)
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

    pub fn lz_receive_types<'a>(
        ctx: Context<'_, '_, 'a, 'a, LzReceiveTypes<'a>>,
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

    pub fn test(ctx: Context<Test>, params: TestParams) -> Result<()> {
        Test::apply(ctx, &params)
    }
}
