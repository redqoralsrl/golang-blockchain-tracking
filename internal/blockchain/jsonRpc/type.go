package jsonRpc

import "github.com/ethereum/go-ethereum/common"

type EthBlockNumberResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  string `json:"result"`
}

type EthGetBlockByNumber struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Block  `json:"result"`
}

type EthGetTransactionReceipt struct {
	Jsonrpc string             `json:"jsonrpc"`
	ID      int                `json:"id"`
	Result  TransactionReceipt `json:"result"`
}

type EthGetDebugTraceTransaction struct {
	Jsonrpc string   `json:"jsonrpc"`
	ID      int      `json:"id"`
	Result  CallInfo `json:"result"`
}

type Block struct {
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"` // DB - number_hex와 number의 값을 변형시켜야함
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	Transactions     []Transaction `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []string      `json:"uncles"`
}

type Transaction struct {
	BlockHash        string  `json:"blockHash"`
	BlockNumber      string  `json:"blockNumber"` // DB - block_numer_hex와 block_number의 값을 변형시켜야함
	From             string  `json:"from"`
	Gas              string  `json:"gas"`
	GasPrice         string  `json:"gasPrice"`
	Hash             string  `json:"hash"`
	Input            string  `json:"input"`
	Nonce            string  `json:"nonce"`
	R                string  `json:"r"`
	S                string  `json:"s"`
	To               *string `json:"to,omitempty"`
	TransactionIndex string  `json:"transactionIndex"`
	Type             string  `json:"type"`
	V                string  `json:"v"`
	Value            string  `json:"value"`
}

type TransactionReceipt struct {
	BlockHash         string  `json:"blockHash"`
	BlockNumber       string  `json:"blockNumber"`
	ContractAddress   *string `json:"contractAddress,omitempty"`
	CumulativeGasUsed string  `json:"cumulativeGasUsed"`
	From              string  `json:"from"`
	GasUsed           string  `json:"gasUsed"`
	Logs              []Log   `json:"logs"`
	LogsBloom         string  `json:"logsBloom"`
	Status            string  `json:"status"`
	To                *string `json:"to,omitempty"`
	TransactionHash   string  `json:"transactionHash"`
	TransactionIndex  string  `json:"transactionIndex"`
	Type              string  `json:"type"`
}

type Log struct {
	Address          string         `json:"address"`
	BlockHash        string         `json:"blockHash"`
	BlockNumber      string         `json:"blockNumber"`
	Data             string         `json:"data"`
	LogIndex         string         `json:"logIndex"`
	Removed          bool           `json:"removed"`
	Topics           []common.Hash  `json:"topics"`
	TransactionHash  string         `json:"transactionHash"`
	TransactionIndex string         `json:"transactionIndex"`
	Timestamp        uint64         `json:"timestamp"`
	From             common.Address `json:"from,omitempty"`
	To               common.Address `json:"to,omitempty"`
}

type CallInfo struct {
	Calls   []CallInfo `json:"calls,omitempty"`
	From    string     `json:"from"`
	Gas     string     `json:"gas"`
	GasUsed string     `json:"gasUsed"`
	Input   string     `json:"input"`
	Output  string     `json:"output"`
	To      string     `json:"to"`
	Type    string     `json:"type"`
	Value   string     `json:"value,omitempty"`
}
