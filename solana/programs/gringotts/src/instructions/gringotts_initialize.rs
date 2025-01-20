use crate::constants::MAX_QUEUE_SIZE;
use crate::{
    state::{Gringotts, LzReceiveTypesAccounts},
    GRINGOTTS_SEED, LZ_RECEIVE_TYPES_SEED, VAULT_SEED,
};
use anchor_lang::prelude::*;
use anchor_lang::system_program;
use oapp::endpoint::{instructions::RegisterOAppParams, ID as ENDPOINT_ID};
use pyth_solana_receiver_sdk::price_update::get_feed_id_from_hex;

#[derive(Accounts)]
#[instruction(params: GringottsInitializeParams)]
pub struct GringottsInitialize<'info> {
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
    #[account(mut, seeds = [VAULT_SEED], bump)]
    pub vault: SystemAccount<'info>,

    pub system_program: Program<'info, System>,
}

impl GringottsInitialize<'_> {
    pub fn apply(ctx: Context<GringottsInitialize>, params: &GringottsInitializeParams) -> Result<()> {
        let user = &ctx.accounts.owner;

        let gringotts = &mut ctx.accounts.gringotts;
        let lz_receive_types_accounts = &mut ctx.accounts.lz_receive_types_accounts;
        let system_program = &ctx.accounts.system_program;
        let vault = &ctx.accounts.vault;

        gringotts.owner = user.key();
        gringotts.bump = ctx.bumps.gringotts;
        gringotts.chain_id = params.chain_id;
        gringotts.lz_endpoint_program = params.lz_endpoint_program;
        gringotts.pyth_price_feed_id = get_feed_id_from_hex(&params.pyth_price_feed_id)?;
        gringotts.commission_micro_bps = params.commission_micro_bps;
        gringotts.stable_coins = params.stable_coins.clone();

        lz_receive_types_accounts.gringotts = gringotts.key();
        lz_receive_types_accounts.orders = [Pubkey::default(); MAX_QUEUE_SIZE];

        system_program::transfer(
            CpiContext::new(
                system_program.to_account_info(),
                system_program::Transfer {
                    from: user.to_account_info(),
                    to: vault.to_account_info(),
                },
            ),
            params.vault_fund,
        )?;

        oapp::endpoint_cpi::register_oapp(
            ENDPOINT_ID,
            gringotts.key(),
            ctx.remaining_accounts,
            &[GRINGOTTS_SEED, &[gringotts.bump]],
            RegisterOAppParams {
                delegate: user.key(),
            },
        )?;

        Ok(())
    }
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct GringottsInitializeParams {
    pub chain_id: u8,
    pub lz_endpoint_program: Pubkey,
    pub pyth_price_feed_id: String,
    pub vault_fund: u64,
    pub commission_micro_bps: u32,
    pub stable_coins: Vec<[u8; 32]>,
}
