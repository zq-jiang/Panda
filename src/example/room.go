/**
* 举例
**/
package example

import (
	"github.com/zq-jiang/Panda/base/logger"
	"github.com/zq-jiang/Panda/actor"
)

type Room struct {
	actor.Service
	name string
	users []uint32
	cmd map[string]task
}

func (u *Room)Jion(msg ...interface{}) interface{}{
	user := msg[0].(uint32)
	logger.SrvPrintf("User ID : %d jion", user)

	return "Jion Success"
}

func (u *Room)register(name string, callback task){
	u.cmd[name] = callback
}

func (u *Room)Start(){
	/*注册*/
	u.cmd = make(map[string]task)
	u.register("Jion", u.Jion)
}

func (u *Room)Dispatch(source uint32, session int32, msg ...interface{}){
	cmd := msg[0]
	cmdstr := cmd.(string)

	logger.SrvPrintf("Room Dispatch : %s", cmd)

	actor.Return(u.Self, u.cmd[cmdstr](source))
}