// @Author : detaohe
// @File   : string.go
// @Description:
// @Date   : 2022/8/30 18:01

package uppercase

import (
	"context"
	"strings"
)

type stringService struct {
}

func NewStringService() *stringService {
	return &stringService{}
}

func (svc stringService) Uppercase(ctx context.Context, src string) (string, error) {
	return strings.ToUpper(src), nil
}
