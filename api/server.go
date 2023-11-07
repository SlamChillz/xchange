package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/utils"
)

type Server struct {
	config utils.Config
	router *gin.Engine
	storage db.Store
}

func NewServer(config utils.Config, storage db.Store) (*Server, error) {
	server := &Server{
		config: config,
		storage: storage,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phonenumber", validatePhoneNumber)
	}
	server.ConfigRouter()
	return server, nil
}

func (server *Server) ConfigRouter() {
	router := gin.Default()
	router.POST("/api/v1/swap", server.CoinSwap)
	router.POST("/api/v1/users/signup", server.CreateCustomer)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
