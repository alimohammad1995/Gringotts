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

const TEST_STABLE = "0x301b022b40d06088fc974e767149f4a3feebbf1a";
const RECEIVER = "0x1c191f62728b1498d779559e9ffb75a849582103";
const SOLANA_EID = 40168;
const ARB_EID = 40231;

const provider = anchor.AnchorProvider.env();
anchor.setProvider(anchor.AnchorProvider.env());

const program = anchor.workspace.Gringotts as Program<Gringotts>;
const wallet = provider.wallet as NodeWallet;

export const GRINGOTTS_SEED = 'Gringotts'
export const PEER_SEED = 'Peer'
export const VAULT_SEED = 'Vault'

const [gringottsPDA] = PublicKey.findProgramAddressSync([Buffer.from(GRINGOTTS_SEED)], program.programId);
const [peerPDA] = PublicKey.findProgramAddressSync([Buffer.from(PEER_SEED), new BN(ARB_EID).toArrayLike(Buffer, 'le', 4)], program.programId);
const [selfPeerPDA] = PublicKey.findProgramAddressSync([Buffer.from(PEER_SEED), new BN(SOLANA_EID).toArrayLike(Buffer, 'le', 4)], program.programId);
const [vaultPDA, x] = PublicKey.findProgramAddressSync([Buffer.from(VAULT_SEED)], program.programId);

const endpointProgram = new EndpointProgram.Endpoint(new PublicKey('76y77prsiCMvXMjuoZ5VRrhG5qYBrUMYTE5WgHqgjEn6'));
const ulnProgram = new UlnProgram.Uln(new PublicKey('7a4WjyR8VZ7yZz5XJAKm39BUGn5iT9CKcv2pmG9tdXVH'));
const USDC = new PublicKey("4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU");

async function bridge(chainID: number) {
    const computeUnitInstruction = ComputeBudgetProgram.setComputeUnitLimit({
        units: 1_000_000_000
    });

    const zeroArray = Array.from({length: 32}, (_) => 0);
    const notZeroArray = Array.from({length: 32}, (_, index) => index * 2);

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
        srcEid: SOLANA_EID,
        dstEid: chainID,
        sender: utils.hexlify(gringottsPDA.toBytes()),
        receiver: utils.hexlify(RECEIVER),
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
                inTransfer: {
                    amountUsdx: new BN(100 * 1000 * 1000),
                    items: [
                        {
                            asset: Array.from(USDC.toBytes()),
                            amount: new BN(1000 * 1000),
                            swap: null,
                        }
                    ]
                },
                outTransfers: [
                    {
                        chainId: 1,
                        items: [
                            {
                                asset: Array.from(utils.arrayify(utils.hexZeroPad(TEST_STABLE, 32))),
                                recipient: Array.from(utils.arrayify(utils.hexZeroPad("0x0D595AE2666a2c5Ae6b99cce4DD428a9Cf20B2c9", 32))),
                                executionGasAmount: new BN(100000),
                                distributionBp: 10_000,
                                swap: null,
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
            await provider.connection.getAddressLookupTable(new PublicKey("BLakf36faA77C5hG3WkzUJzBnPcACLreUNqPYsbhyULh"))
        ).value;

        const messageV0 = new TransactionMessage({
            payerKey: wallet.publicKey,
            recentBlockhash: (await provider.connection.getLatestBlockhash()).blockhash,
            instructions: [computeUnitInstruction, tx],
        }).compileToV0Message([lookupTable]);
        const transaction = new VersionedTransaction(messageV0);

        const txID = await provider.sendAndConfirm(transaction, [wallet.payer]);
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
        srcEid: SOLANA_EID,
        dstEid: ARB_EID,
        sender: utils.hexlify(gringottsPDA.toBytes()),
        receiver: utils.hexlify(RECEIVER),
    }
    const accountsEx = await endpointProgram.getSendIXAccountMetaForCPI(provider.connection, vaultPDA, packetPath, ulnProgram);

    const extendInstruction = AddressLookupTableProgram.extendLookupTable({
        payer: wallet.payer.publicKey,
        authority: wallet.payer.publicKey,
        lookupTable: new PublicKey("BLakf36faA77C5hG3WkzUJzBnPcACLreUNqPYsbhyULh"),
        addresses: accountsEx.map((a) => a.pubkey),
    });

    const t2 = new Transaction().add(extendInstruction);
    const s2 = await sendAndConfirmTransaction(provider.connection, t2, [wallet.payer], {
        skipPreflight: true,
    });

    console.log(s2);
}

async function main() {
    // await createAddressLookup();
    await bridge(ARB_EID);
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
    });