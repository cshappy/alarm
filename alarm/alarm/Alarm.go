package alarm

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type ALMTAGNAME struct {
	ALM_TAGNAME string
}

type ALMSTATUS struct {
	ALM_TAGNAME   string
	ALM_DATEIN    string
	ALM_TIMEIN    string
	ALM_DATELAST  string
	ALM_TIMELAST  string
	ALM_ALMSTATUS string
}

func Sqlalarm(data map[string]string) (interface{}, error) {
	var server = "10.1.1.31"
	var port = 1433
	var user = "sa"
	var password = "Sealu2019"
	var database = "ceshi"
	var timein = data["timein"]
	var timeout = data["timeout"]
	// var timein = "2019/10/17 00:12:44.010"
	// var timeout = "2019/10/20 10:09:53.962"
	var i = 1
	var count = 0
	var lastarrange = ""
	var sqlstr = ""
	var alarmstart int
	var alarmend int
	var alarmtime int
	var timesum int
	var resultmap = make(map[string]interface{})
	var tagnamemap = make(map[string]interface{})
	connString := fmt.Sprintf("server=%s;port%d;database=%s;user id=%s;password=%s", server, port, database, user, password)

	//建立连接
	db, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open Connection failed:", err.Error())
	}
	defer db.Close()
	timeindate := timein[0:10]
	timeintime := timein[11:]
	timeoutdate := timeout[0:10]
	timeouttime := timeout[11:]
	sqlstring := `select ALM_TAGNAME from(select [ALM_TAGNAME],[ALM_DATEIN],[ALM_TIMEIN],[ALM_DATELAST],[ALM_TIMELAST],[ALM_ALMSTATUS] from (select * from [ceshi].[dbo].[FIXALARMS] where alm_almstatus!='         ') as list group by [ALM_DATELAST],[ALM_TIMELAST],[ALM_TAGNAME],[ALM_DATEIN],[ALM_TIMEIN],[ALM_ALMSTATUS])as b 
	where (ALM_DATELAST>'` + timeindate + `' and ALM_TIMELAST>'` + timeintime + `' and ALM_DATELAST<='` + timeoutdate + `' and ALM_TIMELAST<='` + timeouttime + `') 
	or (ALM_DATEIN>'` + timeindate + `' and ALM_TIMEIN>'` + timeintime + `' and ALM_DATEIN<='` + timeoutdate + `' and ALM_TIMEIN<='` + timeouttime + `') group by ALM_TAGNAME`
	rows, err := db.Query(sqlstring)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()
	var rowsData []*ALMTAGNAME

	for rows.Next() {
		var row = new(ALMTAGNAME)
		rows.Scan(&row.ALM_TAGNAME)
		rowsData = append(rowsData, row)
	}
	for _, ar := range rowsData {
		tagnamemap = make(map[string]interface{})
		i = 1
		count = 0
		timesum = 0
		sqlstr = `select * from (select * from(select [ALM_TAGNAME],[ALM_DATEIN],[ALM_TIMEIN],[ALM_DATELAST],[ALM_TIMELAST],[ALM_ALMSTATUS] from (select * from [ceshi].[dbo].[FIXALARMS] where alm_almstatus!='         ') as list group by [ALM_DATELAST],[ALM_TIMELAST],[ALM_TAGNAME],[ALM_DATEIN],[ALM_TIMEIN],[ALM_ALMSTATUS])as b 
		where (ALM_DATELAST>'` + timeindate + `' and ALM_TIMELAST>'` + timeintime + `' and ALM_DATELAST<='` + timeoutdate + `' and ALM_TIMELAST<='` + timeouttime + `') 
		or (ALM_DATEIN>'` + timeindate + `' and ALM_TIMEIN>'` + timeintime + `' and ALM_DATEIN<='` + timeoutdate + `' and ALM_TIMEIN<='` + timeouttime + `')) as c` + " where ALM_TAGNAME='" + ar.ALM_TAGNAME + "'"
		sqlstrrows, err := db.Query(sqlstr)
		if err != nil {
			log.Fatal("Query failed:", err.Error())
		}
		defer sqlstrrows.Close()
		var alarmstatusData []*ALMSTATUS
		for sqlstrrows.Next() {
			var sqlstrrowrow = new(ALMSTATUS)
			sqlstrrows.Scan(&sqlstrrowrow.ALM_TAGNAME, &sqlstrrowrow.ALM_DATEIN, &sqlstrrowrow.ALM_TIMEIN, &sqlstrrowrow.ALM_DATELAST, &sqlstrrowrow.ALM_TIMELAST, &sqlstrrowrow.ALM_ALMSTATUS)
			alarmstatusData = append(alarmstatusData, sqlstrrowrow)
		}
		for _, arrange := range alarmstatusData {
			//判断第一条数据状态是否为OK
			if i == 1 {
				if arrange.ALM_ALMSTATUS == "OK       " {
					alarmstart = timetotimestamp(timeindate, timeintime)
					alarmend = timetotimestamp(strings.Replace(arrange.ALM_DATELAST, " ", "", -1), strings.Replace(arrange.ALM_TIMELAST, " ", "", -1))
					alarmtime = alarmend - alarmstart
					lastarrange = "OK"
					count = count + 1
					i = i + 1
					continue
				}
			}
			//判断最后一条数据状态是否为OK
			if i == len(alarmstatusData) {
				if arrange.ALM_ALMSTATUS != "OK       " {
					alarmstart = timetotimestamp(strings.Replace(arrange.ALM_DATEIN, " ", "", -1), strings.Replace(arrange.ALM_TIMEIN, " ", "", -1))
					alarmend = timetotimestamp(timeoutdate, timeouttime)
					alarmtime = alarmend - alarmstart
					count = count + 1
				}
			}
			//判断上一条数据状态是否为OK
			if lastarrange == "OK" {
				alarmstart = timetotimestamp(strings.Replace(arrange.ALM_DATEIN, " ", "", -1), strings.Replace(arrange.ALM_TIMEIN, " ", "", -1))
				lastarrange = ""
			}
			//判断此条数据状态是否为OK
			if arrange.ALM_ALMSTATUS == "OK       " {
				alarmend = timetotimestamp(strings.Replace(arrange.ALM_DATELAST, " ", "", -1), strings.Replace(arrange.ALM_TIMELAST, " ", "", -1))
				alarmtime = alarmend - alarmstart
				lastarrange = "OK"
				count = count + 1
			}
			timesum = timesum + alarmtime
			i = i + 1
		}
		fmt.Println(strings.Replace(ar.ALM_TAGNAME, " ", "", -1), count, "次")
		fmt.Println(strings.Replace(ar.ALM_TAGNAME, " ", "", -1), timesum, "秒")
		tagnamemap["alarmcount"] = count
		tagnamemap["alarmtime"] = timesum
		resultmap[strings.Replace(ar.ALM_TAGNAME, " ", "", -1)] = tagnamemap
	}
	return resultmap, nil
}

func timetotimestamp(datestring string, timestring string) int {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tt, _ := time.ParseInLocation("2006/01/02 15:04:05.000", datestring+" "+timestring, loc)
	strInt64 := strconv.FormatInt(tt.Unix(), 10)
	id16, _ := strconv.Atoi(strInt64)
	return id16
}
