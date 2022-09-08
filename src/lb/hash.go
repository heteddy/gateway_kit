// @Author : detaohe
// @File   : hash.go
// @Description:
// @Date   : 2022/9/3 18:41

package lb

import (
	"errors"
	"sort"
	"sync"
)

type Hash func([]byte) uint32
type Int32Slice []uint32

type ConsistentHash struct {
	mux      sync.Mutex        // guards
	hashFunc Hash              // hash函数
	replicas int               // 复制多少份
	keys     Int32Slice        // 节点key
	hashMap  map[uint32]string // 节点hashmap
}

// Get 获取key所在节点信息
// todo 待测试
func (ch *ConsistentHash) Get(key string) (string, error) {
	if len(ch.hashMap) == 0 {
		return "", errors.New("ch map 为空")
	}
	value := ch.hashFunc([]byte(key))
	idx := sort.Search(len(ch.keys), func(i int) bool {
		return ch.keys[i] >= value
	})
	// hash节点
	if idx == len(ch.keys) {
		idx = 0
	}
	ch.mux.Lock()
	defer ch.mux.Unlock()
	return ch.hashMap[ch.keys[idx]], nil
}
