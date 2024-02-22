package auth

import (
	"context"

	authv1 "github.com/tolseone/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

)

type Auth interface {
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	Login(ctx context.Context, email string, password string, app_id int) (token string, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	Logout(ctx context.Context, token string) (bool, error)
}

type serverAPI struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{
		auth: auth,
	})
}

const (
	emptyValue = 0
)

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if err := validateLogout(req); err != nil {
		return nil, err
	}

	success, err := s.auth.Logout(ctx, req.GetToken())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LogoutResponse{
		Success: success,
	}, nil
}

func validateRegister(req *authv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateLogin(req *authv1.LoginRequest) error {
	if req.GetEmail() == "" {
		status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func validateIsAdmin(req *authv1.IsAdminRequest) error {
	if req.GetUserId() == emptyValue {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}

func validateLogout(req *authv1.LogoutRequest) error {
	if req.GetToken() == "" {
		return status.Error(codes.InvalidArgument, "token is required")
	}

	return nil
}
