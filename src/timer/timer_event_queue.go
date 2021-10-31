/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 定时器事件, 定时器一般会频繁注册超时事件,
*              所以事件采用对象池的方式管理
************************************************************
* Date              Author            Description
*
*/
package timer

import (
	"sync"
	"unsafe"
)

type TimerEvent struct{
	session int32
	handle uint32
	expire uint32
}

/*定时事件节点池*/
var EventPool = sync.Pool{New: func() interface{} { return new(TimerEvent) }}

func GetEvent() *TimerEvent {
	return EventPool.Get().(*TimerEvent)
}

func PutEvent(event *TimerEvent) {
	event.session, event.handle = 0, 0
	EventPool.Put(event)
}

type TimerEventQueue interface {
	Enqueue(unsafe.Pointer)
	Dequeue() unsafe.Pointer
	IsEmpty() bool
}