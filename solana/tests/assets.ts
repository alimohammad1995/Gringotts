import * as anchor from "@coral-xyz/anchor";
import {Program} from "@coral-xyz/anchor";
import {Gringotts} from "../target/types/gringotts";
import {SystemProgram, PublicKey} from '@solana/web3.js';
import {expect} from "chai";

describe("Asset", () => {
    const provider = anchor.AnchorProvider.env();
    anchor.setProvider(anchor.AnchorProvider.env());

    const program = anchor.workspace.Gringotts as Program<Gringotts>;
    const PDA_SEED = "Gringotts";
    const [gringottsPDA, bump] = PublicKey.findProgramAddressSync(
        [Buffer.from(PDA_SEED)],
        program.programId
    );

    beforeEach(async () => {
        try {
            await program.methods
                .destroy()
                .accounts({
                    gringotts: gringottsPDA,
                    owner: provider.wallet.publicKey,
                    systemProgram: SystemProgram.programId,
                })
                .rpc();
        } catch (err) {
        }

        try {
            await program.methods
                .initialize()
                .accounts({
                    gringotts: gringottsPDA,
                    owner: provider.wallet.publicKey,
                    systemProgram: SystemProgram.programId,
                })
                .rpc();
        } catch (e) {
        }
    });

    it("Add asset", async () => {
        const assetId = 1;

        const [assetPda, assetBump] = anchor.web3.PublicKey.findProgramAddressSync(
            [Buffer.from("asset"), Buffer.from(new Uint16Array([assetId]).buffer)],
            program.programId
        );
        const priceFeed = anchor.web3.Keypair.generate().publicKey;
        const chains = Buffer.from([1, 2, 3]);

        try {
            await program.methods
                .addAsset(assetId, priceFeed, chains)
                .accounts({
                    asset: assetPda,
                    gringotts: gringottsPDA,
                    owner: provider.wallet.publicKey,
                    systemProgram: SystemProgram.programId,
                })
                .rpc();
        } catch (err) {
            expect.fail("Add asset transaction failed")
        }

        // Fetch the updated asset account to verify the changes
        const assetAccount = await program.account.asset.fetch(assetPda);

        expect(assetAccount.id).to.equal(assetId);
        expect(assetAccount.priceFeed.toBase58()).to.equal(priceFeed.toBase58());
        expect(assetAccount.chains).to.deep.equal(chains);
    });

    it("Update asset", async () => {
        const assetId = 2;
        const [assetPda, assetBump] = anchor.web3.PublicKey.findProgramAddressSync(
            [Buffer.from("asset"), Buffer.from(new Uint16Array([assetId]).buffer)],
            program.programId
        );
        const oldPriceFeed = anchor.web3.Keypair.generate().publicKey;
        const oldChains = Buffer.from([1, 2, 3]);

        try {
            await program.methods
                .addAsset(assetId, oldPriceFeed, oldChains)
                .accounts({
                    asset: assetPda,
                    gringotts: gringottsPDA,
                    owner: provider.wallet.publicKey,
                    systemProgram: SystemProgram.programId,
                })
                .rpc();
        } catch (err) {
            expect.fail("Add asset transaction failed")
        }

        const priceFeed = anchor.web3.Keypair.generate().publicKey;
        const chains = Buffer.from([1, 2, 3, 4]);

        try {
            await program.methods
                .updateAsset(priceFeed, chains)
                .accounts({
                    asset: assetPda,
                    gringotts: gringottsPDA,
                    owner: provider.wallet.publicKey,
                    systemProgram: SystemProgram.programId,
                })
                .rpc();
        } catch (err) {
            expect.fail("Add asset transaction failed")
        }

        // Fetch the updated asset account to verify the changes
        const assetAccount = await program.account.asset.fetch(assetPda);

        expect(assetAccount.id).to.equal(assetId);
        expect(assetAccount.priceFeed.toBase58()).to.equal(priceFeed.toBase58());
        expect(assetAccount.chains).to.deep.equal(chains);
    });

    it("Delete asset", async () => {
        const assetId = 3;
        const [assetPda, assetBump] = anchor.web3.PublicKey.findProgramAddressSync(
            [Buffer.from("asset"), Buffer.from(new Uint16Array([assetId]).buffer)],
            program.programId
        );

        try {
            await program.methods
                .addAsset(assetId, anchor.web3.Keypair.generate().publicKey, Buffer.from([1, 2, 3]))
                .accounts({
                    asset: assetPda,
                    gringotts: gringottsPDA,
                    owner: provider.wallet.publicKey,
                    systemProgram: SystemProgram.programId,
                })
                .rpc();
        } catch (err) {
            expect.fail("Add asset transaction failed")
        }

        const initialAssetBalance = await provider.connection.getBalance(assetPda);
        const initialGringottsBalance = await provider.connection.getBalance(gringottsPDA);

        try {
            await program.methods
                .removeAsset()
                .accounts({
                    asset: assetPda,
                    gringotts: gringottsPDA,
                    owner: provider.wallet.publicKey,
                    programAccount: program.programId,
                    systemProgram: SystemProgram.programId,
                })
                .rpc();
        } catch (err) {
            console.log(err);
            expect.fail("failed")
        }

        try {
            await program.account.asset.fetch(assetPda);
            expect.fail('The Gringotts account should have been closed');
        } catch (err) {
        }

        const finalVaultBalance = await provider.connection.getBalance(gringottsPDA);
        expect(finalVaultBalance).to.equal(initialGringottsBalance + initialAssetBalance);
    });
});
