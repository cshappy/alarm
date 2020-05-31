//testalarmcurrent.go为测试文件
package main

import (
	"fmt"

	alarmcurrent "sealu/alarmcurrent/alarmcurrent"
)

func main() {
	finalresult := alarmcurrent.Alarmcurrent()
	fmt.Print(finalresult)
}
