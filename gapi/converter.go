package gapi

import (
	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
	"github.com/ryannguyen1105/Simplepayment/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(User db.User) *pb.User {
	return &pb.User{
		Username:          User.Username,
		Fullname:          User.FullName,
		Email:             User.Email,
		PasswordChangedAt: timestamppb.New(User.PasswordChangedAt),
		CreatedAt:         timestamppb.New(User.CreatedAt),
	}
}
