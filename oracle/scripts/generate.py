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
    os.system("abigen --abi={} --pkg=connection --type=GringottsEVM --out=./connection/GringottsEVM.go".format(file_name))

    os.remove(file_name)


if __name__ == '__main__':
    generate_evm()
