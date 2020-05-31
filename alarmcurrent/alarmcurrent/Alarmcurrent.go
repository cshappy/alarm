// Package alarmcurrent包提供正在报警的数据
package alarmcurrent

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

// ALMTAGNAME 获取tagname的表
type ALMTAGNAME struct {
	ALM_TAGNAME string
}

// ALM_NATIVETIMEIN 获取每个tagname的timein的表
type ALM_NATIVETIMEIN struct {
	ALM_NATIVETIMEIN string
}

// ALM_ALMSTATUS 根据tagname和timein的点查出所有的数据然后根据timelast进行排序
type ALM_ALMSTATUS struct {
	ALM_NATIVETIMEIN   string
	ALM_NATIVETIMELAST string
	ALM_TAGNAME        string
	ALM_VALUE          string
	ALM_TAGDESC        string
	ALM_ALMSTATUS      string
}

// Alarmcurrent 是一个无参函数，将正在进行的报警放进interface里，并且返回
func Alarmcurrent() []map[string]string {
	var server = "10.1.1.31"
	var port = 1433
	var user = "sa"
	var password = "Sealu2019"
	var database = "ceshi"
	var table = "[ceshi].[dbo].[FIXALARMS4]"
	var sqlstring string
	var line int
	var resultmap = make(map[string]string)
	var resultinterface []map[string]string
	connString := fmt.Sprintf("server=%s;port%d;database=%s;user id=%s;password=%s", server, port, database, user, password)

	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
	}
	defer db.Close()
	sqlstring = "select ALM_TAGNAME from " + table + " group by ALM_TAGNAME"
	tagnamerows, err := db.Query(sqlstring)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer tagnamerows.Close()
	var tagnamerowsData []*ALMTAGNAME

	for tagnamerows.Next() {
		var tagnamerow = new(ALMTAGNAME)
		tagnamerows.Scan(&tagnamerow.ALM_TAGNAME)
		tagnamerowsData = append(tagnamerowsData, tagnamerow)
	}
	line = 1
	for _, artagname := range tagnamerowsData {
		sqlstring = "select ALM_NATIVETIMEIN from " + table + " " + "where ALM_TAGNAME='" + artagname.ALM_TAGNAME + "' group by ALM_NATIVETIMEIN"
		timeinrows, err := db.Query(sqlstring)
		if err != nil {
			log.Fatal("Query failed:", err.Error())
		}
		defer timeinrows.Close()
		var timeinrowsData []*ALM_NATIVETIMEIN

		for timeinrows.Next() {
			var timeinrow = new(ALM_NATIVETIMEIN)
			timeinrows.Scan(&timeinrow.ALM_NATIVETIMEIN)
			timeinrowsData = append(timeinrowsData, timeinrow)
		}
		for _, artimein := range timeinrowsData {
			sqlstring = "select ALM_NATIVETIMEIN,ALM_NATIVETIMELAST,ALM_TAGNAME,ALM_VALUE,ALM_TAGDESC string,ALM_ALMSTATUS string from " + table + " " +
				"where ALM_TAGNAME='" + artagname.ALM_TAGNAME + "' and ALM_NATIVETIMEIN='" + artimein.ALM_NATIVETIMEIN + "' " + "order by ALM_NATIVETIMELAST"
			statusrows, err := db.Query(sqlstring)
			if err != nil {
				log.Fatal("Query failed:", err.Error())
			}
			defer statusrows.Close()
			var statusrowsData []*ALM_ALMSTATUS

			for statusrows.Next() {
				var statusrow = new(ALM_ALMSTATUS)
				statusrows.Scan(&statusrow.ALM_NATIVETIMEIN, &statusrow.ALM_NATIVETIMELAST, &statusrow.ALM_TAGNAME, &statusrow.ALM_VALUE, &statusrow.ALM_TAGDESC, &statusrow.ALM_ALMSTATUS)
				statusrowsData = append(statusrowsData, statusrow)
			}
			resultmap = make(map[string]string)
			for _, arstatus := range statusrowsData {
				if len(statusrowsData) == line {
					if arstatus.ALM_ALMSTATUS != "OK       " {
						if strings.Replace(arstatus.ALM_TAGNAME, " ", "", -1) != "" {
							resultmap["TimeIn"] = strings.Replace(arstatus.ALM_NATIVETIMEIN, " ", "", -1)
							resultmap["TagName"] = strings.Replace(arstatus.ALM_TAGNAME, " ", "", -1)
							resultmap["Status"] = strings.Replace(arstatus.ALM_ALMSTATUS, " ", "", -1)
							resultmap["Value"] = strings.Replace(arstatus.ALM_VALUE, " ", "", -1)
							resultmap["Description"] = strings.Replace(arstatus.ALM_TAGDESC, " ", "", -1)
							resultinterface = append(resultinterface, resultmap)
						}

					}
					line = 1
					break
				}
				line = line + 1
			}
		}
	}
	return resultinterface
}
