// @Author : detaohe
// @File   : loadbalance.go
// @Description:
// @Date   : 2022/4/23 8:32 PM

package lb

type Listener interface {
	Update(string, []string)
}

type LoadBalancer interface {
	// Next 通过一个svcName 获取真实的地址，后面改成不需要参数，
	Next(string) (string, error)
	Listener
}
