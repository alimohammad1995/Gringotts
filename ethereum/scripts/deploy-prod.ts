import hre from "hardhat";
import {ethers, utils} from "ethers";
import bs58 from "bs58";
import {parseAbi} from 'viem';
import {bytesToHex} from "web3-utils"

const ARB_EID = 30110;
const SOL_EID = 30168;

const Address = "0x7794d4260bf7c0c975dd0df59c4f67c1631eea51";
const SOL_PDA = "BKbhvxVKdSrKnPTKiustSo4Ne9qnNu3gJYeLYh7TTMEk"

const ARB_USDC = "0xaf88d065e77c8cc2239327c5edb3a432268e5831";
const SOL_USDC = "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v";

const LZ_ENDPOINT = "0x1a44076050125825900e736c501f859c50fE728c";

const provider = new ethers.providers.JsonRpcProvider("https://arbitrum-mainnet.infura.io/v3/463ba1f0f1c34f1da7403e6f4e81ebef");
const signer = new ethers.Wallet(process.env.PRIVATE_KEY, provider);

async function deploy() {
    const gringotts = await hre.viem.deployContract(
        "Gringotts",
        [1, 18, LZ_ENDPOINT],
    );

    console.log(gringotts.address);
    return gringotts.address;
}

async function setDataFeed() {
    const [wallet] = await hre.viem.getWalletClients();
    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: wallet,
            },
        }
    );

    const tx = await gringotts.write.setChainlinkPriceFeed(['0x639Fe6ab55C921f74e7fac1ee960C0B6293ba612']);
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
    const endpointAddress = bs58.decode(SOL_PDA);
    const solanaStableCoin = bs58.decode(SOL_USDC);
    const tx1 = await gringotts.write.updateAgent([
        2, {
            endpoint: utils.hexZeroPad(bytesToHex(Array.from(endpointAddress)), 32) as `0x${string}`,
            chainID: 2,
            lzEID: SOL_EID,
            stableCoins: [utils.hexZeroPad(bytesToHex(Array.from(solanaStableCoin)), 32) as `0x${string}`],
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
            stableCoins: [ethers.utils.hexZeroPad(ARB_USDC, 32) as `0x${string}`],
            baseGasEstimate: BigInt(30_000),
        }
    ]);
    console.log(tx2);
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
    const jupiter = bs58.decode("JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4");

    const x = await gringotts.write.bridge([
        {
            inbound: {
                amountUSDX: BigInt(3 * 1000 * 1000),
                items: [
                    {
                        asset: ethers.utils.hexZeroPad(ARB_USDC, 32) as `0x${string}`,
                        amount: BigInt(3 * 1000 * 1000),
                        swap: {
                            executor: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                            command: bytesToHex(Array.from([0])) as `0x${string}`,
                            metadata: bytesToHex(Array.from([0])) as `0x${string}`,
                            stableToken: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                        }
                    }
                ]
            },
            outbounds: [
                {
                    chainId: 2,
                    metadata: "0x12c6fa7af3bedbad3a3d65f36aabc97431b1bbe4c2d2f6e0e47ca60203452f5d61005bca11acc7a4adc5407b711476383212d9637789a0a1bf69ae857609e8a196610128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce0301aa3f371e63ade99128ff53f5fc6b229ff8ec58274bf9e98bb39b9ada262fe7f901b43ffa27f5d7f64a74c09b1f295879de4b09ab36dfc9dd514b321aa7b38ce5e800046d84465266a1d3008798a5b499bac12326c4b15bfb47fff0a549eb42fe2b680106a7d517187bd16635dad40455fdc2c0c124c68f215675a5dbbacb5f0800000000ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc66010479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f007103b1d6bd9f22e5a2fcf19bc626584e55472227358076a52bcb1e9a1d1ae39a01069b8857feab8184fb687f634618c035dac439dc1aeb3b5598a0f0000000000100069be86ec9af65eb4a614fd99b8e92547da0145fab5e804adb594db3e73a271b00a5d8744c3d0ad845e0d44689b8ea407cc2f491d1bb7931f55cd40810e557e1aa0100010203040c050e060d07000103000e060d0f0d090d100011080a060e030b0c050e07",
                    items: [
                        {
                            asset: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                            recipient: utils.hexZeroPad(bytesToHex(Array.from(solReceive)), 32) as `0x${string}`,
                            executionGas: BigInt(400_000),
                            distributionBP: 10_000,
                            swap: {
                                executor: utils.hexZeroPad(bytesToHex(Array.from(jupiter)), 32) as `0x${string}`,
                                command: '0xe517cb977ae3ad2a010000001964000180841e00000000006d6b050000000000640000',
                                metadata: '0x12',
                                stableToken: utils.hexZeroPad(bytesToHex(Array.from(solStableCoin)), 32) as `0x${string}`,
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
    const approveAbi = parseAbi(['function approve(address spender, uint256 amount) returns (bool)']);

    const [wallet] = await hre.viem.getWalletClients();

    const hash = await wallet.writeContract({
        address: ARB_USDC,
        abi: approveAbi,
        functionName: 'approve',
        args: [Address, BigInt(10 * 1000 * 1000)],
    });

    console.log('Transaction Hash:', hash);
}

async function setConfig() {
    const {ethers} = require('ethers');

    const oappAddress = Address;
    const receiveLibAddress = '0x975bcD720be66659e3EB3C0e4F1866a3020E493A';
    // const receiveLibAddress = '0x7B9E184e07a6EE1aC23eAe0fe8D6Be2f663f05e6';

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
    const encodedExecutorConfig = ethers.utils.defaultAbiCoder.encode(
        [configTypeExecutorStruct],
        [executorConfig],
    );
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
        [setConfigParam, setConfigParamExecutor],
        {
            gasLimit: BigInt(200_000),
        }
    );

    console.log('Transaction sent:', tx.hash);
}

async function widthdraw_token() {
    const [wallet] = await hre.viem.getWalletClients();
    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: wallet,
            },
        }
    );

    const x = await gringotts.write.withdrawNative();
    console.log(x)
}

async function main() {
    // const address = await deploy();
    // await setDataFeed();
    // await setAgent();

    // await setConfig();

    // await bridge();
    // await widthdraw_token();
    // await approve();
    // await getConfigAndDecode();
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });