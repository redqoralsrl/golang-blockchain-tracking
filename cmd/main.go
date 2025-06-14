package main

import (
	"blockchain-tracking/config"
	"blockchain-tracking/internal/blockchain/evm"
	"blockchain-tracking/internal/blockchain/jsonRpc"
	"blockchain-tracking/internal/core/domain/blockchain"
	"blockchain-tracking/internal/database/postgresql"
	"blockchain-tracking/internal/logger"
	"blockchain-tracking/internal/logger/zerolog"
	"blockchain-tracking/internal/smtp"
	"context"
	"fmt"
	"sync"
)

func main() {
	config := config.LoadConfig()

	stage := config.STAGE
	if stage == "" {
		stage = "dev"
	}

	l := zerolog.NewLogger(stage)
	defer func() {
		if err := l.Close(); err != nil {
			l.Warn("failed to connect to logger", logger.Field{Key: "error", Value: err.Error()})
		}
	}()

	db, err := postgresql.NewDB(config, l)
	if err != nil {
		l.Fatal("failed to connect to postgresql", logger.Field{Key: "error", Value: err.Error()})
	}
	defer func() {
		if err = db.Close(); err != nil {
			l.Fatal("failed to close db", logger.Field{Key: "error", Value: err.Error()})
		}
	}()

	transactionManager := postgresql.NewManager(db)

	blockchainService := blockchain.NewService(db, transactionManager, l)

	jsonRpcAdapter := jsonRpc.NewJsonRpc(l)
	smtpAdapter := smtp.NewSmtp(config.SMTPID, config.SMTPPassword)

	chainList := []string{"Ethereum"}
	// chainList := []string{"Ethereum", "Biance", "GiantMammoth"}

	var wg sync.WaitGroup

	for _, chain := range chainList {
		wg.Add(1)
		go func(chainName, rpc string) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					l.Error(fmt.Sprintf("go routine recover %s chain", chain), logger.Field{Key: "error", Value: r.(error).Error()})
				}
			}()

			ctx := context.Background()

			l.Info(fmt.Sprintf("go routine start %s chain", chainName), logger.Field{Key: "chain", Value: chainName})

			err := evm.StartTrack(ctx, chainName, rpc, jsonRpcAdapter, smtpAdapter, blockchainService, l)
			if err != nil {
				l.Error(fmt.Sprintf("%s chain error tracking", chainName), logger.Field{Key: "error", Value: err.Error()})
			}
		}(chain, config.RPC[chain])
	}

	l.Info("blockchain tracking starting...")

	wg.Wait()

	l.Error("server is dead")
}
