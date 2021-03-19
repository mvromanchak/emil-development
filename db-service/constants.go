package polynom

import (
	"math"
	"time"
)

const (
	GrpcMaxRecvMsgSize        = 1024 * 1024 * 500
	GrpcMaxSendMsgSize        = 1024 * 1024 * 500
	GrpcMaxConnectionIdle     = time.Minute * 5
	GrpcMaxConnectionAge      = time.Duration(math.MaxInt64)
	GrpcMaxConnectionAgeGrace = time.Duration(math.MaxInt64)
	GrpcPingTime              = time.Hour * 2
	GrpcTimeout               = time.Second * 30
	MaxConcurrentStreams      = 100
)
