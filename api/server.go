package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/token"
	"github.com/slamchillz/xchange/utils"
	"github.com/slamchillz/xchange/redisdb"
)

type Server struct {
	config utils.Config
	token *token.JWT
	router *gin.Engine
	storage db.Store
	redisClient redisdb.RedisClient
}

func NewServer(
	config utils.Config,
	storage db.Store,
	redisClient redisdb.RedisClient,
) (*Server, error) {
	jwt, err := token.NewJWT(config.JWT_SECRET)
	if err != nil {
		return nil, err
	}
	server := &Server{
		config: config,
		token: jwt,
		storage: storage,
		redisClient: redisClient,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("phonenumber", validatePhoneNumber)
		v.RegisterValidation("coinname", validateCoinName)
		v.RegisterValidation("network", validateNetwork)
		v.RegisterValidation("oldpassword", validateOldPassword)
	}
	server.ConfigRouter()
	return server, nil
}

func (server *Server) ConfigRouter() {
	router := gin.Default()

	apiRouter := router.Group("/api/v1")

	logApiRouter := apiRouter.Use(server.HTTPLogger)
	logApiRouter.POST("/user/signup", server.CreateCustomer)
	logApiRouter.POST("/user/login", server.LoginCustomer)
	logApiRouter.POST("/user/google/signup", server.GoogleSignUp)
	logApiRouter.POST("/user/google/signin", server.GoogleSignIn)

	authEndpoints := logApiRouter.Use(server.Authenticate)
	authEndpoints.POST("/token/swap", server.CoinSwap)
	authEndpoints.PATCH("/token/swap/:ref", server.CoinSwapStatusUpdate)
	authEndpoints.GET("/token/swap/history", server.ListCoinSwapHistory)
	authEndpoints.POST("/token/rate/calculate/ngn", server.GetCoinNGNEquivalent)
	authEndpoints.GET("/user/bank/details", server.GetBankDetails)
	authEndpoints.POST("/user/bank/details", server.AddBankDetails)
	authEndpoints.PATCH("/user/password/change", server.ChangePassword)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
