package ui

import (
	"terminal_monitor_ui/config"
	"terminal_monitor_ui/ui/components/chart"
	"time"

	"github.com/IBM/sarama"
	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewUI(cfg *config.AppConfig, logChan chan string) UiModel {
	config := sarama.NewConfig()
	config.Version, _ = sarama.ParseKafkaVersion(cfg.Kafka.Version)
	config.ClientID = cfg.Kafka.GroupId

	config.Consumer.Group.Heartbeat.Interval = 5 * time.Second
	config.Consumer.Group.Session.Timeout = 15 * time.Second
	config.Consumer.MaxProcessingTime = 300 * time.Millisecond
	config.Consumer.Return.Errors = true

	if cfg.Kafka.Newest {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	client, err := sarama.NewClient(
		cfg.Kafka.Addrs,
		config,
	)
	if err != nil {
		panic(err)
	}

	messageBehindTable := table.New(
		table.WithColumns([]table.Column{
			{Title: "Topic", Width: 30},
			{Title: "Partition", Width: 12},
			{Title: "Latest", Width: 12},
			{Title: "Committed", Width: 12},
			{Title: "Lag", Width: 12},
		}),
	)
	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tableStyle.Selected = tableStyle.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	messageBehindTable.SetStyles(tableStyle)

	return UiModel{
		Errors:             []string{},
		ThroughputChart:    chart.InitChart(),
		KafkaClient:        client,
		cfg:                cfg,
		MessageBehindTable: messageBehindTable,
		Stopwatch:          stopwatch.NewWithInterval(time.Millisecond),
		LogChan:            logChan,
	}
}

func waitForLog(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return LogMsg(<-sub)
	}
}

func (m UiModel) Init() tea.Cmd {
	dataSet := []float64{0, 2, 4, 6, 8, 10, 8, 6, 4, 2, 0}
	for i, v := range dataSet {
		date := time.Now().Add(time.Hour * time.Duration(24*i))
		m.ThroughputChart.Push(tslc.TimePoint{
			Time:  date,
			Value: v,
		})
	}

	if m.KafkaClient == nil {
		panic("Kafka client is nil")
	}

	return tea.Batch(
		m.Stopwatch.Init(),
		pollLagCmd(m.KafkaClient, m.cfg.Kafka.Group, m.cfg.Kafka.TopicNames.ToList(), time.Second*10),
		waitForLog(m.LogChan),
	)
}

func pollLagCmd(
	client sarama.Client,
	group string,
	topics []string,
	interval time.Duration,
) tea.Cmd {
	return tea.Tick(interval, func(time.Time) tea.Msg {
		lag, err := fetchLag(client, group, topics)
		if err != nil {
			return KafkaErrorUIMsg{
				Err: err,
			}
		}
		return LagMetricsMsg(lag)
	})
}

func throughputTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return ThroughputTickMsg(t)
	})
}
