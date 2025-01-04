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

const ARB_ADDRESS = "0x1c191f62728b1498d779559e9ffb75a849582103";
const TEST_STABLE = "0x301b022b40d06088fc974e767149f4a3feebbf1a";

const provider = anchor.AnchorProvider.env();
anchor.setProvider(anchor.AnchorProvider.env());

const ARB_EID = 40231;
const SOL_EID = 40168;

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
    [Buffer.from(GRINGOTTS_SEED)],
    program.programId
);

const [vaultPDA] = PublicKey.findProgramAddressSync(
    [Buffer.from(VAULT_SEED)],
    program.programId
);

async function init() {
    const [lzReceiveTypesPDA] = PublicKey.findProgramAddressSync(
        [Buffer.from(LZ_RECEIVE_TYPES_SEED), gringottsPDA.toBytes()], program.programId
    );

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
            chainId: 2,
            lzEid: SOL_EID,
            vaultFund: new BN(10 * 1000 * 1000),
            lzEndpointProgram: endpointProgram.program,
            pythPriceFeedId: '0xef0d8b6fda2ceba41da15d4095d1da392a0d2f8ed0c6c7bc0f4cfac8c280b56d',
            commissionMicroBps: 5000,
        })
        .accounts({
            gringotts: gringottsPDA,
            vault: vaultPDA,
            lzReceiveTypesPDA: lzReceiveTypesPDA,
            owner: wallet.publicKey,
            systemProgram: SystemProgram.programId,
        })
        .remainingAccounts(registerOAppAccounts)
        .rpc();

    console.log(tx);

    console.log(
        "Gringotts PDA", gringottsPDA.toBase58(),
        "Vault PDA", vaultPDA.toBase58(),
        "\nValue", JSON.stringify(await program.account.gringotts.fetch(gringottsPDA))
    );
}

async function addPeer(
    chain_id: number,
    remote: number,
    remotePeer: Uint8Array,
) {
    const [chainPDA] = PublicKey.findProgramAddressSync(
        [Buffer.from(PEER_SEED), new BN(remote).toArrayLike(Buffer, 'be', 4)],
        program.programId
    );

    const zeroArray = Array.from({length: 32}, (_) => 0);
    const oneArray = Array.from({length: 32}, (_) => 1);

    const tx = await program.methods
        .peerAdd({
            chainId: chain_id,
            lzEid: remote,
            stableCoins: [
                zeroArray,
                oneArray,
                Array.from(utils.arrayify(new PublicKey("4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU").toBytes())),
                Array.from(utils.arrayify(utils.hexZeroPad(TEST_STABLE, 32))),
            ],
            baseGasEstimate: new BN(50000),
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


async function main() {
    const chainID = ARB_EID;
    const peer = utils.arrayify(utils.hexZeroPad(ARB_ADDRESS, 32));

    await init();
    await addPeer(1, chainID, peer); // ARB
    await addPeer(2, SOL_EID, utils.arrayify(program.programId.toBytes())); // SOL
    await handleSend(chainID);
    await handleReceive(chainID);
    await handleNonce(chainID, peer);
    await initUlnConfig(chainID);
    await setOappExecutor(chainID);
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