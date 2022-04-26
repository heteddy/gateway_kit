// @Author : detaohe
// @File   : reverse.go
// @Description:
// @Date   : 2022/4/23 8:31 PM

package http

import (
	"bytes"
	"compress/gzip"
	"gateway_kit/service"
	"gateway_kit/util"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 初始化一个全局的transport
var transport = &http.Transport{
	Proxy:               nil,
	DialContext:         (&net.Dialer{}).DialContext,
	TLSHandshakeTimeout: 10 * time.Second,
	MaxIdleConns:        500,
	MaxIdleConnsPerHost: 50,
	MaxConnsPerHost:     100,
	IdleConnTimeout:     60 * time.Second,
}

//func main() {
//	url1, err1 := url.Parse("http://teddy3:9199")
//	if err1 != nil {
//		log.Println(err1)
//		return
//	}
//	p := NewReverseProxy([]*url.URL{url1})
//
//	log.Fatalln(http.ListenAndServe(":9079", p))
//}

func NewReverseProxy(l service.LoadBalance) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// 随机选择一个url
			// note api gateway的功能： url 改写 rewrite
			//r := rand.New(rand.NewSource(time.Now().UnixNano()))
			//n := r.Int63()
			//idx := n % int64(len(urls))
			//host := urls[idx]
			// /api-gateway/server2
			re, _ := regexp.Compile("^/api-gateway/(.*)")
			tmpPath := re.ReplaceAllString(req.URL.Path, "$1")
			// 是否修改host
			log.Println("req url=", req.URL.Path)
			if strings.Contains(req.URL.Path, "server3") {
				req.URL.Host = "192.168.64.7:9193"
				req.URL.Path = util.JoinUrl("/api/v1", tmpPath)
			} else {
				req.URL.Host = "192.168.64.6:9192"
				req.URL.Path = util.JoinUrl("/api/v1", tmpPath)
			}

			// 是否修改scheme
			req.URL.Scheme = "http"
			//log.Println("url:", host)
			if _, ok := req.Header["User-Agent"]; !ok {
				req.Header.Set("User-Agent", "teddy-api-gateway")
			}
		},
		Transport: transport,
		ModifyResponse: func(resp *http.Response) error {
			// todo 兼容websocket,这里返回nil 不支持
			if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
				return nil
			}
			var payload []byte
			var readErr error

			//todo 兼容gzip压缩
			if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
				bodyReader, err := gzip.NewReader(resp.Body)
				if err != nil {
					return err
				}
				payload, readErr = ioutil.ReadAll(bodyReader)
				resp.Header.Del("Content-Encoding")
			} else {
				payload, readErr = ioutil.ReadAll(resp.Body)
			}
			if readErr != nil {
				return readErr
			}

			//note api gateway的功能，错误码统一处理，异常请求时设置StatusCode
			if resp.StatusCode != 200 {
				payload = []byte("StatusCode error:" + string(payload))
			}

			//todo 因为预读了数据所以内容重新回写
			payload2 := []byte(string(payload) + " from api gateway")
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload2))
			resp.ContentLength = int64(len(payload2))
			resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload2)), 10))
			return nil
		},
		ErrorHandler: func(writer http.ResponseWriter, request *http.Request, err error) {

			http.Error(writer, "ErrorHandler error:"+err.Error(), 500)
		},
	}
}
