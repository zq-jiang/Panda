/***********************************************************
* Date: 2021/10/15
* Author: Arno
* Description: 系统日志
************************************************************
* Date              Author            Description
*
*/
package logger

import (
    "log"
)


func SrvPrintlen(str string, args ...interface{}){
	log.Println(str, args)
}