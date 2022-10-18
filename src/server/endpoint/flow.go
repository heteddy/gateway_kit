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

type FlowSvcSumReq struct {
	From string `json:"from" form:"from" binding:"required"`
	End  string `json:"end" form:"end" binding:"required"`
}

func (ctrl *FlowCtrl) GetServicesSum(c *gin.Context) {
	response := util.NewGinResponse(c)
	var req FlowSvcSumReq
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

type FlowSvcDetailReq struct {
	Service string `json:"service" form:"service"`
	From    string `json:"from" form:"from" binding:"required"`
	End     string `json:"end" form:"end" binding:"required"`
}

func (ctrl *FlowCtrl) GetServiceDetail(c *gin.Context) {
	response := util.NewGinResponse(c)
	var req FlowSvcDetailReq
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
	if entities, err := ctrl.svc.GetServiceDetail(c.Request.Context(), req.Service, from, end); err != nil {
		config.Logger.Error("svc GetServiceDetail", zap.Error(err))
		response.ToError(err)
	} else {
		if len(entities) == 0 {
			config.Logger.Warn("no service flow found")
		}
		for _, e := range entities {
			fmt.Println(e.Hour, e.Date, e.Name, e.Count)
		}
		response.ToResp(entities)
	}
}

func HttpSvcFlowRegister(group *gin.RouterGroup, prefixOptions ...string) {
	prefix := getPrefix(prefixOptions...)
	ctrl := &FlowCtrl{
		svc: service.NewFlowService(dao.NewServiceHourDao(), dao.NewRequestHourDao()),
	}
	prefixRouter := group
	if len(prefix) > 0 {
		prefixRouter = group.Group(prefix)
	}
	prefixRouter.GET("/flow-services", ctrl.GetServicesSum)
	prefixRouter.GET("/flow-service-details", ctrl.GetServiceDetail)
	//prefixRouter.DELETE("/services/:id", ctrl.Delete)
	//prefixRouter.GET("/services", ctrl.List)
	//prefixRouter.GET("/services/:id", ctrl.Get)
}
