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

    const x = await gringotts.write.bridge([
        {
            inbound: {
                amountUSDX: BigInt(10 * 1000 * 1000),
                items: [
                    {
                        asset: ethers.utils.hexZeroPad(ARB_USDC, 32) as `0x${string}`,
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
            outbounds: [
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
        address: ARB_USDC,
        abi,
        functionName: 'decimals',
    });

    console.log('Token Decimals:', decimals);

    const [wallet] = await hre.viem.getWalletClients();

    const hash = await wallet.writeContract({
        address: ARB_USDC,
        abi: approveAbi,
        functionName: 'approve',
        args: [Address, BigInt(10) * (BigInt(10) ** BigInt(decimals))],
    });

    console.log('Transaction Hash:', hash);
}

async function setConfig() {
    const {ethers} = require('ethers');

    const oappAddress = Address;
    const receiveLibAddress = '0x7B9E184e07a6EE1aC23eAe0fe8D6Be2f663f05e6';

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
        configType: 2, // RECEIVE_CONFIG_TYPE
        config: encodedUlnConfig,
    };

    const endpointAbi = [
        'function setConfig(address oappAddress, address receiveLibAddress, tuple(uint32 eid, uint32 configType, bytes config)[] setConfigParams) external',
    ];

    const endpointContract = new ethers.Contract(LZ_ENDPOINT, endpointAbi, signer);

    const tx = await endpointContract.setConfig(
        oappAddress,
        receiveLibAddress,
        [setConfigParam],
    );

    console.log('Transaction sent:', tx.hash);
}

async function setLibraries() {
    const endpointAbi = [
        'function setSendLibrary(address oapp, uint32 eid, address sendLib) external',
        'function setReceiveLibrary(address oapp, uint32 eid, address receiveLib) external',
    ];
    const endpointContract = new ethers.Contract(LZ_ENDPOINT, endpointAbi, signer);

    try {
        // Set the send library
        const sendTx = await endpointContract.setSendLibrary(
            Address,
            SOL_EID,
            '0x975bcD720be66659e3EB3C0e4F1866a3020E493A',
            {
                gasLimit: BigInt(200_000),
            }
        );
        console.log('Send library transaction sent:', sendTx.hash);
        await sendTx.wait();
        console.log('Send library set successfully.');

        // Set the receive library
        const receiveTx = await endpointContract.setReceiveLibrary(
            Address,
            SOL_EID,
            '0x7B9E184e07a6EE1aC23eAe0fe8D6Be2f663f05e6',
            {
                gasLimit: BigInt(200_000),
            }
        );
        console.log('Receive library transaction sent:', receiveTx.hash);
        await receiveTx.wait();
        console.log('Receive library set successfully.');
    } catch (error) {
        console.error('Transaction failed:', error);
    }
}

async function main() {
    // const address = await deploy();
    // await setDataFeed();
    // await setAgent();

    // await setLibraries();
    // await setConfig();

    // await bridge();
    // await approve();
    await getConfigAndDecode();
}

// Define the smart contract address and ABI
const ethereumLzEndpointAddress = LZ_ENDPOINT;
const ethereumLzEndpointABI = [
    'function getConfig(address _oapp, address _lib, uint32 _eid, uint32 _configType) external view returns (bytes memory config)',
];

// Create a contract instance
const contract = new ethers.Contract(ethereumLzEndpointAddress, ethereumLzEndpointABI, provider);

// Define the addresses and parameters
const oappAddress = Address;
const sendLibAddress = '0x975bcD720be66659e3EB3C0e4F1866a3020E493A';
const receiveLibAddress = '0x7B9E184e07a6EE1aC23eAe0fe8D6Be2f663f05e6';
const remoteEid = SOL_EID;
const executorConfigType = 1; // 1 for executor
const ulnConfigType = 2; // 2 for UlnConfig

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

async function getConfigAndDecode() {
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
