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

	tableRows := []table.Row{}
	for _, lag := range m.MessageBehind {
		tableRows = append(tableRows, table.Row{
			lag.Topic, strconv.FormatInt(lag.Partition, 10), strconv.FormatInt(lag.Latest, 10), strconv.FormatInt(lag.Committed, 10), strconv.FormatInt(lag.Lag, 10),
		})
	}

	// Send the UI for rendering
	m.MessageBehindTable.SetRows(tableRows)

	view := lipgloss.JoinVertical(lipgloss.Left,
		m.Stopwatch.View(),
		m.MessageBehindTable.View(),
		m.LogViewport.View(),
	)

	mainView := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")). // purple
		Render(view)

	return mainView + s.String()
}
