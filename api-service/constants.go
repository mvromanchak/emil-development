package polynom

import "time"

const (
	//HTTPServerReadTimeout server read timeout
	HTTPServerReadTimeout = time.Second * 60
	//HTTPServerWriteTimeout server write timeout
	HTTPServerWriteTimeout = time.Second * 60

	JWT           = "jwt-token"
	Authorization = "authorization"
)
