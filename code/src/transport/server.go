// @Author : detaohe
// @File   : server.go
// @Description:
// @Date   : 2022/4/23 9:01 PM

package transport

import (
	"gateway_kit/endpoint/svr"
	"github.com/gorilla/mux"
	"net/http"
)

// MakeHttpHandler
// @Description:
// @return http.Handler
//
func MakeHttpHandler() http.Handler {
	router := mux.NewRouter()
	svr.AddRoute(router)
	router.Methods(http.MethodGet).Path("/healthz").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("ok"))
	})
	return router
}
