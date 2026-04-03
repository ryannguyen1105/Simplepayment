package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
	"github.com/ryannguyen1105/Simplepayment/token"
	"github.com/ryannguyen1105/Simplepayment/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	router.DELETE("/users/delete", server.deleteUser)
	authRoutes.POST("/wallets", server.createWallet)
	authRoutes.GET("/wallets/:id", server.getWallet)
	authRoutes.GET("/wallets", server.listWallet)
	authRoutes.PATCH("/wallets/:id", server.updateWallet)
	authRoutes.DELETE("/wallets/:id", server.deleteWallet)
	authRoutes.POST("/payments", server.createPayment)

	server.router = router

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
