package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	// "os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/slamchillz/xchange/api"
	db "github.com/slamchillz/xchange/db/sqlc"
	"github.com/slamchillz/xchange/gapi"
	"github.com/slamchillz/xchange/pb"
	"github.com/slamchillz/xchange/utils"
	"github.com/slamchillz/xchange/redisdb"
	// "github.com/rs/zerolog"
	// "github.com/rs/zerolog/logger"
	log "github.com/slamchillz/xchange/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var logger = log.GetLogger()

func main() {
	config, err := utils.LoadConfig("/home/ubuntu/xchange")
	if err != nil {
		logger.Fatal().Err(err).Stack().Msg("cannot load config")
	}
	if config.Env == "dev" {
		// logger.Logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		logger.Info().Msg("starting server in development mode")
		// logger.Info().Msgf("config: %+v", config)
	}
	conn, err := sql.Open(config.DBDriver, config.DBURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot connect to db")
	}
	err = runDatabaseMigration(config.MigrationURL, config.DBURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot run database migration")
	}
	store := db.NewStorage(conn)
	rdConn := redis.NewClient(&redis.Options{
		Addr:	  config.RedisAddress,
		Password: config.RedisPassword, // no password set
		DB:		  config.RedisDB,  // use default DB
	})
	redisClient := redisdb.NewRedisClient(rdConn)
	// go runGRPCServer(config, store)
	// runGatewayServer(config)
	runHTTPServer(config, store, redisClient)
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
	logger.Info().Msg("database migration successful")
	return nil
}

func runGRPCServer(config utils.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create gapi server")
	}
	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterXchangeServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create listener for GRPC server")
	}
	logger.Info().Msgf("starting gRPC server at %s", config.GRPCServerAddress)
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot start gRPC server")
	}
}

func runHTTPServer(config utils.Config, store db.Store, redisClient redisdb.RedisClient) {
	server, err := api.NewServer(config, store, redisClient)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create HTTP server")
	}
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot start HTTP server")
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
		logger.Fatal().Err(err).Msg("cannot register gateway server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcServerMux)

	fileServer := http.FileServer(http.Dir("./docs/gateway/swagger"))
	mux.Handle("/docs/", http.StripPrefix("/docs/", fileServer))

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create listener for gateway server")
	}
	logger.Info().Msgf("starting gateway server at %s", config.HTTPServerAddress)
	muxLogWrapper := gapi.HttpLogger(mux)
	err = http.Serve(listener, muxLogWrapper)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot start gateway server")
	}
}
