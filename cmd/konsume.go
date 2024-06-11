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
	databaseMap, err := initDatabases(cfg.Databases)
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
	var handler slog.Handler
	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	if cfg.Log == "json" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	slog.SetDefault(slog.New(handler))
}

// initProviders initializes the consumers for each provider
func initProviders(
	cfg *config.Config,
	consumers map[string]queue.MessageQueueConsumer,
	providerMap map[string]*config.ProviderConfig,
) {
	factories := map[string]queue.Factory{
		common.QueueSourceRabbitMQ: rabbitmq.NewConsumerFactory,
		common.QueueSourceKafka:    kafka.NewConsumerFactory,
		common.QueueSourceActiveMQ: activemq.NewConsumerFactory,
	}

	for _, provider := range cfg.Providers {
		factory, exists := factories[provider.Type]
		if !exists {
			slog.Error("Unknown provider type", "type", provider.Type)
			continue
		}

		consumer, err := factory(provider)
		if err != nil {
			slog.Error("Failed to initialize provider", "error", err)
			continue
		}

		consumers[provider.Name] = consumer
		providerMap[provider.Name] = provider
	}
}

// initDatabases initializes the databases based on the configuration
func initDatabases(cfg []*config.DatabaseConfig) (map[string]database.Database, error) {
	dbMap := make(map[string]database.Database)
	for _, dbConfig := range cfg {
		db, err := database.LoadDatabasePlugin(dbConfig.Type)
		if err != nil {
			slog.Error("Failed to load database plugin", "type", dbConfig.Type, "error", err)
			return nil, err
		}
		err = connectDbRetry(db, *dbConfig)
		if err != nil {
			slog.Error("Failed to connect to database after retries", "type", dbConfig.Type, "error", err)
			return nil, err
		}
		dbMap[dbConfig.Name] = db
	}
	return dbMap, nil
}

// connectDbRetry attempts to connect to the database with retry logic
func connectDbRetry(db database.Database, dbConfig config.DatabaseConfig) error {
	var err error
	for attempt := 0; attempt <= dbConfig.Retry; attempt++ {
		if attempt > 0 {
			slog.Info("Retrying to connect to database", "type", dbConfig.Type, "attempt", attempt)
			time.Sleep(5 * time.Second)
		}
		err = db.Connect(dbConfig.ConnectionString, dbConfig.Database)
		if err == nil {
			return nil
		}
	}
	return err
}

// setupSignalHandling configures the signal handling for os.Interrupt and syscall.SIGTERM
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

// waitForShutdown blocks until a shutdown signal is received
func waitForShutdown(done chan bool) {
	<-done
}
