package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
)

type CreateWalletRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

func (server *Server) CreateWallet(ctx *gin.Context) {
	var req CreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorReponse(err))
		return
	}

	arg := db.CreateWalletParams{
		UserID:  req.UserID,
		Balance: 0,
	}
	wallet, err := server.store.CreateWallet(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorReponse(err))
		return
	}
	ctx.JSON(http.StatusOK, wallet)
}
