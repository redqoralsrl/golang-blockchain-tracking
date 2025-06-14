package evmType

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Block struct {
	ChainID          *big.Int      `json:"chainID"` // custom
	Difficulty       string        `json:"difficulty"`
	Hash             string        `json:"hash"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Miner            string        `json:"miner"`
	Number           string        `json:"number"`
	NumberInt        *big.Int      `json:"NumberInt"` // custom
	ParentHash       string        `json:"parentHash"`
	Timestamp        string        `json:"timestamp"`
	TimestampInt     *big.Int      `json:"timestampInt"` // custom
	CreatedAt        time.Time     `json:"createdAt"`    // custom
	TotalDifficulty  string        `json:"totalDifficulty"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Transaction      []Transaction `json:"transaction"`
}

type Transaction struct {
	ChainID          *big.Int      `json:"chainID"` // custom
	BlockHash        string        `json:"blockHash"`
	BlockNumber      string        `json:"blockNumber"`
	BlockNumberInt   *big.Int      `json:"blockNumberInt"` // custom
	From             string        `json:"from"`
	To               string        `json:"to,omitempty"`
	Gas              string        `json:"gas"`
	GasInt           *big.Int      `json:"gasInt"` // custom
	GasPrice         string        `json:"gasPrice"`
	GasPriceInt      *big.Int      `json:"gasPriceInt"` // custom
	Hash             string        `json:"hash"`
	R                string        `json:"r"`
	S                string        `json:"s"`
	V                string        `json:"v"`
	TransactionIndex string        `json:"transactionIndex"`
	Value            string        `json:"value"`
	ValueInt         *big.Int      `json:"valueInt"` // custom
	Nonce            string        `json:"nonce"`
	Input            string        `json:"input"`
	ContractAddress  *string       `json:"contractAddress,omitempty"`
	GasUsed          string        `json:"gasUsed"`
	GasUsedInt       *big.Int      `json:"gasUsedInt"` // custom
	Logs             []*Log        `json:"log"`
	Status           string        `json:"string"`
	Type             string        `json:"type"`
	Timestamp        string        `json:"timestamp"`              // custom
	TimestampInt     *big.Int      `json:"timestampInt"`           // custom
	CreatedAt        time.Time     `json:"createdAt"`              // custom
	Contract         *Contract     `json:"contract,omitempty"`     // custom
	CoinLogs         *CoinLog      `json:"coinLogs,omitempty"`     // custom
	Erc20Logs        []*Erc20Log   `json:"erc20Logs,omitempty"`    // custom
	Erc721Logs       []*Erc721Log  `json:"erc721Logs,omitempty"`   // custom
	Erc1155Logs      []*Erc1155Log `json:"erc1155Logs,omitempty"`  // custom
	CoinCount        *big.Int      `json:"coinCount,omitempty"`    // custom
	NftCount         *big.Int      `json:"nftCount,omitempty"`     // custom
	Erc20Count       *big.Int      `json:"erc20Count,omitempty"`   // custom
	Erc721Count      *big.Int      `json:"erc721Count,omitempty"`  // custom
	Erc1155Count     *big.Int      `json:"erc1155Count,omitempty"` // custom
}

type Log struct {
	ChainID          *big.Int       `json:"chainID"` // custom
	Address          string         `json:"address"`
	BlockHash        string         `json:"blockHash"`      // custom
	BlockNumber      string         `json:"blockNumber"`    // custom
	BlockNumberInt   *big.Int       `json:"blockNumberInt"` // custom
	Data             string         `json:"data"`
	LogIndex         string         `json:"logIndex"`
	Removed          bool           `json:"removed"`
	Topics           []common.Hash  `json:"topics"`
	TransactionHash  string         `json:"transactionHash"`
	TransactionIndex string         `json:"transactionIndex"`
	From             common.Address `json:"from,omitempty"`
	To               common.Address `json:"to,omitempty"`
	Timestamp        string         `json:"timestamp"`    // custom
	TimestampInt     *big.Int       `json:"timestampInt"` // custom
	CreatedAt        time.Time      `json:"createdAt"`    // custom
}

type Contract struct {
	ChainID     *big.Int `json:"chainID"`
	Hash        string   `json:"hash"`
	Name        string   `json:"name,omitempty"`
	Symbol      string   `json:"symbol,omitempty"`
	Decimals    int      `json:"decimals"`
	TotalSupply *big.Int `json:"totalSupply,omitempty"`
	Type        string   `json:"type"`
	Creator     string   `json:"creator"`
}

type CoinLog struct {
	ChainID         *big.Int  `json:"chainID"`
	Timestamp       string    `json:"timestamp"`
	TimestampInt    *big.Int  `json:"timestampInt"`
	CreatedAt       time.Time `json:"createdAt"`
	TransactionHash string    `json:"transactionHash"`
	From            string    `json:"from"`
	To              string    `json:"to,omitempty"`
	Amount          *big.Int  `json:"amount"`
	Gas             string    `json:"gas"`
	GasInt          *big.Int  `json:"gasInt"`
	GasPrice        string    `json:"gasPrice"`
	GasPriceInt     *big.Int  `json:"gasPriceInt"`
	GasUsed         string    `json:"gasUsed"`
	GasUsedInt      *big.Int  `json:"gasUsedInt"`
}

type Erc20Log struct {
	ChainID         *big.Int  `json:"chainID"`
	Timestamp       string    `json:"timestamp"`
	TimestampInt    *big.Int  `json:"timestampInt"`
	CreatedAt       time.Time `json:"createdAt"`
	TransactionHash string    `json:"transactionHash"`
	ContractAddress string    `json:"contractAddress"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	Amount          *big.Int  `json:"amount"`
	Function        string    `json:"function"`
	Name            string    `json:"name,omitempty"`
	Symbol          string    `json:"symbol,omitempty"`
}

type Erc721Log struct {
	ChainID         *big.Int  `json:"chainID"`
	Timestamp       string    `json:"timestamp"`
	TimestampInt    *big.Int  `json:"timestampInt"`
	CreatedAt       time.Time `json:"createdAt"`
	TransactionHash string    `json:"transactionHash"`
	ContractAddress string    `json:"contractAddress"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	TokenId         *big.Int  `json:"tokenID"`
	Function        string    `json:"function"`
	Name            string    `json:"name,omitempty"`
	Symbol          string    `json:"symbol,omitempty"`
}

type Erc1155Log struct {
	ChainID         *big.Int  `json:"chainID"`
	Timestamp       string    `json:"timestamp"`
	TimestampInt    *big.Int  `json:"timestampInt"`
	CreatedAt       time.Time `json:"createdAt"`
	TransactionHash string    `json:"transactionHash"`
	ContractAddress string    `json:"contractAddress"`
	From            string    `json:"from"`
	To              string    `json:"to"`
	TokenId         *big.Int  `json:"tokenID"`
	Amount          *big.Int  `json:"amount"`
	Function        string    `json:"function"`
	Name            string    `json:"name,omitempty"`
	Symbol          string    `json:"symbol,omitempty"`
}
