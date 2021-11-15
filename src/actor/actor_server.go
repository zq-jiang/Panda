/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 服务
* 1. actor状态只由自己维护
* 2. 不同actor的状态只能通过eamil通信影响
* 3. 通过通信去共享内存, 而不是通过共享内存去通信
* ACTOR-A             ACTOR-B
*  eamil  --channel->  eamil
*                       |
*                    channel
*                       |
* ACTOR_A <--channel-- work
************************************************************
* Date              Author            Description
*
*/

package actor

import (
	"sync"
	"unsafe"
	"reflect"
	"github.com/zq-jiang/Panda/base/logger"
	"github.com/zq-jiang/Panda/base/queue"
)

type (

	/*所有的接口都由业务层具体实现*/
	ServiceProHandle interface {
		/*分发函数：业务处理*/
		Dispatch(source uint32, session int32, msg ...interface{})

		/*注册服务名: 可通过名字通信*/
		RegName()(string)

		/*注册协议: 协议解析的格式*/
		RegProtocol()(uint32)

		Start()
	}

	/*自身操作接口*/
	Self interface {
		GetHandle() uint32
	
		GetSession() int32

		IncreaseSession() int32

		Ret(event actorEvent)

		Yiel(event *actorEvent) interface{}

		CreateTimeoutGo(event *actorEvent)
	}

	/*业务服务必须包含*/
	Service struct{
		Self Self
	}
)

type actorServer struct {
	/*服务名称*/
	name string
	/*actor的句柄*/
	handle uint32
	/*会话*/
	session int32
	/*actor消息邮箱*/
	email ActorEventQueue
	/*写通道: 用于对外交互的通道*/
	write chan actorEvent
	/*内部通道：用于与内部协程进行共享同步*/
	inner chan chan actorEvent
	/*actor执行接口*/
	ex *executor
	/*注册业务层处理函数*/
	process ServiceProHandle
}

/*服务者池*/
var actorServerPool = sync.Pool{New: func() interface{} { return new(actorServer) }}

func GetActor() *actorServer {
	actorSRV := actorServerPool.Get().(*actorServer)
	actorSRV.write = make(chan actorEvent)
	actorSRV.inner = make(chan chan actorEvent)
	actorSRV.email = queue.NewQueue()
	actorSRV.ex = new(executor)
	actorSRV.ex.Init()
	return actorSRV
}

func PutActor(actorSRV *actorServer) {
	actorSRV.handle = 0
	for !actorSRV.email.IsEmpty() {
		node := (*actorEvent)(actorSRV.email.Dequeue())
		PutActorEvent(node)
	}	
	actorServerPool.Put(actorSRV)
}

func (actorSRV *actorServer)Run() {
	var event actorEvent
	var channel chan actorEvent
	var idel *chan actorEvent = nil

	actorSRV.Work()

	for {
		select{
			case event =<- actorSRV.write :
				/*消息入队*/
				if nil == idel {
					actorSRV.email.Enqueue(unsafe.Pointer(&event))
				} else {
					channel = *idel
					idel = nil
					if actorSRV.email.IsEmpty() {
						channel <- event
					} else {
						channel <- *(*actorEvent)(actorSRV.email.Dequeue())
					}
				}
			case channel =<- actorSRV.inner :
				/*消息出队*/
				if actorSRV.email.IsEmpty() {
					idel = &channel
				} else {
					idel = nil
					channel <- *(*actorEvent)(actorSRV.email.Dequeue())
				}
		}
	}

}

func (actorSRV *actorServer)Work() {
	var idel chan actorEvent = make(chan actorEvent)
	var event actorEvent

	go func(){
		for {
			/*任务空闲*/
			actorSRV.inner <- idel
			
			select{
				case event =<- idel:
					if ActorEventTypeResp == event.eventType {
						/*response协议已经处理*/
					} else {
						/*非response报文, 目前暂时支持send报文分发*/
						event.process = actorSRV.withAcotrProc(event.source, event.session, event.msg...)
					}
					actorSRV.ex.Resume(event)
			}
		}
	}()
}

func (actorSRV *actorServer)GetName() string{
	/*获取handle*/
	return actorSRV.name
}

func (actorSRV *actorServer)GetHandle() uint32{
	/*获取handle*/
	return actorSRV.handle
}

func (actorSRV *actorServer)GetSession() int32{
	/*获取session*/
	return actorSRV.session
}

func (actorSRV *actorServer)IncreaseSession() int32{
	/*获取session*/
	actorSRV.session++
	return actorSRV.session
}

func (actorSRV *actorServer)Ret(event actorEvent){
	actorSRV.ex.Ret(event)
}

func (actorSRV *actorServer)Yiel(event *actorEvent) interface{}{
	return actorSRV.ex.Yiel(event)
}

func (actorSRV *actorServer)CreateTimeoutGo(event *actorEvent){
	actorSRV.ex.CreateTimeoutGo(event)
}

func (actorSRV *actorServer)withAcotrProc(source uint32, session int32, msg ...interface{}) func() {
	return func(){
		actorSRV.process.Dispatch(source, session, msg...)
	}
}

func (actorSRV *actorServer)Start(timeout uint32, callback func()) {
	event := GetActorEvent()
	event.source = 0
	event.session = actorSRV.IncreaseSession()
	event.eventType = ActorEventTypeSend
	event.process = callback
	event.msg = nil
	actorSRV.ex.Resume(*event)
}

func (actorSRV *actorServer)Init(srvHandle interface{}){

	logger.SrvPrintf("actor init")
	
	/*step1: 检查业务必须提供的接口, 保证服务的基本功能*/

	/*step2：业务提供的接口*/
	if value,ok := srvHandle.(ServiceProHandle); ok {
		actorSRV.process = value
	}else{
		logger.SrvPrintf("ServiceProHandle interface init error")
	}

	/*step3: 通过反射向业务反向注入actor的接口依赖*/
	value := reflect.ValueOf(srvHandle).Elem().FieldByName("Self")
	if value.CanAddr() && value.CanSet() {
		value.Set(reflect.ValueOf(actorSRV))
	}else{
		logger.SrvPrintf("self interface init error")
	}

	/*step3: 非必要接口初始化, 业务可以不提供, 由actor提供默认接口*/

	/*dispatch 分发函数处理消息, 不需要初始化*/

	/*regName: 注册名字, 业务没有注册名字, 使用handle生成*/
	//name := actorSRV.process.regName()

	/*RegProtocol: 注册解析协议， 如果没有注册则使用默认的*/
	//protocaol := actorSRV.process.RegProtocol()
	
	/*step4: 服务启动*/
	actorSRV.Start(0, actorSRV.process.Start)
}