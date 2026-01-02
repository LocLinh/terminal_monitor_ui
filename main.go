package main

import (
	"fmt"
	"log"
	"os"

	"terminal_monitor_ui/application/worker"
	"terminal_monitor_ui/config"
	"terminal_monitor_ui/logger"
	ui "terminal_monitor_ui/ui"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

func main() {
	logChan := make(chan string, 100)
	zapLogger := logger.NewLogger(logChan)
	defer zapLogger.Sync()
	zap.ReplaceGlobals(zapLogger)

	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	go func() {
		errRestart := make(chan error)
		worker := worker.NewWorker(config)
		worker.Start(errRestart)
	}()

	p := tea.NewProgram(ui.NewUI(config, logChan))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
