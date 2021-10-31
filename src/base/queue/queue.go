/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 队列, 禁止在不安全的场景中使用
************************************************************
* Date              Author            Description
*
*/

package queue

import (
	"unsafe"
)

type node struct {
	value unsafe.Pointer
	next  *node
}

type Queue struct {
	head   *node
	tail   *node
	length int32
}

func NewQueue() *Queue {
	return &Queue{head: nil, tail: nil, length: 0} //unsafe.Pointer(&)
}

func (q *Queue) Enqueue(data unsafe.Pointer) {
	n := &node{value: data, next: nil}
	if nil == q.head {
		q.head = n
		q.tail = n
	} else {
		q.tail.next = n
		q.tail = n
	}
	q.length++	
}

func (q *Queue) Dequeue() unsafe.Pointer {
	n := q.head
	if q.head == q.tail {

		if nil == n {
			/*空队列*/
			return nil
		}
		q.head = nil
		q.tail = nil
		q.length--
		return n.value
	} else {
		q.head = n.next
		q.length--
		return n.value
	}
}

func (q *Queue) IsEmpty() bool {
	if (0 == q.length){
		return true
	}
	return false
}