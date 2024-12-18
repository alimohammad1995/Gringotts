use crate::{
    state::{Gringotts, LzReceiveTypesAccounts},
    GRINGOTTS_SEED, LZ_RECEIVE_TYPES_SEED,
};
use anchor_lang::prelude::*;
use oapp::endpoint::{instructions::RegisterOAppParams, ID as ENDPOINT_ID};
use pyth_solana_receiver_sdk::price_update::get_feed_id_from_hex;

#[derive(Accounts)]
#[instruction(params: InitializeParams)]
pub struct Initialize<'info> {
    #[account(mut)]
    pub owner: Signer<'info>,
    #[account(
        init,
        payer = owner,
        space = 8 + Gringotts::INIT_SPACE,
        seeds = [GRINGOTTS_SEED],
        bump
    )]
    pub gringotts: Account<'info, Gringotts>,
    #[account(
        init,
        payer = owner,
        space = 8 + LzReceiveTypesAccounts::INIT_SPACE,
        seeds = [LZ_RECEIVE_TYPES_SEED, &gringotts.key().to_bytes()],
        bump
    )]
    pub lz_receive_types_accounts: Account<'info, LzReceiveTypesAccounts>,
    pub system_program: Program<'info, System>,
}

impl Initialize<'_> {
    pub fn apply(ctx: Context<Initialize>, params: &InitializeParams) -> Result<()> {
        let gringotts = &mut ctx.accounts.gringotts;
        let user = &ctx.accounts.owner;
        let lz_receive_types_accounts = &mut ctx.accounts.lz_receive_types_accounts;

        gringotts.owner = user.key();
        gringotts.bump = ctx.bumps.gringotts;

        gringotts.chain_id = params.chain_id;

        gringotts.lz_eid = params.lz_eid;
        gringotts.lz_endpoint_program = params.lz_endpoint_program;

        gringotts.pyth_price_feed_id = get_feed_id_from_hex(&params.pyth_price_feed_id)?;
        gringotts.commission_micro_bps = params.commission_micro_bps;

        lz_receive_types_accounts.gringotts = gringotts.key();

        let register_params = RegisterOAppParams {
            delegate: user.key(),
        };
        let seeds: &[&[u8]] = &[GRINGOTTS_SEED, &[gringotts.bump]];
        oapp::endpoint_cpi::register_oapp(
            ENDPOINT_ID,
            gringotts.key(),
            ctx.remaining_accounts,
            seeds,
            register_params,
        )?;

        Ok(())
    }
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct InitializeParams {
    pub chain_id: u8,
    pub lz_eid: u32,
    pub lz_endpoint_program: Pubkey,
    pub pyth_price_feed_id: String,
    pub commission_micro_bps: u32,
}
