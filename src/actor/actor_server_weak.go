/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 仅用于打桩
************************************************************
* Date              Author            Description
*
*/

package actor

import (
	// "github.com/zq-jiang/Panda/base/logger"
)

func (s *Service)RegName() string{
	return ""
}

func (s *Service)RegProtocol() uint32{
	return 0
}

func (s *Service)Dispatch(source uint32, session int32, msg ...interface{}){
	return
}

func (s *Service)Start(){
	return
}

func (s *Service)GetName() string{
	return ""
}

func (s *Service)GetHandle() uint32{
	return 0
}

func (s *Service)GetSession() int32{
	return 0
}

func (s *Service)IncreaseSession() int32{
	return 0
}

func (s *Service)Ret(event actorEvent){
	return
}

func (s *Service)Yiel(event *actorEvent) interface{}{
	return nil
}

func (s *Service)CreateTimeoutGo(event *actorEvent){
	return
}