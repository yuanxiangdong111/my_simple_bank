package api

import (
    "database/sql"
    "errors"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    db "techschool/simplebank/db/sqlc"
)

type transferRequest struct {
    FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
    ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
    Amount        int64  `json:"amount" binding:"required,gt=0"`
    Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
    var req transferRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    // 检查转账者的账户货币类型是否正确
    if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
        return
    }

    // 检查收款者的账户货币类型是否正确
    if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
        return
    }

    arg := db.TransferTxParams{
        FromAccountID: req.FromAccountID,
        ToAccountID:   req.ToAccountID,
        Amount:        req.Amount,
    }

    result, err := server.store.TransferTx(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
    account, err := server.store.GetAccount(ctx, accountID)
    if err != nil {

        if errors.Is(err, sql.ErrNoRows) {
            ctx.JSON(http.StatusNotFound, ErrorResponse(err))
            return false
        }

        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return false
    }

    if account.Currency != currency {
        err := fmt.Errorf("account [%d] currency is %s, but the transfer currency is %s", accountID, account.Currency, currency)
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return false
    }

    return true
}
