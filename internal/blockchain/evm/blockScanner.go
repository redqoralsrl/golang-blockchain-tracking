package evm

import (
	"blockchain-tracking/internal/blockchain/evm/evmType"
	"blockchain-tracking/internal/blockchain/jsonRpc"
	"blockchain-tracking/internal/core/domain/blockchain"
	"blockchain-tracking/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func BlockScanner(ctx context.Context, name, rpc string, jrAdapter *jsonRpc.JsonRpc, blockchainService *blockchain.Service, l logger.Logger) error {
	// client 연결 확인 및 체인아이디 가져오기
	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		l.Error("eth client dial context", logger.Field{Key: "error", Value: err.Error()})
		return err
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		l.Error(fmt.Sprintf("%s get chain id", name), logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	// 최신 블록 높이 요청
	ethBlockNumberRes, err := jrAdapter.CreateRequest(rpc, "eth_blockNumber", []interface{}{})
	if err != nil {
		l.Error("eth_blockNumber", logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	// JSON -> struct 변환
	var latestBlockHeight jsonRpc.EthBlockNumberResponse
	err = json.Unmarshal(ethBlockNumberRes, &latestBlockHeight)
	if err != nil {
		l.Error("json unmarshal block height", logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	lastScanedBlockHeight, err := blockchainService.GetBlockHeight(ctx, chainID.String())
	if err != nil {
		l.Error("get block height", logger.Field{Key: "error", Value: err.Error()})
		return err
	}

	start := new(big.Int)
	start.SetString(lastScanedBlockHeight, 10)
	start.Add(start, big.NewInt(1))

	end := HexToBigInt(latestBlockHeight.Result[2:])

	// TODO: 트랜잭션 / 블록 상황에 따라 유동적으로 조절
	BatchSize := big.NewInt(10)

	for batchStart := new(big.Int).Set(start); batchStart.Cmp(end) <= 0; batchStart.Add(batchStart, BatchSize) {
		// TODO: 트랜잭션 / 블록 상황에 따라 유동적으로 조절
		batchEnd := new(big.Int).Add(batchStart, big.NewInt(9))
		if batchEnd.Cmp(end) > 0 {
			batchEnd.Set(end)
		}

		batchCount := new(big.Int).Add(new(big.Int).Sub(batchEnd, batchStart), big.NewInt(1))

		var wg sync.WaitGroup
		errCh := make(chan error, batchCount.Int64())

		// 성공한 블록 추적 값 데이터 넣는 곳
		idx := -1
		blockList := make([]*evmType.Block, batchCount.Int64())
		var mu sync.Mutex

		for block := new(big.Int).Set(batchStart); block.Cmp(batchEnd) <= 0; block.Add(block, big.NewInt(1)) {
			wg.Add(1)

			idx++

			go func(bn *big.Int, list []*evmType.Block, i int) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						errMsg := fmt.Sprintf("panic: %v", r)
						l.Warn(fmt.Sprintf("block scanner recover %s", name), logger.Field{Key: "error", Value: errMsg})
					}
				}()
				// 복사해서 넘기기
				blockHex := "0x" + bn.Text(16)

				blockResult, err := FetchBlockData(rpc, blockHex, jrAdapter, l)
				if err != nil {
					errCh <- err
					return
				}

				blockTimestampInt := HexToBigInt(blockResult.Timestamp[2:])
				blockTime := time.Unix(blockTimestampInt.Int64(), 0).UTC()
				blockNumberInt := HexToBigInt(blockResult.Number[2:])

				input := &evmType.Block{
					ChainID:          chainID,
					Difficulty:       blockResult.Difficulty,
					Hash:             strings.ToLower(blockResult.Hash),
					GasLimit:         blockResult.GasLimit,
					GasUsed:          blockResult.GasUsed,
					Miner:            strings.ToLower(blockResult.Miner),
					Number:           blockResult.Number,
					NumberInt:        blockNumberInt,
					ParentHash:       strings.ToLower(blockResult.ParentHash),
					Timestamp:        blockResult.Timestamp,
					TimestampInt:     blockTimestampInt,
					CreatedAt:        blockTime,
					TotalDifficulty:  blockResult.TotalDifficulty,
					TransactionsRoot: strings.ToLower(blockResult.TransactionsRoot),
					Transaction:      make([]evmType.Transaction, len(blockResult.Transactions)),
				}

				if len(blockResult.Transactions) > 0 {
					err = FetchTransactionData(ctx, rpc, blockResult, input, jrAdapter, l)
					if err != nil {
						errCh <- err
						return
					}
				}

				mu.Lock()
				list[i] = input
				mu.Unlock()
			}(new(big.Int).Set(block), blockList, idx)
		}

		wg.Wait()
		close(errCh)

		if len(errCh) > 0 {
			for errResult := range errCh {
				if errResult != nil {
					l.Error("blockchain decode", logger.Field{Key: "error", Value: errResult})
					return errResult
				}
			}
		}

		for _, data := range blockList {
			err = blockchainService.Create(ctx, data)
			if err != nil {
				l.Error("blockchain service create", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "chain_id", Value: data.ChainID}, logger.Field{Key: "block number", Value: data.NumberInt})
				return err
			}
		}
	}

	return nil
}
