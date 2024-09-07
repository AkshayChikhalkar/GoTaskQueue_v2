package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	pb "github.com/akshaychikhalkar/GoTaskQueue_v2/tasks" // Import generated package
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051" // Consumer service address
)

// Declare version variable
var version = "development" // Default version

// Define metrics
var (
	taskCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "producer_task_count",
			Help: "Total number of tasks sent by the producer.",
		},
		[]string{"status"},
	)
	taskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "producer_task_duration_seconds",
			Help:    "Histogram of task send durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(taskCount)
	prometheus.MustRegister(taskDuration)
}

func main() {
	// Handle version flag
	versionFlag := flag.Bool("version", false, "Show the version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Version:", version)
		os.Exit(0)
	}

	// Set up logging
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	log.Info("Producer application started")

	// Set up metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	// Connect to gRPC server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a client for ProducerService
	client := pb.NewProducerServiceClient(conn)

	// Generate random tasks and send to the consumer
	for {
		taskType := rand.Intn(10)
		taskValue := rand.Intn(100)

		req := &pb.TaskRequest{
			TaskType:  int32(taskType),
			TaskValue: int32(taskValue),
		}

		start := time.Now()
		res, err := client.ProduceTask(context.Background(), req)
		duration := time.Since(start).Seconds()

		if err != nil {
			log.Errorf("Error sending task: %v", err)
			taskCount.WithLabelValues("error").Inc()
			taskDuration.WithLabelValues("error").Observe(duration)
		} else {
			log.Infof("Task sent: %v, Success: %v", req, res.Success)
			taskCount.WithLabelValues("success").Inc()
			taskDuration.WithLabelValues("success").Observe(duration)
		}

		time.Sleep(2 * time.Second)
	}
}
