package evm

import (
	"blockchain-tracking/internal/blockchain/evm/abi/erc20"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

func HexToBigInt(hexStr string) *big.Int {
	bi := new(big.Int)
	bi.SetString(hexStr, 16)
	return bi
}

func ContractType(address common.Address, client *ethclient.Client) (string, *big.Int, int, error) {

	callOpts := &bind.CallOpts{}

	// ERC165 interface IDs
	erc721ID := [4]byte{0x80, 0xac, 0x58, 0xcd}
	erc1155ID := [4]byte{0xd9, 0xb6, 0x7a, 0x26}

	contractAbi, err := abi.JSON(strings.NewReader(erc20.Erc20MetaData.ABI))
	if err != nil {
		fmt.Print("error parsing abi")
		return "", nil, -1, err
	}

	// Bind to generic contract
	bound := bind.NewBoundContract(address, contractAbi, client, client, client)

	// Check ERC721
	var out []interface{}
	if err := bound.Call(callOpts, &out, "supportsInterface", erc721ID); err == nil {
		if len(out) > 0 {
			if supports, ok := out[0].(bool); ok && supports {
				return "721", nil, 0, nil
			}
		}
	}

	// Check ERC1155
	out = nil
	if err := bound.Call(callOpts, &out, "supportsInterface", erc1155ID); err == nil {
		if len(out) > 0 {
			if supports, ok := out[0].(bool); ok && supports {
				return "1155", nil, 0, nil
			}
		}
	}

	// Fallback to ERC20
	erc20Instance, err := erc20.NewErc20(address, client)
	if err == nil {
		supply, err1 := erc20Instance.TotalSupply(callOpts)
		decimals, err2 := erc20Instance.Decimals(callOpts)
		if err1 == nil && err2 == nil {
			return "20", supply, int(decimals), nil
		}
	}

	return "", nil, 0, nil
}
