package auth

import (
	"context"
	"fmt"
	"strings"

	"grpc-server/database"
	"grpc-server/ent"
	"grpc-server/ent/user"
	"grpc-server/proto-generated/auth"
	"grpc-server/proto-generated/auth/authconnect"
	"grpc-server/registry"

	"net/http"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	registry.Register(func(db *database.DB, mux *http.ServeMux) {
		server := NewAuthServer(db)
		path, handler := authconnect.NewAuthServiceHandler(server)
		mux.Handle(path, handler)
	})

}

type Server struct {
	db *database.DB
}

func NewAuthServer(db *database.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Login(ctx context.Context, req *connect.Request[auth.LoginRequest]) (*connect.Response[auth.LoginResponse], error) {
	if err := ValidateEmail(req.Msg.Email); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if err := ValidatePassword(req.Msg.Password); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	entUser, err := s.db.Client.User.Query().Where(user.EmailEQ((req.Msg.Email))).Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
		}
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("failed to query user : %w", err))
	}

	if verifyErr := VerifyPassword(entUser.PasswordHash, req.Msg.Password); verifyErr != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("invalid password"))
	}

	tokenPair, err := GenerateTokenPair(entUser.ID, entUser.Email)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generated tokens : %w", err))
	}
	return connect.NewResponse(&auth.LoginResponse{
		User: entUserToProto(entUser),
		Tokens: &auth.TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresAt:    timestamppb.New(tokenPair.AccessTokenExpiry),
		},
	}), nil
}

func (s *Server) Logout(context.Context, *connect.Request[auth.LogoutRequest]) (*connect.Response[auth.LogoutResponse], error) {
	return connect.NewResponse(&auth.LogoutResponse{}), nil
}

func (s *Server) RefreshToken(ctx context.Context, req *connect.Request[auth.RefreshTokenRequest]) (*connect.Response[auth.RefreshTokenResponse], error) {
	claims, err := ValidateToken(req.Msg.RefreshToken)

	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("failed to generated tokens : %w", err))
	}

	tokenPair, err := GenerateTokenPair(claims.UserID, claims.Email)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generated tokens : %w", err))
	}

	return connect.NewResponse(&auth.RefreshTokenResponse{
		Tokens: &auth.TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresAt:    timestamppb.New(tokenPair.AccessTokenExpiry),
		},
	}), nil
}

func (s *Server) Register(ctx context.Context, req *connect.Request[auth.RegisterRequest]) (*connect.Response[auth.RegisterResponse], error) {
	if err := ValidateEmail(req.Msg.Email); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := ValidatePassword(req.Msg.Password); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	hashedPassword, err := HashPassword(req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to hash password: %w", err))
	}

	entUser, err := s.db.Client.User.
		Create().
		SetEmail(req.Msg.Email).
		SetName(req.Msg.Name).
		SetPasswordHash(hashedPassword).
		Save(ctx)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "UNIQUE constraint") {
			return nil, connect.NewError(connect.CodeAlreadyExists, fmt.Errorf("user with this email already exists"))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create user: %w", err))
	}

	tokenPair, err := GenerateTokenPair(entUser.ID, entUser.Email)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate tokens: %w", err))
	}

	return connect.NewResponse(&auth.RegisterResponse{
		User: entUserToProto(entUser),
		Tokens: &auth.TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
			ExpiresAt:    timestamppb.New(tokenPair.AccessTokenExpiry),
		},
	}), nil
}
