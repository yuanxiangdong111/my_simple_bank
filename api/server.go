package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "techschool/simplebank/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// 注册gin的验证标签currency validCurrency是自定义的验证函数
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validate.RegisterValidation("currency", validCurrency)
	}

	// 创建用户
	router.POST("/users", server.createUser)
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)
	server.router = router
	return server
}

// Start starts the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func ErrorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
