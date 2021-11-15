/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: actor内部协程
* 1. 每次只有一个协程在工作, actor内部无竞态
* 2. 每个请求独立协程处理, 崩溃只会影响当前求. 
* 3. actor可重入, 异步业务可以用同步的写法
************************************************************
* Date              Author            Description
*
*/
package actor

// import (
// 	"github.com/zq-jiang/Panda/base/logger"
// )

type stat = int

const (
	_ stat = iota
	STATCOMPELETE
	STATSUSPEND
	STATYILE
	STATERROR
)

type TaskData struct {
	event actorEvent
	s stat
}

type executor struct {
	session2routing	map[int32]*chan TaskData
	routing2session	map[chan TaskData]int32
	routing2address	map[chan TaskData]uint32
	curRouting *chan TaskData
}

func routingTask(routing chan TaskData, data *TaskData) {
	defer func() {
		data.s = STATERROR
		if recover() != nil {
			/*TODO: core dump printf*/
			routing <- *data
		}
	}()

	/*唤醒*/
	_, ok :=<- routing
	if !ok {
		/*TDOD: 匪夷所思的错误*/
		return
	}
	/*执行*/
	data.event.process()

	data.s = STATCOMPELETE

	routing <- *data
}

func (ex *executor) Init() {
	ex.session2routing = make(map[int32]*chan TaskData)
	ex.routing2session = make(map[chan TaskData]int32)
	ex.routing2address = make(map[chan TaskData]uint32)
}

func (ex *executor) Resume(event actorEvent) {

	taskData := TaskData{event:event}
	var routing chan TaskData
	var pRouting *chan TaskData

	if ActorEventTypeResp == event.eventType {
		pRouting = ex.session2routing[event.session]

		ex.session2routing[event.session] = nil
		if nil == pRouting {
			/*TODO: 该会话已经失效, 打印告警信息即可*/
			return
		}
		
		/*TODO: 校验响应的服务句柄是否与会话一致*/

		ex.curRouting = pRouting
		routing  = *pRouting
		/*激活协程*/
		routing  <- taskData

		taskData =<- routing
	} else {

		/* 创建处理协程 */
		pRouting = ex.createGo(&taskData)
		ex.curRouting = pRouting
		routing = *pRouting
		ex.routing2session[routing] = event.session
		ex.routing2address[routing] = event.source
		ex.session2routing[event.session] = pRouting
		/*激活协程*/
		routing  <- taskData
		taskData =<- routing
	}

	if STATCOMPELETE == taskData.s {
		session := ex.routing2session[routing]
		ex.routing2session[routing] = 0
		ex.routing2address[routing] = 0
		PutActorRouting(*ex.curRouting)
		ex.curRouting = nil
		ex.session2routing[session] = nil
	} else if STATERROR == taskData.s {
		/*TODO: 需要一些善后工作比如告知其待回复的服务发生系统错误*/

		ex.routing2session[routing] = 0
		ex.routing2address[routing] = 0
		PutActorRouting(*ex.curRouting)
		ex.curRouting = nil
	} else {
		/*上个请求未完成*/
		ex.session2routing[taskData.event.session] = &routing
	}

	/*TODO: actor上的业务不能直接使用go, 需要重新实现为actor服务提供并发处理的接口*/
}


func (ex *executor) createGo(data *TaskData) *chan TaskData{
	routing := GetActorRouting()
	go routingTask(routing, data)

	ex.session2routing[data.event.session] = &routing

	return &routing
}

func (ex *executor) Yiel(event *actorEvent) interface{} {
	data := TaskData{event:*event}
	/*暂时挂起协程*/
	data.s = STATYILE
	routing := *ex.curRouting
	routing <- data
	/*唤醒*/
	data =<- routing

	return data.event.msg[0]
}

func (ex *executor) Ret(event actorEvent) {
	data := TaskData{event:event}
	routing := *ex.curRouting
	session := ex.routing2session[routing]
	address := ex.routing2address[routing]	
	data.event.session = session

	ex.Send(address, data.event)
}

func (ex *executor) Send(dest uint32, event actorEvent) {
	actorSRV := gLoadr.findActorByHandle(dest)
	if nil == actorSRV {
		/*TODO: 目标服务已经死亡, 告警信息提示*/
		return
	}

	actorSRV.write <- event
}

func (ex *executor) CreateTimeoutGo(event *actorEvent){
	taskData := TaskData{event:*event}
	pRouting := ex.createGo(&taskData)
	routing := *pRouting
	ex.routing2session[routing] = event.session
	ex.routing2address[routing] = event.source
	ex.session2routing[taskData.event.session] = pRouting
}

func (ex *executor) Done(){
	/*协程退出服务*/
}