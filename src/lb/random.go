// @Author : detaohe
// @File   : random.go
// @Description:
// @Date   : 2022/9/3 18:55

package lb

import (
	"math/rand"
	"time"
)

type randomLB struct {
}

func NewRandomLB() LoadBalancer {
	return &randomLB{}
}
func (lb *randomLB) Add(params ...string) error {
	return nil
}
func (lb *randomLB) Update() {

}
func (lb *randomLB) Next(hosts []string) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Int63()
	idx := n % int64(len(hosts))
	return hosts[idx], nil
}
