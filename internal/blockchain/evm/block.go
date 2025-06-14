package evm

import (
	"blockchain-tracking/internal/blockchain/jsonRpc"
	"blockchain-tracking/internal/logger"
	"encoding/json"
	"fmt"
)

func FetchBlockData(rpc string, blockHeight string, jrAdapter *jsonRpc.JsonRpc, l logger.Logger) (*jsonRpc.Block, error) {
	res, err := jrAdapter.CreateRequest(rpc, "eth_getBlockByNumber", []interface{}{blockHeight, true})
	if err != nil {
		l.Error("eth_getBlockByNumber", logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	var blockData jsonRpc.EthGetBlockByNumber
	err = json.Unmarshal(res, &blockData)
	if err != nil {
		l.Error("json unmarshal block data", logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	if len(blockData.Result.Timestamp) <= 2 {
		err := fmt.Errorf("invalid timestamp in block %s", blockHeight)
		l.Error("block timestamp parse error", logger.Field{Key: "error", Value: err.Error()})
		return nil, err
	}

	return &blockData.Result, nil
}
