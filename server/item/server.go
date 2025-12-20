package item

import (
	"context"
	"fmt"
	"time"

	"grpc-server/database"
	"grpc-server/ent"
	"grpc-server/ent/item"

	itemv1 "grpc-server/proto-generated/item/v1"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	db *database.DB
}

func NewItemServer(db *database.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) CreateItem(
	ctx context.Context,
	req *connect.Request[itemv1.CreateItemRequest],
) (*connect.Response[itemv1.CreateItemResponse], error) {
	status := req.Msg.Status
	if status == itemv1.ItemStatus_ITEM_STATUS_UNSPECIFIED {
		status = itemv1.ItemStatus_ITEM_STATUS_DRAFT
	}

	entItem, err := s.db.Client.Item.
		Create().
		SetName(req.Msg.Name).
		SetDescription(req.Msg.Description).
		SetStatus(int32(status)).
		Save(ctx)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create item: %w", err))
	}

	return connect.NewResponse(&itemv1.CreateItemResponse{
		Item: entItemToProto(entItem),
	}), nil
}

func (s *Server) GetItem(
	ctx context.Context,
	req *connect.Request[itemv1.GetItemRequest],
) (*connect.Response[itemv1.GetItemResponse], error) {
	entItem, err := s.db.Client.Item.
		Query().
		Where(item.IDEQ(req.Msg.Id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("item not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get item: %w", err))
	}

	return connect.NewResponse(&itemv1.GetItemResponse{
		Item: entItemToProto(entItem),
	}), nil
}

func (s *Server) ListItems(
	ctx context.Context,
	req *connect.Request[itemv1.ListItemsRequest],
) (*connect.Response[itemv1.ListItemsResponse], error) {
	entItems, err := s.db.Client.Item.
		Query().
		Order(ent.Desc(item.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to list items: %w", err))
	}

	allItems := make([]*itemv1.Item, 0, len(entItems))
	for _, entItem := range entItems {
		allItems = append(allItems, entItemToProto(entItem))
	}

	filter := ApplyItemFilters(req.Msg.Filters)
	filteredItems := FilterItems(allItems, filter)

	return connect.NewResponse(&itemv1.ListItemsResponse{
		Items:      filteredItems,
		TotalCount: int32(len(filteredItems)),
	}), nil
}

func (s *Server) UpdateItem(
	ctx context.Context,
	req *connect.Request[itemv1.UpdateItemRequest],
) (*connect.Response[itemv1.UpdateItemResponse], error) {
	update := s.db.Client.Item.
		UpdateOneID(req.Msg.Id)

	if req.Msg.Name != nil {
		update = update.SetName(*req.Msg.Name)
	}

	if req.Msg.Description != nil {
		update = update.SetDescription(*req.Msg.Description)
	}

	if req.Msg.Status != nil {
		update = update.SetStatus(int32(*req.Msg.Status))
	}

	entItem, err := update.Save(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("item not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update item: %w", err))
	}

	return connect.NewResponse(&itemv1.UpdateItemResponse{
		Item: entItemToProto(entItem),
	}), nil
}

func (s *Server) DeleteItem(
	ctx context.Context,
	req *connect.Request[itemv1.DeleteItemRequest],
) (*connect.Response[itemv1.DeleteItemResponse], error) {
	err := s.db.Client.Item.
		DeleteOneID(req.Msg.Id).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("item not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete item: %w", err))
	}

	return connect.NewResponse(&itemv1.DeleteItemResponse{}), nil
}

func (s *Server) WatchItems(
	ctx context.Context,
	req *connect.Request[itemv1.WatchItemsRequest],
	stream *connect.ServerStream[itemv1.WatchItemsResponse],
) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			entItems, err := s.db.Client.Item.
				Query().
				Order(ent.Desc(item.FieldCreatedAt)).
				All(ctx)

			if err != nil {
				return fmt.Errorf("failed to query items: %w", err)
			}

			for _, entItem := range entItems {
				if err := stream.Send(&itemv1.WatchItemsResponse{
					Item:      entItemToProto(entItem),
					EventType: "UPDATE",
				}); err != nil {
					return err
				}
			}
		}
	}
}

func entItemToProto(entItem *ent.Item) *itemv1.Item {
	return &itemv1.Item{
		Id:          entItem.ID,
		Name:        entItem.Name,
		Description: entItem.Description,
		Status:      itemv1.ItemStatus(entItem.Status),
		CreatedAt:   timestamppb.New(entItem.CreatedAt),
		UpdatedAt:   timestamppb.New(entItem.UpdatedAt),
	}
}
