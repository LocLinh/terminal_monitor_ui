package chart

import (
	tslc "github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/lipgloss"
)

func InitChart() tslc.Model {
	width := 30
	height := 12

	chart := tslc.New(width, height)

	chart.SetStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("9"))) // red

	return chart
}
