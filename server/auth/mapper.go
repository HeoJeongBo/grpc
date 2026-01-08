package auth

import (
	"grpc-server/ent"
	"grpc-server/proto-generated/user"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func entUserToProto(u *ent.User) *user.User {
	return &user.User{
		Id:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}
