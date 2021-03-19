package infrastructure

import (
	"context"
	"encoding/json"
	"encoding/xml"
	constants "github.com/mvromanchak/emil-development/api-service"
	"github.com/mvromanchak/emil-development/api-service/entities"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/mvromanchak/emil-development/api-service/gps"
	"github.com/mvromanchak/emil-development/api-service/sdk"
	"github.com/pkg/errors"
)

// NewHTTPHandler creates a new HTTP handler for service endpoints.
func NewHTTPHandler(r *mux.Router, endpoints gps.Endpoints) {
	options := sdk.HTTPServerOptions()

	r.Methods(http.MethodPost).Path("/api/v1/gps").Handler(httptransport.NewServer(
		endpoints[gps.CreateGPSEndpoint],
		decodeHTTPBlockRequest,
		encodeHTTPBlockResponse,
		options...,
	))
}

func decodeHTTPBlockRequest(ctx context.Context, r *http.Request) (req interface{}, err error) {
	auth := r.Header.Get(constants.Authorization)
	if auth == "" {
		return nil, errors.New("error Authorization token is missing")
	}

	gpxRequest := entities.GPXRequestBody{}
	b := r.Body
	if b != nil {
		defer func() {
			if closeErr := b.Close(); closeErr != nil {
				err = errors.Wrap(closeErr, "failed to close request body")
			}
		}()
	}
	err = xml.NewDecoder(b).Decode(&gpxRequest)
	if err != nil {
		return nil, err
	}
	gpsr := entities.GPSDataRequest{
		DeviceID:   gpxRequest.Creator,
		JWTSecured: sdk.NewJWTSecured(),
	}
	return gpsr, nil
}

func encodeHTTPBlockResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(entities.GPXResponse)
	if !ok {
		panic("type error GPXResponse")
	}
	return encodeResponse(resp, w, http.StatusOK)
}

func encodeResponse(payload interface{}, w http.ResponseWriter, status int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	data, err := json.Marshal(&payload)
	if err != nil {
		return errors.Wrap(err, "encode gps response")
	}
	if _, err := w.Write(data); err != nil {
		return errors.Wrap(err, "gps writer write")
	}
	return nil
}
