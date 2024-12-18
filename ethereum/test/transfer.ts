import {expect} from "chai";
import hre from "hardhat";
import {parseEther, formatEther} from "viem";
import {Contract, ContractFactory, ethers, utils} from 'ethers';
import {bytesToHex} from "web3-utils"
import bs58 from "bs58";

const Address = "0x22ec28a484065dad34cd7600ec1319ce25406a93";
const SOL_PDA = "AkYxtYn81jkDS7hEceyq4ibCjzVagzDuStbPN6R6Ua5x"
const TEST_STABLE = "0x301b022b40d06088fc974e767149f4a3feebbf1a";
const SOL_USDC = "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU";
const ARB_EID = 40231;
const SOL_EID = 40168;

describe("Transfer contract", function () {

    let layerZero;
    let priceFeed;
    // @ ts-ignore
    let gringotts;

    beforeEach(async () => {
        const [main] = await hre.viem.getWalletClients();

        layerZero = await hre.viem.deployContract("MockLayerZero", [40217])
        priceFeed = await hre.viem.deployContract("MockChainLink");

        gringotts = await hre.viem.deployContract(
            "Gringotts",
            [1, 1, layerZero.address],
        );

        const decodedBytes = bs58.decode(SOL_PDA);
        const stableCoin = bs58.decode(SOL_USDC);
        const tx1 = await gringotts.write.updateAgent([
            2, {
                endpoint: utils.hexZeroPad(bytesToHex(Array.from(decodedBytes)), 32) as `0x${string}`,
                chainID: 2,
                lzEID: SOL_EID,
                stableCoins: [utils.hexZeroPad(bytesToHex(Array.from(stableCoin)), 32) as `0x${string}`],
                baseGasEstimate: BigInt(30_000),
            }
        ]);
        console.log(tx1);

        // self
        const tx2 = await gringotts.write.updateAgent([
            1, {
                endpoint: ethers.utils.hexZeroPad(Address, 32) as `0x${string}`,
                chainID: 1,
                lzEID: ARB_EID,
                stableCoins: [ethers.utils.hexZeroPad(TEST_STABLE, 32) as `0x${string}`],
                baseGasEstimate: BigInt(30_000),
            }
        ]);
        await main.sendTransaction({to: gringotts.address, value: parseEther("1")});
    });

    it("One to One Bridge", async function () {
        const [main] = await hre.viem.getWalletClients();
        const blockchain = await hre.viem.getPublicClient();

        const solStableCoin = bs58.decode(SOL_USDC);
        const solReceive = bs58.decode("H6ieTjyWqcRFv1RwD9LghEFDxVtMaEMgVbgnYrdHMjr5");

        const tx = await gringotts.write.bridge([
            {
                inTransfer: {
                    amountUSDX: BigInt(1000 * 1000),
                    items: [
                        {
                            asset: ethers.utils.hexZeroPad(TEST_STABLE, 32) as `0x${string}`,
                            amount: BigInt(1) * (BigInt(10) ** BigInt(18)),
                            swap: {
                                executor: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                                command: bytesToHex(Array.from([0])) as `0x${string}`,
                                metadata: bytesToHex(Array.from([0])) as `0x${string}`,
                                stableToken: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                            }
                        }
                    ]
                },
                outTransfers: [
                    {
                        chainId: 2,
                        items: [
                            {
                                asset: utils.hexZeroPad(bytesToHex(Array.from(solStableCoin)), 32) as `0x${string}`,
                                recipient: utils.hexZeroPad(bytesToHex(Array.from(solReceive)), 32) as `0x${string}`,
                                executionGasAmount: BigInt(1_000_000),
                                distributionBP: 10_000,
                                swap: {
                                    executor: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                                    command: bytesToHex(Array.from([0])) as `0x${string}`,
                                    metadata: bytesToHex(Array.from([0])) as `0x${string}`,
                                    stableToken: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                                },
                            }
                        ]
                    }
                ]
            }
        ]);

        const txR = await blockchain.getTransactionReceipt({hash: tx});
    });
});