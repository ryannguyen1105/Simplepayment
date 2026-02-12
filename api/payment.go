package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
)

type paymentRequest struct {
	FromWalletID int64 `json:"from_wallet_id" binding:"required,min=1"`
	ToWalletID   int64 `json:"to_wallet_id" binding:"required,min=1"`
	Amount       int64 `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createPayment(ctx *gin.Context) {
	var req paymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.PaymentTxParams{
		FromWalletID: req.FromWalletID,
		ToWalletID:   req.ToWalletID,
		Amount:       req.Amount,
	}
	result, err := server.store.PaymentTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}
