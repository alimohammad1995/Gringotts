use crate::state::Gringotts;
use crate::*;

#[derive(Accounts)]
#[instruction(params: GringottsUpdateParams)]
pub struct GringottsUpdate<'info> {
    #[account(mut)]
    pub owner: Signer<'info>,

    #[account(has_one = owner, seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,

    pub system_program: Program<'info, System>,
}

impl GringottsUpdate<'_> {
    pub fn apply(ctx: Context<GringottsUpdate>, params: &GringottsUpdateParams) -> Result<()> {
        let gringotts = &mut ctx.accounts.gringotts;

        gringotts.commission_micro_bps = params.commission_micro_bps;
        gringotts.commission_discount_bps = params.commission_discount_bps;
        gringotts.gas_discount_bps = params.gas_discount_bps;

        Ok(())
    }
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct GringottsUpdateParams {
    pub commission_micro_bps: u32,
    pub commission_discount_bps: u32,
    pub gas_discount_bps: u32,
}
