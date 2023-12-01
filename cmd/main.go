package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"konsume/pkg/config"
)

func main() {
	slog.Info("Starting konsume")

	_, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	//signalChannel := setupSignalHandling()
	//waitForShutdown(signalChannel)

	slog.Info("Shut down gracefully")
}

// setupSignalHandling configures the signal handling for os.Interrupt and syscall.SIGTERM.
func setupSignalHandling() chan bool {
	done := make(chan bool, 1)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChannel
		slog.Info("Received signal", "signal", sig)
		done <- true
	}()

	return done
}

// waitForShutdown blocks until a shutdown signal is received.
func waitForShutdown(done chan bool) {
	<-done
}
