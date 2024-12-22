use crate::msg_codec::{ChainTransfer, ChainTransferItem, Message, CHAIN_TRANSFER_TYPE};
use crate::state::{Gringotts, Peer};
use crate::utils::{bps, change_decimals, OptionsBuilder};
use crate::*;
use anchor_lang::context::Context;
use anchor_lang::prelude::{Account, Program, Result, System};
use anchor_lang::solana_program::entrypoint::ProgramResult;
use anchor_lang::solana_program::instruction::Instruction;
use anchor_lang::solana_program::program::invoke_signed;
use anchor_lang::{system_program, Accounts, AnchorDeserialize, AnchorSerialize};
use anchor_spl::associated_token;
use anchor_spl::associated_token::AssociatedToken;
use anchor_spl::token::{self, Mint, Token, TokenAccount, Transfer};
use oapp::endpoint::instructions::SendParams;
use pyth_solana_receiver_sdk::price_update::PriceUpdateV2;

#[derive(Accounts)]
pub struct Bridge<'info> {
    #[account(mut)]
    pub user: Signer<'info>,

    pub price_feed: Account<'info, PriceUpdateV2>,

    #[account(seeds = [GRINGOTTS_SEED], bump = gringotts.bump)]
    pub gringotts: Account<'info, Gringotts>,

    #[account(mut, seeds = [VAULT_SEED], bump)]
    pub vault: SystemAccount<'info>,

    #[account(seeds = [PEER_SEED, &self_peer.lz_eid.to_le_bytes()], bump = self_peer.bump)]
    pub self_peer: Account<'info, Peer>,

    #[account(
        init_if_needed, payer = user,
        associated_token::mint = stable_coin_mint,
        associated_token::authority = gringotts,
    )]
    pub gringotts_stable_coin: Account<'info, TokenAccount>,
    pub stable_coin_mint: Account<'info, Mint>,

    /// CHECK: This is the swap program
    #[account(executable)]
    pub swap_program: AccountInfo<'info>,
    pub associated_token_program: Program<'info, AssociatedToken>,
    pub token_program: Program<'info, Token>,
    pub system_program: Program<'info, System>,
}

impl<'info> Bridge<'info> {
    pub fn apply(
        ctx: Context<'_, '_, 'info, 'info, Bridge<'info>>,
        params: &BridgeRequest,
    ) -> Result<BridgeResponse> {
        let signer = &ctx.accounts.user;

        let gringotts = &ctx.accounts.gringotts;
        let vault_account = &ctx.accounts.vault;
        let vault_bump: u8 = ctx.bumps.vault;

        let price_feed = &ctx.accounts.price_feed;

        let swap_program = &ctx.accounts.swap_program;
        let system_program = &ctx.accounts.system_program;
        let token_program = &ctx.accounts.token_program;
        let associated_token_program = &ctx.accounts.associated_token_program;

        let gringotts_stable_coin = &mut ctx.accounts.gringotts_stable_coin;
        let stable_coin_mint = &ctx.accounts.stable_coin_mint;
        let self_peer = &ctx.accounts.self_peer;

        let remaining_accounts = ctx.remaining_accounts;
        let mut r = 0;
        let mut m = 0;

        let mut net_usdx = 0u64;

        require!(
            self_peer.chain_id == gringotts.chain_id,
            BridgeErrorCode::InvalidParams
        );

        /*********** [Inbound transaction] ***********/
        for i in 0..params.inbound.items.len() {
            let item = &params.inbound.items[i];
            let first_stable_coin_amount = gringotts_stable_coin.amount;

            if self_peer.stable_coins.contains(&item.asset) && stable_coin_mint.key().to_bytes() == item.asset {
                let user_token_account = &remaining_accounts[r];

                let cpi_accounts = Transfer {
                    from: user_token_account.to_account_info(),
                    to: gringotts_stable_coin.to_account_info(),
                    authority: signer.to_account_info(),
                };

                let cpi_ctx = CpiContext::new(token_program.to_account_info(), cpi_accounts);
                token::transfer(cpi_ctx, item.amount)?;

                r += 1;
            } else {
                require!(!item.swap.is_none(), BridgeErrorCode::InvalidParams);

                let token_mint: &Account<Mint> = &Account::try_from(&remaining_accounts[r])?;
                let token_account = &remaining_accounts[r + 1];

                init_token_account_if_needed(
                    token_account,
                    &gringotts.to_account_info(),
                    token_mint,
                    signer,
                    system_program,
                    token_program,
                    associated_token_program,
                )?;

                let mut gringotts_token: Account<TokenAccount> = Account::try_from(token_account)?;
                require!(gringotts_token.owner.key() == gringotts.key(), BridgeErrorCode::InvalidParams);

                if item.asset == [0; 32] {
                    let seeds = &[GRINGOTTS_SEED, &[gringotts.bump]];
                    let signer_seeds = &[&seeds[..]];

                    system_program::transfer(
                        CpiContext::new(
                            system_program.to_account_info(),
                            system_program::Transfer {
                                from: signer.to_account_info(),
                                to: gringotts_token.to_account_info(),
                            },
                        ),
                        item.amount,
                    )?;

                    token::sync_native(CpiContext::new_with_signer(
                        gringotts_token.to_account_info(),
                        token::SyncNative {
                            account: gringotts_token.to_account_info(),
                        },
                        signer_seeds,
                    ))?;
                } else {
                    let user_token_account = remaining_accounts[r + 2].to_account_info();
                    let cpi_ctx = CpiContext::new(
                        token_program.to_account_info(),
                        Transfer {
                            from: user_token_account,
                            to: gringotts_token.to_account_info(),
                            authority: signer.to_account_info(),
                        },
                    );

                    token::transfer(cpi_ctx, item.amount)?;
                }

                let swap = item.swap.as_ref().unwrap();
                let swap_account_counts = swap.metadata[m] as usize;

                r += if item.asset == [0; 32] { 2 } else { 3 };
                m += 1;

                let before_swap_amount = gringotts_token.amount;

                swap_on_jupiter(
                    gringotts,
                    swap_program,
                    &remaining_accounts[r..r + swap_account_counts],
                    swap.command.clone(),
                )?;

                gringotts_token.reload()?;

                require!(
                    gringotts_token.amount - before_swap_amount <= item.amount,
                    BridgeErrorCode::InvalidSwapAmount
                );

                r += swap_account_counts;
            }

            gringotts_stable_coin.reload()?;

            net_usdx += change_decimals(
                gringotts_stable_coin.amount - first_stable_coin_amount,
                stable_coin_mint.decimals,
                CHAIN_TRANSFER_DECIMALS,
            );
        }

        require!(net_usdx >= params.inbound.amount_usdx, BridgeErrorCode::InvalidSwapAmount);

        /*********** [Estimate transaction] ***********/
        let mut estimate_outbounds = Vec::with_capacity(params.outbounds.len());
        let mut peers = Vec::with_capacity(params.outbounds.len());

        for i in 0..params.outbounds.len() {
            let mut items = Vec::with_capacity(params.outbounds[i].items.len());

            for j in 0..params.outbounds[i].items.len() {
                let mut command_length = 0;
                let mut metadata_length = 0;

                if let Some(swap) = params.outbounds[i].items[j].swap.as_ref() {
                    command_length = swap.command.len() as u16;
                    metadata_length = swap.metadata.len() as u16;
                }

                let item = EstimateOutboundTransferItem {
                    asset: params.outbounds[i].items[j].asset,
                    execution_gas_amount: params.outbounds[i].items[j].execution_gas_amount,
                    execution_command_length: command_length,
                    execution_metadata_length: metadata_length,
                };

                items.push(item);
            }

            estimate_outbounds.push(EstimateOutboundTransfer {
                chain_id: params.outbounds[i].chain_id,
                items: items,
            });

            let mut data = &remaining_accounts[r + i].try_borrow_data()?[..];
            let peer = Peer::try_deserialize(&mut data)?;
            peers.push(peer);
        }

        r = r + params.outbounds.len();

        let estimate_request = &EstimateRequest {
            inbound: EstimateInboundTransfer {
                amount_usdx: net_usdx,
            },
            outbounds: estimate_outbounds,
        };

        let estimate_result = estimate_marketplace(
            estimate_request,
            gringotts,
            peers.as_slice(),
            price_feed,
            &remaining_accounts[r..],
            true,
        )?;

        net_usdx = net_usdx - estimate_result.transfer_gas_price_usdx;
        net_usdx = net_usdx - estimate_result.commission_usdx;

        /*********** [Send transfers] ***********/
        let mut message_ids = Vec::with_capacity(params.outbounds.len());
        let vault_seeds: &[&[u8]] = &[VAULT_SEED, &[vault_bump]];
        let gringotts_seeds: &[&[u8]] = &[GRINGOTTS_SEED, &[gringotts.bump]];

        for i in 0..params.outbounds.len() {
            let mut chain_transfers = Vec::with_capacity(params.outbounds[i].items.len());

            for j in 0..params.outbounds[i].items.len() {
                let amount_usdx = bps(
                    net_usdx,
                    params.outbounds[i].items[j].distribution_bp as u32,
                );

                let mut executor = [0; 32];
                let mut stable_token = [0; 32];
                let mut command: Vec<u8> = Vec::new();
                let mut metadata: Vec<u8> = Vec::new();

                if let Some(swap) = params.outbounds[i].items[j].swap.as_ref() {
                    executor = swap.executor;
                    stable_token = swap.stable_token;
                    command = swap.command.clone();
                    metadata = swap.metadata.clone();
                }

                chain_transfers.push(ChainTransferItem {
                    amount_usdx: amount_usdx,
                    asset: params.outbounds[i].items[j].asset,
                    recipient: params.outbounds[i].items[j].recipient,
                    executor: executor,
                    stable_token: stable_token,
                    command: command,
                    metadata: metadata,
                });
            }

            let chain_transfer = ChainTransfer {
                items: chain_transfers,
            };

            let mut builder = OptionsBuilder::new();
            builder.add_executor_lz_receive_option(
                estimate_result.outbound_metadata[i].execution_gas_amount,
                0,
            );

            let message = Message::new(CHAIN_TRANSFER_TYPE, chain_transfer.encode());

            let send_params = SendParams {
                dst_eid: peers[i].lz_eid,
                receiver: peers[i].address,
                message: message.encode(),
                options: builder.options(),
                native_fee: estimate_result.outbound_metadata[i].transfer_gas_amount * 2,
                lz_token_fee: 0,
            };

            let lz_send_accounts: Vec<AccountInfo> = remaining_accounts
                [r + (i * LZ_SEND_ACCOUNTS_LEN)..r + ((i + 1) * LZ_SEND_ACCOUNTS_LEN)]
                .iter()
                .map(|acc| AccountInfo {
                    is_signer: acc.is_signer || acc.key() == vault_account.key(),
                    ..acc.clone()
                })
                .collect();

            let send_response = oapp::endpoint_cpi::send(
                gringotts.lz_endpoint_program,
                gringotts.key(),
                lz_send_accounts.as_slice(),
                &[vault_seeds, gringotts_seeds],
                send_params,
            )?;

            message_ids.push(
                send_response
                    .guid
                    .iter()
                    .map(|byte| format!("{:02x}", byte))
                    .collect::<String>(),
            );
        }

        Ok(BridgeResponse {
            message_ids: message_ids,
        })
    }
}

pub fn swap_on_jupiter<'info>(
    gringotts: &Account<'info, Gringotts>,
    swap_program: &AccountInfo<'info>,
    remaining_accounts: &[AccountInfo],
    data: Vec<u8>,
) -> ProgramResult {
    let accounts: Vec<AccountMeta> = remaining_accounts
        .iter()
        .map(|acc| AccountMeta {
            pubkey: *acc.key,
            is_signer: *acc.key == gringotts.key(),
            is_writable: acc.is_writable,
        })
        .collect();

    let accounts_infos: Vec<AccountInfo> = remaining_accounts
        .iter()
        .map(|acc| AccountInfo { ..acc.clone() })
        .collect();

    let seeds = &[GRINGOTTS_SEED, &[gringotts.bump]];
    let signer_seeds = &[&seeds[..]];

    invoke_signed(
        &Instruction {
            program_id: swap_program.key(),
            accounts,
            data,
        },
        &accounts_infos,
        signer_seeds,
    )
}

fn init_token_account_if_needed<'info>(
    token_account: &AccountInfo<'info>,
    authority: &AccountInfo<'info>,
    token_mint: &Account<'info, Mint>,
    signer: &Signer<'info>,
    system_program: &Program<'info, System>,
    token_program: &Program<'info, Token>,
    associated_token_program: &Program<'info, AssociatedToken>,
) -> Result<()> {
    let ata = associated_token::get_associated_token_address(&authority.key(), &token_mint.key());

    require!(
        ata.key() == token_account.key(),
        BridgeErrorCode::InvalidParams
    );

    if token_account.data_is_empty() {
        let cpi_accounts = associated_token::Create {
            payer: signer.to_account_info(),
            associated_token: token_account.clone(),
            authority: authority.to_account_info(),
            mint: token_mint.to_account_info(),
            system_program: system_program.to_account_info(),
            token_program: token_program.to_account_info(),
        };

        let cpi_context = CpiContext::new(associated_token_program.to_account_info(), cpi_accounts);
        associated_token::create(cpi_context)?;
    }

    Ok(())
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct Swap {
    executor: [u8; 32],
    command: Vec<u8>,
    metadata: Vec<u8>,
    stable_token: [u8; 32],
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct BridgeInboundTransferItem {
    pub asset: [u8; 32],
    pub amount: u64,
    pub swap: Option<Swap>,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct BridgeInboundTransfer {
    pub amount_usdx: u64,
    pub items: Vec<BridgeInboundTransferItem>,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct BridgeOutboundTransferItem {
    pub asset: [u8; 32],
    pub recipient: [u8; 32],
    pub execution_gas_amount: u64,
    pub distribution_bp: u16,
    pub swap: Option<Swap>,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct BridgeOutboundTransfer {
    pub chain_id: u8,
    pub items: Vec<BridgeOutboundTransferItem>,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct BridgeRequest {
    pub inbound: BridgeInboundTransfer,
    pub outbounds: Vec<BridgeOutboundTransfer>,
}

#[derive(Clone, AnchorSerialize, AnchorDeserialize)]
pub struct BridgeResponse {
    pub message_ids: Vec<String>,
}

#[error_code]
pub enum BridgeErrorCode {
    #[msg("Invalid params.")]
    InvalidParams,
    #[msg("Invalid swap amount.")]
    InvalidSwapAmount,
}
