/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 定时器实体对象
* 1) 5级时间轮, 32位从高位到低位, 时间周长依次|6|6|6|6|8|
* 2) 每格时间单位1ms, 时间滴答每走过一格即表示时间格内所有事件超时
* 3) 时间轮每走过一圈, 上一级时间轮拨动一格, 依次类推直至最高级时间轮
************************************************************
* Date              Author            Description
*
*/
package timer

import (
	"time"
	"unsafe"
	"github.com/zq-jiang/Panda/base/logger"
	"github.com/zq-jiang/Panda/base/queue"
)

/*定义时间轮的尺度：8:6:6:6:6*/
const(
	TIME_NEAR_SHIFT = 8
	TIME_NEAR = (1 << TIME_NEAR_SHIFT)
	TIME_LEVEL_SHIFT = 6
	TIME_LEVEL = (1 << TIME_LEVEL_SHIFT)
	TIME_NEAR_MASK = (TIME_NEAR-1)
	TIME_LEVEL_MASK = (TIME_LEVEL-1)
)

type timeEntity struct {
	ticks			uint32
	current			uint64
	near			[TIME_NEAR] TimerEventQueue
	t				[4][TIME_LEVEL] TimerEventQueue
}

func ServiceTimerEntityNew() (t *timeEntity){
	var i uint32
	var j uint32
	t = new(timeEntity)
	t.ticks = 0
	t.current = uint64(time.Now().UnixNano() / 1e6)

	for i = 0; i < TIME_NEAR; i++ {
		t.near[i] = queue.NewFreeQueue()
	}

	for i = 0; i < 4; i++ {
		for j = 0; j < TIME_LEVEL; j++ {
			t.t[i][j] = queue.NewFreeQueue()
		}
	}

	return t
}

func (t *timeEntity)serviceTimerEntityDispatch(eventNode *TimerEvent){
	logger.SrvPrintlen("test timer: %d,%d,%d", eventNode.session, eventNode.handle, eventNode.expire)
}

func (t *timeEntity)ServiceTimerEntityAdd(eventNode *TimerEvent){
	var i uint32
	expire := eventNode.expire
	ticks := t.ticks
	
	if (expire | TIME_NEAR_MASK) == (ticks | TIME_NEAR_MASK) {
		t.near[expire&TIME_NEAR_MASK].Enqueue(unsafe.Pointer(eventNode))
	} else {
		mask := uint32(TIME_NEAR << TIME_LEVEL_SHIFT)
		for i = 0 ; i < 3 ; i++ {
			if (expire | (mask-1)) == (ticks | (mask-1)) {
				break
			}
			mask = mask << TIME_LEVEL_SHIFT
		}
		t.t[i][((expire >> (TIME_NEAR_SHIFT + i * TIME_LEVEL_SHIFT)) & TIME_LEVEL_MASK)].Enqueue(unsafe.Pointer(eventNode))	
	}
}

func (t *timeEntity)serviceTimerEntityMoveList(level uint32, idx uint32){
	for !t.t[level][idx].IsEmpty() {
		eventNode := (*TimerEvent)(t.t[level][idx].Dequeue())
		t.ServiceTimerEntityAdd(eventNode)
	}
}

func (t *timeEntity)serviceTimerEntityExec(){
	idx := t.ticks & TIME_NEAR_MASK
	
	for !t.near[idx].IsEmpty() {
		eventNode := (*TimerEvent)(t.near[idx].Dequeue())
		t.serviceTimerEntityDispatch(eventNode)
	}
}

func (t *timeEntity)serviceTimerEntityShift(){
	var idx uint32
	mask := uint32(TIME_NEAR)
	t.ticks++
	ticks := t.ticks
	if (ticks == 0) {
		t.serviceTimerEntityMoveList(3, 0)
	} else {
		times := uint32(ticks >> TIME_NEAR_SHIFT)
		i := uint32(0)
		for (ticks & (mask-1)) == 0 {
			idx = times & TIME_LEVEL_MASK
			if (idx != 0) {
				t.serviceTimerEntityMoveList(i, idx)
				break
			}
			mask <<= TIME_LEVEL_SHIFT
			times >>= TIME_LEVEL_SHIFT
			i++
		}
	}	
}

func (t *timeEntity)serviceTimerEntityUpdate(){
	t.serviceTimerEntityExec()

	t.serviceTimerEntityShift()

	t.serviceTimerEntityExec()
}

func (t *timeEntity)ServiceTimerEntityWalk(){
	var i uint64
	cp := uint64(time.Now().UnixNano() / 1e6)
	if cp < t.current {
		//logger.SrvPrintlen("time diff error: change from %lld to %lld", cp, t.current)
		t.current = cp
	} else if cp != t.current {
		diff := cp - t.current
		t.current = cp
		for i = 0; i < diff; i++ {
			t.serviceTimerEntityUpdate()
		}
	}
}

func (t *timeEntity)ServiceTimeout(handle uint32, timeout uint32, session int32) int32{
	eventNode := GetEvent()
	eventNode.handle = handle
	eventNode.expire = timeout + t.ticks
	eventNode.session = session
	t.ServiceTimerEntityAdd(eventNode)

	return session
}