package auth

import (
	"context"

	"google.golang.org/grpc"

	authv1 "github.com/tolseone/protos/gen/go/auth"

)

type serverAPI struct {
	authv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	panic("implement me")
}
