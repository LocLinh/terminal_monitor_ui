package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

func (m UiModel) View() string {
	// The header
	var s strings.Builder
	s.WriteString("\n\n")

	// Iterate over our choices
	for _, err := range m.Errors {
		// Render the row
		fmt.Fprintf(&s, "%s\n", err)
	}

	// The footer
	s.WriteString("\nPress q to quit.\n")

	// Logs
	s.WriteString("\nLogs:\n")
	for _, log := range m.Logs {
		fmt.Fprintf(&s, "%s\n", log)
	}

	tableRows := []table.Row{}
	for _, lag := range m.MessageBehind {
		fmt.Fprintf(&s, "topic: %s, behind: %d\n", lag.Topic, lag.Lag)
		tableRows = append(tableRows, table.Row{
			lag.Topic, strconv.FormatInt(lag.Partition, 10), strconv.FormatInt(lag.Latest, 10), strconv.FormatInt(lag.Committed, 10), strconv.FormatInt(lag.Lag, 10),
		})
	}

	// Send the UI for rendering
	m.MessageBehindTable.SetRows(tableRows)

	mainView := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")). // purple
		Render(m.Stopwatch.View(), "\n", m.MessageBehindTable.View())

	return mainView + s.String()
}
