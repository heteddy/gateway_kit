// @Author : detaohe
// @File   : server.go
// @Description:
// @Date   : 2022/4/23 9:50 PM

package svr

import (
	"context"
	"encoding/json"
	"gateway_kit/svc"
	"gateway_kit/svc/uppercase"
	"gateway_kit/util"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
)

type UppercaseRequest struct {
	S string `json:"s"`
}

type UppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

// MakeUpperCaseEndpoint
// @Description: 这里是每个controller
// @param uppercase
// @return endpoint.Endpoint
func MakeUpperCaseEndpoint(uppercase svc.Upper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*UppercaseRequest)
		return uppercase.Uppercase(ctx, req.S)
	}
}

func AddStringRoute(r *mux.Router, options ...httptransport.ServerOption) {
	v1 := r.PathPrefix("/api/v1/").Subrouter()
	v1.Methods(http.MethodPost).Path("/strings").Handler(
		// 这里可以继承这种server，然后通过自定义的route
		httptransport.NewServer(
			MakeUpperCaseEndpoint(uppercase.NewStringService()),
			decodeStringReq,
			util.EncodeHttpResp,
			options...,
		),
	)
}

func decodeStringReq(ctx context.Context, req *http.Request) (interface{}, error) {
	decoder := json.NewDecoder(req.Body)
	defer func() {
		req.Body.Close()
	}()
	var sReq UppercaseRequest
	if err := decoder.Decode(&sReq); err != nil {
		return nil, err
	}
	return &sReq, nil

}
