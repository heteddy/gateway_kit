// @Author : detaohe
// @File   : string.go
// @Description:
// @Date   : 2022/4/23 9:01 PM

package svr

import (
	"context"
	"errors"
	"strings"
)

// StringService provides operations on strings.
type StringService interface {
	Uppercase(context.Context, string) (string, error)
	Count(context.Context, string) int
}

type stringService struct{}

func NewStringSvc() StringService {
	return stringService{}
}

func (stringService) Uppercase(ctx context.Context, s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(ctx context.Context, s string) int {
	return len(s)
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// ServiceMiddleware is a chainable behavior modifier for StringService.
type ServiceMiddleware func(StringService) StringService
