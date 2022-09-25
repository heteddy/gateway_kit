// @Author : detaohe
// @File   : reverse_proxy.go
// @Description:
// @Date   : 2022/9/9 16:40

package gateway

import (
	"bytes"
	"compress/gzip"
	"gateway_kit/core/lb"
	"gateway_kit/util"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
)

type ProxyBuilder struct {
	transports *TransportPool
}

func NewProxyBuilder() *ProxyBuilder {
	return &ProxyBuilder{
		transports: NewTransportPool(),
	}
}

func (builder *ProxyBuilder) BuildDirector(balancer lb.LoadBalancer, svcName, scheme string) func(req *http.Request) {
	return func(req *http.Request) {
		// 随机选择一个url
		// note api gateway的功能： url 改写 rewrite
		// /api-gateway/server2
		re, _ := regexp.Compile("^/api-gateway/(.*)")
		urlPath := re.ReplaceAllString(req.URL.Path, "$1")
		//svcName := ServiceName(urlPath)
		// 是否修改host
		//hosts, err := repo.GetServices(svcName)
		log.Printf("servicename=%s \n", svcName)
		host, err := balancer.Next(svcName) // 这里需要返回一个对象，scheme， host都是从这里获取
		//log.Printf("loadbalance host=%s, error=%v\n", host, err)
		if err != nil {

		} else {
			// 是否修改scheme
			req.URL.Scheme = scheme // 通过配置获取scheme；是否可以转换https-> http wss->ws
			req.URL.Host = host
			// todo 这里可以做path改写，当降级的时候，直接改地址就可以了
			req.URL.Path = util.JoinPath("", urlPath)
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// 这里增加一个前缀，如果请求header不包括x-request-id 可以增加一个header
			req.Header.Set("User-Agent", "teddy-api-gateway")
		} else {
			req.Header.Add("User-Agent", "api-gateway")
		}
		for k, v := range req.Header {
			log.Printf("k=%s,v=%s\n", k, v)
		}
		log.Printf("new request to %s", req.URL.String())
	}
}

func (builder *ProxyBuilder) BuildModifier() func(req *http.Response) error {
	return func(resp *http.Response) error {
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
		// note api gateway的功能，错误码统一处理，异常请求时设置StatusCode
		if resp.StatusCode != 200 {
			payload = []byte("StatusCode error:" + string(payload))
		}
		//todo 因为预读了数据所以内容重新回写
		payload2 := []byte(string(payload) + " from api gateway\n")
		resp.Body = ioutil.NopCloser(bytes.NewBuffer(payload2))
		resp.ContentLength = int64(len(payload2))
		resp.Header.Set("Content-Length", strconv.FormatInt(int64(len(payload2)), 10))
		return nil
	}
}

func (builder *ProxyBuilder) BuildErrorHandler() func(http.ResponseWriter, *http.Request, error) {
	return func(writer http.ResponseWriter, request *http.Request, err error) {
		http.Error(writer, "ErrorHandler error:"+err.Error(), 500)
	}
}

func (builder *ProxyBuilder) BuildHttpProxy(balancer lb.LoadBalancer, svcName, scheme string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director:       builder.BuildDirector(balancer, svcName, scheme),
		Transport:      builder.transports.Get(svcName),
		ModifyResponse: builder.BuildModifier(),
		ErrorHandler:   builder.BuildErrorHandler(),
	}
}

func ServiceName(path string) string {
	idx := strings.Index(path, "/")
	var service = ""
	defer func() {
		log.Printf("path=%s, idx=%d, servicename=%s\n", path, idx, service)
	}()
	switch {
	case idx > 0:
		service = path[:idx]
	case idx == 0:
		idx := strings.Index(path[1:], "/")
		service = path[:idx]
	case idx < 0:
		service = path
	default:
		service = path
	}
	return service
}
