package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	itemv1 "github.com/heojeongbo/grpc/server/proto-generated/item/v1"
	"github.com/heojeongbo/grpc/server/proto-generated/item/v1/itemv1connect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type itemServer struct {
	mu    sync.RWMutex
	items map[string]*itemv1.Item
}

func newItemServer() *itemServer {
	return &itemServer{
		items: make(map[string]*itemv1.Item),
	}
}

func (s *itemServer) CreateItem(
	ctx context.Context,
	req *connect.Request[itemv1.CreateItemRequest],
) (*connect.Response[itemv1.CreateItemResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := timestamppb.New(time.Now())
	item := &itemv1.Item{
		Id:          uuid.New().String(),
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.items[item.Id] = item

	return connect.NewResponse(&itemv1.CreateItemResponse{
		Item: item,
	}), nil
}

func (s *itemServer) GetItem(
	ctx context.Context,
	req *connect.Request[itemv1.GetItemRequest],
) (*connect.Response[itemv1.GetItemResponse], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.items[req.Msg.Id]
	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("item not found"))
	}

	return connect.NewResponse(&itemv1.GetItemResponse{
		Item: item,
	}), nil
}

func (s *itemServer) ListItems(
	ctx context.Context,
	req *connect.Request[itemv1.ListItemsRequest],
) (*connect.Response[itemv1.ListItemsResponse], error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]*itemv1.Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return connect.NewResponse(&itemv1.ListItemsResponse{
		Items: items,
	}), nil
}

func (s *itemServer) UpdateItem(
	ctx context.Context,
	req *connect.Request[itemv1.UpdateItemRequest],
) (*connect.Response[itemv1.UpdateItemResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item, exists := s.items[req.Msg.Id]
	if !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("item not found"))
	}

	item.Name = req.Msg.Name
	item.Description = req.Msg.Description
	item.UpdatedAt = timestamppb.New(time.Now())

	return connect.NewResponse(&itemv1.UpdateItemResponse{
		Item: item,
	}), nil
}

func (s *itemServer) DeleteItem(
	ctx context.Context,
	req *connect.Request[itemv1.DeleteItemRequest],
) (*connect.Response[itemv1.DeleteItemResponse], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[req.Msg.Id]; !exists {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("item not found"))
	}

	delete(s.items, req.Msg.Id)

	return connect.NewResponse(&itemv1.DeleteItemResponse{}), nil
}

func (s *itemServer) WatchItems(
	ctx context.Context,
	req *connect.Request[itemv1.WatchItemsRequest],
	stream *connect.ServerStream[itemv1.WatchItemsResponse],
) error {
	// Demo: Send current items every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			s.mu.RLock()
			for _, item := range s.items {
				if err := stream.Send(&itemv1.WatchItemsResponse{
					Item:      item,
					EventType: "UPDATE",
				}); err != nil {
					s.mu.RUnlock()
					return err
				}
			}
			s.mu.RUnlock()
		}
	}
}

func main() {
	server := newItemServer()
	mux := http.NewServeMux()

	path, handler := itemv1connect.NewItemServiceHandler(server)
	mux.Handle(path, handler)

	// Setup CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}).Handler(mux)

	addr := ":8080"
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, h2c.NewHandler(corsHandler, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
