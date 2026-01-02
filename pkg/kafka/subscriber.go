package kafka

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"terminal_monitor_ui/logger"
	"terminal_monitor_ui/pkg/contxt"
	"terminal_monitor_ui/pkg/graceful"

	"terminal_monitor_ui/config"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var _ Subscriber = (*subscriber)(nil)

type subscriber struct {
	consumer sarama.ConsumerGroup
	topics   []string
	group    string
}

// NewSubscriber ...
func NewSubscriber(cfg config.KafkaConfig) Subscriber {
	config := sarama.NewConfig()
	config.Version, _ = sarama.ParseKafkaVersion(cfg.Version)
	config.ClientID = cfg.Group

	config.Consumer.Group.Heartbeat.Interval = 5 * time.Second
	config.Consumer.Group.Session.Timeout = 15 * time.Second
	config.Consumer.MaxProcessingTime = 300 * time.Millisecond
	config.Consumer.Return.Errors = true

	if cfg.Newest {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	//start consumer group
	consumer, err := sarama.NewConsumerGroup(cfg.Addrs, cfg.Group, config)
	if err != nil {
		zap.S().Errorf("init consumer group fail err: %v", err)
		_ = logger.WriteFile("./logger/error.txt", err.Error()+"\n")
		panic(err)
	}

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			zap.S().Errorf("consume error: %v", err)
			_ = logger.WriteFile("./logger/error.txt", err.Error()+"\n")
		}
	}()

	return &subscriber{
		consumer: consumer,
		topics:   cfg.Topics,
		group:    cfg.Group,
	}
}

func (s *subscriber) Read(callback CallBack, errRestart chan error) {
	if s.consumer == nil {
		zap.S().Error("consumer nil -> missing consumer")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer s.Close()

	c := &consumerHandler{
		fc:    callback,
		ready: make(chan bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := s.consumer.Consume(ctx, s.topics, c); err != nil {
				zap.S().Errorf("kafka consume topics err: %v", err)
				_ = logger.WriteFile("./logger/error.txt", err.Error()+"\n")
				errRestart <- err
				return
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				zap.S().Info(ctx.Err())
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready // Await till the consumer has been set up
	// zap.S().Debug("kafka consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case <-ctx.Done():
		zap.S().Info("terminating: context cancelled")
	case <-sigterm:
		zap.S().Info("terminating: via signal")
	}
	wg.Wait()
}

func (s subscriber) Close() {
	if err := s.consumer.Close(); err != nil {
		zap.S().Errorf("Error closing client: %v", err)
		_ = logger.WriteFile("./logger/error.txt", err.Error()+"\n")
	}
}

type consumerHandler struct {
	fc    func(context.Context, string, []byte) error
	ready chan bool
}

func (c *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	// zap.S().Info("setup consumer group handler")
	close(c.ready)
	return nil
}
func (c *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	zap.S().Info("cleanup consumer group handler")
	return nil
}

func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	stop := make(chan bool, 1)
	graceful.Stop(0, func() {
		stop <- true
	})

	logger := zap.L()
	for {
		select {
		case msg, ok := <-claim.Messages():
			if ok {
				graceful.AddProcess()
				requestID := uuid.New().String()
				ctx := ctxzap.ToContext(contxt.Background(), logger.With(zap.String("prefix", "consumer"), zap.String("request_id", requestID)))

				err := c.fc(ctx, msg.Topic, msg.Value)
				if err != nil {
					continue
				}
				session.MarkMessage(msg, "")
				graceful.DoneProcess()
			}
		case <-session.Context().Done():
			zap.S().Debugf("session timeout")
			return nil
		case <-stop:
			graceful.ShutDown()
			return nil
		}
	}
}

func NewSubscriberV2(cfg config.KafkaConfig) Subscriber {
	config := sarama.NewConfig()
	config.Version, _ = sarama.ParseKafkaVersion(cfg.Version)
	config.ClientID = cfg.GroupId

	config.Consumer.Group.Heartbeat.Interval = 5 * time.Second
	config.Consumer.Group.Session.Timeout = 15 * time.Second
	config.Consumer.MaxProcessingTime = 300 * time.Millisecond
	config.Consumer.Return.Errors = true

	if cfg.Newest {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	//start consumer group
	consumer, err := sarama.NewConsumerGroup(cfg.Addrs, cfg.Group, config)
	if err != nil {
		zap.S().Errorf("init consumer group fail err: %v", err)
		panic(err)
	}

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			zap.S().Errorf("consume error: %v", err)
		}
	}()

	return &subscriber{
		consumer: consumer,
		topics:   cfg.Topics,
		group:    cfg.Group,
	}
}
