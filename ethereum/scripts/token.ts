// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import hre from "hardhat";
import {TOKEN_ABI} from "./token_abi";
import {getContract} from 'viem';

const TOKEN_ADDRESS = "0xC19B63081BbD08fbc0900422CAa1D64c0EB01b3b";
const BRIDGE_ADDRESS = "0xec97e7eeb8fc6b162ef75e478a44471825f96fea";

async function main() {
    const [main, guest] = await hre.viem.getWalletClients();
    console.log(main.account.address, guest.account.address);

    // @ts-ignore
    const token = getContract({
        address: TOKEN_ADDRESS,
        abi: TOKEN_ABI,
        client: {
            public: await hre.viem.getPublicClient(),
            wallet: main,
        },
    });

    const decimals = await token.read.decimals() as bigint;
    const fraction = BigInt(10) ** BigInt(decimals);
    console.log("Token decimals", decimals);
    const mainBalance = await token.read.balanceOf(["0x1d5B493882BDDe19D082426bF32aE3918Cb36916"]) as bigint;
    console.log("Main account balance:", mainBalance / fraction);
    const guestBalance = await token.read.balanceOf(["0x0D595AE2666a2c5Ae6b99cce4DD428a9Cf20B2c9"]) as bigint
    console.log("Guest account balance:", guestBalance / fraction);

    const bridge = await hre.viem.getContractAt(
        "Gringotts",
        BRIDGE_ADDRESS,
        {
            client: {
                wallet: main,
            },
        }
    );

    const amount = BigInt(100) * fraction;
    const approve = await token.write.approve([bridge.address, amount]);
    console.log("Main approves hash:", approve);

    await new Promise(f => setTimeout(f, 1000));
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });