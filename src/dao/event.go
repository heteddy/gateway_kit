// @Author : detaohe
// @File   : event.go
// @Description:
// @Date   : 2022/9/22 18:19

package dao

const (
	EventInvalid = -1
	EventCreate  = iota
	EventUpdate
	EventDelete
)

type SvcEvent struct {
	EventType int
	Entity    *HttpSvcEntity
}

//func (e *SvcEvent) Empty() bool {
//	return len(e.Entity) == 0
//}

type GwEvent struct {
	EventType int
	Entity    *GatewayEntity
}

//func (e *GwEvent) Empty() bool {
//	return len(e.Entity) == 0
//}
