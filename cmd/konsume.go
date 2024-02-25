package konsume

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bugrakocabay/konsume/pkg/common"
	"github.com/bugrakocabay/konsume/pkg/config"
	"github.com/bugrakocabay/konsume/pkg/database"
	"github.com/bugrakocabay/konsume/pkg/metrics"
	"github.com/bugrakocabay/konsume/pkg/queue"
	"github.com/bugrakocabay/konsume/pkg/queue/activemq"
	"github.com/bugrakocabay/konsume/pkg/queue/kafka"
	"github.com/bugrakocabay/konsume/pkg/queue/rabbitmq"
	"github.com/bugrakocabay/konsume/pkg/runner"
)

func Execute() {
	slog.Info("Starting konsume")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}

	setupLogger(cfg)

	consumerMap := make(map[string]queue.MessageQueueConsumer)
	providerMap := make(map[string]*config.ProviderConfig)

	initProviders(cfg, consumerMap, providerMap)
	databaseMap, err := initDatabases(cfg)
	if err != nil {
		slog.Error("Failed to initialize databases", "error", err)
		return
	}
	if cfg.Metrics != nil && cfg.Metrics.Enabled {
		metrics.InitMetrics(cfg.Metrics)
	}

	go func() {
		if err = runner.StartConsumers(cfg, consumerMap, providerMap, databaseMap); err != nil {
			slog.Error("Failed to start consumers", "error", err)
		}
	}()

	signalChannel := setupSignalHandling()
	waitForShutdown(signalChannel)

	runner.StopConsumers(consumerMap, databaseMap)

	slog.Info("Shut down gracefully")
}

// setupLogger configures the logger based on the configuration
func setupLogger(cfg *config.Config) {
	if cfg.Log == "json" {
		if cfg.Debug {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
		} else {
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
		}
	} else {
		if cfg.Debug {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
		} else {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
		}
	}
}

// initProviders initializes the consumers for each provider
func initProviders(cfg *config.Config, consumers map[string]queue.MessageQueueConsumer, providerMap map[string]*config.ProviderConfig) {
	for _, provider := range cfg.Providers {
		slog.Debug("Initializing provider", "provider", provider.Name, "type", provider.Type)
		switch provider.Type {
		case common.QueueSourceRabbitMQ:
			consumer := rabbitmq.NewConsumer(provider.AMQPConfig)
			consumers[provider.Name] = consumer
			providerMap[provider.Name] = provider
		case common.QueueSourceKafka:
			consumer := kafka.NewConsumer(provider.KafkaConfig)
			consumers[provider.Name] = consumer
			providerMap[provider.Name] = provider
		case common.QueueSourceActiveMQ:
			consumer := activemq.NewConsumer(provider.StompMQConfig)
			consumers[provider.Name] = consumer
			providerMap[provider.Name] = provider
		default:
			log.Fatalf("Unknown queue source: %s", provider.Type)
		}
	}
}

func initDatabases(cfg *config.Config) (map[string]database.Database, error) {
	dbMap := make(map[string]database.Database)
	for _, dbConfig := range cfg.Databases {
		db, err := database.LoadDatabasePlugin(dbConfig.Type)
		if err != nil {
			return nil, err
		}
		if err = db.Connect(dbConfig.ConnectionString); err != nil {
			slog.Error("Failed to connect to database", "type", dbConfig.Type, "error", err)
			if dbConfig.Retry > 0 {
				for i := 1; i <= dbConfig.Retry; i++ {
					time.Sleep(time.Duration(5) * time.Second)
					slog.Info("Retrying to connect to database", "retry", i)
					if err = db.Connect(dbConfig.ConnectionString); err == nil {
						break
					}
				}
			}
			return nil, err
		}
		dbMap[dbConfig.Name] = db
	}
	return dbMap, nil
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
