// @Author : detaohe
// @File   : svc
// @Description:
// @Date   : 2022/8/30 18:04

package svc

import "context"

type Upper interface {
	Uppercase(context.Context, string) (string, error)
}
