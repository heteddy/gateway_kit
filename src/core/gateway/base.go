// @Author : detaohe
// @File   : base.go
// @Description:
// @Date   : 2022/9/13 15:49

package gateway

import "sync"

type base struct {
	mutex sync.RWMutex
	stopC chan struct{}
}
