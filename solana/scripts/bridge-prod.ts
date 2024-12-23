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
    console.log("Gringotts USDC", vaultPDA.toBase58());

    // const tx1 = await program.methods.tokenWithdraw({
    //     amount: new BN(7 * 1000 * 1000)
    // }).accounts({
    //     gringotts: gringottsPDA,
    //     recipient: wallet.publicKey,
    //     tokenMint: USDC,
    //     gringottsTokenAccount: gringottsUSDC,
    //     recipientTokenAccount: userUSDC,
    //     associatedTokenProgram: ASSOCIATED_TOKEN_PROGRAM_ID,
    //     tokenProgram: TOKEN_PROGRAM_ID,
    //     systemProgram: SystemProgram.programId
    // }).rpc()
    //
    // console.log(tx1);

    // const tx2 = await program.methods.vaultWithdraw({
    //     amount: new BN(10 * 1000 * 1000)
    // }).accounts({
    //     gringotts: gringottsPDA,
    //     recipient: wallet.publicKey,
    //     vault: vaultPDA,
    //     systemProgram: SystemProgram.programId
    // }).rpc()
    //
    // console.log(tx2);
}

async function lzReceive() {
    console.log(gringottsPDA.toBase58());
    console.log("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL");
    console.log("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA");
    console.log("11111111111111111111111111111111");
    console.log("H6ieTjyWqcRFv1RwD9LghEFDxVtMaEMgVbgnYrdHMjr5");
    console.log(NATIVE_MINT.toBase58());
    console.log(getAssociatedTokenAddressSync(
        NATIVE_MINT,
        gringottsPDA,
        true
    ).toBase58())
    console.log("JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4");
    console.log("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v");
    console.log(getAssociatedTokenAddressSync(
        USDC,
        gringottsPDA,
        true
    ).toBase58())
    console.log(getAssociatedTokenAddressSync(
        USDC,
        wallet.publicKey,
        true
    ).toBase58())

    const tx = await program.methods.lzReceiveTypes({
        srcEid: ARB_EID,
        sender: Array.from(utils.arrayify(utils.hexZeroPad(ARB_ADDRESS, 32))),
        nonce: new BN(2),
        guid: Array.from(utils.arrayify(utils.hexZeroPad("0x0", 32))),
        // message: Buffer.from("010100000000002bfb2f0000000000000000000000000000000000000000000000000000000000000000ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc660479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138fc6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d610023e517cb977ae3ad2a010000001964000180841e00000000006d6b050000000000640000033a1906ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a900ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc6600429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd650128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce03010479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d61000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00b43ffa27f5d7f64a74c09b1f295879de4b09ab36dfc9dd514b321aa7b38ce5e8000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f000e03685f8e909053e458121c66f5a76aedc7706aa11c82f8aa952a8f2b7879a90006ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a90006ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a900054a535a992921064d24e87160da387c7c35b5ddbc92bb81e41fa8404105448d00ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc6600538745d236517bba822fb87c8a0ffd7a981125acd6bf793f2f07529551a37c2201069b8857feab8184fb687f634618c035dac439dc1aeb3b5598a0f0000000000100c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d6100429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd65014eacd57702d35bcb7e52c58caf01c4bdcfebf35e0e99547dc8b55066e8db5afa0128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce0301c528bd00d0774c59320cb908f8c41b8fddcb0bc86cd53bf317b4fbbac28355880102333e1d4d67788300eae9a294feca87813ce579503b6d308ecb9f6448e3bed00143e4192edb5ae4ffdb2882b9a003b63fefe9df8ac324e3e3a39033632322aea10127aac7", 'hex'),
        // message: Buffer.from("010100000000002bfa410000000000000000000000000000000000000000000000000000000000000000ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc660479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138fc6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d610023e517cb977ae3ad2a010000001964000180841e00000000006d6b05000000000064000002531206ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a900ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc6600429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd650128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce03010479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d61000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00b43ffa27f5d7f64a74c09b1f295879de4b09ab36dfc9dd514b321aa7b38ce5e8000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00069be86ec9af65eb4a614fd99b8e92547da0145fab5e804adb594db3e73a271b00ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc6600a5d8744c3d0ad845e0d44689b8ea407cc2f491d1bb7931f55cd40810e557e1aa01aa3f371e63ade99128ff53f5fc6b229ff8ec58274bf9e98bb39b9ada262fe7f901046d84465266a1d3008798a5b499bac12326c4b15bfb47fff0a549eb42fe2b6801429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd650128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce030106ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a90006a7d517187bd16635dad40455fdc2c0c124c68f215675a5dbbacb5f0800000000", 'hex'),
        message: Buffer.from("010100000000002bfa410000000000000000000000000000000000000000000000000000000000000000ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc660479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138fc6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d610023e517cb977ae3ad2a010000001964000180841e00000000006d6b05000000000064000002531206ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a900ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc6600429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd650128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce03010479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d61000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00b43ffa27f5d7f64a74c09b1f295879de4b09ab36dfc9dd514b321aa7b38ce5e8000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00069be86ec9af65eb4a614fd99b8e92547da0145fab5e804adb594db3e73a271b00ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc6600a5d8744c3d0ad845e0d44689b8ea407cc2f491d1bb7931f55cd40810e557e1aa01aa3f371e63ade99128ff53f5fc6b229ff8ec58274bf9e98bb39b9ada262fe7f901046d84465266a1d3008798a5b499bac12326c4b15bfb47fff0a549eb42fe2b6801429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd650128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce030106ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a90006a7d517187bd16635dad40455fdc2c0c124c68f215675a5dbbacb5f0800000000", 'hex'),
        extraData: Buffer.from([]),
    }).accounts({
        gringotts: gringottsPDA,
    }).simulate()

    console.log(tx);
}

async function main() {
    // await createAddressLookup();
    // await bridge(ARB_EID);
    await widthdraw_token();
    // await lzReceive();
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
    });