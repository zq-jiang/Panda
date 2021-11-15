/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 事件队列
************************************************************
* Date              Author            Description
*
*/

package actor

import (
	"sync"
	"unsafe"
)

const(
	ActorEventTypeSend = iota /*系统默认提供的处理*/
	ActorEventTypeResp
)

type actorEvent struct{
	source uint32
	session int32
	eventType int32
	process func()
	msg []interface{}
}

/*事件节点池*/
var EventPool = sync.Pool{New: func() interface{} { return new(actorEvent) }}

func GetActorEvent() *actorEvent {
	return EventPool.Get().(*actorEvent)
}

func PutActorEvent(event *actorEvent) {
	event.session, event.session, event.eventType, event.msg = 0, 0, 0, nil
	EventPool.Put(*event)
}

type ActorEventQueue interface {
	Enqueue(unsafe.Pointer)
	Dequeue() unsafe.Pointer
	IsEmpty() bool
}