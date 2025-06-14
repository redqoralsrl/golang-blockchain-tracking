package evm

import (
	"blockchain-tracking/internal/blockchain/jsonRpc"
	"blockchain-tracking/internal/core/domain/blockchain"
	"blockchain-tracking/internal/logger"
	"blockchain-tracking/internal/smtp"
	"context"
	"fmt"
	"time"
)

func StartTrack(ctx context.Context, name, rpc string, jrAdapter *jsonRpc.JsonRpc, smtpAdapter *smtp.Smtp, blockchainService *blockchain.Service, l logger.Logger) error {
	defer func() {
		if r := recover(); r != nil {
			l.Warn(fmt.Sprintf("panic occurred in StartTracking %s", name), logger.Field{
				Key:   "recovered",
				Value: r,
			})
		}
	}()

	// 추적 주기
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// tracking 중복 방지 채널
	isTracking := make(chan struct{}, 1)

	failCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			select {
			case isTracking <- struct{}{}:

				err := BlockScanner(ctx, name, rpc, jrAdapter, blockchainService, l)
				if err != nil {
					failCount++
					l.Error(fmt.Sprintf("tracking error on %s (failCount: %d)", name, failCount), logger.Field{
						Key:   "error",
						Value: err.Error(),
					})

					if failCount >= 5 {
						// TODO:
						subject := fmt.Sprintf("[%s] %s chain tracking failed 5 times!", time.Now().Format(time.RFC3339), name)
						smtpAdapter.SendEmail("test@gmail.com", subject, "check please!")
						return err
					}
				} else {
					failCount = 0 // 성공 시 카운터 초기화
				}
				<-isTracking
			default:
				l.Info(fmt.Sprintf("%s is still tracking, skipping", name))
			}
		}
	}
}
