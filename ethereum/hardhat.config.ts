import type {HardhatUserConfig} from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox-viem";
import dotenv from 'dotenv';

dotenv.config();


const config: HardhatUserConfig = {
    networks: {
        hardhat: {},
        holesky: {
            // url: "https://holesky.drpc.org",
            url: "https://holesky.infura.io/v3/463ba1f0f1c34f1da7403e6f4e81ebef",
            chainId: 17000,
            accounts: [process.env.PRIVATE_KEY || '', process.env.PRIVATE_KEY2 || ''],
        },
        sepolia: {
            url: "https://sepolia.drpc.org",
            chainId: 11155111,
            accounts: [process.env.PRIVATE_KEY || '', process.env.PRIVATE_KEY2 || ''],
        },
        arb: {
            url: "https://arbitrum-mainnet.infura.io/v3/463ba1f0f1c34f1da7403e6f4e81ebef",
            chainId: 42161,
            accounts: [process.env.PRIVATE_KEY || '', process.env.PRIVATE_KEY2 || ''],
        },
        arb_test: {
            url: "https://arbitrum-sepolia.infura.io/v3/463ba1f0f1c34f1da7403e6f4e81ebef",
            chainId: 421614,
            accounts: [process.env.PRIVATE_KEY || '', process.env.PRIVATE_KEY2 || ''],
        }
    },
    solidity: {
        version: "0.8.27",
        settings: {
            optimizer: {
                enabled: true,
                runs: 200
            },
            viaIR: true,
        }
    }
};

export default config;
