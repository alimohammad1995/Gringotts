use crate::state::Gringotts;
use crate::*;
use anchor_lang::context::Context;
use anchor_lang::prelude::{Account, Program, Signer, System, SystemAccount};
use anchor_lang::{Accounts, AnchorDeserialize, AnchorSerialize};
use anchor_spl::associated_token::AssociatedToken;
use anchor_spl::token;
use anchor_spl::token::{Mint, Token, TokenAccount, Transfer};

#[derive(Accounts)]
pub struct TokenWithdraw<'info> {
    #[account( seeds = [b"Gringotts"], bump, has_one = owner)]
    pub gringotts: Account<'info, Gringotts>,
    #[account( mut)]
    pub owner: Signer<'info>,

    /// CHECK: recipient address
    pub recipient: SystemAccount<'info>,

    pub token_mint: Account<'info, Mint>,

    #[account(
        mut,
        associated_token::mint = token_mint,
        associated_token::authority = gringotts,
    )]
    pub gringotts_token_account: Account<'info, TokenAccount>,

    #[account(
        init_if_needed,
        payer = owner,
        associated_token::mint = token_mint,
        associated_token::authority = recipient,
    )]
    pub recipient_token_account: Account<'info, TokenAccount>,

    pub associated_token_program: Program<'info, AssociatedToken>,
    pub token_program: Program<'info, Token>,
    pub system_program: Program<'info, System>,
}

impl TokenWithdraw<'_> {
    pub fn apply(ctx: Context<TokenWithdraw>, params: &TokenWithdrawParams) -> Result<()> {
        let gringotts = &ctx.accounts.gringotts;

        let seeds = &[GRINGOTTS_SEED, &[gringotts.bump]];
        let signer_seeds = &[&seeds[..]];

        let cpi_accounts = Transfer {
            from: ctx.accounts.gringotts_token_account.to_account_info(),
            to: ctx.accounts.recipient_token_account.to_account_info(),
            authority: gringotts.to_account_info(),
        };
        let cpi_program = ctx.accounts.token_program.to_account_info();

        let cpi_ctx = CpiContext::new_with_signer(cpi_program, cpi_accounts, signer_seeds);
        token::transfer(cpi_ctx, params.amount)?;

        Ok(())
    }
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct TokenWithdrawParams {
    pub amount: u64,
}
