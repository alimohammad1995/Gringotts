use crate::state::Gringotts;
use crate::*;

#[derive(Accounts)]
pub struct Test<'info> {
    #[account(has_one = owner, seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,
    #[account(mut, seeds = [VAULT_SEED], bump)]
    pub vault: SystemAccount<'info>,
    #[account(mut)]
    /// CHECK: will be init
    pub fc_item: AccountInfo<'info>,
    #[account(mut)]
    pub owner: Signer<'info>,
    pub system_program: Program<'info, System>,
}

impl Test<'_> {
    pub fn apply(ctx: Context<Test>, params: &TestParams) -> Result<()> {
        if params.close {
            close_order_account(
                &ctx.accounts.fc_item,
                &ctx.accounts.vault,
            )?;
        } else {
            create_order(
                &ctx.accounts.fc_item,
                &ctx.accounts.vault,
                ctx.bumps.vault,
                &ctx.accounts.system_program,
                Vec::new(),
                &Pubkey::default().to_bytes(),
                1,
            )?;
        }

        Ok(())
    }
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct TestParams {
    pub fc_bump: u8,
    pub close: bool,
}
