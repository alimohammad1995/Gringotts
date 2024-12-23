import * as anchor from "@coral-xyz/anchor";
import {BN, Program} from "@coral-xyz/anchor";
import {SystemProgram, PublicKey} from '@solana/web3.js';
import {Gringotts} from '../target/types/gringotts';
import {expect} from "chai";
import {utils} from "ethers";

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
    });

    it("Init", async () => {
        try {
            const tx = await program.methods.lzReceiveTypes({
                srcEid: 1,
                sender: Array.from(utils.arrayify(utils.hexZeroPad("0x0", 32))),
                nonce: new BN(2),
                guid: Array.from(utils.arrayify(utils.hexZeroPad("0x0", 32))),
                message: Buffer.from("010100000000002bfa410000000000000000000000000000000000000000000000000000000000000000ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc660479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138fc6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d610023e517cb977ae3ad2a010000001964000180841e00000000006d6b05000000000064000002531206ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a900ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc6600429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd650128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce03010479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d61000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00b43ffa27f5d7f64a74c09b1f295879de4b09ab36dfc9dd514b321aa7b38ce5e8000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00069be86ec9af65eb4a614fd99b8e92547da0145fab5e804adb594db3e73a271b00ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc6600a5d8744c3d0ad845e0d44689b8ea407cc2f491d1bb7931f55cd40810e557e1aa01aa3f371e63ade99128ff53f5fc6b229ff8ec58274bf9e98bb39b9ada262fe7f901046d84465266a1d3008798a5b499bac12326c4b15bfb47fff0a549eb42fe2b6801429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd650128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce030106ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a90006a7d517187bd16635dad40455fdc2c0c124c68f215675a5dbbacb5f0800000000", 'hex'),
                extraData: Buffer.from([]),
            }).accounts({
                gringotts: gringottsPDA,
            }).rpc();

            console.log(tx);
        } catch (err) {
            expect.fail('Init transaction failed:', err)
        }
    });
});