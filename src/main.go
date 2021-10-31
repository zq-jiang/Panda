/**
* Date: 2021/10/13
* Author: Arno
* Description: 主函数
**/
package main

import (
	"github.com/zq-jiang/Panda/timer"
	"sync"
)
func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		timer.ServiceTimerThread()
	}()
	wg.Wait()
}