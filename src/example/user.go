/**
* 举例
**/
package example

import (
	"github.com/zq-jiang/Panda/base/logger"
	"github.com/zq-jiang/Panda/actor"
)

type task = func(...interface{}) interface{}

type User struct {
	actor.Service
	name string
	arg uint32
	cmd map[string]task
}

func (u *User)Hello(msg ...interface{}) interface{}{
	logger.SrvPrintf("Hello")

	return "hello!"
}

func (u *User)register(name string, callback task){
	u.cmd[name] = callback
}


func (u *User)Start(){

	/*注册*/
	u.cmd = make(map[string]task)
	u.register("Hello", u.Hello)

	handle := actor.New(&Room{})

	result := actor.Call(u.Self, handle, "Jion")
	logger.SrvPrintf(result.(string))

	actor.CreateTimeout(u.Self, 3000, func(){
		logger.SrvPrintf("3 sec timeout")
	})
}

func (u *User)Dispatch(source uint32, session int32, msg ...interface{}){
	logger.SrvPrintf("User Dispatch")

	cmd := msg[0]
	args := msg[1:]

	actor.Return(u.Self, u.cmd[cmd.(string)](args))
}