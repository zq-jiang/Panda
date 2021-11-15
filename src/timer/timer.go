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

var gTime *timeEntity


func ServiceTimerInit(){
	gTime = ServiceTimerEntityNew()
}

func ServiceTimerThread(){
	logger.SrvPrintf("timer start")
	for {
		gTime.ServiceTimerEntityWalk()
		time.Sleep(time.Duration(2)*time.Millisecond)
	}
}

func ServiceGetEntity() *timeEntity{
	return gTime
}