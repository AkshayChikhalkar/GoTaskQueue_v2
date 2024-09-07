package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"

	pb "github.com/akshaychikhalkar/GoTaskQueue_v2/tasks" // Correct import path

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	port = ":50051" // Port for gRPC server
)

var (
	version string // Version variable to be set at build time

	taskCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "consumer_task_count",
			Help: "Total number of tasks consumed by the consumer.",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(taskCounter)
}

type server struct {
	pb.UnimplementedConsumerServiceServer // Embed the correct unimplemented server struct
}

func (s *server) ConsumeTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	log.Printf("Consumed task: Type: %v, Value: %v", req.TaskType, req.TaskValue)

	// Update metrics
	taskCounter.WithLabelValues("success").Inc()

	return &pb.TaskResponse{Success: true}, nil
}

func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	flag.Parse()

	if *versionFlag {
		log.Printf("Consumer Service Version: %s", version)
		return
	}

	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	log.Info("Consumer application started")

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":9091", nil); err != nil {
			log.Fatalf("Failed to start metrics server: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterConsumerServiceServer(grpcServer, &server{}) // Register ConsumerService

	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
