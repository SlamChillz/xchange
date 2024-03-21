package gapi

import (
	"github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/token"
	"github.com/slamchillz/xchange/utils"
	"github.com/slamchillz/xchange/pb"
)

type Server struct {
	pb.UnimplementedXchangeServer
	config utils.Config
	token *token.JWT
	storage db.Store
}

func NewServer(config utils.Config, storage db.Store) (*Server, error) {
	jwt, err := token.NewJWT(config.JWT_SECRET)
	if err != nil {
		return nil, err
	}
	server := &Server{
		UnimplementedXchangeServer: pb.UnimplementedXchangeServer{},
		config: config,
		token: jwt,
		storage: storage,
	}
	return server, nil
}
