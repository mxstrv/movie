package main

import (
	"google.golang.org/grpc"
	"log"
	"movieapp/gen"
	"movieapp/metadata/internal/controller/metadata"
	grpchandler "movieapp/metadata/internal/handler/grpc"
	"movieapp/metadata/internal/repository/memory"
	"net"
)

//const serviceName = "metadata"
//
//func main() {
//	var port int
//	flag.IntVar(&port, "port", 8081, "API handler port")
//	flag.Parse()
//	log.Printf("Starting the movie metadata service on port %d", port)
//	registry, err := consul.NewRegistry("localhost:8500")
//	if err != nil {
//		panic(err)
//	}
//	ctx := context.Background()
//	instanceID := discovery.GenerateInstanceID(serviceName)
//	if err := registry.Register(
//		ctx,
//		instanceID,
//		serviceName,
//		fmt.Sprintf("localhost:%d", port),
//	); err != nil {
//		panic(err)
//	}
//	go func() {
//		for {
//			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
//				log.Println("Failed to report healthy state: " + err.Error())
//			}
//			time.Sleep(1 * time.Second)
//		}
//	}()
//	defer func(registry *consul.Registry, ctx context.Context, instanceID string, _ string) {
//		_ = registry.Deregister(ctx, instanceID, serviceName)
//	}(registry, ctx, instanceID, serviceName)
//
//	repo := memory.New()
//	ctrl := metadata.New(repo)
//	h := httphandler.New(ctrl)
//	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
//	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil); err != nil {
//		panic(err)
//	}
//}

func main() {
	log.Println("Starting movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	srv.Serve(lis)
}
