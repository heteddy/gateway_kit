// @Author : detaohe
// @File   : gateway.go
// @Description:
// @Date   : 2022/9/11 09:52

package endpoint

import (
	"gateway_kit/admin/service"
	"gateway_kit/dao"
	"gateway_kit/util"
	"github.com/gin-gonic/gin"
)

type gatewayCtrl struct {
	svc *service.GatewaySvc
}

type GatewayRequest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name" binding:"name,required" validator:"min=3,max=10"`    // gateway name
	Description string   `json:"description" binding:"required" validator:"min=0,max=127"` //描述
	BlockList   []string `json:"block_list" binding:"required"`                            // 网关黑名单，所有的服务通用
	AllowList   []string `json:"allow_list" binding:"required"`
}

// List godoc
// @Summary 列表
// @Tags gateway
// @version 1.0
// @Accept application/json
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/gateway/ [get]
func (ctrl *gatewayCtrl) List(c *gin.Context) {
	response := util.NewGinResponse(c)
	if entities, err := ctrl.svc.List(c.Request.Context()); err != nil {
		response.ToError(err)
	} else {
		response.ToResp(entities)
	}
}

// Create godoc
// @Summary 创建
// @Tags gateway
// @version 1.0
// @Accept application/json
// @Param gw body GatewayRequest true "网关信息"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/gateway [post]
func (ctrl *gatewayCtrl) Create(c *gin.Context) {
	response := util.NewGinResponse(c)
	req := GatewayRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ToError(err)
		return
	}
	if err2 := util.DefaultValidator.Struct(req); err2 != nil {
		response.ToError(err2, "参数错误")
		return
	}

	newEntity := dao.GatewayEntity{
		Name:        req.Name,
		Description: req.Description,
		BlockList:   req.BlockList,
		AllowList:   req.AllowList,
	}
	if entity, err := ctrl.svc.Create(c.Request.Context(), &newEntity); err != nil {
		response.ToError(err)
	} else {
		response.ToResp(entity)
	}
}

// Update godoc
// @Summary 更新
// @Tags gateway
// @version 1.0
// @Accept application/json
// @Param id path string true "网关id"
// @Param gw body GatewayRequest true "网关信息"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/gateway/{id} [put]
func (ctrl *gatewayCtrl) Update(c *gin.Context) {
	response := util.NewGinResponse(c)
	req := GatewayRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ToError(err)
		return
	}
	if err2 := util.DefaultValidator.Struct(req); err2 != nil {
		response.ToError(err2, "参数错误")
		return
	}
	id := c.Param("id")
	newEntity := dao.GatewayEntity{
		Name:        req.Name,
		Description: req.Description,
		BlockList:   req.BlockList,
		AllowList:   req.AllowList,
	}
	if entity, err := ctrl.svc.Update(c.Request.Context(), id, &newEntity); err != nil {
		response.ToError(err)
	} else {
		response.ToResp(entity)
	}
}

// Delete godoc
// @Summary 删除
// @Tags gateway
// @version 1.0
// @Accept application/json
// @Param id path string true "网关id"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/gateway/{id} [delete]
func (ctrl *gatewayCtrl) Delete(c *gin.Context) {
	id := c.Param("id")
	response := util.NewGinResponse(c)
	err := ctrl.svc.Delete(c.Request.Context(), id)
	if err != nil {
		response.ToError(err, "")
	} else {
		response.ToResp("ok")
	}
}

func getPrefix(prefixOptions ...string) string {
	prefix := ""
	if len(prefixOptions) > 0 {
		prefix = prefixOptions[0]
	}
	return prefix
}

func GatewayRouteRegister(group *gin.RouterGroup, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)
	ctrl := &gatewayCtrl{
		svc: service.NewGatewaySvc(),
	}
	prefixRouter := group
	if len(prefix) > 0 {
		prefixRouter = group.Group(prefix)
	}
	prefixRouter.POST("/gateway", ctrl.Create)
	prefixRouter.PUT("/gateway/:id", ctrl.Update)
	prefixRouter.DELETE("/gateway/:id", ctrl.Delete)
}
