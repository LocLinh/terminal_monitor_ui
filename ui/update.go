package ui

import (
	"strings"
	"time"

	"github.com/IBM/sarama"
	tea "github.com/charmbracelet/bubbletea"
)

func (m UiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	// case KafkaUIMsg:
	case KafkaErrorUIMsg:
		m.Errors = append(m.Errors, msg.Err.Error())
		return m, nil
	case LogMsg:
		m.Logs = append(m.Logs, string(msg))
		if len(m.Logs) > 100 {
			m.Logs = m.Logs[1:]
		}
		m.LogViewport.SetContent(strings.Join(m.Logs, "\n"))
		m.LogViewport.GotoBottom()
		return m, waitForLog(m.LogChan)
	case LagMetricsMsg:
		m.MessageBehind = msg
		return m, pollLagCmd(
			m.KafkaClient,
			m.cfg.Kafka.Group,
			m.cfg.Kafka.TopicNames.ToList(),
			5*time.Second,
		)
	}

	m.ThroughputChart, _ = m.ThroughputChart.Update(msg)
	m.ThroughputChart.DrawBrailleAll()

	var cmd tea.Cmd
	m.Stopwatch, cmd = m.Stopwatch.Update(msg)
	m.LogViewport, _ = m.LogViewport.Update(msg)

	return m, cmd

}

func fetchLag(
	client sarama.Client,
	group string,
	topics []string,
) ([]PartitionLag, error) {
	om, err := sarama.NewOffsetManagerFromClient(group, client)
	if err != nil {
		return nil, err
	}
	defer om.Close()

	var result []PartitionLag

	for _, topic := range topics {
		partitions, err := client.Partitions(topic)
		if err != nil {
			return nil, err
		}

		for _, p := range partitions {
			pom, err := om.ManagePartition(topic, p)
			if err != nil {
				return nil, err
			}

			committed, _ := pom.NextOffset()
			latest, err := client.GetOffset(topic, p, sarama.OffsetNewest)
			if err != nil {
				pom.Close()
				return nil, err
			}

			lag := latest - committed
			if committed < 0 {
				lag = latest
			}

			result = append(result, PartitionLag{
				Topic:     topic,
				Partition: int64(p),
				Latest:    latest,
				Committed: committed,
				Lag:       lag,
			})

			pom.Close()
		}
	}

	return result, nil
}
