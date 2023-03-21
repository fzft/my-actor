package internel

import (
	"go.uber.org/zap"
	"testing"
)

func TestDefaultMailbox(t *testing.T) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()
	sugarLogger := logger.Sugar()
	// Create a new mailbox
	mailbox := NewDefaultMailbox(sugarLogger)

	// Add some messages to the mailbox
	for i := 1; i <= 10; i++ {
		mailbox.Source(i)
	}

	//// Consume messages from the mailbox
	//consumed := make([]any, 0)
	//for msg := range mailbox.Consume() {
	//	consumed = append(consumed, msg)
	//	if len(consumed) == 10 {
	//		break
	//	}
	//}
}
