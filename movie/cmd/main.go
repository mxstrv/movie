package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"movieapp/movie/internal/controller/movie"
	metadatagateway "movieapp/movie/internal/gateway/metadata/http"
	ratinggateway "movieapp/movie/internal/gateway/rating/http"
	httphandler "movieapp/movie/internal/handler/http"
	"movieapp/pkg/discovery"
	"movieapp/pkg/discovery/consul"
	"net/http"
	"time"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
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
	svc := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(svc)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
