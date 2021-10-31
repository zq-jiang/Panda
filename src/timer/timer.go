/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 定时器
************************************************************
* Date              Author            Description
*
*/
package timer

import (
	"time"
	"github.com/zq-jiang/Panda/base/logger"
)

func ServiceTimerThread(){
	logger.SrvPrintlen("timer start")

	t := ServiceTimerEntityNew()
	t.ServiceTimeout(uint32(1), uint32(1000), 1)
	t.ServiceTimeout(uint32(2), uint32(2000), 1)
	t.ServiceTimeout(uint32(3), uint32(3000), 1)
	t.ServiceTimeout(uint32(4), uint32(4000), 1)
	t.ServiceTimeout(uint32(5), uint32(3000), 1)
	t.ServiceTimeout(uint32(6), uint32(2000), 1)
	t.ServiceTimeout(uint32(7), uint32(1000), 1)
	t.ServiceTimeout(uint32(8), uint32(10000), 1)
	t.ServiceTimeout(uint32(9), uint32(20000), 1)
	for true {
		t.ServiceTimerEntityWalk()
		time.Sleep(time.Duration(2)*time.Millisecond)
	}
}