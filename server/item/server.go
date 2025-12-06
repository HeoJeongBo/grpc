package item

import (
	"context"
	"fmt"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	itemv1 "github.com/heojeongbo/grpc/server/proto-generated/item/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	mu    sync.RWMutex
	items map[string]*itemv1.Item
}

func NewServer() *Server {
	return &Server{
		items: make(map[string]*itemv1.Item),
	}
}

func (s *Server) CreateItem(
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

func (s *Server) GetItem(
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

func (s *Server) ListItems(
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

func (s *Server) UpdateItem(
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

func (s *Server) DeleteItem(
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

func (s *Server) WatchItems(
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
