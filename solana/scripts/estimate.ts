import * as anchor from "@coral-xyz/anchor";
import {BN, Program} from "@coral-xyz/anchor";
import {Gringotts} from "../target/types/gringotts";
import {ComputeBudgetProgram, Connection, PublicKey} from '@solana/web3.js';
import {utils} from "ethers";
import {EndpointProgram, UlnProgram} from '@layerzerolabs/lz-solana-sdk-v2'
import {PacketPath} from '@layerzerolabs/lz-v2-utilities'
import NodeWallet from "@coral-xyz/anchor/dist/cjs/nodewallet";

const provider = anchor.AnchorProvider.env();
anchor.setProvider(anchor.AnchorProvider.env());

const program = anchor.workspace.Gringotts as Program<Gringotts>;
const wallet = provider.wallet as NodeWallet;

export const GRINGOTTS_SEED = 'Gringotts'
export const PEER_SEED = 'Peer'
const ARB_EID = 40231;
const SOLANA_EID = 40168;
const RECEIVER = "0x1c191f62728b1498d779559e9ffb75a849582103";

const [gringottsPDA] = PublicKey.findProgramAddressSync([Buffer.from(GRINGOTTS_SEED)], program.programId);
const [peerPDA] = PublicKey.findProgramAddressSync([Buffer.from(PEER_SEED), new BN(ARB_EID).toArrayLike(Buffer, 'le', 4)], program.programId);

const endpointProgram = new EndpointProgram.Endpoint(new PublicKey('76y77prsiCMvXMjuoZ5VRrhG5qYBrUMYTE5WgHqgjEn6'));
const ulnProgram = new UlnProgram.Uln(new PublicKey('7a4WjyR8VZ7yZz5XJAKm39BUGn5iT9CKcv2pmG9tdXVH'));

async function estimate(chainID: number) {
    const connection = new Connection(provider.connection.rpcEndpoint, 'confirmed')

    const packetPath: PacketPath = {
        srcEid: SOLANA_EID,
        dstEid: chainID,
        sender: utils.hexlify(gringottsPDA.toBytes()),
        receiver: utils.hexlify(RECEIVER),
    }

    const accounts = [
        {pubkey: peerPDA, isSigner: false, isWritable: true},
    ];

    let accountsEx = await endpointProgram.getQuoteIXAccountMetaForCPI(connection, wallet.publicKey, packetPath, ulnProgram);
    const zeroArray = Array.from({length: 32}, (_, index) => index);

    // let x = [];
    // for (let i = 0; i < accountsEx.length; i++) {
    //     x.push({'address': accountsEx[i].pubkey.toBase58(), 'is_writable': accountsEx[i].isWritable, 'is_signer': accountsEx[i].isSigner});
    // }
    // console.log(JSON.stringify(x));

    const computeUnitInstruction = ComputeBudgetProgram.setComputeUnitLimit({
        units: 1_000_000_000
    });

    try {
        const tx = await program.methods
            .estimate({
                inbound: {
                    amountUsdx: new BN(1000000000),
                },
                outbounds: [
                    {
                        chainId: 1,
                        items: [
                            {
                                asset: zeroArray,
                                executionGasAmount: new BN(100000),
                                executionCommandLength: 1000,
                                executionMetadataLength: 0,
                            }
                        ],
                    }
                ],
            }).accounts({
                priceFeed: new PublicKey("7UVimffxr9ow1uXYxsr4LHAcV58mLzhmwaeKvJ1pjLiE"),
                gringotts: gringottsPDA,
            })
            .remainingAccounts(accounts.concat(accountsEx))
            .preInstructions([computeUnitInstruction])
            .simulate();

        console.log(peerPDA.toBase58());

        const returnPrefix = `Program return: ${program.programId} `
        const returnLog = tx.raw?.find((l) => l.startsWith(returnPrefix))
        const buffer = Buffer.from(returnLog.slice(returnPrefix.length), 'base64')
        console.log(returnLog);

        const estimates = program.coder.types.decode('estimateResponse', buffer);
        console.log(estimates);
    } catch (err) {
        console.log(err);
        console.log(await err.getLogs());
    }
}

async function main() {
    const chainID = 40231;
    await estimate(chainID);
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
    });