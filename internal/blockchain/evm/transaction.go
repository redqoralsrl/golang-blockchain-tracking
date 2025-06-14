package evm

import (
	"blockchain-tracking/internal/blockchain/evm/abi/contract"
	"blockchain-tracking/internal/blockchain/evm/abi/erc1155"
	"blockchain-tracking/internal/blockchain/evm/abi/erc20"
	"blockchain-tracking/internal/blockchain/evm/abi/erc721"
	"blockchain-tracking/internal/blockchain/evm/evmType"
	"blockchain-tracking/internal/blockchain/jsonRpc"
	"blockchain-tracking/internal/logger"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
	"sync"
)

func FetchTransactionData(ctx context.Context, rpc string, block *jsonRpc.Block, result *evmType.Block, jrAdapter *jsonRpc.JsonRpc, l logger.Logger) error {
	tx := make([]jsonRpc.Transaction, len(block.Transactions))

	payloads := make([]jsonRpc.Payload, 0, len(block.Transactions))
	for index, transaction := range block.Transactions {
		p := jsonRpc.Payload{
			Jsonrpc: "2.0",
			Method:  "eth_getTransactionReceipt",
			Params:  []interface{}{fmt.Sprintf("%s", transaction.Hash)},
			ID:      index,
		}
		tx[index] = transaction
		payloads = append(payloads, p)
	}

	ethGetTransactionReceiptRes, err := jrAdapter.CreateRequestMultiple(rpc, payloads)
	if err != nil {
		l.Error("eth_getTransactionReceipt", logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	var ethTransactionReceipt []*jsonRpc.EthGetTransactionReceipt
	err = json.Unmarshal(ethGetTransactionReceiptRes, &ethTransactionReceipt)
	if err != nil {
		l.Error("json unmarshal transaction data", logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	client, err := ethclient.DialContext(ctx, rpc)
	defer client.Close()
	if err != nil {
		l.Error("eth client dial context", logger.Field{
			Key:   "error",
			Value: err.Error(),
		})
		return err
	}

	erc20Erc721TransferSig := []byte("Transfer(address,address,uint256)")
	erc1155TransferSingleSig := []byte("TransferSingle(address,address,address,uint256,uint256)")
	erc1155TransferBatchSig := []byte("TransferBatch(address,address,address,uint256[],uint256[])")
	erc20Erc721TransferSigHash := crypto.Keccak256Hash(erc20Erc721TransferSig)
	erc1155TransferSingleSigHash := crypto.Keccak256Hash(erc1155TransferSingleSig)
	erc1155TransferBatchSigHash := crypto.Keccak256Hash(erc1155TransferBatchSig)

	goroutineErr := make(chan error, len(ethTransactionReceipt))

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, res := range ethTransactionReceipt {
		wg.Add(1)

		go func(txReceipt *jsonRpc.TransactionReceipt, idx int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					errMsg := fmt.Sprintf("panic: %v", r)
					l.Warn(fmt.Sprintf("fetch transaction data recover %s", result.ChainID.String()), logger.Field{
						Key:   "error",
						Value: errMsg,
					})
				}
			}()
			gasInt := big.NewInt(0)
			if tx[idx].Gas != "0x0" {
				gasInt = HexToBigInt(tx[idx].Gas[2:])
			}
			gasPriceInt := big.NewInt(0)
			if tx[idx].GasPrice != "0x0" {
				gasPriceInt = HexToBigInt(tx[idx].GasPrice[2:])
			}
			valueInt := big.NewInt(0)
			if tx[idx].Value != "0x0" {
				valueInt = HexToBigInt(tx[idx].Value[2:])
			}
			gasUsedInt := big.NewInt(0)
			if txReceipt.GasUsed != "0x0" {
				gasUsedInt = HexToBigInt(txReceipt.GasUsed[2:])
			}
			if txReceipt.ContractAddress != nil {
				lowered := strings.ToLower(*txReceipt.ContractAddress)
				txReceipt.ContractAddress = &lowered
			}

			mu.Lock()
			result.Transaction[idx] = evmType.Transaction{
				ChainID:          result.ChainID,
				BlockHash:        strings.ToLower(result.Hash),
				BlockNumber:      result.Number,
				BlockNumberInt:   result.NumberInt,
				From:             strings.ToLower(txReceipt.From),
				Gas:              tx[idx].Gas,
				GasInt:           gasInt,
				GasPrice:         tx[idx].GasPrice,
				GasPriceInt:      gasPriceInt,
				Hash:             strings.ToLower(txReceipt.TransactionHash),
				R:                tx[idx].R,
				S:                tx[idx].S,
				V:                tx[idx].V,
				TransactionIndex: tx[idx].TransactionIndex,
				Value:            tx[idx].Value,
				ValueInt:         valueInt,
				Nonce:            tx[idx].Nonce,
				Input:            tx[idx].Input,
				ContractAddress:  txReceipt.ContractAddress,
				GasUsed:          txReceipt.GasUsed,
				GasUsedInt:       gasUsedInt,
				Logs:             nil,
				Status:           txReceipt.Status,
				Type:             txReceipt.Type,
				Timestamp:        result.Timestamp,
				TimestampInt:     result.TimestampInt,
				CreatedAt:        result.CreatedAt,
				Contract:         nil,
				CoinLogs:         nil,
				Erc20Logs:        nil,
				Erc721Logs:       nil,
				Erc1155Logs:      nil,
				CoinCount:        big.NewInt(0),
				NftCount:         big.NewInt(0),
				Erc20Count:       big.NewInt(0),
				Erc721Count:      big.NewInt(0),
				Erc1155Count:     big.NewInt(0),
			}
			if txReceipt.To != nil {
				*txReceipt.To = strings.ToLower(*txReceipt.To)
				result.Transaction[idx].To = *txReceipt.To
			}
			// coin tracking
			if valueInt.Cmp(big.NewInt(0)) > 0 {
				result.Transaction[idx].CoinLogs = &evmType.CoinLog{
					ChainID:         result.ChainID,
					Timestamp:       result.Timestamp,
					TimestampInt:    result.TimestampInt,
					CreatedAt:       result.CreatedAt,
					TransactionHash: strings.ToLower(txReceipt.TransactionHash),
					From:            strings.ToLower(txReceipt.From),
					To:              *txReceipt.To,
					Amount:          valueInt,
					Gas:             tx[idx].Gas,
					GasInt:          gasInt,
					GasPrice:        tx[idx].GasPrice,
					GasPriceInt:     gasPriceInt,
					GasUsed:         txReceipt.GasUsed,
					GasUsedInt:      gasUsedInt,
				}
				result.Transaction[idx].CoinCount = new(big.Int).Add(result.Transaction[idx].CoinCount, big.NewInt(1))
			}
			mu.Unlock()

			// is contract create
			if txReceipt.To == nil && txReceipt.ContractAddress != nil {
				contractInput := evmType.Contract{
					ChainID:     result.ChainID,
					Hash:        *txReceipt.ContractAddress,
					Name:        "",
					Symbol:      "",
					Decimals:    0,
					TotalSupply: big.NewInt(0),
					Type:        "",
					Creator:     "",
				}

				contractAddress := common.HexToAddress(*txReceipt.ContractAddress)

				typeStr, _, _, err := ContractType(contractAddress, client)
				if err == nil {
					contractInput.Type = typeStr
				}

				contractInstance, err := contract.NewContract(contractAddress, client)
				if err != nil {
					l.Error("create new instance", logger.Field{Key: "error", Value: err.Error()})
					goroutineErr <- err
					return
				}

				name, err := contractInstance.Name(&bind.CallOpts{})
				if err == nil {
					contractInput.Name = name
				}
				symbol, err := contractInstance.Symbol(&bind.CallOpts{})
				if err == nil {
					contractInput.Symbol = symbol
				}

				erc20Instance, err := erc20.NewErc20(contractAddress, client)
				if err == nil {
					totalSupply, err := erc20Instance.TotalSupply(&bind.CallOpts{})
					if err == nil {
						contractInput.TotalSupply = totalSupply
					}
					decimals, err := erc20Instance.Decimals(&bind.CallOpts{})
					if err == nil {
						contractInput.Decimals = int(decimals)
					}
				}

				mu.Lock()
				result.Transaction[idx].Contract = &contractInput
				mu.Unlock()
			}

			mu.Lock()
			result.Transaction[idx].Logs = make([]*evmType.Log, len(txReceipt.Logs))
			mu.Unlock()

			for index, txLogs := range txReceipt.Logs {
				if len(txLogs.Topics) > 0 {
					switch txLogs.Topics[0].Hex() {
					case erc20Erc721TransferSigHash.Hex():
						if txLogs.Data != "0x" {
							// erc20
							function := "transfer"

							from := common.HexToAddress(txLogs.Topics[1].Hex())
							to := common.HexToAddress(txLogs.Topics[2].Hex())

							if from == common.HexToAddress("0x0000000000000000000000000000000000000000") {
								function = "mint"
							} else if to == common.HexToAddress("0x0000000000000000000000000000000000000000") || to == common.HexToAddress("0x000000000000000000000000000000000000dead") {
								function = "burn"
							}

							contractAbi, _ := abi.JSON(strings.NewReader(string(erc20.Erc20MetaData.ABI)))
							byteData, _ := hex.DecodeString(txLogs.Data[2:])
							results, _ := contractAbi.Unpack("Transfer", byteData)

							tokenValue := big.NewInt(0)
							if intValue, ok := results[0].(*big.Int); ok {
								tokenValue = intValue
							}

							name := ""
							symbol := ""

							erc20Instance, err := erc20.NewErc20(common.HexToAddress(txLogs.Address), client)
							if err != nil {
								l.Error("create new instance", logger.Field{Key: "error", Value: err.Error()})
								goroutineErr <- err
								return
							}

							tokenName, err := erc20Instance.Name(&bind.CallOpts{})
							if err == nil {
								name = tokenName
							}
							tokenSymbol, err := erc20Instance.Symbol(&bind.CallOpts{})
							if err == nil {
								symbol = tokenSymbol
							}

							if len(txLogs.Topics) == 3 && len(results) == 1 {
								erc20Input := &evmType.Erc20Log{
									ChainID:         result.ChainID,
									Timestamp:       result.Timestamp,
									TimestampInt:    result.TimestampInt,
									CreatedAt:       result.CreatedAt,
									TransactionHash: strings.ToLower(txReceipt.TransactionHash),
									ContractAddress: strings.ToLower(txLogs.Address),
									From:            strings.ToLower(from.String()),
									To:              strings.ToLower(to.String()),
									Amount:          tokenValue,
									Function:        function,
									Name:            name,
									Symbol:          symbol,
								}

								mu.Lock()
								result.Transaction[idx].Erc20Logs = append(result.Transaction[idx].Erc20Logs, erc20Input)
								result.Transaction[idx].Erc20Count = new(big.Int).Add(result.Transaction[idx].Erc20Count, big.NewInt(1))
								mu.Unlock()
							}
						} else if txLogs.Data == "0x" && len(txLogs.Topics) == 4 {
							// erc721
							function := "transfer"
							from := common.HexToAddress(txLogs.Topics[1].Hex())
							to := common.HexToAddress(txLogs.Topics[2].Hex())

							if from == common.HexToAddress("0x0000000000000000000000000000000000000000") {
								function = "mint"
							} else if to == common.HexToAddress("0x0000000000000000000000000000000000000000") || to == common.HexToAddress("0x000000000000000000000000000000000000dead") {
								function = "burn"
							}

							tokenID := new(big.Int)
							tokenIdByte, err := hex.DecodeString(txLogs.Topics[3].Hex()[2:])
							if err == nil {
								tokenID.SetBytes(tokenIdByte)
							}

							name := ""
							symbol := ""

							erc721Instance, err := erc721.NewErc721(common.HexToAddress(txLogs.Address), client)
							if err != nil {
								l.Error("create new instance", logger.Field{Key: "error", Value: err.Error()})
								goroutineErr <- err
								return
							}

							tokenName, err := erc721Instance.Name(&bind.CallOpts{})
							if err == nil {
								name = tokenName
							}
							tokenSymbol, err := erc721Instance.Symbol(&bind.CallOpts{})
							if err == nil {
								symbol = tokenSymbol
							}

							erc721Input := &evmType.Erc721Log{
								ChainID:         result.ChainID,
								Timestamp:       result.Timestamp,
								TimestampInt:    result.TimestampInt,
								CreatedAt:       result.CreatedAt,
								TransactionHash: strings.ToLower(txReceipt.TransactionHash),
								ContractAddress: strings.ToLower(txLogs.Address),
								From:            strings.ToLower(from.String()),
								To:              strings.ToLower(to.String()),
								TokenId:         tokenID,
								Function:        function,
								Name:            name,
								Symbol:          symbol,
							}

							mu.Lock()
							result.Transaction[idx].Erc721Logs = append(result.Transaction[idx].Erc721Logs, erc721Input)
							result.Transaction[idx].Erc721Count = new(big.Int).Add(result.Transaction[idx].Erc721Count, big.NewInt(1))
							mu.Unlock()
						}
					case erc1155TransferSingleSigHash.Hex():
						contractAbi, _ := abi.JSON(strings.NewReader(string(erc1155.Erc1155MetaData.ABI)))
						byteData, _ := hex.DecodeString(txLogs.Data[2:])
						results, _ := contractAbi.Unpack("TransferSingle", byteData)

						if txLogs.Data != "0x" && len(results) == 2 {
							function := "transfer"

							from := common.HexToAddress(txLogs.Topics[2].Hex())
							to := common.HexToAddress(txLogs.Topics[3].Hex())

							if from == common.HexToAddress("0x0000000000000000000000000000000000000000") {
								function = "mint"
							} else if to == common.HexToAddress("0x0000000000000000000000000000000000000000") || to == common.HexToAddress("0x000000000000000000000000000000000000dead") {
								function = "burn"
							}

							var tokenID *big.Int
							var tokenValue *big.Int

							if intValue, ok := results[0].(*big.Int); ok {
								tokenID = intValue
							}
							if intValue, ok := results[1].(*big.Int); ok {
								tokenValue = intValue
							}

							name := ""
							symbol := ""

							erc1155Instance, err := erc1155.NewErc1155(common.HexToAddress(txLogs.Address), client)
							if err != nil {
								l.Error("create new instance", logger.Field{Key: "error", Value: err.Error()})
								goroutineErr <- err
								return
							}

							tokenName, err := erc1155Instance.Name(&bind.CallOpts{})
							if err == nil {
								name = tokenName
							}
							tokenSymbol, err := erc1155Instance.Symbol(&bind.CallOpts{})
							if err == nil {
								symbol = tokenSymbol
							}

							erc1155Input := &evmType.Erc1155Log{
								ChainID:         result.ChainID,
								Timestamp:       result.Timestamp,
								TimestampInt:    result.TimestampInt,
								CreatedAt:       result.CreatedAt,
								TransactionHash: strings.ToLower(txReceipt.TransactionHash),
								ContractAddress: strings.ToLower(txLogs.Address),
								From:            strings.ToLower(from.String()),
								To:              strings.ToLower(to.String()),
								TokenId:         tokenID,
								Amount:          tokenValue,
								Function:        function,
								Name:            name,
								Symbol:          symbol,
							}

							mu.Lock()
							result.Transaction[idx].Erc1155Logs = append(result.Transaction[idx].Erc1155Logs, erc1155Input)
							result.Transaction[idx].Erc1155Count = new(big.Int).Add(result.Transaction[idx].Erc1155Count, big.NewInt(1))
							mu.Unlock()
						}
					case erc1155TransferBatchSigHash.Hex():
						contractAbi, _ := abi.JSON(strings.NewReader(string(erc1155.Erc1155MetaData.ABI)))
						byteData, _ := hex.DecodeString(txLogs.Data[2:])
						results, _ := contractAbi.Unpack("TransferBatch", byteData)

						if txLogs.Data != "0x" && len(results) == 2 {
							function := "transfer"

							from := common.HexToAddress(txLogs.Topics[2].Hex())
							to := common.HexToAddress(txLogs.Topics[3].Hex())

							if from == common.HexToAddress("0x0000000000000000000000000000000000000000") {
								function = "mint"
							} else if to == common.HexToAddress("0x0000000000000000000000000000000000000000") || to == common.HexToAddress("0x000000000000000000000000000000000000dead") {
								function = "burn"
							}

							name := ""
							symbol := ""

							erc1155Instance, err := erc1155.NewErc1155(common.HexToAddress(txLogs.Address), client)
							if err != nil {
								l.Error("create new instance", logger.Field{Key: "error", Value: err.Error()})
								goroutineErr <- err
								return
							}

							tokenName, err := erc1155Instance.Name(&bind.CallOpts{})
							if err == nil {
								name = tokenName
							}
							tokenSymbol, err := erc1155Instance.Symbol(&bind.CallOpts{})
							if err == nil {
								symbol = tokenSymbol
							}

							tokenIDs := results[0].([]*big.Int)
							tokenValues := results[1].([]*big.Int)

							for i, tokenID := range tokenIDs {
								erc1155Input := &evmType.Erc1155Log{
									ChainID:         result.ChainID,
									Timestamp:       result.Timestamp,
									TimestampInt:    result.TimestampInt,
									CreatedAt:       result.CreatedAt,
									TransactionHash: strings.ToLower(txReceipt.TransactionHash),
									ContractAddress: strings.ToLower(txLogs.Address),
									From:            strings.ToLower(from.String()),
									To:              strings.ToLower(to.String()),
									TokenId:         tokenID,
									Amount:          tokenValues[i],
									Function:        function,
									Name:            name,
									Symbol:          symbol,
								}

								mu.Lock()
								result.Transaction[idx].Erc1155Logs = append(result.Transaction[idx].Erc1155Logs, erc1155Input)
								result.Transaction[idx].Erc1155Count = new(big.Int).Add(result.Transaction[idx].Erc1155Count, big.NewInt(1))
								mu.Unlock()
							}
						}
					default:
						// TODO: custom 추적 이벤트
					}
				}

				mu.Lock()
				result.Transaction[idx].Logs[index] = &evmType.Log{
					ChainID:          result.ChainID,
					Address:          strings.ToLower(txLogs.Address),
					BlockHash:        strings.ToLower(result.Hash),
					BlockNumber:      result.Number,
					BlockNumberInt:   result.NumberInt,
					Data:             txLogs.Data,
					LogIndex:         txLogs.LogIndex,
					Removed:          txLogs.Removed,
					Topics:           txLogs.Topics,
					TransactionHash:  strings.ToLower(txLogs.TransactionHash),
					TransactionIndex: txLogs.TransactionIndex,
					From:             txLogs.From,
					To:               txLogs.To,
					Timestamp:        result.Timestamp,
					TimestampInt:     result.TimestampInt,
					CreatedAt:        result.CreatedAt,
				}
				mu.Unlock()
			}
		}(&res.Result, res.ID)
	}

	wg.Wait()

	close(goroutineErr)

	for len(goroutineErr) > 0 {
		for errResult := range goroutineErr {
			if errResult != nil {
				l.Error("fetch transaction data", logger.Field{Key: "error", Value: errResult})
				return errResult
			}
		}
	}

	return nil
}
