/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: actor内部协程池(用channel间接控制协程，即每个channle相当于一个协程)
* &后续考虑ants库
************************************************************
* Date              Author            Description
*
*/
package actor

import (
	"sync"
)

/*事件节点池*/
var RoutingPool = sync.Pool{New: func() interface{} { return make(chan TaskData) }}

func GetActorRouting() chan TaskData {
	return RoutingPool.Get().(chan TaskData)
}

func PutActorRouting(channel chan TaskData) {
	RoutingPool.Put(channel)
}