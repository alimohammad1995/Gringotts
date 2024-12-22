import * as anchor from "@coral-xyz/anchor";
import {BN, Program} from "@coral-xyz/anchor";
import {Gringotts} from "../target/types/gringotts";
import {
    AddressLookupTableProgram,
    ComputeBudgetProgram,
    PublicKey,
    sendAndConfirmTransaction,
    SystemProgram,
    Transaction,
    TransactionMessage,
    VersionedTransaction,
} from '@solana/web3.js';
import {utils} from "ethers";
import {EndpointProgram, UlnProgram} from '@layerzerolabs/lz-solana-sdk-v2'
import {PacketPath} from '@layerzerolabs/lz-v2-utilities'
import NodeWallet from "@coral-xyz/anchor/dist/cjs/nodewallet";
import {ASSOCIATED_TOKEN_PROGRAM_ID, getAssociatedTokenAddressSync, NATIVE_MINT, TOKEN_PROGRAM_ID} from "@solana/spl-token";

const ARB_EID = 30110;
const SOL_EID = 30168;
const ARB_CHAIN_ID = 1;
const SOL_CHAIN_ID = 2;

const ARB_ADDRESS = "0x9c4E6e7e2f2387c3fd9fccc499c18D6c98528931";
const ARB_USDC = "0xaf88d065e77c8cc2239327c5edb3a432268e5831";
const SOL_USDC = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v";

const ALT = "CHgiodkqryp7jqKTrETT45MLG28ukg42NFBZWp26SCJ5";

const provider = anchor.AnchorProvider.env();
anchor.setProvider(anchor.AnchorProvider.env());

const program = anchor.workspace.Gringotts as Program<Gringotts>;
const wallet = provider.wallet as NodeWallet;

export const GRINGOTTS_SEED = 'Gringotts'
export const PEER_SEED = 'Peer'
export const VAULT_SEED = 'Vault'

const [gringottsPDA] = PublicKey.findProgramAddressSync([Buffer.from(GRINGOTTS_SEED)], program.programId);
const [peerPDA] = PublicKey.findProgramAddressSync([Buffer.from(PEER_SEED), new BN(ARB_EID).toArrayLike(Buffer, 'le', 4)], program.programId);
const [selfPeerPDA] = PublicKey.findProgramAddressSync([Buffer.from(PEER_SEED), new BN(SOL_EID).toArrayLike(Buffer, 'le', 4)], program.programId);
const [vaultPDA, x] = PublicKey.findProgramAddressSync([Buffer.from(VAULT_SEED)], program.programId);

const endpointProgram = new EndpointProgram.Endpoint(new PublicKey('76y77prsiCMvXMjuoZ5VRrhG5qYBrUMYTE5WgHqgjEn6'));
const ulnProgram = new UlnProgram.Uln(new PublicKey('7a4WjyR8VZ7yZz5XJAKm39BUGn5iT9CKcv2pmG9tdXVH'));
const USDC = new PublicKey(SOL_USDC);

async function bridge(chainID: number) {
    const computeUnitInstruction = ComputeBudgetProgram.setComputeUnitLimit({
        units: 1_000_000_000
    });
    const computeUnitPriceInstruction = ComputeBudgetProgram.setComputeUnitPrice({
        microLamports: 600_000
    });


    const gringottsUSDC = getAssociatedTokenAddressSync(
        USDC,
        gringottsPDA,
        true
    );

    const gringottsWSOL = getAssociatedTokenAddressSync(
        NATIVE_MINT,
        gringottsPDA,
        true
    );

    const userUSDC = getAssociatedTokenAddressSync(
        USDC,
        wallet.publicKey
    );

    console.log("gringottsPDA", gringottsPDA.toBase58());
    console.log("vaultPDA", vaultPDA.toBase58());
    console.log("gringottsUSDC", gringottsUSDC.toBase58());
    console.log("gringottsWSOL", gringottsWSOL.toBase58());
    console.log("userUSDC", userUSDC.toBase58());
    // return;

    const packetPath: PacketPath = {
        srcEid: SOL_EID,
        dstEid: chainID,
        sender: utils.hexlify(gringottsPDA.toBytes()),
        receiver: utils.hexlify(ARB_ADDRESS),
    }

    const accountsEx = await endpointProgram.getSendIXAccountMetaForCPI(provider.connection, vaultPDA, packetPath, ulnProgram);

    let accounts = [
        {pubkey: userUSDC, isSigner: false, isWritable: true},
        {pubkey: peerPDA, isSigner: false, isWritable: true},
    ];
    accounts = accounts.concat(accountsEx);
    console.log(accounts.length);

    try {
        const tx = await program.methods
            .bridge({
                inbound: {
                    amountUsdx: new BN(5 * 1000 * 1000),
                    items: [
                        {
                            asset: Array.from(USDC.toBytes()),
                            amount: new BN(5 * 1000 * 1000),
                            swap: null,
                        }
                    ]
                },
                outbounds: [
                    {
                        chainId: ARB_CHAIN_ID,
                        items: [
                            {
                                asset: Array.from([]),
                                recipient: Array.from(utils.arrayify(utils.hexZeroPad("0x0D595AE2666a2c5Ae6b99cce4DD428a9Cf20B2c9", 32))),
                                executionGasAmount: new BN(500_000),
                                distributionBp: 10_000,
                                swap: {
                                    executor: Array.from(utils.arrayify(utils.hexZeroPad("0x6352a56caadC4F1E25CD6c75970Fa768A3304e64", 32))),
                                    command: Buffer.from("bc80f1a80000000000000000000000000d595ae2666a2c5ae6b99cce4dd428a9cf20b2c900000000000000000000000000000000000000000000000000000000002dc6c000000000000000000000000000000000000000000000000000032e0c9554103e00000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001a000640000000000000000006f38e884725a116c9c7fbf208e79fe8828a2595f", 'hex'),
                                    metadata: Buffer.from([]),
                                    stableToken: Array.from(utils.arrayify(utils.hexZeroPad(ARB_USDC, 32))),
                                },
                            }
                        ]
                    }
                ]
            }).accounts({ // TEST IF WE CAN MOVE THEM ALL TO ALT
                priceFeed: new PublicKey("7UVimffxr9ow1uXYxsr4LHAcV58mLzhmwaeKvJ1pjLiE"),
                gringotts: gringottsPDA,
                selfPeer: selfPeerPDA,
                gringottsStableCoin: gringottsUSDC,
                stableCoinMint: USDC,
                vaultPDA: vaultPDA,
                swapProgram: program.programId,
                associatedTokenProgram: ASSOCIATED_TOKEN_PROGRAM_ID,
                tokenProgram: TOKEN_PROGRAM_ID,
                systemProgram: SystemProgram.programId
            })
            .remainingAccounts(accounts)
            .instruction();


        const lookupTable = (
            await provider.connection.getAddressLookupTable(new PublicKey(ALT))
        ).value;

        const messageV0 = new TransactionMessage({
            payerKey: wallet.publicKey,
            recentBlockhash: (await provider.connection.getLatestBlockhash()).blockhash,
            instructions: [computeUnitInstruction, computeUnitPriceInstruction, tx],
        }).compileToV0Message([lookupTable]);
        const transaction = new VersionedTransaction(messageV0);

        const txID = await provider.sendAndConfirm(transaction, [wallet.payer], {
            skipPreflight: false,
        });
        console.log(txID)
    } catch (err) {
        console.log("err", err);
    }
}

async function createAddressLookup() {
    // const [lookupTableInst, lookupTableAddress] =
    //     AddressLookupTableProgram.createLookupTable({
    //         authority: wallet.payer.publicKey,
    //         payer: wallet.payer.publicKey,
    //         recentSlot: await provider.connection.getSlot(),
    //     });
    //
    // console.log("lookup table address:", lookupTableAddress.toBase58());
    // const t1 = new Transaction().add(lookupTableInst);
    // const s1 = await sendAndConfirmTransaction(provider.connection, t1, [wallet.payer], {
    //     skipPreflight: true,
    // });
    //
    // console.log(s1);

    const packetPath: PacketPath = {
        srcEid: SOL_EID,
        dstEid: ARB_EID,
        sender: utils.hexlify(gringottsPDA.toBytes()),
        receiver: utils.hexlify(ARB_ADDRESS),
    }
    const accountsEx = await endpointProgram.getSendIXAccountMetaForCPI(provider.connection, vaultPDA, packetPath, ulnProgram);

    const extendInstruction = AddressLookupTableProgram.extendLookupTable({
        payer: wallet.payer.publicKey,
        authority: wallet.payer.publicKey,
        lookupTable: new PublicKey(ALT),
        addresses: accountsEx.map((a) => a.pubkey),
    });

    const t2 = new Transaction().add(extendInstruction);
    const s2 = await sendAndConfirmTransaction(provider.connection, t2, [wallet.payer], {
        skipPreflight: true,
    });

    console.log(s2);
}

async function widthdraw_token() {
    const gringottsUSDC = getAssociatedTokenAddressSync(
        USDC,
        gringottsPDA,
        true
    );

    const userUSDC = getAssociatedTokenAddressSync(
        USDC,
        wallet.publicKey,
    );

    console.log("User USDC", userUSDC.toBase58());
    console.log("Gringotts USDC", gringottsUSDC.toBase58());

    const tx = await program.methods.tokenWithdraw({
        amount: new BN(7 * 1000 * 1000)
    }).accounts({
        gringotts: gringottsPDA,
        recipient: wallet.publicKey,
        tokenMint: USDC,
        gringottsTokenAccount: gringottsUSDC,
        recipientTokenAccount: userUSDC,
        associatedTokenProgram: ASSOCIATED_TOKEN_PROGRAM_ID,
        tokenProgram: TOKEN_PROGRAM_ID,
        systemProgram: SystemProgram.programId
    }).rpc()

    console.log(tx);
}

async function main() {
    // await createAddressLookup();
    await bridge(ARB_EID);
    await widthdraw_token();
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
    });