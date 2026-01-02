package graceful

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	stop, shutdown sync.WaitGroup
)

func init() {
	stop.Add(1)
	go func() {
		q := make(chan os.Signal, 1)
		signal.Notify(q, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
		zap.S().Debug("receive signal: ", <-q)
		stop.Done()
	}()
}

// AddProcess ...
func AddProcess() {
	stop.Add(1)
}

// DoneProcess ...
func DoneProcess() {
	stop.Done()
}

// Stop ...
func Stop(t time.Duration, f func()) {
	shutdown.Add(1)
	go func() {
		stop.Wait()
		zap.S().Debug("Stop")
		time.Sleep(t)
		f()
		shutdown.Done()
	}()
}

// ShutDown ...
func ShutDown() {
	stop.Wait()
	shutdown.Wait()
	zap.S().Debug("Shutdown")
	os.Exit(0)
}

// ShutDownSlowly ...
func ShutDownSlowly(delay time.Duration) {
	stop.Wait()
	shutdown.Wait()
	zap.S().Debug("ShutDownSlowly")
	time.Sleep(delay)
	os.Exit(0)
}

func WaitShutdown() {
	stop := make(chan bool, 1)
	Stop(0, func() {
		stop <- true
	})
	<-stop
	ShutDown()
}
