package sdk

import (
	"context"
	httpkit "github.com/go-kit/kit/transport/http"
	constants "github.com/mvromanchak/emil-development/api-service"
	"net/http"
)

const Token = "api-token"

// authorization mdv
func PopulateTokenRequestContext() httpkit.RequestFunc {
	return func(ctx context.Context, req *http.Request) context.Context {
		auth := req.Header.Get(constants.Authorization)
		token := ExtractTokenFromBearer(auth)
		ctx = context.WithValue(ctx, interface{}(constants.JWT), token)
		return ctx
	}
}
