package main

import (
	"fmt"

	alarmselect "sealu/alarmcurrent/alarmselect"
)

func main() {
	var getmap = make(map[string]string)
	getmap["tagname"] = "TAG1"
	getmap["perffull"] = ""
	getmap["alarmarea"] = "ALL"
	getmap["timein"] = "1571295856000"
	getmap["timelast"] = "1571457432000"
	finalresult, err := alarmselect.Alarmselect(getmap)
	fmt.Print(finalresult, err)
}
