Openzeppelin
https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v5.0.1/contracts/token/ERC20/IERC20.sol

1. make .sol file
2. solc --abi contract.sol -o ./
3. solc --bin contract.sol -o ./
3. abigen --bin=Contract.bin --abi=Contract.abi --pkg=contracts --out=contract.go

https://gist.github.com/qbig/4708b468014fb84fdbbcb70453132a4a