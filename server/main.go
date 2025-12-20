package main

import (
	"log"
	"net/http"

	"grpc-server/database"
	"grpc-server/item"
	"grpc-server/proto-generated/item/v1/itemv1connect"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	connString := "host=localhost port=5432 user=postgres password=postgres dbname=grpc_dev sslmode=disable"

	db, err := database.New(connString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	server := item.NewItemServer(db)
	mux := http.NewServeMux()

	path, handler := itemv1connect.NewItemServiceHandler(server)
	mux.Handle(path, handler)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"*",
			"Connect-Protocol-Version",
			"Content-Type",
		},
		ExposedHeaders: []string{
			"Connect-Protocol-Version",
		},
	}).Handler(mux)

	addr := ":8080"
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, h2c.NewHandler(corsHandler, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
