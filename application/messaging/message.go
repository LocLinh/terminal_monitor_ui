package messaging

import (
	"context"
	"fmt"
	no "terminal_monitor_ui/application/domains/now_order"
	"terminal_monitor_ui/config"
	"terminal_monitor_ui/constant"
	"terminal_monitor_ui/logger"
	"terminal_monitor_ui/pkg/kafka"
	"time"

	"go.uber.org/zap"
)

type UseCase struct {
	nowOrder no.INowOrderUsecase
}

type MessageHandler struct {
	UseCase
	cfg    *config.AppConfig
	reader []kafka.Subscriber
}

func NewMessageHandler(
	cfg *config.AppConfig,
	reader []kafka.Subscriber,
	nowOrder no.INowOrderUsecase,
) *MessageHandler {
	return &MessageHandler{
		cfg:    cfg,
		reader: reader,
		UseCase: UseCase{
			nowOrder: nowOrder,
		},
	}
}

func (m *MessageHandler) Run(errRestart chan error) {
	if len(m.reader) > 0 {
		for _, sub := range m.reader {
			go sub.Read(m.processMessage, errRestart)
		}
	}
}

func (m *MessageHandler) Stop() {
	if len(m.reader) > 0 {
		for _, sub := range m.reader {
			sub.Close()
		}
	}
}

func (m *MessageHandler) processMessage(ctx context.Context, topic string, vals []byte) error {
	defer func() error {
		if err := recover(); err != nil {
			zap.S().Errorf("runtime error: %v", err)
			_ = logger.WriteFile("./logger/error.txt", fmt.Sprintf("runtime error: %v", err))
			return fmt.Errorf("runtime error: %v", err)
		}
		return nil
	}()

	if len(vals) == 0 || string(vals) == "" {
		return nil
	}

	var err error

	switch topic {
	case m.cfg.Kafka.TopicNames.NowOrder:
		err = m.nowOrder.ConsumeMessage(ctx, vals)
	case m.cfg.Kafka.TopicNames.ProductTester:
		time.Sleep(10 * time.Second)
	}

	if err != nil {
		_ = logger.WriteFile(constant.ERROR_FILE_PATH, fmt.Sprintf("processMessage on topic [%s] has err:[%v] --msg:[%s]\n\n", topic, err, string(vals)))
		return err
	}

	return nil
}
