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

const ARB_ADDRESS = "0x7794d4260bf7c0c975dd0df59c4f67c1631eea51";
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
export const LZ_RECEIVE_TYPES_SEED = 'LzReceiveTypes'

const [gringottsPDA] = PublicKey.findProgramAddressSync([Buffer.from(GRINGOTTS_SEED)], program.programId);
const [peerPDA] = PublicKey.findProgramAddressSync([Buffer.from(PEER_SEED), new BN(ARB_EID).toArrayLike(Buffer, 'be', 4)], program.programId);
const [selfPeerPDA] = PublicKey.findProgramAddressSync([Buffer.from(PEER_SEED), new BN(SOL_EID).toArrayLike(Buffer, 'be', 4)], program.programId);
const [vaultPDA] = PublicKey.findProgramAddressSync([Buffer.from(VAULT_SEED)], program.programId);

const endpointProgram = new EndpointProgram.Endpoint(new PublicKey('76y77prsiCMvXMjuoZ5VRrhG5qYBrUMYTE5WgHqgjEn6'));
const ulnProgram = new UlnProgram.Uln(new PublicKey('7a4WjyR8VZ7yZz5XJAKm39BUGn5iT9CKcv2pmG9tdXVH'));
const USDC = new PublicKey(SOL_USDC);

const [lzReceiveTypesPDA] = PublicKey.findProgramAddressSync(
    [Buffer.from(LZ_RECEIVE_TYPES_SEED), gringottsPDA.toBytes()], program.programId
);

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
                        message: Buffer.from([]),
                        executionGas: new BN(50_000),
                        items: [
                            {

                                distributionBp: 10_000,
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
    console.log("Vault", vaultPDA.toBase58());

    const tx1 = await program.methods.tokenWithdraw({
        amount: new BN(0)
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

    console.log(tx1);

    const tx2 = await program.methods.vaultWithdraw({
        amount: new BN(0)
    }).accounts({
        gringotts: gringottsPDA,
        recipient: wallet.publicKey,
        vault: vaultPDA,
        systemProgram: SystemProgram.programId
    }).rpc()

    console.log(tx2);
}

async function fund_token() {
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

    const tx1 = await program.methods.tokenFund({
        amount: new BN(7 * 1000 * 1000)
    }).accounts({
        gringotts: gringottsPDA,
        tokenMint: USDC,
        gringottsTokenAccount: gringottsUSDC,
        ownerTokenAccount: userUSDC,
        associatedTokenProgram: ASSOCIATED_TOKEN_PROGRAM_ID,
        tokenProgram: TOKEN_PROGRAM_ID,
        systemProgram: SystemProgram.programId
    }).rpc()

    console.log(tx1);
}


async function lzReceive() {
    const remaining_accounts = [
        {pubkey: new PublicKey('5Wn9HhqZPrr6hps4s6DPxSboHmjVwHPpW9CrX4gmT8tj'), isSigner: false, isWritable: false},
        {pubkey: new PublicKey('2HB1S1fSkPQuwA94rcAF7r7FYkFh1wcMZwBhUGWhGVsT'), isSigner: false, isWritable: false},
        {pubkey: new PublicKey('11111111111111111111111111111111'), isSigner: false, isWritable: false},
        {pubkey: new PublicKey('11111111111111111111111111111111'), isSigner: false, isWritable: false},
        {pubkey: new PublicKey('11111111111111111111111111111111'), isSigner: false, isWritable: false},
        {pubkey: new PublicKey('11111111111111111111111111111111'), isSigner: false, isWritable: false},
    ]

    const tx = await program.methods.lzReceiveTypes({
        srcEid: ARB_EID,
        sender: Array.from(utils.arrayify(utils.hexZeroPad(ARB_ADDRESS, 32))),
        nonce: new BN(2),
        guid: Array.from(utils.arrayify(utils.hexZeroPad("0x0", 32))),
        message: Buffer.from("01b4e4388606dac6dd544f6f1bcf5386772830e1f45429d29adc29392d25474b6d", 'hex'),
        extraData: Buffer.from([]),
    }).accounts({
        gringotts: gringottsPDA,
    }).remainingAccounts(
        remaining_accounts
    ).simulate()

    console.log(tx);

    const buffer = Buffer.from("HAAAAPidux6zXxY4wQQz/JjPjTmY5wuiK+13E+mcWLqo2RAfAAADdKiL7XsEzQw2puUEDCgopGiJ+j38MWZdWbAzwTlFbgABCg2Qzzj8c3oYttcm55/JCOJTxdFq5nPpdI457i5EqbsAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABDD6JoVOqAhHN5xZasSwKIYb6aPxDcrChAvnB+btSOpAABjJclj04kifG7PRApFI4NgwtaE5na/xCEBI572Nvp+FkAAAbd9uHXZaGT2cvhRs7reawctIXtX1s3kTqM9YV+/wCpAAAEedVb8jHAbu50xW7OaBUH/bGy3qP0jlECsc2iVrwTjwAAxvp6877brTo9ZfNqq8l0MbG75MLS9uDkfKYCA0UvXWEAAHiVAvl9H7GFZ7d4HP1MTDYNprZTOh2wB7E0RYTcYwJ0AAGqPzceY63pkSj/U/X8ayKf+OxYJ0v56Yuzm5raJi/n+QABBG2ERlJmodMAh5iltJm6wSMmxLFb+0f/8KVJ60L+K2gAAQan1RcYe9FmNdrUBFX9wsDBJMaPIVZ1pdu6y18IAAAAAADvMfvoVRFcYk6bNdPqPhrFgaSyE0hzRYaeD2D1mdjMZgABBpuIV/6rgYT7aH9jRhjANdrEOdwa6ztVmKDwAAAAAAEAALQ/+if11/ZKdMCbHylYed5LCas238ndUUsyGqezjOXoAAAGm+huya9l60phT9mbjpJUfaAUX6tegErbWU2z5zonGwAApdh0TD0K2EXg1EaJuOpAfML0kdG7eTH1XNQIEOVX4aoAASjXk4wmhSaJKChJh92bRkCbvx4i6lYoXWU6f7mr7M4DAAFCn3DU4BdktyVZNlRs+/Du29imqIfOoBhRUJYyxcnNZQABWq122lFLbh3PEQN+kE2sPTdfUlyfuvyxlQe3iQfYwYsAAPidux6zXxY4wQQz/JjPjTmY5wuiK+13E+mcWLqo2RAfAADA6ARO8LAojYEzfBFngaGxAmrxhJH6SygwygggizoFvgAAKZ+ew2wLxxDmgf2nFKEqMg9uTp9DeB9MowVT4FfPsFEAAUfUA2RchS6eyAyTDh9jUzO8JojWHbxO3ZasmZ/B9JYAAAEcXq5fqohHjocpvEUFfrlRVPD37Eo+yAUDxiuLfN2XzAAB0d2GrDYbYlLEBsKB85EqsTuSQSbAEbWHJ46K8LCO8JsAAFqtdtpRS24dzxEDfpBNrD03X1Jcn7r8sZUHt4kH2MGLAAA=", 'base64')
    // const lzAccounts = program.coder.types.decode('lzAccount', buffer);

    let offset = 4;
    const lzAccounts = [];
    for (let i = 0; i < (buffer.length - 4) / 34; i++) {
        const chunk = buffer.slice(offset, offset + 34);
        const decoded = program.coder.types.decode('lzAccount', chunk);
        lzAccounts.push(decoded);
        offset += 34;
    }

    for(let i = 0; i < lzAccounts.length; i++) {
        console.log(lzAccounts[i].pubkey.toBase58(), lzAccounts[i].isWritable);
    }
}

async function close(pda: string) {
    const tx = await program.methods.accountDestroy({})
        .accounts({
            gringottsPDA: gringottsPDA,
            account: new PublicKey(pda),
        }).rpc();
    console.log(tx);
}

async function main() {
    console.log(gringottsPDA.toBase58());
    console.log(vaultPDA.toBase58());

    // await createAddressLookup();
    // await bridge(ARB_EID);
    await widthdraw_token();
    // await lzReceive();

    // console.log(
    //     "Value", JSON.stringify(await program.account.lzReceiveTypesAccounts.fetch(lzReceiveTypesPDA))
    // );
    await close("HjVcDEcpjVvzwEnaHkZgdrV2pV2DoFcg7TyemtMovuZg");
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
    });