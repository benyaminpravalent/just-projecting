package hello

import (
	"context"

	pb "github.com/mine/just-projecting/proto/v1"
	"google.golang.org/grpc"
)

type Server struct {
	grpc *grpc.Server
}

func InitServer(grpc *grpc.Server) *Server {
	return &Server{grpc}
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: in.Name + " world"}, nil
}
