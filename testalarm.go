package main

import (
	"fmt"

	alarm "github.com/gitea/alarm/alarm"
)

func main() {
	var timevalue = make(map[string]string)
	timevalue["timein"] = "2019/10/17 00:12:44.010"
	timevalue["timeout"] = "2019/10/20 10:09:53.962"
	finalresult, _ := alarm.Sqlalarm(timevalue)
	fmt.Println(finalresult)
}
