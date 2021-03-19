package sdk

import (
	httptransport "github.com/go-kit/kit/transport/http"
)

// HTTPServerOptions provides the standard set of http server options to be used for any AddHTTPHandlers block.
func HTTPServerOptions() []httptransport.ServerOption {
	options := []httptransport.ServerOption{
		httptransport.ServerBefore(
			PopulateTokenRequestContext(),
		),
	}
	return options
}
