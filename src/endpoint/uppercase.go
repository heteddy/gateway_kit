// @Author : detaohe
// @File   : http.go
// @Description:
// @Date   : 2022/4/23 9:50 PM

package endpoint

import (
	"gateway_kit/svr/uppercase"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UppercaseRequest struct {
	S string `json:"s"`
}

type UppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"`
}

type uppercaseCtrl struct {
	svc *uppercase.StringService
}

func (ctrl uppercaseCtrl) Post(c *gin.Context) {
	var req UppercaseRequest
	if e := c.ShouldBindJSON(&req); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	ret, err := ctrl.svc.Uppercase(c.Request.Context(), req.S)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	} else {
		resp := UppercaseResponse{
			V:   ret,
			Err: "",
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
	}
}

func StringRouteReg(rg *gin.RouterGroup, prefix ...string) {
	var prefixGroup *gin.RouterGroup
	if len(prefix) > 0 {
		pre := prefix[0]
		prefixGroup = rg.Group(pre)
	} else {
		prefixGroup = rg
	}
	uc := uppercaseCtrl{
		svc: uppercase.NewStringService(),
	}

	prefixGroup.POST("uppercase", uc.Post)
}
