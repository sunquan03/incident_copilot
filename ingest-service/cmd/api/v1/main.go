package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/sunquan03/ingest-service/internal/brokers"
	"github.com/sunquan03/ingest-service/internal/config"
	"github.com/sunquan03/ingest-service/internal/database"
	"github.com/sunquan03/ingest-service/internal/handlers"
	"github.com/sunquan03/ingest-service/internal/repositories"
	"github.com/sunquan03/ingest-service/internal/routers"
	"github.com/sunquan03/ingest-service/internal/services"
	"github.com/valyala/fasthttp"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	kafkaConf := brokers.DefaultProducerConfig(cfg.Kafka.Brokers)
	kafkaCli, err := sarama.NewClient(cfg.Kafka.Brokers, sarama.NewConfig())
	if err != nil {
		log.Fatal(err)
	}
	producer, err := brokers.NewProducer(kafkaConf)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDB(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}

	repository := repositories.NewRepository(db)
	service := services.NewService(repository, producer)
	reqHandler := handlers.NewHandler(service)
	healthHandler := handlers.NewHealthHandler(db, kafkaCli)

	r := routers.NewRouter(reqHandler, healthHandler)
	handler := r.Setup()
	server := fasthttp.Server{
		Handler:      handler,
		WriteTimeout: 7 * time.Second,
		ReadTimeout:  7 * time.Second,
		IdleTimeout:  7 * time.Second,
	}

	server.ListenAndServe(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
}
