/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 服务引导程序
************************************************************
* Date              Author            Description
*
*/

package actor

import (
	"github.com/zq-jiang/Panda/timer"
)

var gLoadr *actorBootLoader

const(
	/*预留高8bit作为remote id*/
	AcHandleMMask uint32 = 0xffffff
	AcHandleRemoteShift int32 = 24
)

type actorBootLoader struct {
	handleIndex int
	harbor uint32
	slot []*actorServer
}

func (loader *actorBootLoader)registerActorHandle(actorSRV *actorServer) uint32{
	handle := loader.handleIndex

	for index, value := range loader.slot[handle:] {
		if nil == value {
			loader.slot[index] = actorSRV
			loader.handleIndex = index
			return loader.harbor | uint32(index)
		}
	}

	for index, value := range loader.slot[:handle] {
		if nil == value {
			loader.slot[index] = actorSRV
			loader.handleIndex = index
			return loader.harbor | uint32(index)
		}
	}

	/*扩容*/
	loader.handleIndex = 0
	loader.slot = append(loader.slot, actorSRV)
	handle = cap(loader.slot) - 1

	return  loader.harbor | uint32(handle)
}

func (loader *actorBootLoader)findActorByHandle(handle uint32) *actorServer{
	index := handle & AcHandleMMask
	actorSRV := loader.slot[index]

	return actorSRV
}

func (loader *actorBootLoader)findActorByName(name string) *actorServer{
	/*TODO*/
	return nil
}

func (loader *actorBootLoader)sendMsgByHandle(self Self, dest uint32, args ...interface{}) interface{}{
	actorSRV := loader.findActorByHandle(dest)
	if nil == actorSRV {
		return nil
	}

	event := GetActorEvent()
	event.source = self.GetHandle()
	event.session = self.IncreaseSession()
	event.eventType = ActorEventTypeSend
	event.msg = args

	actorSRV.write <- *event

	return self.Yiel(event)
}

func (loader *actorBootLoader)sendMsgByHandleNoYiel(self Self, dest uint32, args ...interface{}){
	actorSRV := loader.findActorByHandle(dest)
	if nil == actorSRV {
		return
	}

	event := GetActorEvent()
	event.source = self.GetHandle()
	event.session = self.IncreaseSession()
	event.eventType = ActorEventTypeSend
	event.msg = args

	actorSRV.write <- *event
}

func (loader *actorBootLoader)sendMsgByName(self Self, dest string, args ...interface{}) interface{}{
	/*TODO*/
	return nil
}

func (loader *actorBootLoader)sendMsgByNameNoYiel(self Self, dest string, args ...interface{}){
	/*TODO*/
	return
}

func (loader *actorBootLoader)sendTimeout2Srv(dest uint32, session int32) {
	actorSRV := loader.findActorByHandle(dest)
	if nil == actorSRV {
		return
	}

	event := GetActorEvent()
	event.source = 0
	event.session = session
	event.eventType = ActorEventTypeResp

	actorSRV.write <- *event
}

func (loader *actorBootLoader)timeout(handle uint32, session int32, timeout uint32){
	timer := timer.ServiceGetEntity()
	timer.ServiceTimeout(handle, session, timeout, gLoadr.sendTimeout2Srv)
}

func ActorBootLoadrInit() {
	gLoadr = new(actorBootLoader)
}

func ActorBootLoadrFinit() {
	/*暂时不需要清理*/
}