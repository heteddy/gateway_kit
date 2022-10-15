// @Author : detaohe
// @File   : loadbalance.go
// @Description:
// @Date   : 2022/4/23 8:32 PM

package lb

type Node struct {
	Svc       string
	EventType int
	Addr      string
	Weight    int64
}

func (node *Node) Update(other Node) {
	//node.Addr = other.Addr
	node.Weight = other.Weight
}

// IsSameNode 判断是否是同一个node，
// note 如果是更新只能更新weight，不支持修改地址，如果必须修改地址，先删除再添加
func (node Node) IsSameNode(other Node) bool {
	//switch other.EventType {
	//case dao.EventDelete:
	//	return node.Svc == other.Svc && node.Addr == other.Addr && node.Weight == node.Weight
	//case dao.EventUpdate:
	//	return node.Svc == other.Svc
	//default:
	//
	//}
	return node.Svc == other.Svc && node.Addr == other.Addr
}

type LoadBalancer interface {
	// Next 通过一个svcName 获取真实的地址，后面改成不需要参数，
	Next(string) (string, error)
	UpdateNode(*Node)
}
