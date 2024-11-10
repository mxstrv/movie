package main

import (
	"google.golang.org/grpc"
	"log"
	"movieapp/gen"
	"movieapp/rating/internal/controller/rating"
	grpchandler "movieapp/rating/internal/handler/grpc"
	"movieapp/rating/internal/repository/memory"
	"net"
)

//const serviceName = "rating"
//
//func main() {
//	var port int
//	flag.IntVar(&port, "port", 8082, "API handler port")
//	flag.Parse()
//	log.Printf("Starting Rating service on port %d", port)
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
//	repo := memory.New()
//	ctrl := rating.New(repo)
//	h := httphandler.New(ctrl)
//	http.Handle("/ratings", http.HandlerFunc(h.Handle))
//	if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil); err != nil {
//		panic(err)
//	}
//}

func main() {
	log.Println("Starting movie rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	srv.Serve(lis)
}
