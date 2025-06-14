package blockchain

import (
	"blockchain-tracking/internal/blockchain/evm/evmType"
	"blockchain-tracking/internal/database/gen"
	"blockchain-tracking/internal/database/postgresql"
	"blockchain-tracking/internal/logger"
	"context"
	"database/sql"
	"errors"
)

type Service struct {
	db        *postgresql.Database
	txManager postgresql.DBTransactionManager
	l         logger.Logger
}

func NewService(d *postgresql.Database, tx postgresql.DBTransactionManager, l logger.Logger) *Service {
	return &Service{
		db:        d,
		txManager: tx,
		l:         l,
	}
}

func (s *Service) GetBlockHeight(ctx context.Context, chainId string) (string, error) {
	height, err := s.db.Queries.GetBlockHeight(ctx, chainId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "-1", nil
		}
		return "", err
	}

	return height, nil
}

func (s *Service) Create(ctx context.Context, block *evmType.Block) error {
	err := s.txManager.WithTransaction(ctx, sql.LevelRepeatableRead, false, func(ctx context.Context) error {

		err := s.db.Queries.InsertBlock(ctx, gen.InsertBlockParams{
			ChainID:          block.ChainID.String(),
			Difficulty:       sql.NullString{String: block.Difficulty, Valid: block.Difficulty != ""},
			Hash:             block.Hash,
			GasLimit:         sql.NullString{String: block.GasLimit, Valid: block.GasLimit != ""},
			GasUsed:          sql.NullString{String: block.GasUsed, Valid: block.GasUsed != ""},
			Miner:            sql.NullString{String: block.Miner, Valid: block.Miner != ""},
			Number:           block.Number,
			NumberInt:        block.NumberInt.String(),
			ParentHash:       sql.NullString{String: block.ParentHash, Valid: block.ParentHash != ""},
			Timestamp:        block.Timestamp,
			TimestampInt:     block.TimestampInt.String(),
			CreatedAt:        block.CreatedAt,
			TotalDifficulty:  sql.NullString{String: block.TotalDifficulty, Valid: block.TotalDifficulty != ""},
			TransactionsRoot: sql.NullString{String: block.TransactionsRoot, Valid: block.TransactionsRoot != ""},
		})
		if err != nil {
			s.l.Error("create block", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "block", Value: block.Hash})
			return err
		}

		if block.Transaction != nil && len(block.Transaction) > 0 {
			for _, tx := range block.Transaction {
				txInput := gen.InsertTransactionParams{
					ChainID:          tx.ChainID.String(),
					BlockHash:        tx.BlockHash,
					BlockNumber:      tx.BlockNumber,
					BlockNumberInt:   tx.BlockNumberInt.String(),
					From:             sql.NullString{String: tx.From, Valid: tx.From != ""},
					To:               sql.NullString{String: tx.To, Valid: tx.To != ""},
					Gas:              sql.NullString{String: tx.Gas, Valid: tx.Gas != ""},
					GasInt:           sql.NullString{String: tx.GasInt.String(), Valid: tx.GasInt.String() != ""},
					GasPrice:         sql.NullString{String: tx.GasPrice, Valid: tx.GasPrice != ""},
					GasPriceInt:      sql.NullString{String: tx.GasPriceInt.String(), Valid: tx.GasPriceInt.String() != ""},
					Hash:             tx.Hash,
					R:                sql.NullString{String: tx.R, Valid: tx.R != ""},
					S:                sql.NullString{String: tx.S, Valid: tx.S != ""},
					V:                sql.NullString{String: tx.V, Valid: tx.V != ""},
					TransactionIndex: sql.NullString{String: tx.TransactionIndex, Valid: tx.TransactionIndex != ""},
					Value:            sql.NullString{String: tx.Value, Valid: tx.Value != ""},
					ValueInt:         sql.NullString{String: tx.ValueInt.String(), Valid: tx.ValueInt != nil},
					Nonce:            sql.NullString{String: tx.Nonce, Valid: tx.Nonce != ""},
					Input:            sql.NullString{String: tx.Input, Valid: tx.Input != ""},
					GasUsed:          sql.NullString{String: tx.GasUsed, Valid: tx.GasUsed != ""},
					GasUsedInt:       sql.NullString{String: tx.GasUsedInt.String(), Valid: tx.GasUsedInt != nil},
					Status:           sql.NullString{String: tx.Status, Valid: tx.Status != ""},
					Type:             sql.NullString{String: tx.Type, Valid: tx.Type != ""},
					Timestamp:        tx.Timestamp,
					TimestampInt:     tx.TimestampInt.String(),
					CreatedAt:        tx.CreatedAt,
					CoinCount:        tx.CoinCount.String(),
					NftCount:         tx.NftCount.String(),
					Erc20Count:       tx.Erc20Count.String(),
					Erc721Count:      tx.Erc721Count.String(),
					Erc1155Count:     tx.Erc1155Count.String(),
				}
				if tx.ContractAddress != nil {
					txInput.ContractAddress = sql.NullString{String: *tx.ContractAddress, Valid: *tx.ContractAddress != ""}
				} else {
					txInput.ContractAddress = sql.NullString{Valid: false}
				}
				err = s.db.Queries.InsertTransaction(ctx, txInput)
				if err != nil {
					s.l.Error("create transaction", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "transaction", Value: tx.Hash})
					return err
				}

				if tx.Contract != nil {
					err = s.db.Queries.InsertContract(ctx, gen.InsertContractParams{
						ChainID:     tx.Contract.ChainID.String(),
						Hash:        tx.Contract.Hash,
						Name:        sql.NullString{String: tx.Contract.Name, Valid: tx.Contract.Name != ""},
						Symbol:      sql.NullString{String: tx.Contract.Symbol, Valid: tx.Contract.Symbol != ""},
						Decimals:    sql.NullInt32{Int32: int32(tx.Contract.Decimals), Valid: true},
						TotalSupply: sql.NullString{String: tx.Contract.TotalSupply.String(), Valid: tx.Contract.TotalSupply.String() != ""},
						Type:        sql.NullString{String: tx.Contract.Type, Valid: tx.Contract.Type != ""},
						Creator:     sql.NullString{String: tx.Contract.Creator, Valid: tx.Contract.Creator != ""},
					})
					if err != nil {
						s.l.Error("create contract", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "contract", Value: tx.Contract.Hash})
						return err
					}
				}
				if tx.CoinLogs != nil {
					err = s.db.Queries.InsertCoinLog(ctx, gen.InsertCoinLogParams{
						ChainID:         tx.CoinLogs.ChainID.String(),
						Timestamp:       tx.CoinLogs.Timestamp,
						TimestampInt:    tx.CoinLogs.TimestampInt.String(),
						CreatedAt:       tx.CoinLogs.CreatedAt,
						TransactionHash: sql.NullString{String: tx.CoinLogs.TransactionHash, Valid: tx.CoinLogs.TransactionHash != ""},
						From:            sql.NullString{String: tx.CoinLogs.From, Valid: tx.CoinLogs.From != ""},
						To:              sql.NullString{String: tx.CoinLogs.To, Valid: tx.CoinLogs.To != ""},
						Amount:          sql.NullString{String: tx.CoinLogs.Amount.String(), Valid: tx.CoinLogs.Amount.String() != ""},
						Gas:             sql.NullString{String: tx.CoinLogs.Gas, Valid: tx.CoinLogs.Gas != ""},
						GasInt:          sql.NullString{String: tx.CoinLogs.GasInt.String(), Valid: tx.CoinLogs.GasInt.String() != ""},
						GasPrice:        sql.NullString{String: tx.CoinLogs.GasPrice, Valid: tx.CoinLogs.GasPrice != ""},
						GasPriceInt:     sql.NullString{String: tx.CoinLogs.GasPriceInt.String(), Valid: tx.CoinLogs.GasPriceInt.String() != ""},
						GasUsed:         sql.NullString{String: tx.CoinLogs.GasUsed, Valid: tx.CoinLogs.GasUsed != ""},
						GasUsedInt:      sql.NullString{String: tx.CoinLogs.GasUsedInt.String(), Valid: tx.CoinLogs.GasUsedInt.String() != ""},
					})
					if err != nil {
						s.l.Error("create coin logs", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "coinLog", Value: tx.CoinLogs.TransactionHash})
						return err
					}

					err = s.db.Queries.UpsertWalletBalance(ctx, gen.UpsertWalletBalanceParams{
						ChainID: tx.CoinLogs.ChainID.String(),
						Address: tx.CoinLogs.From,
						Balance: "-" + tx.CoinLogs.Amount.String(),
					})
					if err != nil {
						s.l.Error("create wallet", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "wallet", Value: tx.CoinLogs.From})
						return err
					}
					err = s.db.Queries.UpsertWalletBalance(ctx, gen.UpsertWalletBalanceParams{
						ChainID: tx.CoinLogs.ChainID.String(),
						Address: tx.CoinLogs.To,
						Balance: tx.CoinLogs.Amount.String(),
					})
					if err != nil {
						s.l.Error("create wallet", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "wallet", Value: tx.CoinLogs.To})
						return err
					}
				}
				if tx.Erc20Logs != nil && len(tx.Erc20Logs) > 0 {
					for _, erc20Data := range tx.Erc20Logs {
						err = s.db.Queries.InsertERC20Log(ctx, gen.InsertERC20LogParams{
							ChainID:         erc20Data.ChainID.String(),
							Timestamp:       erc20Data.Timestamp,
							TimestampInt:    erc20Data.TimestampInt.String(),
							CreatedAt:       erc20Data.CreatedAt,
							TransactionHash: erc20Data.TransactionHash,
							ContractAddress: erc20Data.ContractAddress,
							From:            erc20Data.From,
							To:              erc20Data.To,
							Amount:          erc20Data.Amount.String(),
							Function:        erc20Data.Function,
							Name:            sql.NullString{String: erc20Data.Name, Valid: erc20Data.Name != ""},
							Symbol:          sql.NullString{String: erc20Data.Symbol, Valid: erc20Data.Symbol != ""},
						})
						if err != nil {
							s.l.Error("create erc20 log", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc20Log", Value: erc20Data.TransactionHash})
							return err
						}

						err = s.db.Queries.InsertWallet(ctx, gen.InsertWalletParams{
							ChainID: erc20Data.ChainID.String(),
							Address: erc20Data.From,
						})
						if err != nil {
							s.l.Error("create erc20 log wallet", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc20LogWallet", Value: erc20Data.From})
							return err
						}
						err = s.db.Queries.InsertWallet(ctx, gen.InsertWalletParams{
							ChainID: erc20Data.ChainID.String(),
							Address: erc20Data.To,
						})
						if err != nil {
							s.l.Error("create erc20 log wallet", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc20LogWallet", Value: erc20Data.To})
							return err
						}

						if erc20Data.Function != "mint" {
							err = s.db.Queries.UpsertERC20Balance(ctx, gen.UpsertERC20BalanceParams{
								ChainID: erc20Data.ChainID.String(),
								Balance: "-" + erc20Data.Amount.String(),
								Hash:    erc20Data.ContractAddress,
								Address: erc20Data.From,
							})
							if err != nil {
								s.l.Error("create erc20 balance", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc20Balance", Value: erc20Data.From})
								return err
							}
						}
						err = s.db.Queries.UpsertERC20Balance(ctx, gen.UpsertERC20BalanceParams{
							ChainID: erc20Data.ChainID.String(),
							Balance: erc20Data.Amount.String(),
							Hash:    erc20Data.ContractAddress,
							Address: erc20Data.To,
						})
						if err != nil {
							s.l.Error("create erc20 balance", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc20Balance", Value: erc20Data.To})
							return err
						}
					}
				}
				if tx.Erc721Logs != nil && len(tx.Erc721Logs) > 0 {
					for _, erc721Data := range tx.Erc721Logs {
						err = s.db.Queries.InsertERC721Log(ctx, gen.InsertERC721LogParams{
							ChainID:         erc721Data.ChainID.String(),
							Timestamp:       erc721Data.Timestamp,
							TimestampInt:    erc721Data.TimestampInt.String(),
							CreatedAt:       erc721Data.CreatedAt,
							TransactionHash: erc721Data.TransactionHash,
							ContractAddress: erc721Data.ContractAddress,
							From:            erc721Data.From,
							To:              erc721Data.To,
							TokenID:         erc721Data.TokenId.String(),
							Function:        sql.NullString{String: erc721Data.Function, Valid: erc721Data.Function != ""},
							Name:            sql.NullString{String: erc721Data.Name, Valid: erc721Data.Name != ""},
							Symbol:          sql.NullString{String: erc721Data.Symbol, Valid: erc721Data.Symbol != ""},
						})
						if err != nil {
							s.l.Error("create erc721 log", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc721Log", Value: erc721Data.TransactionHash})
							return err
						}

						err = s.db.Queries.UpsertERC721Balance(ctx, gen.UpsertERC721BalanceParams{
							ChainID: erc721Data.ChainID.String(),
							Hash:    erc721Data.ContractAddress,
							TokenID: erc721Data.TokenId.String(),
							Address: erc721Data.To,
						})
						if err != nil {
							s.l.Error("create erc721 balance", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc721Balance", Value: erc721Data.To})
							return err
						}

						err = s.db.Queries.InsertWallet(ctx, gen.InsertWalletParams{
							ChainID: erc721Data.ChainID.String(),
							Address: erc721Data.From,
						})
						if err != nil {
							s.l.Error("create erc721 log wallet", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc721LogWallet", Value: erc721Data.From})
							return err
						}

						err = s.db.Queries.InsertWallet(ctx, gen.InsertWalletParams{
							ChainID: erc721Data.ChainID.String(),
							Address: erc721Data.To,
						})
						if err != nil {
							s.l.Error("create erc721 log wallet", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc721LogWallet", Value: erc721Data.To})
							return err
						}

						err = s.db.Queries.CreateErc721(ctx, gen.CreateErc721Params{
							ChainID: erc721Data.ChainID.String(),
							Hash:    erc721Data.ContractAddress,
							TokenID: erc721Data.TokenId.String(),
						})

						err = s.db.Queries.UpdateContractType(ctx, gen.UpdateContractTypeParams{
							Type:    sql.NullString{String: "721", Valid: true},
							ChainID: erc721Data.ChainID.String(),
							Hash:    erc721Data.ContractAddress,
						})
					}
				}
				if tx.Erc1155Logs != nil && len(tx.Erc1155Logs) > 0 {
					for _, erc1155Data := range tx.Erc1155Logs {
						err = s.db.Queries.InsertERC1155Log(ctx, gen.InsertERC1155LogParams{
							ChainID:         erc1155Data.ChainID.String(),
							Timestamp:       erc1155Data.Timestamp,
							TimestampInt:    erc1155Data.TimestampInt.String(),
							CreatedAt:       erc1155Data.CreatedAt,
							TransactionHash: erc1155Data.TransactionHash,
							ContractAddress: sql.NullString{String: erc1155Data.ContractAddress, Valid: erc1155Data.ContractAddress != ""},
							From:            sql.NullString{String: erc1155Data.From, Valid: erc1155Data.From != ""},
							To:              sql.NullString{String: erc1155Data.To, Valid: erc1155Data.To != ""},
							TokenID:         sql.NullString{String: erc1155Data.TokenId.String(), Valid: erc1155Data.TokenId.String() != ""},
							Amount:          sql.NullString{String: erc1155Data.Amount.String(), Valid: erc1155Data.Amount.String() != ""},
							Function:        sql.NullString{String: erc1155Data.Function, Valid: erc1155Data.Function != ""},
							Name:            sql.NullString{String: erc1155Data.Name, Valid: erc1155Data.Name != ""},
							Symbol:          sql.NullString{String: erc1155Data.Symbol, Valid: erc1155Data.Symbol != ""},
						})
						if err != nil {
							s.l.Error("create erc1155 log", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc1155Log", Value: erc1155Data.TransactionHash})
							return err
						}

						err = s.db.Queries.InsertWallet(ctx, gen.InsertWalletParams{
							ChainID: erc1155Data.ChainID.String(),
							Address: erc1155Data.From,
						})
						if err != nil {
							s.l.Error("create erc1155 log wallet", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc1155LogWallet", Value: erc1155Data.From})
							return err
						}
						err = s.db.Queries.InsertWallet(ctx, gen.InsertWalletParams{
							ChainID: erc1155Data.ChainID.String(),
							Address: erc1155Data.To,
						})
						if err != nil {
							s.l.Error("create erc1155 log wallet", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc1155LogWallet", Value: erc1155Data.To})
							return err
						}

						if erc1155Data.Function != "mint" {
							err = s.db.Queries.SubtractERC1155Balance(ctx, gen.SubtractERC1155BalanceParams{
								ChainID: erc1155Data.ChainID.String(),
								Hash:    erc1155Data.ContractAddress,
								TokenID: erc1155Data.TokenId.String(),
								Address: erc1155Data.From,
								Amount:  erc1155Data.Amount.String(),
							})
							if err != nil {
								s.l.Error("create erc1155 balance", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc1155Balance", Value: erc1155Data.From})
								return err
							}
						}

						err = s.db.Queries.UpsertERC1155Balance_Add(ctx, gen.UpsertERC1155Balance_AddParams{
							ChainID: erc1155Data.ChainID.String(),
							Hash:    erc1155Data.ContractAddress,
							TokenID: erc1155Data.TokenId.String(),
							Address: erc1155Data.To,
							Amount:  erc1155Data.Amount.String(),
						})
						if err != nil {
							s.l.Error("create erc1155 balance", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "erc1155Balance", Value: erc1155Data.To})
							return err
						}

						err = s.db.Queries.CreateErc1155(ctx, gen.CreateErc1155Params{
							ChainID: erc1155Data.ChainID.String(),
							Hash:    erc1155Data.ContractAddress,
							TokenID: erc1155Data.TokenId.String(),
						})

						err = s.db.Queries.UpdateContractType(ctx, gen.UpdateContractTypeParams{
							Type:    sql.NullString{String: "1155", Valid: true},
							ChainID: erc1155Data.ChainID.String(),
							Hash:    erc1155Data.ContractAddress,
						})
					}
				}
				if tx.Logs != nil && len(tx.Logs) > 0 {
					for _, txLogData := range tx.Logs {
						topics := make([]string, len(txLogData.Topics))
						for i, topic := range txLogData.Topics {
							topics[i] = topic.Hex()
						}

						err = s.db.Queries.InsertLog(ctx, gen.InsertLogParams{
							ChainID:          txLogData.ChainID.String(),
							Address:          sql.NullString{String: txLogData.Address, Valid: txLogData.Address != ""},
							BlockHash:        txLogData.BlockHash,
							BlockNumber:      txLogData.BlockNumber,
							BlockNumberInt:   sql.NullString{String: txLogData.BlockNumberInt.String(), Valid: txLogData.BlockNumberInt.String() != ""},
							Data:             sql.NullString{String: txLogData.Data, Valid: txLogData.Data != ""},
							LogIndex:         sql.NullString{String: txLogData.LogIndex, Valid: txLogData.LogIndex != ""},
							Removed:          txLogData.Removed,
							Topics:           topics,
							TransactionHash:  txLogData.TransactionHash,
							TransactionIndex: txLogData.TransactionIndex,
							From:             sql.NullString{String: txLogData.From.String(), Valid: txLogData.From.String() != ""},
							To:               sql.NullString{String: txLogData.To.String(), Valid: txLogData.To.String() != ""},
							Timestamp:        txLogData.Timestamp,
							TimestampInt:     txLogData.TimestampInt.String(),
							CreatedAt:        txLogData.CreatedAt,
						})
						if err != nil {
							s.l.Error("create log", logger.Field{Key: "error", Value: err.Error()}, logger.Field{Key: "log", Value: txLogData.TransactionHash})
							return err
						}
					}
				}

			}
		}

		return nil
	})

	return err
}
