import hre from "hardhat";
import {ethers, utils} from "ethers";
import bs58 from "bs58";
import {parseAbi} from 'viem';
import {bytesToHex} from "web3-utils"

const ARB_EID = 40231;
const SOL_EID = 40168;

const Address = "0x0b481d55839b6118a917e89dd53eb35e5359fae9";
const SOL_PDA = "HjVcDEcpjVvzwEnaHkZgdrV2pV2DoFcg7TyemtMovuZg"

const TEST_STABLE = "0x301b022b40d06088fc974e767149f4a3feebbf1a";
const SOL_USDC = "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU";

async function deploy() {
    const gringotts = await hre.viem.deployContract(
        "Gringotts",
        [1, 18, "0x6EDCE65403992e310A62460808c4b910D972f10f"],
    );

    console.log(gringotts.address);
    const [wallet] = await hre.viem.getWalletClients();

    const tx = await wallet.sendTransaction({
        to: gringotts.address,
        value: BigInt(1e18 * 0.01),
    });

    console.log("Fund -> ", tx);
}

async function setDataFeed() {
    const [wallet] = await hre.viem.getWalletClients();
    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: wallet,
            },
        }
    );

    const tx = await gringotts.write.setChainlinkPriceFeed(['0xd30e2101a97dcbAeBCBC04F14C3f624E67A35165']);
    console.log(tx);
}

async function setAgent() {
    const [wallet] = await hre.viem.getWalletClients();
    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: wallet,
            },
        }
    );

    //solana
    const decodedBytes = bs58.decode(SOL_PDA);
    const stableCoin = bs58.decode(SOL_USDC);
    const tx1 = await gringotts.write.updateAgent([
        2, {
            endpoint: utils.hexZeroPad(bytesToHex(Array.from(decodedBytes)), 32) as `0x${string}`,
            chainID: 2,
            lzEID: SOL_EID,
            multiSend: true,
            baseGasEstimate: BigInt(100_000),
        }
    ]);
    console.log(tx1);
}

async function bridge() {
    const [wallet] = await hre.viem.getWalletClients();
    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: wallet,
            },
        }
    );

    const solStableCoin = bs58.decode(SOL_USDC);
    const solReceive = bs58.decode("H6ieTjyWqcRFv1RwD9LghEFDxVtMaEMgVbgnYrdHMjr5");

    const x = await gringotts.write.bridge([
        {
            inTransfer: {
                amountUSDX: BigInt(10 * 1000 * 1000),
                items: [
                    {
                        asset: ethers.utils.hexZeroPad(TEST_STABLE, 32) as `0x${string}`,
                        amount: BigInt(10) * (BigInt(10) ** BigInt(18)),
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
                            executionGasAmount: BigInt(60_000),
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
    ], {
        gas: BigInt(1_000_000_000),
    });

    console.log(x);
}

async function approve() {
    const abi = parseAbi(['function decimals() public view returns (uint8)']);
    const approveAbi = parseAbi(['function approve(address spender, uint256 amount) returns (bool)']);

    const decimals = await (await hre.viem.getPublicClient()).readContract({
        address: TEST_STABLE,
        abi,
        functionName: 'decimals',
    });

    console.log('Token Decimals:', decimals);

    const [wallet] = await hre.viem.getWalletClients();

    const hash = await wallet.writeContract({
        address: TEST_STABLE,
        abi: approveAbi,
        functionName: 'approve',
        args: [Address, BigInt(100) * (BigInt(10) ** BigInt(decimals))],
    });

    console.log('Transaction Hash:', hash);
}

async function test() {
    const [wallet] = await hre.viem.getWalletClients();
    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: wallet,
            },
        }
    );

    const x = await gringotts.write.testSend(
        [
            2,
            2,
            1,
            "0x1245",
            BigInt(400_000),
            false
        ]);
    console.log(x)
}

async function widthdraw_token() {
    const [wallet] = await hre.viem.getWalletClients();
    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: wallet,
            },
        }
    );

    // const x = await gringotts.write.withdrawERC20(["0xaf88d065e77c8cc2239327c5edb3a432268e5831", BigInt(0)]);
    // console.log(x)

    const y = await gringotts.write.withdrawNative([BigInt(0)]);
    console.log(y)
}


async function main() {
    // const address = await deploy();
    // await setDataFeed();
    // await setAgent();
    // await bridge();
    await test();
    // await approve();
    // await widthdraw_token()
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });