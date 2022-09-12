// @Author : detaohe
// @File   : gateway.go
// @Description:
// @Date   : 2022/9/11 09:52

package endpoint

import (
	"gateway_kit/admin"
	"github.com/gin-gonic/gin"
)

type gatewayCtrl struct {
	svc *admin.GatewaySvc
}

type GatewayRequest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`        // gateway name
	Description string   `json:"description"` //描述
	BlockList   []string `json:"block_list"`  // 网关黑名单，所有的服务通用
	AllowList   []string `json:"allow_list"`
}

func (gateway *gatewayCtrl) Create(c *gin.Context) {

}

func (gateway *gatewayCtrl) Update(c *gin.Context) {

}

func (gateway *gatewayCtrl) Delete(c *gin.Context) {

}
