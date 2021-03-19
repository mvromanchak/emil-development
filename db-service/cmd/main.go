package main

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	constant "github.com/mvromanchak/emil-development/db-service"
	"github.com/mvromanchak/emil-development/db-service/application"
	"github.com/mvromanchak/emil-development/db-service/application/repository"
	"github.com/mvromanchak/emil-development/db-service/config"
	"github.com/mvromanchak/emil-development/db-service/gps"
	"github.com/mvromanchak/emil-development/db-service/gps/transport"
	"github.com/mvromanchak/emil-development/db-service/sdk"
	"github.com/pressly/goose"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

func main() {
	// load config.
	v, err := sdk.LoadConfig()
	if err != nil {
		log.Printf("unable to load config %v", err)
		return
	}
	var settings config.Settings
	err = v.Unmarshal(&settings)
	if err != nil {
		log.Printf("unable to unmarshal config %v", err)
		return
	}

	// log panic if it happens.
	defer func() {
		if r := recover(); r != nil {
			log.Println("panic happened in the main", "panic info", r)
		}
	}()

	var (
		logger kitlog.Logger
		db     *sqlx.DB
	)

	// db connection data.
	dbConnection := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		settings.DbHost,
		settings.DbUsername,
		settings.DbPassword,
		settings.DbDataBase,
		settings.DbPort,
		settings.DbSSLMode)
	fmt.Println(settings.DbType, dbConnection)

	db, err = sqlx.Connect(settings.DbType, dbConnection)
	if err != nil {
		log.Println(err, "sqlx init failed")
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("error: failed to close DB: %v\n", err)
		}
	}()

	// run migration.
	log.Println("starting migration")
	gooseDB, err := goose.OpenDBWithDriver(settings.DbType, dbConnection)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}
	defer func() {
		if err := gooseDB.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()
	if goose.Up(gooseDB, "./migration/"); err != nil {
		log.Fatalf("goose migration err", err)
	}

	logger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)

	// prepare GRPC listener and server.
	grpcHost := "0.0.0.0"
	listener, err := net.Listen("tcp", net.JoinHostPort(grpcHost, "30000"))
	if err != nil {
		logger.Log("error while trying to listen %s port", "30000")
		return
	}

	grpcServer := initGRPCServer()
	rep := repository.NewRepository(db)
	app := application.NewGroupsService(rep)
	svc := gps.NewGroupsService(app)
	endpoints := gps.NewEndpoints(svc)
	transport.AddGRPCHandler(grpcServer, endpoints, logger)

	// Add server reflection
	reflection.Register(grpcServer)
	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			logger.Log("error while serving tcp listener")
			return
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	logger.Log("received system interruption", "signal", <-c)

	// close grpc server.
	if err := listener.Close(); err != nil && err != http.ErrServerClosed {
		log.Println(errors.Wrap(err, "gRPC server close"))
		return
	}

	log.Println("stop db-service")
}
func initGRPCServer() *grpc.Server {
	serverOptions := []grpc.ServerOption{
		grpc.MaxSendMsgSize(constant.GrpcMaxSendMsgSize),
		grpc.MaxRecvMsgSize(constant.GrpcMaxRecvMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     constant.GrpcMaxConnectionIdle,
			MaxConnectionAge:      constant.GrpcMaxConnectionAge,
			MaxConnectionAgeGrace: constant.GrpcMaxConnectionAgeGrace,
			Time:                  constant.GrpcPingTime,
			Timeout:               constant.GrpcTimeout,
		}),
		grpc.MaxConcurrentStreams(constant.MaxConcurrentStreams),
	}
	return grpc.NewServer(serverOptions...)
}
