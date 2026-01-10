package user

import (
	"context"
	"grpc-server/database"
	"grpc-server/proto-generated/user"
	"grpc-server/proto-generated/user/userconnect"
	"net/http"

	"connectrpc.com/connect"
)

func Register(db *database.DB, mux *http.ServeMux) {
	server := NewUserServer(db)
	path, handle := userconnect.NewUserServiceHandler(server)
	mux.Handle(path, handle)
}

type Server struct {
	db *database.DB
}

func NewUserServer(db *database.DB) *Server {
	return &Server{
		db: db,
	}
}

// DeleteUser implements userconnect.UserServiceHandler.
func (s *Server) DeleteUser(context.Context, *connect.Request[user.DeleteUserRequest]) (*connect.Response[user.DeleteUserResponse], error) {
	panic("unimplemented")
}

// GetUser implements userconnect.UserServiceHandler.
func (s *Server) GetUser(context.Context, *connect.Request[user.GetUserRequest]) (*connect.Response[user.GetUserResponse], error) {
	panic("unimplemented")
}

// UpdateUser implements userconnect.UserServiceHandler.
func (s *Server) UpdateUser(context.Context, *connect.Request[user.UpdateUserRequest]) (*connect.Response[user.UpdateUserResponse], error) {
	panic("unimplemented")
}
