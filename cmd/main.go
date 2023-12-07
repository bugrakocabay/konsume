package konsume

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bugrakocabay/konsume/pkg/common"
	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/queue"
	"github.com/bugrakocabay/konsume/pkg/queue/kafka"
	"github.com/bugrakocabay/konsume/pkg/queue/rabbitmq"
	"github.com/bugrakocabay/konsume/pkg/runner"
)

func Execute() {
	slog.Info("Starting konsume")

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	consumers := make(map[string]queue.MessageQueueConsumer)
	for _, provider := range cfg.Providers {
		switch provider.Type {
		case common.QueueSourceRabbitMQ:
			consumer := rabbitmq.NewConsumer(provider.AMQPConfig)
			consumers[provider.Name] = consumer
		case common.QueueSourceKafka:
			consumer := kafka.NewConsumer(provider.KafkaConfig)
			consumers[provider.Name] = consumer
		default:
			log.Fatalf("Unknown queue source: %s", provider.Type)
		}
	}

	if err = runner.StartConsumers(cfg, consumers); err != nil {
		log.Fatalf("Failed to start consumers: %s", err)
	}

	signalChannel := setupSignalHandling()
	waitForShutdown(signalChannel)

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
