// import artifacts from "../artifacts/contracts/TestTron.sol/TestTron.json";
import artifacts from "../artifacts/contracts/Gringotts.sol/Gringotts.json";
import {TronWeb} from "tronweb";
import {ethers, BigNumber} from "ethers";


const url = "https://api.trongrid.io"
const urlTest = "https://api.nileex.io"
// const tronURL = url;
const tronURL = urlTest;
const tronWeb = new TronWeb(tronURL, tronURL, tronURL, process.env.PRIVATE_KEY);

const contractABI = artifacts.abi;
const contractBytecode = artifacts.bytecode;

const contractAddress = "TLemeSqqQkCMvykYcF4KY1jU2VHLDsKVg9"

async function deployContract() {
    const contract = await tronWeb.contract().new({
        abi: contractABI,
        bytecode: contractBytecode,
        parameters: [6, 'TYsbWxNnyTgsZaTFaue9hqpxkU3Fkco94a'],
        feeLimit: 15_000_000_000
    });

    const hexAddress = contract.address;
    const base58Address = tronWeb.address.fromHex(hexAddress);
    console.log('Contract deployed at address:');
    console.log('Hexadecimal: ', hexAddress);
    console.log('Base58Check: ', base58Address);
}

async function bridge() {
    try {
        const contract = await tronWeb.contract().at(contractAddress);

        const x = await contract.methods.bridge(
            '45000000000000000000',
            "TF17BgPaZYbz8oxbjhriubPDsA7ArKoLX3"
        ).send();

        // const x = await contract.methods.swapFromTRX(
        //     '100000',
        //     "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf"
        // ).send();

        console.log(x);
    } catch (error) {
        console.error(error);
    }
}

async function execute() {
    const command = "zvlSKQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAJAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACYloAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAFc0+Nc96aHV4UECeug52GsCpwJtAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGc9MkMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA+zsxNPE8zSyB9AEuUwJOgTXVj+4AAAAAAAAAAAAAAADsqbyCijAFuaO5CfLMXCpUeU3gXwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAnYyAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA";
    const contract = await tronWeb.contract().at(contractAddress);

    const x = await contract.methods.execute(
        tronWeb.address.toHex("TB6xBCixqRPUSKiXb45ky1GhChFJ7qrfFj"),
        Buffer.from(command, 'base64'),
        10 * 1000000,
        50 * 1000000
    ).send({
        // feeLimit: 200 * 1000000,
    });
    console.log(x);
}

async function main() {
    await deployContract();
    // await execute();
    // await swapFromETHV2();
    // await directTest();
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });

