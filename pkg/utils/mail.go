package utils

import (
	"time"

	"go.uber.org/zap"
)

func SendEmail(logger *zap.Logger) {
	time.Sleep(2 * time.Second)
	logger.Info("send email success")
}
