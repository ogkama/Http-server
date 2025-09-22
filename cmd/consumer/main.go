package main

import (
	"encoding/base64"
	"encoding/json"
	"http_server/cmd/config"
	"http_server/domain"
	"http_server/repository/postgres"
	"log"
	"net/http"
	"os"
	"time"

	rabbitmq "http_server/repository/rabbit_mq"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// http://localhost:15672/api/healthchecks/node

type Response struct {
    domain.Task  
}

type Result struct {
	Result string `json:"result"`
}

var (
	taskProcessed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tasks_processed_total",
		Help: "Total number of processed tasks",
	})
	taskFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tasks_failed_total",
		Help: "Total number of failed tasks",
	})
	taskInProgress = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasks_in_progress",
		Help: "Total number of tasks in progress",
	})
	taskDuration = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "task_processing_duration_time",
		Help: "Task processing duration time in seconds",
		Objectives: map[float64]float64{
            0.5:  0.05, // 50th percentile
            0.9:  0.01, // 90th percentile
            0.99: 0.001, // 99th percentile
        },
	},
	[]string{"task"},
	)
)

func init(){
	prometheus.MustRegister(taskProcessed, taskFailed, taskInProgress, taskDuration)
}

func main() {
	go func(){
		http.Handle("/metrics",promhttp.Handler())
		http.ListenAndServe(":2112",nil)
	}()


	appFlags := config.ParseFlags()
	var cfg config.AppConfig
	config.MustLoad(appFlags.ConfigPath, &cfg)

	connStr := os.Getenv("DB_CONN_STR")
	if connStr == "" {
		log.Fatal("DB_CONN_STR environment: variable is required")
	}

	storage, err := postgres.NewPostgresStorage(connStr)
	if err != nil {
		log.Fatalf("Can't connect to the db %s", err.Error())
	}

	taskConsumer, err := rabbitmq.NewRabbitMQConsumer(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("Failed creating rabbitmq consumer: %s", err)
	}

	defer taskConsumer.Close()

	msgs, err := taskConsumer.Consume()
	if err != nil {
		log.Fatalf("Failed to register consumer %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			startTime := time.Now()

			taskInProgress.Inc()

			log.Println("Got task from queue")

			var resp Response
			if err := json.Unmarshal(d.Body, &resp); err != nil {
				log.Println("json unmarshal error:", err)
			}

			log.Println("Ino Data Len : ", len(resp.Data))

			result, err := RunPythonContainer(resp.Data)
			if err != nil {
				log.Println("Error processing:", err)
				storage.SetStatus(resp.Task_id, "ERROR")
				taskFailed.Inc()
				taskInProgress.Dec()
				d.Nack(false, false)
				continue
			}
			
			taskProcessed.Inc()
			taskInProgress.Dec()

			elapsed := time.Since(startTime).Seconds()
			taskDuration.WithLabelValues(resp.Task_id).Observe(elapsed)

			log.Println("Processed file, result size:", len(result))

			encoded := base64.StdEncoding.EncodeToString(result)
			
			resultJSON := Result{Result: encoded}

    		jsonBytes, err := json.Marshal(resultJSON)
    		if err != nil {
        		log.Fatal(err)
    		}

			storage.SetStatus(resp.Task_id, "done")
			storage.SetResult(resp.Task_id, jsonBytes)

			// подтверждаем задачу
			d.Ack(false)
		}
	}()

	log.Println(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
