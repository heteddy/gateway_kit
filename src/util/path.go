// @Author : detaohe
// @File   : path.go
// @Description:
// @Date   : 2022/4/23 8:36 PM

package util

import "strings"

func JoinPath(a, b string) string {
	suffixA := strings.HasSuffix(a, "/")
	suffixB := strings.HasPrefix(b, "/")
	switch {
	case suffixA && suffixB:
		return a + b[1:]
	case !suffixA && !suffixB:
		return a + "/" + b
	}
	return a + b
}
