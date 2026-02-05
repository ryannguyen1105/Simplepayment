package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/wallets", server.createWallet)
	router.GET("/wallets/:id", server.getWallet)
	router.GET("/wallets", server.listWallet)
	router.PATCH("/wallets/:id", server.updateWallet)
	router.DELETE("/wallets/:id", server.deleteWallet)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
