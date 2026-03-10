package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/ryannguyen1105/Simplepayment/db/sqlc"
	"github.com/ryannguyen1105/Simplepayment/token"
)

type paymentRequest struct {
	FromWalletID int64 `json:"from_wallet_id" binding:"required,min=1"`
	ToWalletID   int64 `json:"to_wallet_id" binding:"required,min=1"`
	Amount       int64 `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createPayment(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	var req paymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fromWallet, err := server.store.GetWallet(ctx, req.FromWalletID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if fromWallet.Owner != authPayload.Username {
		err := errors.New("from wallet doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	_, err = server.store.GetWallet(ctx, req.ToWalletID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

func (server *Server) validWallet(ctx *gin.Context, walletID int64, status string) bool {
	wallet, err := server.store.GetPayment(ctx, walletID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}
	if wallet.Status != status {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", wallet.ID, wallet.Status, status)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
