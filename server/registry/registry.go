package registry

import (
	"net/http"
	"sync"

	"grpc-server/database"
)

type HandlerRegistrar func(db *database.DB, mux *http.ServeMux)

var (
	mu         sync.RWMutex
	registrars []HandlerRegistrar
)

func Register(registrar HandlerRegistrar) {
	mu.Lock()
	defer mu.Unlock()
	registrars = append(registrars, registrar)
}

func RegisterAll(db *database.DB, mux *http.ServeMux) {
	mu.RLock()
	defer mu.RUnlock()

	for _, registrar := range registrars {
		registrar(db, mux)
	}
}
