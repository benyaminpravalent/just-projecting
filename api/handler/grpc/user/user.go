package user

import (
	"context"
	"fmt"

	pb "github.com/mine/just-projecting/proto/v1"
	"github.com/mine/just-projecting/service"
	"google.golang.org/grpc"
)

type Server struct {
	grpc    *grpc.Server
	service service.UserService
}

func InitServer(grpc *grpc.Server, service service.UserService) *Server {
	return &Server{grpc, service}
}

func (s *Server) Get(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	fmt.Printf("\n Get Users: %v\n", in)
	users, err := s.service.Find()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("data not found")
	}

	var listUser []*pb.ModelUser
	for _, user := range users {
		userResp := &pb.ModelUser{
			Id:   int32(user.ID),
			Age:  int32(user.Age),
			Name: user.Name,
		}
		listUser = append(listUser, userResp)
	}

	return &pb.GetUserResponse{
		Status: 200,
		Data: &pb.GetUserResponse_Data{
			Users: listUser,
		},
	}, nil
}

func (s *Server) GetById(ctx context.Context, in *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	fmt.Printf("\n Get User By Id: %v\n", in)
	user, err := s.service.FindOne()
	if err != nil {
		return nil, err
	}
	return &pb.GetUserByIdResponse{
		Status: 200,
		Data: &pb.GetUserByIdResponse_Data{User: &pb.ModelUser{
			Id:   int32(user.ID),
			Age:  int32(user.Age),
			Name: user.Name,
		}},
	}, nil
}
