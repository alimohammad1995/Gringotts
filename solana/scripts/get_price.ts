import * as anchor from "@coral-xyz/anchor";
import {BN, Program} from "@coral-xyz/anchor";
import {Gringotts} from "../target/types/gringotts";
import {SystemProgram, PublicKey} from '@solana/web3.js';


const provider = anchor.AnchorProvider.env();
anchor.setProvider(anchor.AnchorProvider.env());

const program = anchor.workspace.Gringotts as Program<Gringotts>;
const wallet = provider.wallet

async function main() {
    console.log("Main wallet", wallet.publicKey.toBase58())
    console.log("Program ID", program.programId.toBase58())

    const [gringottsPDA,] = PublicKey.findProgramAddressSync(
        [Buffer.from('Gringotts')],
        program.programId
    );

    console.log("Gringotts account initialized at:", gringottsPDA.toBase58());


    const assetId = 1;
    const [assetPDA,] = PublicKey.findProgramAddressSync(
        [Buffer.from('Asset'), Buffer.from(new Uint16Array([assetId]).buffer)],
        program.programId
    );

    // try {
    //     await program.methods
    //         .addAsset(
    //             assetId,
    //             new PublicKey("669U43LNHx7LsVj95uYksnhXUfWKDsdzVqev3V4Jpw3P"),
    //             Buffer.from([1, 2, 3]),
    //         ).accounts({
    //             asset: assetPDA,
    //             gringotts: gringottsPDA,
    //             owner: wallet.publicKey,
    //             systemProgram: SystemProgram.programId,
    //         })
    //         .rpc();
    // } catch (err) {
    //     // console.log(err);
    // }
    //
    // try {
    //     await program.methods
    //         .updateAsset(
    //             new PublicKey("669U43LNHx7LsVj95uYksnhXUfWKDsdzVqev3V4Jpw3P"),
    //             Buffer.from([1, 2, 3]),
    //         ).accounts({
    //             asset: assetPDA,
    //             gringotts: gringottsPDA,
    //             owner: wallet.publicKey,
    //         })
    //         .rpc();
    // } catch (err) {
    //     console.log(err);
    // }


    try {
        const x = await program.methods
            .getAssetPrice()
            .accounts({
                asset: assetPDA,
                priceUpdate: "HAm5DZhrgrWa12heKSxocQRyJWGCtXegC77hFQ8F5QTH"
            })
            .rpc();

        console.log(x);
    } catch (err) {
        console.log(err);
    }

    console.log('Asset account initialized at:', assetPDA.toBase58());
}

main()
    .then(() => process.exit(0))
    .catch((err) => {
        console.error(err);
    });