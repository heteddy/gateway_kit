// @Author : detaohe
// @File   : http.go
// @Description:
// @Date   : 2022/9/13 20:11

package endpoint

import (
	"gateway_kit/dao"
	"gateway_kit/server/service"
	"gateway_kit/util"
	"github.com/gin-gonic/gin"
)

type httpSvcCtrl struct {
	svc *service.ServiceSvr
}

type HttpSvcRequest struct {
	ID             string   `json:"id"`
	Name           string   `json:"name" binding:"name,required" validator:"min=3,max=10"`    // gateway name
	Description    string   `json:"description" binding:"required" validator:"min=0,max=120"` //描述
	BlockList      []string `json:"block_list" binding:"required"`                            // 网关黑名单，所有的服务通用
	AllowList      []string `json:"allow_list" binding:"required"`
	Addrs          []string `json:"addrs" binding:"required" validator:"min=1"`
	ClientQps      int      `json:"client_qps" binding:"required" validator:"min=1"` // 客户端流量控制
	ServerQps      int      `json:"server_qps" binding:"required" validator:"min=1"` // 服务端流量控制
	Category       int      `json:"category"  binding:"required" `                   // 如果gateway绑定多个域名，可以通过访问的host，来进行重定向
	MatchRule      string   `json:"match_rule" binding:"required" validator:"min=0"` // 匹配的项目与category结合使用，如果是domain，host==match_rule，否则是url前缀匹配
	IsHttps        bool     `json:"need_https" binding:"required" `
	IsWebsocket    bool     `json:"need_websocket" binding:"required" `
	StripUri       []string `json:"strip_uri" binding:"required" `   // 如果修改url可以通过gateway修改
	UrlRewrite     []string `json:"url_rewrite" binding:"required" ` // todo 需要支持正则表达式？,当修改了uri，可以对客户端保持兼容
	HeaderTransfer []string `json:"header_transfer" binding:"required"`
}

// List godoc
// @Summary 列表
// @Tags 服务
// @version 1.0
// @Accept application/json
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/services [get]
func (ctrl *httpSvcCtrl) List(c *gin.Context) {
	response := util.NewGinResponse(c)
	if entities, err := ctrl.svc.All(c.Request.Context()); err != nil {
		response.ToError(err)
	} else {
		response.ToResp(entities)
	}
}

// Get godoc
// @Summary 列表
// @Tags 服务
// @version 1.0
// @Accept application/json
// @Param id path string true "服务id"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/services/{id} [get]
func (ctrl *httpSvcCtrl) Get(c *gin.Context) {
	response := util.NewGinResponse(c)
	id := c.Param("id")
	if entities, err := ctrl.svc.GetByID(c.Request.Context(), id); err != nil {
		response.ToError(err)
	} else {
		response.ToResp(entities)
	}
}

// Create godoc
// @Summary 创建
// @Tags 服务
// @version 1.0
// @Accept application/json
// @Param gw body HttpSvcRequest true "服务信息"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/services [post]
func (ctrl *httpSvcCtrl) Create(c *gin.Context) {
	response := util.NewGinResponse(c)
	req := HttpSvcRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ToError(err)
		return
	}
	if err2 := util.DefaultValidator.Struct(req); err2 != nil {
		response.ToError(err2, "参数错误")
		return
	}

	newEntity := dao.HttpSvcEntity{
		Name:           req.Name,
		Description:    req.Description,
		BlockList:      req.BlockList,
		AllowList:      req.AllowList,
		Addr:           req.Addrs,
		ClientQps:      req.ClientQps,
		ServerQps:      req.ServerQps,
		Category:       req.Category,
		MatchRule:      req.MatchRule,
		IsHttps:        req.IsHttps,
		IsWebsocket:    req.IsWebsocket,
		StripUri:       req.StripUri,
		UrlRewrite:     req.UrlRewrite,
		HeaderTransfer: req.HeaderTransfer,
	}
	if entity, err := ctrl.svc.Create(c.Request.Context(), &newEntity); err != nil {
		response.ToError(err)
	} else {
		response.ToResp(entity)
	}
}

// Update godoc
// @Summary 更新
// @Tags 服务
// @version 1.0
// @Accept application/json
// @Param id path string true "服务id"
// @Param gw body HttpSvcRequest true "服务信息"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/services/{id} [put]
func (ctrl *httpSvcCtrl) Update(c *gin.Context) {
	response := util.NewGinResponse(c)
	req := HttpSvcRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ToError(err)
		return
	}
	if err2 := util.DefaultValidator.Struct(req); err2 != nil {
		response.ToError(err2, "参数错误")
		return
	}
	id := c.Param("id")
	newEntity := dao.HttpSvcEntity{
		Name:           req.Name,
		Description:    req.Description,
		BlockList:      req.BlockList,
		AllowList:      req.AllowList,
		Addr:           req.Addrs,
		ClientQps:      req.ClientQps,
		ServerQps:      req.ServerQps,
		Category:       req.Category,
		MatchRule:      req.MatchRule,
		IsHttps:        req.IsHttps,
		IsWebsocket:    req.IsWebsocket,
		StripUri:       req.StripUri,
		UrlRewrite:     req.UrlRewrite,
		HeaderTransfer: req.HeaderTransfer,
	}
	if entity, err := ctrl.svc.Update(c.Request.Context(), id, &newEntity); err != nil {
		response.ToError(err)
	} else {
		response.ToResp(entity)
	}
}

// Delete godoc
// @Summary 删除
// @Tags 服务
// @version 1.0
// @Accept application/json
// @Param id path string true "服务id"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/services/{id} [delete]
func (ctrl *httpSvcCtrl) Delete(c *gin.Context) {
	id := c.Param("id")
	response := util.NewGinResponse(c)
	err := ctrl.svc.Delete(c.Request.Context(), id)
	if err != nil {
		response.ToError(err, "")
	} else {
		response.ToResp("ok")
	}
}

func HttpSvcRouteRegister(group *gin.RouterGroup, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)
	ctrl := &httpSvcCtrl{
		svc: service.NewServiceSvr(),
	}
	prefixRouter := group
	if len(prefix) > 0 {
		prefixRouter = group.Group(prefix)
	}
	prefixRouter.POST("/services", ctrl.Create)
	prefixRouter.PUT("/services/:id", ctrl.Update)
	prefixRouter.DELETE("/services/:id", ctrl.Delete)
	prefixRouter.GET("/services", ctrl.List)
	prefixRouter.GET("/services/:id", ctrl.Get)
}
