use crate::state::Gringotts;
use anchor_lang::prelude::*;
use anchor_lang::system_program;

#[derive(Accounts)]
pub struct Destroy<'info> {
    #[account(
        mut,
        close = owner,
        seeds = [b"Gringotts"],
        bump,
        has_one = owner
    )]
    pub gringotts: Account<'info, Gringotts>,
    #[account(mut)]
    pub owner: Signer<'info>,
    pub system_program: Program<'info, System>,
}

pub fn destroy(ctx: Context<Destroy>) -> Result<()> {
    let gringotts_account_info = ctx.accounts.gringotts.to_account_info();
    let owner_account_info = ctx.accounts.owner.to_account_info();

    **owner_account_info.try_borrow_mut_lamports()? +=
        **gringotts_account_info.try_borrow_lamports()?;
    **gringotts_account_info.try_borrow_mut_lamports()? = 0;

    gringotts_account_info.assign(&system_program::ID);

    Ok(())
}
