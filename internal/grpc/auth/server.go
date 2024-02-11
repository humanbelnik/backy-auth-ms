package auth

import (
	"context"

	authv1 "github.com/humanbelnik/backy-contracts/codegen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implemented in ./services/auth/auth.go
type Auth interface {
	Login(ctx context.Context, loginString string, password string) (token string, err error)
	Logout(ctx context.Context, token string) (ok bool, err error)
	Register(ctx context.Context, email string, nickname string, password string) (userID int64, err error)
	Unregister(ctx context.Context, unregisterString string, password string, passowordConfirmed string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

type APIServer struct {
	authv1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthServer(gRPC, &APIServer{auth: auth})
}

func (s *APIServer) Unregister(ctx context.Context, req *authv1.UnregisterRequest) (*authv1.UnregisterResponse, error) {
	panic("implement me!")
}

func (s *APIServer) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	panic("implement me!")
}

func (s *APIServer) IsAdmin(ctx context.Context, req *authv1.IsAdminRequest) (*authv1.IsAdminResponse, error) {
	panic("implement me!")
}

func (s *APIServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetLoginString(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *APIServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	id, err := s.auth.Register(ctx, req.GetEmail(), req.GetNickname(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authv1.RegisterResponse{
		UserId: int64(id),
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// String in 'LoginString' can be either Email or Nickname.
// If there's an '@' symbol -> Email.
//
// Deal with LoginString when quering to DB.
func validateLogin(req *authv1.LoginRequest) error {
	if req.GetLoginString() == "" {
		return status.Error(codes.InvalidArgument, "nickname or email required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}

	return nil
}

func validateRegister(req *authv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "Email required")
	}

	if req.GetNickname() == "" {
		return status.Error(codes.InvalidArgument, "Nickname required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "Password required")
	}

	return nil
}
