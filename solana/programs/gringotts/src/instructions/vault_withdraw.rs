use crate::state::Gringotts;
use crate::*;
use anchor_lang::context::Context;
use anchor_lang::prelude::{Account, AccountInfo, Program, Pubkey, Signer, System, SystemAccount};
use anchor_lang::system_program::transfer;
use anchor_lang::{system_program, Accounts, AnchorDeserialize, AnchorSerialize};

#[derive(Accounts)]
pub struct VaultWithdraw<'info> {
    #[account(has_one = owner, seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,
    pub owner: Signer<'info>,

    /// CHECK: recipient address
    pub recipient: SystemAccount<'info>,

    #[account(mut, seeds = [VAULT_SEED], bump)]
    pub vault: SystemAccount<'info>,

    pub system_program: Program<'info, System>,
}

impl VaultWithdraw<'_> {
    pub fn apply(ctx: Context<VaultWithdraw>, params: &VaultWithdrawParams) -> Result<()> {
        let system_program = &ctx.accounts.system_program;
        let vault = &ctx.accounts.vault;
        let recipient = &ctx.accounts.recipient;

        let seeds = &[VAULT_SEED, &[ctx.bumps.vault]];
        let signer_seeds = &[&seeds[..]];

        transfer(
            CpiContext::new_with_signer(
                system_program.to_account_info(),
                system_program::Transfer {
                    from: vault.to_account_info(),
                    to: recipient.to_account_info(),
                },
                signer_seeds,
            ),
            params.amount,
        )?;
        Ok(())
    }
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct VaultWithdrawParams {
    pub amount: u64,
}
