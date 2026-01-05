package auth

import (
	"context"
	"grpc-server/database"
	"grpc-server/proto-generated/auth"
	"grpc-server/proto-generated/auth/authconnect"
	"grpc-server/registry"
	"net/http"

	"connectrpc.com/connect"
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

// Login implements authconnect.AuthServiceHandler.
func (s *Server) Login(context.Context, *connect.Request[auth.LoginRequest]) (*connect.Response[auth.LoginResponse], error) {
	panic("unimplemented")
}

// Logout implements authconnect.AuthServiceHandler.
func (s *Server) Logout(context.Context, *connect.Request[auth.LogoutRequest]) (*connect.Response[auth.LogoutResponse], error) {
	panic("unimplemented")
}

// RefreshToken implements authconnect.AuthServiceHandler.
func (s *Server) RefreshToken(context.Context, *connect.Request[auth.RefreshTokenRequest]) (*connect.Response[auth.RefreshTokenResponse], error) {
	panic("unimplemented")
}

// Register implements authconnect.AuthServiceHandler.
func (s *Server) Register(context.Context, *connect.Request[auth.RegisterRequest]) (*connect.Response[auth.RegisterResponse], error) {
	panic("unimplemented")
}
