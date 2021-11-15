/**
* Date: 2021/10/13
* Author: Arno
* Description: 主函数
**/
package main

import (
	"sync"
	"github.com/zq-jiang/Panda/timer"
	"github.com/zq-jiang/Panda/example"
	"github.com/zq-jiang/Panda/actor"
)

func main() {
	var wg sync.WaitGroup
	
	/*初始化*/
	timer.ServiceTimerInit()
	actor.ActorBootLoadrInit()

	wg.Add(1)
	go func() {
		defer wg.Done()
		timer.ServiceTimerThread()
	}()
	actor.New(&example.User{})
	wg.Wait()
}