use crate::state::Gringotts;
use crate::*;

#[derive(Accounts)]
pub struct AccountDestroy<'info> {
    #[account(has_one = owner, seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,
    /// CHECK: This is a PDA account
    #[account(mut)]
    pub account: AccountInfo<'info>,
    #[account(mut, seeds = [VAULT_SEED], bump)]
    pub vault: SystemAccount<'info>,
    #[account(mut)]
    pub owner: Signer<'info>,
    pub system_program: Program<'info, System>,
}

impl AccountDestroy<'_> {
    pub fn apply(ctx: &Context<AccountDestroy>, _params: &AccountDestroyParams) -> Result<()> {
        let account_info = ctx.accounts.account.to_account_info();
        let vault_account_info = ctx.accounts.vault.to_account_info();

        **vault_account_info.try_borrow_mut_lamports()? += **account_info.try_borrow_lamports()?;
        **account_info.try_borrow_mut_lamports()? = 0;

        Ok(())
    }
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct AccountDestroyParams {}
