import * as anchor from "@coral-xyz/anchor";
import {BN, Program} from "@coral-xyz/anchor";
import {Gringotts} from "../target/types/gringotts";
import {Connection, PublicKey, Signer, SystemProgram, TransactionInstruction} from '@solana/web3.js';
import {utils} from "ethers";
import {
    buildVersionedTransaction,
    EndpointProgram,
    EventPDADeriver,
    ExecutorPDADeriver,
    SetConfigType,
    UlnProgram
} from '@layerzerolabs/lz-solana-sdk-v2'
import NodeWallet from "@coral-xyz/anchor/dist/cjs/nodewallet";
import {ASSOCIATED_TOKEN_PROGRAM_ID, getAssociatedTokenAddressSync, TOKEN_PROGRAM_ID} from "@solana/spl-token";

const provider = anchor.AnchorProvider.env();
anchor.setProvider(anchor.AnchorProvider.env());

const ARB_EID = 30110;
const SOL_EID = 30168;
const ARB_CHAIN_ID = 1;
const SOL_CHAIN_ID = 2;

const ARB_ADDRESS = "0x02fa4ad1fe96e10f81a258bbe29e0392124bf27c";
const SOL_USDC = new PublicKey("EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v");

const program = anchor.workspace.Gringotts as Program<Gringotts>;
const wallet = provider.wallet as NodeWallet;

export const VAULT_SEED = 'Vault'
export const GRINGOTTS_SEED = 'Gringotts'
export const PEER_SEED = 'Peer'
export const LZ_RECEIVE_TYPES_SEED = 'LzReceiveTypes'

const endpointProgram = new EndpointProgram.Endpoint(new PublicKey('76y77prsiCMvXMjuoZ5VRrhG5qYBrUMYTE5WgHqgjEn6'));
const ulnProgram = new UlnProgram.Uln(new PublicKey('7a4WjyR8VZ7yZz5XJAKm39BUGn5iT9CKcv2pmG9tdXVH'));
const executorProgram = new PublicKey('6doghB248px58JSSwG4qejQ46kFMW4AMj7vzJnWZHNZn');

const [gringottsPDA] = PublicKey.findProgramAddressSync(
    [Buffer.from(GRINGOTTS_SEED)], program.programId
);

const [vaultPDA] = PublicKey.findProgramAddressSync(
    [Buffer.from(VAULT_SEED)], program.programId
);

const [lzReceiveTypesPDA] = PublicKey.findProgramAddressSync(
    [Buffer.from(LZ_RECEIVE_TYPES_SEED), gringottsPDA.toBytes()], program.programId
);

async function init() {
    const [oAppRegistry] = endpointProgram.deriver.oappRegistry(gringottsPDA);
    const [eventAuthority] = new EventPDADeriver(endpointProgram.program).eventAuthority()

    const ixAccounts = EndpointProgram.instructions.createRegisterOappInstructionAccounts(
        {
            payer: wallet.publicKey,
            oapp: gringottsPDA,
            oappRegistry: oAppRegistry,
            eventAuthority: eventAuthority,
            program: endpointProgram.program,
        },
        endpointProgram.program
    )
    const registerOAppAccounts = [
        {pubkey: endpointProgram.program, isSigner: false, isWritable: false},
        ...ixAccounts,
    ]
    registerOAppAccounts[1].isSigner = false
    registerOAppAccounts[2].isSigner = false

    const tx = await program.methods
        .gringottsInitialize({
            chainId: SOL_CHAIN_ID,
            lzEid: SOL_EID,
            stableCoins: [Array.from(utils.arrayify(SOL_USDC.toBytes()))],
            vaultFund: new BN(10 * 1000 * 1000),
            lzEndpointProgram: endpointProgram.program,
            pythPriceFeedId: '0xef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d',
            commissionMicroBps: 5000,
        })
        .accounts({
            gringotts: gringottsPDA,
            vault: vaultPDA,
            lzReceiveTypesPDA: lzReceiveTypesPDA,
            systemProgram: SystemProgram.programId,
        })
        .remainingAccounts(registerOAppAccounts)
        .rpc();

    console.log(tx);

    console.log(
        "Gringotts PDA", gringottsPDA.toBase58(),
        "Vault PDA", vaultPDA.toBase58(),
        "Value", JSON.stringify(await program.account.gringotts.fetch(gringottsPDA))
    );
}

async function addPeer(
    remotePeer: Uint8Array,
) {
    const [chainPDA] = PublicKey.findProgramAddressSync(
        [Buffer.from(PEER_SEED), new BN(ARB_EID).toArrayLike(Buffer, 'be', 4)],
        program.programId
    );

    const tx = await program.methods
        .peerAdd({
            chainId: ARB_CHAIN_ID,
            lzEid: ARB_EID,
            multiSend: false,
            baseGasEstimate: new BN(30_000),
            address: Array.from(remotePeer),
        })
        .accounts({
            peer: chainPDA,
            gringotts: gringottsPDA,
            systemProgram: SystemProgram.programId,
        })
        .rpc();

    console.log("tx", tx);
}

async function handleSend(remote: number) {
    const ix1 = endpointProgram.initSendLibrary(wallet.publicKey, gringottsPDA, remote);
    const tx1 = await sendAndConfirm([wallet.payer], [ix1]);
    console.log("tx1", tx1);

    const ix2 = endpointProgram.setSendLibrary(wallet.publicKey, gringottsPDA, ulnProgram.program, remote)
    const tx2 = await sendAndConfirm([wallet.payer], [ix2]);
    console.log("tx2", tx2);
}

async function handleReceive(remote: number) {
    const ix1 = endpointProgram.initReceiveLibrary(wallet.publicKey, gringottsPDA, remote)
    const tx1 = await sendAndConfirm([wallet.payer], [ix1]);
    console.log("tx1", tx1);

    const ix2 = endpointProgram.setReceiveLibrary(wallet.publicKey, gringottsPDA, ulnProgram.program, remote)
    const tx2 = await sendAndConfirm([wallet.payer], [ix2]);
    console.log("tx2", tx2);
}

async function handleNonce(
    remote: number,
    remotePeer: Uint8Array
) {
    const ix = endpointProgram.initOAppNonce(wallet.publicKey, remote, gringottsPDA, remotePeer)
    const tx = await sendAndConfirm([wallet.payer], [ix])
    console.log("tx", tx);
}

async function initUlnConfig(
    remote: number
) {
    const ix = endpointProgram.initOAppConfig(wallet.publicKey, ulnProgram, wallet.publicKey, gringottsPDA, remote)
    const tx = await sendAndConfirm([wallet.payer], [ix]);
    console.log("tx", tx);
}

async function setOappExecutor(remote: number) {
    const defaultOutboundMaxMessageSize = 10000

    const [executorPda] = new ExecutorPDADeriver(executorProgram).config()
    const expected: UlnProgram.types.ExecutorConfig = {
        maxMessageSize: defaultOutboundMaxMessageSize,
        executor: executorPda,
    }

    const connection = new Connection(provider.connection.rpcEndpoint, 'confirmed');
    const ix = await endpointProgram.setOappConfig(connection, wallet.publicKey, gringottsPDA, ulnProgram.program, remote, {
        configType: SetConfigType.EXECUTOR,
        value: expected,
    })

    const tx = await sendAndConfirm([wallet.payer], [ix])
    console.log("tx", tx);
}

async function fund_usdc_token() {
    const gringottsUSDC = getAssociatedTokenAddressSync(
        SOL_USDC,
        gringottsPDA,
        true
    );

    const userUSDC = getAssociatedTokenAddressSync(
        SOL_USDC,
        wallet.publicKey,
    );

    console.log("User USDC", userUSDC.toBase58());
    console.log("Gringotts USDC", gringottsUSDC.toBase58());

    const tx1 = await program.methods.tokenFund({
        amount: new BN(6 * 1000 * 1000)
    }).accounts({
        gringotts: gringottsPDA,
        tokenMint: SOL_USDC,
        gringottsTokenAccount: gringottsUSDC,
        ownerTokenAccount: userUSDC,
        associatedTokenProgram: ASSOCIATED_TOKEN_PROGRAM_ID,
        tokenProgram: TOKEN_PROGRAM_ID,
        systemProgram: SystemProgram.programId
    }).rpc()

    console.log(tx1);
}

async function main() {
    const peer = utils.arrayify(utils.hexZeroPad(ARB_ADDRESS, 32));

    console.log("Gringotts PDA", gringottsPDA.toBase58());
    console.log("Vault PDA", vaultPDA.toBase58());

    // await init();
    // await addPeer(peer); // ARB
    // await handleSend(ARB_EID);
    // await handleReceive(ARB_EID);
    // await handleNonce(ARB_EID, peer);
    // await initUlnConfig(ARB_EID);
    // await setOappExecutor(ARB_EID);
    // await fund_usdc_token();
}

async function sendAndConfirm(
    signers: Signer[],
    instructions: TransactionInstruction[]
) {
    const connection = new Connection(provider.connection.rpcEndpoint, 'confirmed');
    const tx = await buildVersionedTransaction(connection, signers[0].publicKey, instructions, 'confirmed')
    tx.sign(signers)
    const hash = await connection.sendTransaction(tx, {skipPreflight: true})
    await connection.confirmTransaction(hash, 'confirmed')
    return hash
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
    });