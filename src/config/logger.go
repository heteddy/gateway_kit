// @Author : detaohe
// @File   : logger
// @Description:
// @Date   : 2022/9/4 16:12

package config

import (
	"gateway_kit/util"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger(path, name string) (e error) {
	Logger, e = util.InitZapLogger(path, name)
	return
}
