package api

import (
    "database/sql"
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"
    db "techschool/simplebank/db/sqlc"
)

// 创建账户的话, 金额是0
type createAccountRequest struct {
    Owner    string `json:"owner" binding:"required"`
    Currency string `json:"currency" binding:"required,oneof=USD EUR CNY"`
}

func (server *Server) createAccount(ctx *gin.Context) {
    var req createAccountRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    arg := db.CreateAccountParams{
        Owner:    req.Owner,
        Balance:  0,
        Currency: req.Currency,
    }
    account, err := server.store.CreateAccount(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
    ID int64 `uri:"id" binding:"required,min=1"`
}

// 通过id获取用户
func (server *Server) getAccount(ctx *gin.Context) {
    var req getAccountRequest
    if err := ctx.ShouldBindUri(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    account, err := server.store.GetAccount(ctx, req.ID)
    if err != nil {

        // 一种是没找到这个数据的错误
        if errors.Is(err, sql.ErrNoRows) {
            ctx.JSON(http.StatusNotFound, ErrorResponse(err))
            return
        }

        // 一种是其他错误
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
    PageID   int32 `form:"page_id" binding:"required,min=1"`
    PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(ctx *gin.Context) {
    var req listAccountRequest
    if err := ctx.ShouldBindQuery(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
        return
    }

    arg := db.ListAccountsParams{
        Limit:  req.PageSize,
        Offset: (req.PageID - 1) * req.PageSize,
    }
    accounts, err := server.store.ListAccounts(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, accounts)
}
