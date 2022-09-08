// @Author : detaohe
// @File   : string.go
// @Description:
// @Date   : 2022/8/30 18:01

package uppercase

import (
	"context"
	"strings"
)

type StringService struct {
}

func NewStringService() *StringService {
	return &StringService{}
}

func (svc StringService) Uppercase(ctx context.Context, src string) (string, error) {
	return strings.ToUpper(src), nil
}
