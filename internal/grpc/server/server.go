package server

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/bxrne/goforge/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	port string
}

type Config struct {
	Port string
}

func NewServer(config Config) *Server {
	return &Server{
		port: config.Port,
	}
}

func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// Placeholder implementation
	user := &pb.User{
		Id:        req.Id,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	return &pb.GetUserResponse{User: user}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Placeholder implementation
	user := &pb.User{
		Id:        1,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: "2024-01-01T00:00:00Z",
	}

	return &pb.CreateUserResponse{User: user}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Placeholder implementation
	users := []*pb.User{
		{Id: 1, Name: "John Doe", Email: "john@example.com", CreatedAt: "2024-01-01T00:00:00Z"},
		{Id: 2, Name: "Jane Smith", Email: "jane@example.com", CreatedAt: "2024-01-01T00:00:00Z"},
	}

	return &pb.ListUsersResponse{Users: users, Total: int32(len(users))}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, s)

	// Enable reflection for tools like grpcurl
	reflection.Register(grpcServer)

	log.Printf("gRPC server starting on port %s", s.port)
	return grpcServer.Serve(lis)
}

func (s *Server) StartWithContext(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	go func() {
		log.Printf("gRPC server starting on port %s", s.port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("gRPC server stopped")

	return nil
}
