package gapi

import (
	"fmt"

	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
	"github.com/ryannguyen1105/Simplepayment/pb"
	"github.com/ryannguyen1105/Simplepayment/token"
	"github.com/ryannguyen1105/Simplepayment/util"
)

type Server struct {
	pb.UnimplementedSimplePaymentServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server, nil
}
