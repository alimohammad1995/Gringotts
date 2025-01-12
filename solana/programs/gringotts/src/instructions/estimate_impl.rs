use crate::instructions::{
    EstimateErrorCode, EstimateOutboundDetails, EstimateRequest, EstimateResponse,
};
use crate::state::{Gringotts, Peer};
use crate::utils::{bps, change_decimals, micro_bps, OptionsBuilder};
use crate::constants::{CHAIN_TRANSFER_DECIMALS, MAX_PRICE_AGE, NETWORK_DECIMALS};
use anchor_lang::prelude::{Account, AccountInfo, Clock, Result, SolanaSysvar};
use anchor_lang::{require, Key};
use oapp::endpoint::instructions::QuoteParams;
use pyth_solana_receiver_sdk::price_update::PriceUpdateV2;

pub const LZ_QUOTE_ACCOUNTS_LEN: usize = 18;
pub const LZ_SEND_ACCOUNTS_LEN: usize = 26;

pub fn estimate_marketplace(
    request: &EstimateRequest,
    gringotts: &Account<Gringotts>,
    peers: &[Peer],
    price_feed: &PriceUpdateV2,
    quote_accounts: &[AccountInfo],
    use_send_accounts: bool,
) -> Result<EstimateResponse> {
    let mut total_transfer_gas_price = 0u64;

    let mut outbound_metadata = Vec::with_capacity(request.outbounds.len());

    for i in 0usize..request.outbounds.len() {
        let outbound = &request.outbounds[i];
        let peer = &peers[i];

        require!(outbound.chain_id == peer.chain_id, EstimateErrorCode::InvalidParams);

        let chain_execution_gas_price = peer.base_gas_estimate + outbound.execution_gas;

        let mut builder = OptionsBuilder::new();
        builder.add_executor_lz_receive_option(chain_execution_gas_price, 0);

        let quote_params = QuoteParams {
            sender: gringotts.key(),
            dst_eid: peer.lz_eid,
            receiver: peer.address,
            message: vec![0; outbound.message_length as usize],
            pay_in_lz_token: false,
            options: builder.clone().options(),
        };

        let oapp_quote_accounts;

        if use_send_accounts {
            oapp_quote_accounts = generate_quote_remaining_account(
                &quote_accounts[i * LZ_SEND_ACCOUNTS_LEN..(i + 1) * LZ_SEND_ACCOUNTS_LEN],
            );
        } else {
            oapp_quote_accounts = quote_accounts[i * LZ_QUOTE_ACCOUNTS_LEN..(i + 1) * LZ_QUOTE_ACCOUNTS_LEN].to_vec();
        }

        let quote = oapp::endpoint_cpi::quote(
            gringotts.lz_endpoint_program,
            oapp_quote_accounts.as_slice(),
            quote_params,
        )?;

        outbound_metadata.push(EstimateOutboundDetails {
            chain_id: outbound.chain_id,
            execution_gas: chain_execution_gas_price,
            transfer_gas: quote.native_fee,
        });

        total_transfer_gas_price = total_transfer_gas_price + quote.native_fee;
    }

    let commission_usdx = micro_bps(request.inbound.amount_usdx, gringotts.commission_micro_bps);
    let transfer_gas_usdx = get_native_price(total_transfer_gas_price, CHAIN_TRANSFER_DECIMALS, gringotts.pyth_price_feed_id, &price_feed)?;

    let commission_discount_usdx = bps(commission_usdx, gringotts.commission_discount_bps);
    let transfer_gas_discount_usdx = bps(transfer_gas_usdx, gringotts.gas_discount_bps);

    require!(
        request.inbound.amount_usdx >= (commission_usdx + transfer_gas_usdx) - (commission_discount_usdx + transfer_gas_discount_usdx),
        EstimateErrorCode::InvalidParams
    );

    Ok(EstimateResponse {
        commission_usdx: commission_usdx,
        commission_discount_usdx: commission_discount_usdx,
        transfer_gas_usdx: transfer_gas_usdx,
        transfer_gas_discount_usdx: transfer_gas_discount_usdx,
        outbound_details: outbound_metadata,
    })
}

pub fn get_native_price(
    amount: u64,
    decimals: u8,
    feed_id: [u8; 32],
    price_feed: &PriceUpdateV2,
) -> Result<u64> {
    let price = price_feed.get_price_no_older_than(&Clock::get()?, MAX_PRICE_AGE, &feed_id)?;
    let asset_price = amount * (u64::try_from(price.price)? + price.conf) / 10u64.pow(u32::try_from(-price.exponent)?);

    Ok(change_decimals(asset_price, NETWORK_DECIMALS, decimals))
}

pub fn generate_quote_remaining_account<'a>(
    input_list: &[AccountInfo<'a>],
) -> Vec<AccountInfo<'a>> {
    let mut accounts = vec![
        input_list[0].clone(),  // Index 0 -> input_list[0]
        input_list[2].clone(),  // Index 1 -> input_list[2]
        input_list[3].clone(),  // Index 2 -> input_list[3]
        input_list[4].clone(),  // Index 3 -> input_list[4]
        input_list[5].clone(),  // Index 4 -> input_list[5]
        input_list[6].clone(),  // Index 5 -> input_list[6]
        input_list[7].clone(),  // Index 6 -> input_list[7]
        input_list[10].clone(), // Index 7 -> input_list[10]
        input_list[11].clone(), // Index 8 -> input_list[11]
        input_list[12].clone(), // Index 9 -> input_list[12]
        input_list[18].clone(), // Index 10 -> input_list[18]
        input_list[19].clone(), // Index 11 -> input_list[19]
        input_list[20].clone(), // Index 12 -> input_list[20]
        input_list[21].clone(), // Index 13 -> input_list[21]
        input_list[22].clone(), // Index 14 -> input_list[22]
        input_list[23].clone(), // Index 15 -> input_list[23]
        input_list[20].clone(), // Index 16 -> input_list[20]
        input_list[21].clone(), // Index 17 -> input_list[21]
    ];

    for i in 0..accounts.len() {
        accounts[i].is_writable = false;
        accounts[i].is_signer = false;
    }

    accounts
}
