package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/slamchillz/xchange/api"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/gapi"
	"github.com/slamchillz/xchange/pb"
	"github.com/slamchillz/xchange/utils"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := utils.LoadConfig("./")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	if config.Env == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	fmt.Printf("%+v\n", config)
	conn, err := sql.Open(config.DBDriver, config.DBURL)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	err = runDatabaseMigration(config.MigrationURL, config.DBURL)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot run database migration")
	}
	store := db.NewStorage(conn)
	go runGRPCServer(config, store)
	runGatewayServer(config)
	// runHTTPServer(config, store)
}

func runDatabaseMigration(migrationURL, databaseURL string) error {
	m, err := migrate.New(migrationURL, databaseURL)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Info().Msg("database migration successful")
	return nil
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gapi server")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterXchangeServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener for GRPC server")
	}
	log.Printf("starting gRPC server at %s", config.GRPCServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

func runHTTPServer(config utils.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create HTTP server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start HTTP server")
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
		log.Fatal().Err(err).Msg("cannot register gateway server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcServerMux)

	fileServer := http.FileServer(http.Dir("./docs/gateway/swagger"))
	mux.Handle("/docs/", http.StripPrefix("/docs/", fileServer))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener for gateway server")
	}
	log.Printf("starting gateway server at %s", config.HTTPServerAddress)
	muxLogWrapper := gapi.HttpLogger(mux)
	err = http.Serve(listener, muxLogWrapper)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start gateway server")
	}
}
