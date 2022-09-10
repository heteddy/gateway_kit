// @Author : detaohe
// @File   : weighted.go
// @Description:
// @Date   : 2022/9/10 11:56

package lb

import (
	"errors"
	"strconv"
)

type weightedRoundRobinLB struct {
	nodes []*weightNode
}

/*
1. 对于每个请求，遍历集群中的所有可用后端，对于每个后端peer执行：
peer->current_weight += peer->effecitve_weight。
同时累加所有peer的effective_weight，保存为total。
2. 从集群中选出current_weight最大的peer，作为本次选定的后端。
3. 对于本次选定的后端，执行：peer->current_weight -= total。

*/
type weightNode struct {
	Addr   string
	Weight int64 //配置文件中指定的该后端的权重，这个值是固定不变的。
	/*
		后端目前的权重，一开始为0，之后会动态调整。动态调整：
		每次选取后端时，会遍历集群中所有后端，对于每个后端，current_weight增加effective_weight，
		同时累加所有后端的effective_weight，保存为total。
		如果该后端的current_weight是最大的，就选定这个后端，然后把它的current_weight减去total。
		如果该后端没有被选定，那么current_weight不用减小。
	*/
	Current int64
	/*
		后端的有效权重，初始值为weight。
		在释放后端时，如果发现和后端的通信过程中发生了错误，就减小effective_weight。
		此后有新的请求过来时，在选取后端的过程中，再逐步增加effective_weight，最终又恢复到weight。
		之所以增加这个字段，是为了当后端发生错误时，降低其权重
	*/
	Effective int64
}

//func (lb *weightedRoundRobinLB) Add(node *weightNode) {
//	// 构造weight node
//	node.Effective = node.Weight
//	lb.nodes = append(lb.nodes, node)
//}

func (lb *weightedRoundRobinLB) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("参数错误")
	}
	v, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	node := &weightNode{
		Addr:      params[0],
		Weight:    v,
		Current:   0,
		Effective: v,
	}
	// 构造weight node
	//node.Effective = node.Weight
	lb.nodes = append(lb.nodes, node)
	return nil
}
func (lb *weightedRoundRobinLB) Update() {

}

// Next
// 选择currentWeight最大的node，
func (lb *weightedRoundRobinLB) Next() (string, error) {
	var total int64 // total是effective
	var choice *weightNode
	for _, node := range lb.nodes {
		node.Current += node.Effective
		total += node.Effective //
		if node.Effective < node.Weight {
			node.Effective += 1
		}
		if choice == nil {
			choice = node
		} else {
			if node.Current > choice.Current {
				choice = node
			}
		}
	}

	choice.Current -= total
	return choice.Addr, nil
}

func (lb *weightedRoundRobinLB) GetService(hosts []string) (string, error) {
	return hosts[0], nil
}
