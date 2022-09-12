// @Author : detaohe
// @File   : listener.go
// @Description:
// @Date   : 2022/9/12 00:16

package driver

import "gateway_kit/dao"

type HttpSvcListener interface {
	Update([]*dao.HttpSvcEntity)
}
