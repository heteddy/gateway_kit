// @Author : detaohe
// @File   : path.go
// @Description:
// @Date   : 2022/4/23 8:36 PM

package util

import "strings"

func JoinUrl(a, b string) string {
	suffixA := strings.HasSuffix(a, "/")
	suffixB := strings.HasPrefix(b, "/")
	switch {
	case suffixA && suffixB:
		return a + b[1:]
	case !suffixA && !suffixB:
		return a + "/" + b
	}
	//log.Printf("singleJoiningSlash a=%s,b=%s,a+b=%s \n", a, b, a+b)
	return a + b
}
