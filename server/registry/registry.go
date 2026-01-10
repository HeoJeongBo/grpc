package registry

import (
	"net/http"

	"grpc-server/auth"
	"grpc-server/database"
	"grpc-server/item"
	"grpc-server/user"
)

func RegisterAll(db *database.DB, mux *http.ServeMux) {
	auth.Register(db, mux)
	item.Register(db, mux)
	user.Register(db, mux)
}
