package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	constants "github.com/mvromanchak/emil-development/api-service"
	"github.com/mvromanchak/emil-development/api-service/application"
	"github.com/mvromanchak/emil-development/api-service/config"
	"github.com/mvromanchak/emil-development/api-service/gps"
	"github.com/mvromanchak/emil-development/api-service/gps/infrastructure"
	"github.com/mvromanchak/emil-development/api-service/sdk"
	"github.com/mvromanchak/emil-development/api-service/sdk/helper"
	"github.com/openzipkin/zipkin-go"
	zreporter "github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
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
		httpRouter  *mux.Router
		httpHandler http.Handler
	)

	httpRouter = mux.NewRouter()
	httpHandler = httpRouter

	srv := &http.Server{
		ReadTimeout:  constants.HTTPServerReadTimeout,
		WriteTimeout: constants.HTTPServerWriteTimeout,
		Handler:      timeoutHandler(httpHandler, time.Second*time.Duration(60)),
		Addr:         net.JoinHostPort(settings.Host, settings.Port),
	}

	var dbGRPCConn *grpc.ClientConn

	// make createGRPCconnectiont. //if will tot work on localhost.
	dbGRPCConn, err = createGRPCconnection(net.JoinHostPort("db.service", "30000"))
	if err != nil {
		log.Println("failed to init gRPC")
		return
	}

	gpsClient, err := infrastructure.NewGPSGRPCClient(dbGRPCConn)
	if err != nil {
		panic(err)
	}

	hostAddress, err := helper.GetLocalIPaddress()
	if err != nil {
		error.Error(errors.Wrap(err, "can't get ip address"))
	}

	// zipkin distributive tracing
	var zipkinTracer *zipkin.Tracer
	var zipkinReporter zreporter.Reporter
	zipkinTracer, zipkinReporter, err = initZipkinTracer(settings.ZipkinHost, settings.ZipkinPort, hostAddress, "5686")
	if err != nil {
		log.Println(err)
	} else {
		defer func() {
			err = zipkinReporter.Close()
			if err != nil {
				errors.Wrap(err, "zipkin reporter close failed")
			}
		}()
	}
	fmt.Println(zipkinTracer)

	app := application.NewGPSService(gpsClient)
	svc := gps.NewGPSService(app)
	endp := gps.NewEndpoints(svc)
	infrastructure.NewHTTPHandler(httpRouter, endp)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// start http server.
	go func() {
		log.Printf("server starts on %s : %s", settings.Host, settings.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(errors.Wrap(err, "http server failed"))
			return
		}
	}()

	log.Println("received system interruption", "signal", <-c)
	// stop the service.
	if err := srv.Shutdown(context.Background()); err != nil && err != http.ErrServerClosed {
		log.Println(errors.Wrap(err, "server shutdown"))
		return
	}
	// stop the gRPC client.
	if err := dbGRPCConn.Close(); err != nil && err != http.ErrServerClosed {
		log.Println(errors.Wrap(err, "gRPC client close"))
		return
	}
	log.Println("stop api-service")
}

func timeoutHandler(h http.Handler, duration time.Duration) http.Handler {
	m := map[string]interface{}{
		"error": "Service unavailable timeout abort",
		"code":  "errors.TimeoutErrorCode",
	}
	data, err := json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	return http.TimeoutHandler(h, duration, string(data))
}

func createGRPCconnection(address string) (*grpc.ClientConn, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 10,
			Timeout:             time.Second * 10,
			PermitWithoutStream: true,
		}),
	}
	return grpc.Dial(address, dialOptions...)
}

func initZipkinTracer(zipkinHost, zipkinPort, httpHost, httpPort string) (*zipkin.Tracer, zreporter.Reporter, error) {
	var zipkinURL string
	zipkinURL = "http://" + net.JoinHostPort(zipkinHost, zipkinPort) + "/api/v2/spans"

	var (
		err      error
		hostPort = net.JoinHostPort(httpHost, httpPort)
		reporter = zipkinhttp.NewReporter(zipkinURL)
	)
	zEP, err := zipkin.NewEndpoint("api-service", hostPort)
	if err != nil {
		return nil, nil, errors.Wrap(err, "init zipkin endpoint failed")
	}

	// Sampler tells you which traces are going to be sampled or not. 1.00 -> 100%
	sampler, err := zipkin.NewCountingSampler(1.00)
	if err != nil {
		return nil, nil, errors.Wrap(err, "sampler init fail")
	}

	zipkinTracer, err := zipkin.NewTracer(
		reporter, zipkin.WithLocalEndpoint(zEP), zipkin.WithNoopTracer(false), zipkin.WithSampler(sampler),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "zipkin init fail")
	}
	return zipkinTracer, reporter, nil
}
