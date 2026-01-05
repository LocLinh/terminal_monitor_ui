package ui

import (
	"terminal_monitor_ui/config"
	"time"

	"github.com/IBM/sarama"
	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
)

type UiModel struct {
	Errors             []string
	ThroughputChart    tslc.Model
	MessageBehind      []PartitionLag
	KafkaClient        sarama.Client
	cfg                *config.AppConfig
	MessageBehindTable table.Model
	Stopwatch          stopwatch.Model
	Logs               []string
	LogChan            chan string
	LogViewport        viewport.Model
}

type LogMsg string

type KafkaUIMsg struct {
	Key   string
	Value string
}

type KafkaErrorUIMsg struct {
	Err error
}

type PartitionLag struct {
	Topic     string
	Partition int64
	Latest    int64
	Committed int64
	Lag       int64
}

type LagMetricsMsg []PartitionLag

type ThroughputTickMsg time.Time
