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
                inbound: {
                    amountUsdx: new BN(100 * 1000 * 1000),
                    items: [
                        {
                            asset: Array.from(USDC.toBytes()),
                            amount: new BN(1000 * 1000),
                            swap: null,
                        }
                    ]
                },
                outbounds: [
                    {
                        chainId: 1,
                        items: [
                            {
                                asset: Array.from(utils.arrayify(utils.hexZeroPad(TEST_STABLE, 32))),
                                recipient: Array.from(utils.arrayify(utils.hexZeroPad("0x0D595AE2666a2c5Ae6b99cce4DD428a9Cf20B2c9", 32))),
                                executionGasAmount: new BN(100000),
                                distributionBp: 10_000,
                                swap: {
                                    executor: Array.from(utils.arrayify(utils.hexZeroPad("0x6352a56caadC4F1E25CD6c75970Fa768A3304e64", 32))),
                                    command: Buffer.from("bc80f1a80000000000000000000000000d595ae2666a2c5ae6b99cce4dd428a9cf20b2c900000000000000000000000000000000000000000000000000000000002dc6c0000000000000000000000000000000000000000000000000000326ad37a4dc0900000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000002000064000000000000000000be3ad6a5669dc0b8b12febc03608860c31e2eef6a0000000000000000000000042161084d0672e1d3f26a9b53e653be2084ff19c", 'hex'),
                                    metadata: Buffer.from([]),
                                    stableToken: Array.from(utils.arrayify(utils.hexZeroPad(TEST_STABLE, 32))),
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
    const address = accountsEx.map((a) => a.pubkey);
    address.push(new PublicKey("HjVcDEcpjVvzwEnaHkZgdrV2pV2DoFcg7TyemtMovuZg"));
    address.push(new PublicKey("EVPzMgjbbejQEv8AW1kr4CJidApoqMh1RovB3r2ZVi5"));

    const extendInstruction = AddressLookupTableProgram.extendLookupTable({
        payer: wallet.payer.publicKey,
        authority: wallet.payer.publicKey,
        lookupTable: new PublicKey("BLakf36faA77C5hG3WkzUJzBnPcACLreUNqPYsbhyULh"),
        addresses: address,
    });

    const t2 = new Transaction().add(extendInstruction);
    const s2 = await sendAndConfirmTransaction(provider.connection, t2, [wallet.payer], {
        skipPreflight: true,
    });

    console.log(s2);
}

async function lzReceive() {
    const tx = await program.methods.lzReceiveTypes({
        srcEid: ARB_EID,
        sender: Array.from(utils.arrayify(utils.hexZeroPad(ARB_ADDRESS, 32))),
        nonce: new BN(2),
        guid: Array.from(utils.arrayify(utils.hexZeroPad("0x0", 32))),
        message: Buffer.from("010100000000002cad430428d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce03ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc66c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d617d7cdbe036610dfd0b509d8763c05f438474dc88711dd48e5071145043e0a01d90040607080500", 'hex'),
        extraData: Buffer.from([]),
    }).accounts({
        gringotts: gringottsPDA,
    }).simulate()

    console.log(tx);
}

async function main() {
    // await createAddressLookup();
    // await bridge(ARB_EID);
    // await lzReceive();
    console.log(gringottsPDA.toBase58());
    console.log(vaultPDA.toBase58());
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
    });