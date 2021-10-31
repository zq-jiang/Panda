/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 队列, 无锁安全队列
************************************************************
* Date              Author            Description
*
*/

package queue

import (
	"sync/atomic"
	"unsafe"
)

type FreeNode struct {
	value unsafe.Pointer
	next  unsafe.Pointer
}

type FreeQueue struct {
	head   unsafe.Pointer
	tail   unsafe.Pointer
	length int32
}

func NewFreeQueue() *FreeQueue {
	/*创建一个冗余节点, 用于保持队列头永远不变, 原子操作队列尾即可保证队列安全*/
	n := unsafe.Pointer(&node{})
	return &FreeQueue{head: n, tail: n, length: 0} //unsafe.Pointer(&)
}

func (q *FreeQueue) Enqueue(data unsafe.Pointer) {
	n := &FreeNode{value: data}
retry:
	tail := load(&q.tail)
	next := load(&tail.next)
	if tail == load(&q.tail) {
		if nil == next {
			if cas(&tail.next, next, n) {
				cas(&q.tail, tail, n)
				atomic.AddInt32(&q.length, 1)
				return
			}
		} else {
			cas(&q.tail, tail, next)
		}
	}
	goto retry
}

func (q *FreeQueue) Dequeue() unsafe.Pointer {
retry:
	head := load(&q.head)
	tail := load(&q.tail)
	next := load(&head.next)
	if head == load(&q.head) {
		if head == tail {
			if next == nil {
				return nil
			}
			cas(&q.tail, tail, next) 
		} else {
			task := next.value
			if cas(&q.head, head, next) {
				atomic.AddInt32(&q.length, -1)
				return task
			}
		}
	}
	goto retry
}

func (q *FreeQueue) IsEmpty() bool {
	return atomic.LoadInt32(&q.length) == 0
}

func load(p *unsafe.Pointer) (n *FreeNode) {
	return (*FreeNode)(atomic.LoadPointer(p))
}

func cas(p *unsafe.Pointer, old, new *FreeNode) bool {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}