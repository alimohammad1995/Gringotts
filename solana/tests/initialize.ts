import * as anchor from "@coral-xyz/anchor";
import {Program} from "@coral-xyz/anchor";
import {SystemProgram, PublicKey} from '@solana/web3.js';
import {Gringotts} from '../target/types/gringotts';
import {expect} from "chai";

describe("Gringotts", () => {
    const provider = anchor.AnchorProvider.env();
    anchor.setProvider(anchor.AnchorProvider.env());

    const program = anchor.workspace.Gringotts as Program<Gringotts>;
    const PDA_SEED = "Gringotts";
    const [gringottsPDA,] = PublicKey.findProgramAddressSync(
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
    });

    it("Init", async () => {
        try {
            await program.methods
                .initialize()
                .accounts({
                    gringotts: gringottsPDA,
                    owner: provider.wallet.publicKey,
                    systemProgram: SystemProgram.programId,
                })
                .rpc();
        } catch (err) {
            expect.fail('Init transaction failed:', err)
        }

        try {
            const gringottsAccount = await program.account.gringotts.fetch(gringottsPDA);
            const owner = gringottsAccount.owner.toBase58();
            expect(owner).to.equal(provider.wallet.publicKey.toBase58());
        } catch (err) {
            expect.fail('Failed to fetch Gringotts account:', err)
        }
    });

    it('Destroys Gringotts', async () => {
        try {
            await program.methods
                .initialize()
                .accounts({
                    gringotts: gringottsPDA,
                    owner: provider.wallet.publicKey,
                    systemProgram: SystemProgram.programId,
                })
                .rpc();
        } catch (err) {
            expect.fail('Init transaction failed:', err)
        }

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
            expect.fail('Close transaction failed:', err)
        }

        try {
            await program.account.gringotts.fetch(gringottsPDA);
            expect.fail('The Gringotts account should have been closed');
        } catch (err) {
        }
    });
});