// @Author : detaohe
// @File   : event.go
// @Description:
// @Date   : 2022/9/22 18:19

package dao

const (
	EventUpdate = iota
	EventDelete
)

type SvcEvent struct {
	EventType int
	Entities  []*HttpSvcEntity
}

func (e *SvcEvent) Empty() bool {
	return len(e.Entities) == 0
}

type GwEvent struct {
	EventType int
	Entities  []*GatewayEntity
}

func (e *GwEvent) Empty() bool {
	return len(e.Entities) == 0
}
