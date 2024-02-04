package api

import (
    "github.com/gin-gonic/gin"
    db "techschool/simplebank/db/sqlc"
)

type Server struct {
    store  db.Store
    router *gin.Engine
}

func NewServer(store db.Store) *Server {
    server := &Server{store: store}
    router := gin.Default()

    router.POST("/accounts", server.createAccount)
    router.GET("/accounts/:id", server.getAccount)
    router.GET("/accounts", server.listAccount)
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
