import hre from "hardhat";
import {ethers, utils} from "ethers";
import bs58 from "bs58";
import {parseAbi} from 'viem';
import {bytesToHex} from "web3-utils"
import {GetContractReturnType} from "@nomicfoundation/hardhat-viem/types";
import {Gringotts$Type} from "../artifacts/contracts/Gringotts.sol/Gringotts";
import {BytesLike} from "@ethersproject/bytes";

const SOL_EID = 30168;

const Address = "0x02fa4ad1fe96e10f81a258bbe29e0392124bf27c";
const SOL_PDA = "HjVcDEcpjVvzwEnaHkZgdrV2pV2DoFcg7TyemtMovuZg"

const USDC = "0xaf88d065e77c8cc2239327c5edb3a432268e5831";

const LZ_ENDPOINT = "0x1a44076050125825900e736c501f859c50fE728c";

const provider = new ethers.providers.JsonRpcProvider("https://arbitrum-mainnet.infura.io/v3/463ba1f0f1c34f1da7403e6f4e81ebef");
const signer = new ethers.Wallet(process.env.PRIVATE_KEY as BytesLike, provider);

let wallet;
let gringotts: GetContractReturnType<Gringotts$Type["abi"]>;

async function deploy() {
    const gringotts = await hre.viem.deployContract(
        "Gringotts",
        [1, 18, LZ_ENDPOINT],
    );

    console.log(gringotts.address);
    return gringotts.address;
}

async function setDataFeed() {
    const tx = await gringotts.write.setChainlinkPriceFeed(['0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612']);
    console.log(tx);
}

async function setAgent() {
    //solana
    const tx1 = await gringotts.write.updateAgent([
        2, {
            endpoint: utils.hexZeroPad(bytesToHex(Array.from(bs58.decode(SOL_PDA))), 32) as `0x${string}`,
            chainID: 2,
            lzEID: SOL_EID,
            multiSend: true,
            baseGasEstimate: BigInt(30_000),
        }
    ]);
    console.log(tx1);

    const tx2 = await gringotts.write.setConfig([{
        commissionMicroBPS: 1000,
        commissionDiscountBPS: 0,
        gasDiscountBPS: 0,
        stableCoins: [ethers.utils.hexZeroPad(USDC, 32) as `0x${string}`]
    }]);
    console.log(tx2);
}

async function bridge() {
    const x = await gringotts.write.bridge([
        {
            inbound: {
                amountUSDX: BigInt(3 * 1000 * 1000),
                items: [
                    {
                        asset: ethers.utils.hexZeroPad(USDC, 32) as `0x${string}`,
                        amount: BigInt(3 * 1000 * 1000),
                        executor: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                        command: bytesToHex(Array.from([0])) as `0x${string}`,
                        stableToken: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                    }
                ]
            },
            outbounds: [
                {
                    chainId: 2,
                    message: "0x0db66a03ef9666cf81aabb9f5207e73915943bda9b8ffa4a5c387d1704aea5276306a7d517187bd16635dad40455fdc2c0c124c68f215675a5dbbacb5f08000000069b8857feab8184fb687f634618c035dac439dc1aeb3b5598a0f000000000010479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138fb43ffa27f5d7f64a74c09b1f295879de4b09ab36dfc9dd514b321aa7b38ce5e8c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d617d7cdbe036610dfd0b509d8763c05f438474dc88711dd48e5071145043e0a01d069be86ec9af65eb4a614fd99b8e92547da0145fab5e804adb594db3e73a271b451e6d1bc494c362c3035a8ea1115692af3f543f0a4a881711a10b37781df0382495549e89eedc61b20af2a1a5d2354a8198198f76e5ef472e17b2c6e05a1aa8ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc66257527de29a3b62c518bde1e665fcf51b2f5938ae575f0988f5cf4e0970f70b328d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce0382f81b0f0710110a0b08000103000b1008070809080c00050d0e100b030680140024e517cb977ae3ad2a010000003d0164000180841e00000000000ebc8a0000000000640000",
                    executionGas: BigInt(500_000),
                    items: [
                        {
                            distributionBP: 10_000,
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
    const approveAbi = parseAbi(['function approve(address spender, uint256 amount) returns (bool)']);

    const [wallet] = await hre.viem.getWalletClients();

    const hash = await wallet.writeContract({
        address: USDC,
        abi: approveAbi,
        functionName: 'approve',
        args: [Address, BigInt(100 * 1000 * 1000)],
    });

    console.log('Transaction Hash:', hash);
}

async function setConfig() {
    const {ethers} = require('ethers');

    const oappAddress = Address;
    const receiveLibAddress = '0x975bcD720be66659e3EB3C0e4F1866a3020E493A'; // Send
    // const receiveLibAddress = '0x7B9E184e07a6EE1aC23eAe0fe8D6Be2f663f05e6'; // Receive

    const remoteEid = SOL_EID;
    const ulnConfig = {
        confirmations: 32,
        requiredDVNCount: 1,
        optionalDVNCount: 0,
        optionalDVNThreshold: 0,
        requiredDVNs: ['0x2f55c492897526677c5b68fb199ea31e2c126416'],
        optionalDVNs: [],
    };
    const configTypeUlnStruct =
        'tuple(uint64 confirmations, uint8 requiredDVNCount, uint8 optionalDVNCount, uint8 optionalDVNThreshold, address[] requiredDVNs, address[] optionalDVNs)';
    const encodedUlnConfig = ethers.utils.defaultAbiCoder.encode([configTypeUlnStruct], [ulnConfig]);
    const setConfigParam = {
        eid: remoteEid,
        configType: 2,
        config: encodedUlnConfig,
    };

    const executorConfig = {
        maxMessageSize: 10_000,
        executorAddress: '0x31CAe3B7fB82d847621859fb1585353c5720660D',
    };
    const configTypeExecutorStruct = 'tuple(uint32 maxMessageSize, address executorAddress)';
    const encodedExecutorConfig = ethers.utils.defaultAbiCoder.encode([configTypeExecutorStruct], [executorConfig]);
    const setConfigParamExecutor = {
        eid: remoteEid,
        configType: 1,
        config: encodedExecutorConfig,
    };

    const endpointAbi = [
        'function setConfig(address oappAddress, address receiveLibAddress, tuple(uint32 eid, uint32 configType, bytes config)[] setConfigParams) external',
    ];
    const endpointContract = new ethers.Contract(LZ_ENDPOINT, endpointAbi, signer);
    const tx = await endpointContract.setConfig(
        oappAddress,
        receiveLibAddress,
        [setConfigParam, setConfigParamExecutor], // Send
        // [setConfigParam], // Receive
        {
            gasLimit: BigInt(1_000_000),
        }
    );

    console.log('Transaction sent:', tx.hash);
}

async function widthdraw_token() {
    const x = await gringotts.write.withdrawERC20([USDC, BigInt(0)]);
    console.log(x)

    const y = await gringotts.write.withdrawNative([BigInt(0)]);
    console.log(y)
}

async function test() {
    // const x = await gringotts.write.testSend(
    //     [
    //         2,
    //         2,
    //         1,
    //         "0x0100000000001e848000083a4c3b29fdf0ef9860ca2f76e548bd7b9cca00e63f83dd2a432d59fd907a924d06a7d517187bd16635dad40455fdc2c0c124c68f215675a5dbbacb5f08000000ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc66069b8857feab8184fb687f634618c035dac439dc1aeb3b5598a0f00000000001b43ffa27f5d7f64a74c09b1f295879de4b09ab36dfc9dd514b321aa7b38ce5e8069be86ec9af65eb4a614fd99b8e92547da0145fab5e804adb594db3e73a271b2958b2ea4d21ef6b184cfb7f9e1e6df6b11137db070b88f027700cedb9c64719e488dd07eb03e922aae5874961a9d63b72f7c8cc4cf91b65cbf931af61a1c232a3090f0a101202000610040a040b040c000d0e07100602080024e517cb977ae3ad2a010000003d0164000180841e000000000079b6a50000000000f40100",
    //         BigInt(200_000),
    //         false
    //     ]);
    // console.log(x)

    const x = await gringotts.write.testSend(
        [
            2,
            1,
            1,
            "0xf4269880b8a00550d1698545abf89b76f6e1438c951733e5a845c148e78c2f11",
            BigInt(400_000),
            false
        ]);
    console.log(x)
}

async function main() {
    [wallet] = await hre.viem.getWalletClients();
    gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: wallet,
            },
        }
    );

    // const address = await deploy();
    // await setDataFeed();
    // await setAgent();
    // await setConfig();

    // await bridge();
    await widthdraw_token();
    // await approve();
    // await test();
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });