// @Author : detaohe
// @File   : command.go
// @Description:
// @Date   : 2022/10/16 18:07

package flow

type StorageCmd interface {
	Run(*FlowStorage) error
}

type LoadGwFlowCmd struct {
	key string
	gwC chan map[string]int64
}

func (cmd *LoadGwFlowCmd) Run(storage *FlowStorage) error {
	v1, e1 := storage.loadByKey(cmd.key)
	if e1 != nil {
		cmd.gwC <- nil
		return e1
	}
	cmd.gwC <- v1
	return nil
}

type DecreaseGwFlowCmd struct {
	gwKey   string
	gwValue int64
}

func (cmd *DecreaseGwFlowCmd) Run(storage *FlowStorage) error {
	return storage.DecrBy(cmd.gwKey, cmd.gwValue)

}

type IncreaseGwFlowCmd struct {
	gwKey   string
	gwValue int64
}

func (cmd *IncreaseGwFlowCmd) Run(storage *FlowStorage) error {
	return storage.IncrBy(cmd.gwKey, cmd.gwValue)
}

type IncreaseSvcFlowCmd struct {
	svcKey   string
	svcValue int64
	reqKey   string
	reqValue int64
}

func (cmd *IncreaseSvcFlowCmd) Run(storage *FlowStorage) error {
	e1 := storage.IncrBy(cmd.svcKey, cmd.svcValue)
	if e1 != nil {
		return e1
	}
	e2 := storage.IncrBy(cmd.reqKey, cmd.reqValue)
	return e2
}

type DecreaseSvcFlowCmd struct {
	svcKey      string
	svcValue    int64
	requestFlow map[string]int64
}

func (cmd *DecreaseSvcFlowCmd) Run(storage *FlowStorage) error {
	e1 := storage.DecrBy(cmd.svcKey, cmd.svcValue)
	if e1 != nil {
		return e1
	}
	for k, v := range cmd.requestFlow {
		storage.DecrBy(k, v)
	}
	return nil
}

type LoadSvcFlowCmd struct {
	svcKey    string
	reqPrefix string
	svcC      chan map[string]int64
	reqC      chan map[string]int64
}

func (cmd *LoadSvcFlowCmd) Run(storage *FlowStorage) error {
	v1, e1 := storage.loadByKey(cmd.svcKey)
	if e1 != nil {
		return e1
	}
	cmd.svcC <- v1
	v2, e2 := storage.loadByPrefix(cmd.reqPrefix)
	if e2 != nil {
		return e2
	}
	cmd.reqC <- v2
	return nil
}
