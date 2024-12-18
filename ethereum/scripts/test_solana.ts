// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import hre from "hardhat";
import {ethers, utils} from "ethers";
import {parseEther, stringToBytes} from "viem";
import bs58 from "bs58";
import {bytesToHex} from "web3-utils"
import {address} from "hardhat/internal/core/config/config-validation";


const Address = "0xcc9462adf0a45db3e9ab95b52829e638f886e1b3";

async function sleep(seconds: number) {
    return new Promise(f => setTimeout(f, seconds * 1000))
}

async function deploy() {
    const gringotts = await hre.viem.deployContract(
        "Gringotts",
        [1, 18, "0x6EDCE65403992e310A62460808c4b910D972f10f"],
    );

    console.log(gringotts.address);
    return gringotts.address;
}

async function setAgent(addressPDA: string) {
    const [wallet] = await hre.viem.getWalletClients();
    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: wallet,
            },
        }
    );

    const decodedBytes = bs58.decode(addressPDA);

    const tx = await gringotts.write.updateAgent([
        2, {
            endpoint: utils.hexZeroPad(bytesToHex(Array.from(decodedBytes)), 32) as `0x${string}`,
            chainID: 2,
            lzEID: 40217,
            stableCoins: [],
            baseGasEstimate: BigInt(30_000),
        }
    ]);

    console.log(tx);
}

async function send() {
    const [main] = await hre.viem.getWalletClients();

    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: main,
            },
        }
    );

    try {
        const tx = await gringotts.write.testSend([40168, "salam"]);
        console.log('Transaction sent:', tx);
    } catch (error) {
        console.error('Error executing testSend:', error);
    }
}

async function back() {
    const [main] = await hre.viem.getWalletClients();

    const c = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: main,
            },
        }
    );

    console.log(await c.write.withdrawNative());
}

async function main() {
    // await deploy();
    // await setAgent("8C7AaZgP3WKvHDjnUuwD1rsSVKGw7zxqk49iEWSNcd5P");
    // await send();
    // await back();
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });