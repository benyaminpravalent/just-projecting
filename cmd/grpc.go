package cmd

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/mine/just-projecting/api/handler/grpc/hello"
	"github.com/mine/just-projecting/api/handler/grpc/user"
	"github.com/mine/just-projecting/domain/repository"
	"github.com/mine/just-projecting/pkg/config"
	pb "github.com/mine/just-projecting/proto/v1"
	"github.com/mine/just-projecting/service"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func Init() {
	g, _ := errgroup.WithContext(context.Background())
	var servers []*http.Server
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			log.Println("received signal: %s\n", sig)

			for i, s := range servers {
				if err := s.Shutdown(context.Background()); err != nil {
					if err == nil {
						log.Println("error shutting down server %d: %v", i, err)
						panic(err)
					}
				}
			}
			os.Exit(1)
		}
		return nil
	})

	// g.Go(func() error { return NewGrpcServer() })
	// g.Go(func() error { return NewHttpServer() })

	// err := g.Wait()
	// if err != nil {
	// 	panic(err)
	// }

	return
}

func NewGrpcServer() error {
	if err := config.Load(DefaultConfig, configURL); err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":"+config.GetString("grpc_port"))
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	userRepo := repository.NewUserRepository()
	// userRedis := redis.NewUserRedis(context.Background())

	userService := service.NewUserService().
		SetUserRepository(userRepo).
		Validate()

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, hello.InitServer(s))
	pb.RegisterUserServer(s, user.InitServer(s, userService))

	log.Println("Serving gRPC on 0.0.0.0:" + config.GetString("grpc_port"))
	s.Serve(lis)

	return nil
}

func NewHttpServer() error {
	if err := config.Load(DefaultConfig, configURL); err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.DialContext(context.Background(), "0.0.0.0:"+config.GetString("grpc_port"), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	err = pb.RegisterUserHandler(context.Background(), gwmux, conn)

	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":" + config.GetString("grpc_gateway_port"),
		Handler: gwmux,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:" + config.GetString("grpc_gateway_port"))
	log.Fatalln(gwServer.ListenAndServe())

	return nil
}
