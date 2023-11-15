package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/slamchillz/xchange/api"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/gapi"
	"github.com/slamchillz/xchange/pb"
	"github.com/slamchillz/xchange/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStorage(conn)
	go runGRPCServer(config, store)
	runGatewayServer(config)
	// runHTTPServer(config, store)
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create gapi server:", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterXchangeServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener for GRPC server:", err)
	}
	log.Printf("starting gRPC server at %s", config.GRPCServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC server:", err)
	}
}

func runHTTPServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create HTTP server:", err)
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start HTTP server:", err)
	}
}

func runGatewayServer(config utils.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcServerMux := runtime.NewServeMux(opts)
	err := pb.RegisterXchangeHandlerFromEndpoint(ctx, grpcServerMux, config.GRPCServerAddress, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatal("cannot register gateway server:", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", grpcServerMux)
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create listener for gateway server:", err)
	}
	log.Printf("starting gateway server at %s", config.HTTPServerAddress)
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start gateway server:", err)
	}
}
