import json
import os


def generate_evm():
    file_name = 'gringotts.json'

    print('\033[94m' + "Compiling solidity")
    os.system("cd ../ethereum; npx hardhat compile; cd ../oracle")

    lines = open("../ethereum/artifacts/contracts/Gringotts.sol/Gringotts.json", 'r').readlines()
    items = json.loads(''.join(lines))

    f = open(file_name, "w")
    f.write(json.dumps(items["abi"]))
    f.close()

    print('\033[94m' + "Generating Go file")
    os.system(
        "abigen --abi={} --pkg=connection --type=GringottsEVM --out=./connection/GringottsEVM.go".format(file_name))

    os.remove(file_name)


def generate_sol():
    print('\033[94m' + "Compiling solana")
    os.system("cd ../solana; anchor build; cd ../oracle")

    lines = open("../solana/target/idl/gringotts.json", 'r').readlines()
    items = json.loads(''.join(lines))

    estimate_discriminator = bridge_discriminator = ""
    for instruction in items["instructions"]:
        if instruction["name"] == "estimate":
            estimate_discriminator = instruction["discriminator"]
        if instruction["name"] == "bridge":
            bridge_discriminator = instruction["discriminator"]

    lines = open("config/solana.json", 'r').readlines()
    items = json.loads(''.join(lines))
    items["estimate_discriminator"] = estimate_discriminator
    items["bridge_discriminator"] = bridge_discriminator

    f = open("config/solana.json", "w")
    f.write(json.dumps(items, indent=4))
    f.close()


if __name__ == '__main__':
    generate_evm()
    generate_sol()
