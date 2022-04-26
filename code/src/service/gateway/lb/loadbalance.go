// @Author : detaohe
// @File   : loadbalance.go
// @Description:
// @Date   : 2022/4/23 8:32 PM

package lb

type LoadBalancer interface {
	// GetService 通过一个src地址 获取真实的代理
	GetService(string) (string, error)
}

type WeightRoundRobinLB struct {
}

func (lb *WeightRoundRobinLB) GetService(src string) (string, error) {
	return "", nil
}
