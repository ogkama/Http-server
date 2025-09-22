package main

import (
	"http_server/cmd/config"
	"http_server/repository/postgres"
	rabbitmq "http_server/repository/rabbit_mq"
	"http_server/repository/redis"
	"http_server/usecases/services"
	"log"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"http_server/api/http"
	_ "http_server/docs"
	pkgHttp "http_server/pkg/http"
)

// @title My API
// @version 1.0
// @description This is a sample server.

// @host localhost:8080
// @BasePath /

func main() {

	log.SetOutput(os.Stderr)
	
	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)


	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		log.Fatalf("DB_CONN_STR environment: variable is required")
	}

	postgres_storage, err := postgres.NewPostgresStorage(connStr)
	if err != nil {
		log.Fatalf("Can't connect to the db %s", err.Error())
	}


	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Fatalf("Error parsing REDIS_DB: %s", err.Error())
	}

	redis_storage, err := redis.NewRedisClient(redisAddr, redisPassword, redisDB)
	if err != nil {
		log.Fatalf("Error creatig Redis client: %s", err.Error())
	}


	taskSender, err := rabbitmq.NewRabbitMQSender(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("failed creating rabbitMQ: %s", err.Error())
	}
	

	taskService := services.NewTask(postgres_storage, taskSender)
	userService := services.NewUser(postgres_storage, redis_storage)

	TaskHandlers := http.NewTaskHandler(taskService)
	UserHandlers := http.NewUserHandler(userService)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	UserHandlers.WithUserHandlers(r)
	TaskHandlers.WithTaskHandlers(r, UserHandlers)

	log.Printf("Starting server on %s", cfg.Address)
	if err := pkgHttp.CreateAndRunServer(r, cfg.Address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
