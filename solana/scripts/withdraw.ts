import * as anchor from "@coral-xyz/anchor";
import {BN, Program} from "@coral-xyz/anchor";
import {Gringotts} from "../target/types/gringotts";
import {PublicKey, SystemProgram} from '@solana/web3.js';
import NodeWallet from "@coral-xyz/anchor/dist/cjs/nodewallet";
import {ASSOCIATED_TOKEN_PROGRAM_ID, getAssociatedTokenAddressSync, TOKEN_PROGRAM_ID} from "@solana/spl-token";

const SOLANA_EID = 40168;

const provider = anchor.AnchorProvider.env();
anchor.setProvider(anchor.AnchorProvider.env());

const program = anchor.workspace.Gringotts as Program<Gringotts>;
const wallet = provider.wallet as NodeWallet;

export const GRINGOTTS_SEED = 'Gringotts'
export const PEER_SEED = 'Peer'
export const VAULT_SEED = 'Vault'

const [gringottsPDA] = PublicKey.findProgramAddressSync([Buffer.from(GRINGOTTS_SEED)], program.programId);
const [vaultPDA] = PublicKey.findProgramAddressSync([Buffer.from(VAULT_SEED)], program.programId);

const USDC = new PublicKey("4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU");


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


    const tx = await program.methods.tokenWithdraw({
        amount: new BN(1 * 1000 * 1000)
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
    await widthdraw_token()
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
    });