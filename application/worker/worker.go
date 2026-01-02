package worker

import (
	"fmt"
	"strconv"
	noRepository "terminal_monitor_ui/application/domains/now_order/repository"
	noUsecase "terminal_monitor_ui/application/domains/now_order/usecase"
	"terminal_monitor_ui/application/messaging"
	"terminal_monitor_ui/config"
	"terminal_monitor_ui/pkg/database"
	"terminal_monitor_ui/pkg/kafka"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	RETRY_DELAY_TIME = 10 * time.Second
)

type Worker struct {
	cfg         *config.AppConfig
	db          *gorm.DB
	ksubscriber []kafka.Subscriber
}

func NewWorker(cfg *config.AppConfig) *Worker {
	database, err := database.GetDB(cfg)
	if err != nil {
		zap.S().Errorf("Error loading DB %v", err)
		panic("Error loading DB")
	}

	kafkaCfg := config.KafkaConfig{
		Addrs:           cfg.Kafka.Addrs,
		Group:           cfg.Kafka.Group,
		MaxMessageBytes: cfg.Kafka.MaxMessageBytes,
		Compress:        cfg.Kafka.Compress,
		Newest:          cfg.Kafka.Newest,
		Version:         cfg.Kafka.Version,
	}

	kafkaCfg.Topics = cfg.Kafka.TopicNames.ToList()
	ksubscribers := make([]kafka.Subscriber, cfg.Kafka.Partition)
	for i := 0; i < cfg.Kafka.Partition; i++ {
		ksubscriber := kafka.NewSubscriberV2(config.KafkaConfig{
			Addrs:           cfg.Kafka.Addrs,
			Group:           cfg.Kafka.Group,
			GroupId:         cfg.Kafka.Group + "-" + strconv.Itoa(i),
			MaxMessageBytes: cfg.Kafka.MaxMessageBytes,
			Compress:        cfg.Kafka.Compress,
			Newest:          cfg.Kafka.Newest,
			Version:         cfg.Kafka.Version,
			Topics:          cfg.Kafka.TopicNames.ToList(),
		})
		ksubscribers[i] = ksubscriber
	}

	return &Worker{
		cfg:         cfg,
		db:          database,
		ksubscriber: ksubscribers,
	}
}

func (w *Worker) Start(errRestart chan error) {
	defer func() error {
		if err := recover(); err != nil {
			time.Sleep(RETRY_DELAY_TIME)
			for _, sub := range w.ksubscriber {
				sub.Close()
			}
			ksubscribers := make([]kafka.Subscriber, w.cfg.Kafka.Partition)
			for i := 0; i < w.cfg.Kafka.Partition; i++ {
				ksubscriber := kafka.NewSubscriberV2(config.KafkaConfig{
					Addrs:           w.cfg.Kafka.Addrs,
					Group:           w.cfg.Kafka.Group,
					GroupId:         w.cfg.Kafka.Group + "-" + strconv.Itoa(i),
					MaxMessageBytes: w.cfg.Kafka.MaxMessageBytes,
					Compress:        w.cfg.Kafka.Compress,
					Newest:          w.cfg.Kafka.Newest,
					Version:         w.cfg.Kafka.Version,
					Topics:          w.cfg.Kafka.TopicNames.ToList(),
				})
				ksubscribers[i] = ksubscriber
			}
			w.ksubscriber = ksubscribers
			w.Start(errRestart)
			return fmt.Errorf("runtime error: %v", err)
		}
		return nil
	}()

	w.init(errRestart)
	if <-errRestart != nil {
		time.Sleep(RETRY_DELAY_TIME)
		for _, sub := range w.ksubscriber {
			sub.Close()
		}
		ksubscribers := make([]kafka.Subscriber, w.cfg.Kafka.Partition)
		for i := 0; i < w.cfg.Kafka.Partition; i++ {
			ksubscriber := kafka.NewSubscriberV2(config.KafkaConfig{
				Addrs:           w.cfg.Kafka.Addrs,
				Group:           w.cfg.Kafka.Group,
				GroupId:         w.cfg.Kafka.Group + "-" + strconv.Itoa(i),
				MaxMessageBytes: w.cfg.Kafka.MaxMessageBytes,
				Compress:        w.cfg.Kafka.Compress,
				Newest:          w.cfg.Kafka.Newest,
				Version:         w.cfg.Kafka.Version,
				Topics:          w.cfg.Kafka.TopicNames.ToList(),
			})
			ksubscribers[i] = ksubscriber
		}
		w.ksubscriber = ksubscribers
		w.Start(errRestart)
	}
}

func (w *Worker) Stop() {

}

func (w *Worker) init(errRestart chan error) {
	// now order
	noRepo := noRepository.InitNowOrderRepository(w.db)
	noUsecase := noUsecase.InitUsecase(noRepo)

	// inject usecase
	handler := messaging.NewMessageHandler(w.cfg,
		w.ksubscriber,
		noUsecase,
	)
	go handler.Run(errRestart)
	zap.S().Infof("Worker started")
}
