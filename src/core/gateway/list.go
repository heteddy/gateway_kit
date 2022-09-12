// @Author : detaohe
// @File   : temporary
// @Description:
// @Date   : 2022/9/11 16:35

package gateway

import "sync"

type controlLayer struct {
	blockList []string
	allowList []string
}

type AccessControl struct {
	mutex        sync.RWMutex
	gatewayLayer *controlLayer
	serviceLayer *controlLayer
	temporary    *controlLayer
}
