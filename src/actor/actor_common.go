/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 服务对外接口
************************************************************
* Date              Author            Description
*
*/

package actor

/*******************************************************
*   FUNC: CreateTimeout
*   INT: self actor内部接口
*        handle 服务句柄
*        session 会话ID
*        timeout 超时时间, 单位ms
*   OUT: NULL
*   RETURN: NULL
*   DESC: 定时器接口
********************************************************/
func CreateTimeout(self Self, timeout uint32, callback func()) {
	event := GetActorEvent()
	event.source = 0
	event.session = self.IncreaseSession()
	event.eventType = ActorEventTypeSend
	event.process = callback
	event.msg = nil
	self.CreateTimeoutGo(event)

	gLoadr.timeout(self.GetHandle(), event.session, timeout)
}

/*******************************************************
*   FUNC: New
*   INT: interface{} 服务接口
*   OUT: NULL
*   RETURN: UINT32 成功返回服务句柄
*   DESC: 用于生成服务
********************************************************/
func New(srvHandle interface{}) uint32{
	actorSRV := GetActor()
	handle := gLoadr.registerActorHandle(actorSRV)
	actorSRV.handle = handle

	actorSRV.Init(srvHandle)

	go actorSRV.Run()

	return handle
}

/*******************************************************
*   FUNC: Call
*   INT: self actor内部接口
*        handle 服务句柄
*        dest 目标服务(ID/名字)
*        args 可变参数, 服务间通信的参数
*   OUT: NULL
*   RETURN: 任意参数
*   DESC: 服务之间的通信, 会挂起当前服务但是不会阻塞, 当对方服务
*         响应后继续
********************************************************/
func Call(self Self, dest interface{}, args ...interface{}) interface{} {

	switch dest.(type) {
		case uint32:
			return gLoadr.sendMsgByHandle(self, dest.(uint32), args...)
		case string:
			return gLoadr.sendMsgByName(self, dest.(string), args...)
	}

	return nil
}

/*******************************************************
*   FUNC: Send
*   INT: self actor内部接口
*        handle 服务句柄
*        dest 目标服务(ID/名字)
*        args 可变参数, 服务间通信的参数
*   OUT: NULL
*   RETURN: NULL
*   DESC: 服务之间的通信, 不挂起
********************************************************/
func Send(self Self, dest interface{}, args ...interface{}){

	switch dest.(type) {
		case uint32:
			gLoadr.sendMsgByHandleNoYiel(self, dest.(uint32), args...)
		case string:
			gLoadr.sendMsgByNameNoYiel(self, dest.(string), args...)
	}

	return
}

/*******************************************************
*   FUNC: Return
*   INT: self actor内部接口
*        args 可变参数, 服务间通信的参数
*   OUT: NULL
*   RETURN: NULL
*   DESC: 响应服务
********************************************************/
func Return(self Self, args ...interface{}) {
	event := GetActorEvent()
	event.source = self.GetHandle()
	event.session = 0
	event.eventType = ActorEventTypeResp
	event.msg = args
	self.Ret(*event)
}

/*******************************************************
*   FUNC: Exit
*   INT: self actor内部接口
*   OUT: NULL
*   RETURN: NULL
*   DESC: 服务退出, 服务生命已经走到尽头, 停止接收新的消息, 
*         做好善后工作,防止因为未回复导致对方无限期等待
********************************************************/
func Exit(self Self){
	
}