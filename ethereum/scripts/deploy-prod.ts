import hre from "hardhat";
import {ethers, utils} from "ethers";
import bs58 from "bs58";
import {parseAbi} from 'viem';
import {bytesToHex} from "web3-utils"

const ARB_EID = 30110;
const SOL_EID = 30168;

const Address = "0x9c4e6e7e2f2387c3fd9fccc499c18d6c98528931";
const SOL_PDA = "7uQ48SWiZ3PZtahmq1NGgbKpQGG9zY2M7ex1akHCBeUu"

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
                    items: [
                        {
                            asset: utils.hexZeroPad(bytesToHex(Array.from([0])), 32) as `0x${string}`,
                            recipient: utils.hexZeroPad(bytesToHex(Array.from(solReceive)), 32) as `0x${string}`,
                            executionGasAmount: BigInt(400_000),
                            distributionBP: 10_000,
                            swap: {
                                executor: utils.hexZeroPad(bytesToHex(Array.from(jupiter)), 32) as `0x${string}`,
                                command: '0xe517cb977ae3ad2a010000001964000180841e00000000006d6b050000000000640000',
                                metadata: '0x1206ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a900ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc660028d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce0301429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd65010479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00069b8857feab8184fb687f634618c035dac439dc1aeb3b5598a0f00000000001000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00b43ffa27f5d7f64a74c09b1f295879de4b09ab36dfc9dd514b321aa7b38ce5e8000479d55bf231c06eee74c56ece681507fdb1b2dea3f48e5102b1cda256bc138f00069be86ec9af65eb4a614fd99b8e92547da0145fab5e804adb594db3e73a271b00ef31fbe855115c624e9b35d3ea3e1ac581a4b213487345869e0f60f599d8cc66002958b2ea4d21ef6b184cfb7f9e1e6df6b11137db070b88f027700cedb9c6471901e488dd07eb03e922aae5874961a9d63b72f7c8cc4cf91b65cbf931af61a1c232013a4c3b29fdf0ef9860ca2f76e548bd7b9cca00e63f83dd2a432d59fd907a924d01429f70d4e01764b7255936546cfbf0eedbd8a6a887cea01851509632c5c9cd650128d7938c2685268928284987dd9b46409bbf1e22ea56285d653a7fb9abecce030106ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a90006a7d517187bd16635dad40455fdc2c0c124c68f215675a5dbbacb5f0800000000',
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

    const x = await gringotts.write.withdrawERC20([ARB_USDC]);
    console.log(x)
}

async function main() {
    // const address = await deploy();
    // await setDataFeed();
    // await setAgent();

    // await setConfig();

    // await bridge();
    await widthdraw_token();
    // await approve();
    // await getConfigAndDecode();
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

async function getConfigAndDecode() {
    const ethereumLzEndpointAddress = LZ_ENDPOINT;
    const ethereumLzEndpointABI = [
        'function getConfig(address _oapp, address _lib, uint32 _eid, uint32 _configType) external view returns (bytes memory config)',
    ];

    const contract = new ethers.Contract(ethereumLzEndpointAddress, ethereumLzEndpointABI, provider);

    const oappAddress = Address;
    const sendLibAddress = '0x975bcD720be66659e3EB3C0e4F1866a3020E493A';
    const receiveLibAddress = '0x7B9E184e07a6EE1aC23eAe0fe8D6Be2f663f05e6';
    const remoteEid = SOL_EID;
    const executorConfigType = 1; // 1 for executor
    const ulnConfigType = 2; // 2 for UlnConfig

    try {
        // Fetch and decode for sendLib (both Executor and ULN Config)
        const sendExecutorConfigBytes = await contract.getConfig(
            oappAddress,
            sendLibAddress,
            remoteEid,
            executorConfigType,
        );
        const executorConfigAbi = ['tuple(uint32 maxMessageSize, address executorAddress)'];
        const executorConfigArray = ethers.utils.defaultAbiCoder.decode(
            executorConfigAbi,
            sendExecutorConfigBytes,
        );
        console.log('Send Library Executor Config:', executorConfigArray);

        const sendUlnConfigBytes = await contract.getConfig(
            oappAddress,
            sendLibAddress,
            remoteEid,
            ulnConfigType,
        );
        const ulnConfigStructType = [
            'tuple(uint64 confirmations, uint8 requiredDVNCount, uint8 optionalDVNCount, uint8 optionalDVNThreshold, address[] requiredDVNs, address[] optionalDVNs)',
        ];
        const sendUlnConfigArray = ethers.utils.defaultAbiCoder.decode(
            ulnConfigStructType,
            sendUlnConfigBytes,
        );
        console.log('Send Library ULN Config:', sendUlnConfigArray);

        // Fetch and decode for receiveLib (only ULN Config)
        const receiveUlnConfigBytes = await contract.getConfig(
            oappAddress,
            receiveLibAddress,
            remoteEid,
            ulnConfigType,
        );
        const receiveUlnConfigArray = ethers.utils.defaultAbiCoder.decode(
            ulnConfigStructType,
            receiveUlnConfigBytes,
        );
        console.log('Receive Library ULN Config:', receiveUlnConfigArray);
    } catch (error) {
        console.error('Error fetching or decoding config:', error);
    }
}
