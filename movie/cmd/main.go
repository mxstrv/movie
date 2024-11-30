package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"log"
	"movieapp/gen"
	"movieapp/movie/internal/controller/movie"
	metadatagateway "movieapp/movie/internal/gateway/metadata/http"
	ratinggateway "movieapp/movie/internal/gateway/rating/http"
	grpchandler "movieapp/movie/internal/handler/grpc"
	"movieapp/pkg/discovery"
	"movieapp/pkg/discovery/consul"
	"net"
	"os"
	"time"
)

const serviceName = "movie"

func main() {
	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	port := cfg.API.Port
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(
		ctx,
		instanceID,
		serviceName,
		fmt.Sprintf("localhost:%d", port),
	); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state:" + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer func(registry *consul.Registry, ctx context.Context, instanceID string, _ string) {
		_ = registry.Deregister(ctx, instanceID, serviceName)
	}(registry, ctx, instanceID, serviceName)
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := movie.New(ratingGateway, metadataGateway)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
