package main

import (
	"fmt"
	alarmcurrent "github.com/gitea/alarmcurrent/alarmcurrent"
)

func main() {
	finalresult :=alarmcurrent.Alarmcurrent()
	fmt.Print(finalresult)
}