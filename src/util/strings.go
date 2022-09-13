// @Author : detaohe
// @File   : strings.go
// @Description:
// @Date   : 2022/9/12 20:00

package util

import "strings"

type IPSlice []string

func (ss IPSlice) Has(s string) bool {
	for _, _s := range ss {
		if strings.Contains(_s, "*") { //正则匹配
			sList := strings.Split(_s, "*")
			if len(sList) > 0 {
				return strings.HasPrefix(s, sList[0])
			}
		}
		if _s == s {
			return true
		}
	}
	return false
}
