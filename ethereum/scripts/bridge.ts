// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import hre from "hardhat";
import {utils} from "ethers";

const Address = "0xeb44d9cf3df41bc623b1d54eb62971d194fa0385";

async function sleep(seconds: number) {
    return new Promise(f => setTimeout(f, seconds * 1000))
}

async function deploy() {
    const gringotts = await hre.viem.deployContract(
        "Gringotts",
        [1, 1, "0x6EDCE65403992e310A62460808c4b910D972f10f"],
    );

    console.log(gringotts.address);

    await sleep(1);
    const t1 = await gringotts.write.setChainMapping([1, 40217]);
    console.log(t1)

    await sleep(1);
    const t2 = await gringotts.write.setPeer([40217, utils.hexZeroPad(gringotts.address, 32) as `0x${string}`]);
    console.log(t2);

    await sleep(1);
    const t3 = await gringotts.write.setDataFeed([1, "0x9ba2d448c8b9127ef11bf0e1e4bffda7ae7a0fdc"]);
    console.log(t3);
}

async function estimate() {
    const [main] = await hre.viem.getWalletClients();

    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: main,
            },
        }
    );

    const amount = BigInt(100_000_000_000_000_000) //0.10 ETH
    const [commission, gas, finalAmounts] = await gringotts.read.estimate([
        1,
        amount,
        1,
        [1],
        [100],
        [main.account.address],
    ]);

    console.log(commission, gas, finalAmounts, finalAmounts.reduce((acc: bigint, current) => acc + current, 0n));
}

async function bridge() {
    const [main] = await hre.viem.getWalletClients();

    const gringotts = await hre.viem.getContractAt("Gringotts", Address, {
            client: {
                wallet: main,
            },
        }
    );

    const amount = BigInt(100_000_000_000_000_000) //0.10 ETH
    const tx = await gringotts.write.bridge([
        1,
        amount,
        1,
        [1],
        [100],
        ["0x0D595AE2666a2c5Ae6b99cce4DD428a9Cf20B2c9"],
    ], {
        value: amount,
    })

    console.log(tx);
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
    // await estimate();
    await bridge();

    // await back();
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });