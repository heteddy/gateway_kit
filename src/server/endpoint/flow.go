// @Author : detaohe
// @File   : flow
// @Description:
// @Date   : 2022/10/18 10:46

package endpoint

import (
	"errors"
	"fmt"
	"gateway_kit/config"
	"gateway_kit/dao"
	"gateway_kit/server/service"
	"gateway_kit/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

type FlowCtrl struct {
	svc *service.FlowService
}

type FlowSumReq struct {
	From string `json:"from" form:"from" binding:"required"`
	End  string `json:"end" form:"end" binding:"required"`
}

// GetGwSum godoc
// @Summary 列表
// @Tags 统计
// @version 1.0
// @Accept application/json
// @Param from query string true "开始时间"
// @Param end query string true "结束时间"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/flow-gateway/ [get]
func (ctrl *FlowCtrl) GetGwSum(c *gin.Context) {
	response := util.NewGinResponse(c)
	var req FlowSumReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ToError(err)
		return
	}
	if len(req.From) != len("2006-01-02") || len(req.End) != len("2006-01-02") {
		response.ToError(errors.New("时间格式错误"))
		return
	}
	from, err2 := time.Parse("2006-01-02", req.From)
	end, err3 := time.Parse("2006-01-02", req.End)
	if err2 != nil || err3 != nil {
		response.ToError(errors.New("时间格式错误"))
		return
	}
	if entities, err := ctrl.svc.GetGwSum(c.Request.Context(), from, end); err != nil {
		response.ToError(err)
	} else {
		response.ToResp(entities)
	}
}

//type FlowDetailReq struct {
//	Name string `json:"service" form:"service"`
//	From    string `json:"from" form:"from" binding:"required"`
//	End     string `json:"end" form:"end" binding:"required"`
//}

// GetGwDetail godoc
// @Summary 列表
// @Tags 统计
// @version 1.0
// @Accept application/json
// @Param from query string true "开始时间"
// @Param end query string true "结束时间"
// @Param name query string true "服务名"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/flow-gw-details/ [get]
func (ctrl *FlowCtrl) GetGwDetail(c *gin.Context) {
	response := util.NewGinResponse(c)
	var req FlowDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ToError(err)
		return
	}
	if len(req.From) != len("2006-01-02") || len(req.End) != len("2006-01-02") {
		response.ToError(errors.New("时间格式错误"))
		return
	}
	from, err2 := time.Parse("2006-01-02", req.From)
	end, err3 := time.Parse("2006-01-02", req.End)
	if err2 != nil || err3 != nil {
		response.ToError(errors.New("时间格式错误"))
		return
	}
	if entities, err := ctrl.svc.GetGwDetail(c.Request.Context(), req.Name, from, end); err != nil {
		config.Logger.Error(" GetGwDetail", zap.Error(err))
		response.ToError(err)
	} else {
		if len(entities) == 0 {
			config.Logger.Warn("no gw flow found")
		}
		for _, e := range entities {
			fmt.Println(e.Date, e.Name, e.Count)
		}
		response.ToResp(entities)
	}
}

// GetServicesSum godoc
// @Summary 列表
// @Tags 统计
// @version 1.0
// @Accept application/json
// @Param from query string true "开始时间"
// @Param end query string true "结束时间"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/flow-services/ [get]
func (ctrl *FlowCtrl) GetServicesSum(c *gin.Context) {
	response := util.NewGinResponse(c)
	var req FlowSumReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ToError(err)
		return
	}
	if len(req.From) != len("2006-01-02") || len(req.End) != len("2006-01-02") {
		response.ToError(errors.New("时间格式错误"))
		return
	}
	from, err2 := time.Parse("2006-01-02", req.From)
	end, err3 := time.Parse("2006-01-02", req.End)
	if err2 != nil || err3 != nil {
		response.ToError(errors.New("时间格式错误"))
		return
	}
	if entities, err := ctrl.svc.GetServicesSum(c.Request.Context(), from, end); err != nil {
		response.ToError(err)
	} else {
		response.ToResp(entities)
	}
}

type FlowDetailReq struct {
	Name string `json:"name" form:"name"`
	From string `json:"from" form:"from" binding:"required"`
	End  string `json:"end" form:"end" binding:"required"`
}

// GetServiceDetail godoc
// @Summary 列表
// @Tags 统计
// @version 1.0
// @Accept application/json
// @Param from query string true "开始时间"
// @Param end query string true "结束时间"
// @Param name query string true "服务名"
// @Success 200 {object} string
// @Failure 200 {object} string
// @Router /gateway-kit-svr/flow-service-details/ [get]
func (ctrl *FlowCtrl) GetServiceDetail(c *gin.Context) {
	response := util.NewGinResponse(c)
	var req FlowDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ToError(err)
		return
	}
	if len(req.From) != len("2006-01-02") || len(req.End) != len("2006-01-02") {
		response.ToError(errors.New("时间格式错误"))
		return
	}
	from, err2 := time.Parse("2006-01-02", req.From)
	end, err3 := time.Parse("2006-01-02", req.End)
	if err2 != nil || err3 != nil {
		response.ToError(errors.New("时间格式错误"))
		return
	}
	if entities, err := ctrl.svc.GetServiceDetail(c.Request.Context(), req.Name, from, end); err != nil {
		config.Logger.Error("svc GetServiceDetail", zap.Error(err))
		response.ToError(err)
	} else {
		if len(entities) == 0 {
			config.Logger.Warn("no service flow found")
		}
		for _, e := range entities {
			fmt.Println(e.Date, e.Name, e.Count)
		}
		response.ToResp(entities)
	}
}

func HttpSvcFlowRegister(group *gin.RouterGroup, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)
	ctrl := &FlowCtrl{
		svc: service.NewFlowService(dao.NewServiceHourDao(), dao.NewRequestHourDao(), dao.NewServiceDayDao()),
	}
	prefixRouter := group
	if len(prefix) > 0 {
		prefixRouter = group.Group(prefix)
	}
	prefixRouter.GET("/flow-services", ctrl.GetServicesSum)
	prefixRouter.GET("/flow-service-details", ctrl.GetServiceDetail)
	prefixRouter.GET("/flow-gateway", ctrl.GetGwSum)
	prefixRouter.GET("/flow-gw-details", ctrl.GetGwDetail)
	//prefixRouter.DELETE("/services/:id", ctrl.Delete)
	//prefixRouter.GET("/services", ctrl.List)
	//prefixRouter.GET("/services/:id", ctrl.Get)
}
