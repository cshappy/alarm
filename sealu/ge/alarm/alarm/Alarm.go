package alarm

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

// ALM_SELECT为获取的数据结构
type ALM_SELECT struct {
	ALM_TAGNAME        string
	ALM_ALMSTATUS      string
	ALM_NATIVETIMEIN   string
	ALM_NATIVETIMELAST string
	ALM_VALUE          string
	ALM_TAGDESC        string
}

// ALM_ALMSTATUS 根据tagname和timein的点查出所有的数据然后根据timelast进行排序
type ALM_ALMSTATUS struct {
	TAG_TIMEIN            string
	ALM_NATIVETIMEIN      string
	ALM_NATIVETIMELAST    string
	maxALM_NATIVETIMELAST string
	ALM_TAGNAME           string
	ALM_VALUE             string
	ALM_TAGDESC           string
	ALM_ALMSTATUS         string
}

type TAG_TIMEIN struct {
	TAG_TIMEIN string
}

/* Alarmcurrent 函数需要获取一个map，
map包含server,port,user,password,database,table字段
将正在进行的报警放进interface里，并且返回*/
func Alarmcurrent() (map[string]interface{}, error) {

	// server是指sqlserver所在服务器的地址
	var server = "10.1.1.31"

	// port是指sqlserver所在服务器的端口号
	var port = "1443"

	// user是指sqlserver登录的用户名
	var user = "sa"

	// password是指sqlserver登录的密码
	var password = "Sealu2019"

	// database是指sqlserver所使用的库
	var database = "ceshi"

	// table是指sqlserver所要操作的表
	var table = "FIXALARMSTEST"

	// sqlstring是指sql语句
	var sqlstring string

	// resultmap是指每次小循环的结果放进map里面
	var resultmap = make(map[string]string)

	// resultinterface是指每次循环完resultmap放进interface里面
	var resultinterface = []interface{}{}

	// 最后需要返回的map
	var finalresultmap = make(map[string]interface{})

	// 判断上一次timein和这次是否一样
	var lasttimein string

	// connString是指准备连接sqlserver的字符串
	connString := "server=" + server + ";port" + port + ";database=" + database + ";user id=" + user + ";password=" + password
	db, err := sql.Open("mssql", connString)
	if err != nil {
		return finalresultmap, err
	}
	defer db.Close()
	sqlstring = "exec [dbo].[" + table + "tag_timeinselect]"
	tagnamerows, err := db.Query(sqlstring)
	if err != nil {
		return finalresultmap, err
	}
	defer tagnamerows.Close()
	var tagnamerowsData []*TAG_TIMEIN

	for tagnamerows.Next() {
		var tagnamerow = new(TAG_TIMEIN)
		tagnamerows.Scan(&tagnamerow.TAG_TIMEIN)
		tagnamerowsData = append(tagnamerowsData, tagnamerow)
	}
	for _, artagname := range tagnamerowsData {
		sqlstring = "exec [dbo].[" + table + "maxtimelastselect] @tagtimein='" + artagname.TAG_TIMEIN + "'"
		maxrows, err := db.Query(sqlstring)
		if err != nil {
			return finalresultmap, err
		}
		defer maxrows.Close()
		var maxrowsData []*ALM_ALMSTATUS

		for maxrows.Next() {
			var maxrow = new(ALM_ALMSTATUS)
			maxrows.Scan(&maxrow.TAG_TIMEIN, &maxrow.ALM_NATIVETIMEIN, &maxrow.ALM_NATIVETIMELAST, &maxrow.maxALM_NATIVETIMELAST, &maxrow.ALM_TAGNAME, &maxrow.ALM_VALUE, &maxrow.ALM_TAGDESC, &maxrow.ALM_ALMSTATUS)
			maxrowsData = append(maxrowsData, maxrow)
		}
		resultmap = make(map[string]string)
		for _, armax := range maxrowsData {
			if strings.Replace(armax.ALM_TAGNAME, " ", "", -1) != "" {
				if strings.Replace(armax.ALM_ALMSTATUS, " ", "", -1) != "OK" {
					if lasttimein == strings.Replace(armax.ALM_NATIVETIMEIN, " ", "", -1) {
						lasttimein = ""
						break
					}
					resultmap["ALM_NATIVETIMEIN"] = timetotimestamp(strings.Replace(armax.ALM_NATIVETIMEIN, " ", "", -1))
					resultmap["ALM_TAGNAME"] = strings.Replace(armax.ALM_TAGNAME, " ", "", -1)
					resultmap["ALM_VALUE"] = strings.Replace(armax.ALM_VALUE, " ", "", -1)
					resultmap["ALM_DESCR"] = strings.Replace(armax.ALM_TAGDESC, " ", "", -1)
					resultmap["ALM_ALMSTATUS"] = strings.Replace(armax.ALM_ALMSTATUS, " ", "", -1)
					lasttimein = strings.Replace(armax.ALM_NATIVETIMEIN, " ", "", -1)
					resultinterface = append(resultinterface, resultmap)
					finalresultmap["data"] = resultinterface
				}
			}
		}
	}
	return finalresultmap, nil
}

// timetotimestamp 函数是指把时间转换成时间戳格式，返回数据类型为string
func timetotimestamp(timeinput string) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tt, _ := time.ParseInLocation("2006-01-02T15:04:05", timeinput[0:len("2006-01-02T15:04:05")], loc)
	timestring := strconv.FormatInt(tt.UnixNano()/1e6, 10)
	return timestring
}

func Alarmselect(data map[string]string) (map[string]interface{}, error) {
	// server是指sqlserver所在服务器的地址
	var server = "10.1.1.31"

	// port是指sqlserver所在服务器的端口号
	var port = "1443"

	// user是指sqlserver登录的用户名
	var user = "sa"

	// password是指sqlserver登录的密码
	var password = "Sealu2019"

	// database是指sqlserver所使用的库
	var database = "ceshi"

	// table是指sqlserver所要操作的表
	var table = "FIXALARMS"

	// tagname是指标签点名
	var tagname = data["tagname"]

	// perffull是指操作员
	var perffull = data["perffull"]

	// alarmarea是指报警区域
	var alarmarea = data["alarmarea"]

	// timein是指开始时间
	var timein = data["timein"]

	// timelast是指结束时间
	var timelast = data["timelast"]

	// confirmtime是指确认时间
	// var confirmtime = data["confirmtime"]

	// resultmap是指每次小循环的结果放进map里面
	var resultmap = make(map[string]string)

	// resultinterface是指每次循环完resultmap放进interface里面
	var resultinterface = []interface{}{}

	// 最后需要返回的map
	var finalresultmap = make(map[string]interface{})

	connString := "server=" + server + ";port" + port + ";database=" + database + ";user id=" + user + ";password=" + password
	db, err := sql.Open("mssql", connString)
	if err != nil {
		return finalresultmap, err
	}
	db.Exec("CREATE PROCEDURE [dbo].[" + table + "timeselect] @tagname varchar(64),@perffull varchar(64),@alarmarea varchar(64),@timein varchar(64),@timelast varchar(64) AS BEGIN " +
		"select ALM_TAGNAME,ALM_ALMSTATUS,ALM_NATIVETIMEIN,ALM_NATIVETIMELAST,ALM_VALUE,ALM_DESCR from " + table +
		" where ALM_TAGNAME like @tagname and ALM_ALMAREA like @alarmarea and ALM_PERFNAME like @perffull and (ALM_NATIVETIMEIN<@timelast and ALM_NATIVETIMEIN>@timein) or (ALM_TAGNAME like @tagname and ALM_ALMAREA like @alarmarea and ALM_PERFNAME like @perffull and ALM_NATIVETIMELAST<@timelast and ALM_NATIVETIMELAST>@timein) order by ALM_TAGNAME,ALM_NATIVETIMEIN,ALM_NATIVETIMELAST" +
		" END")

	if strings.Replace(tagname, " ", "", -1) == "" {
		tagname = "%"
	}
	if strings.Replace(perffull, " ", "", -1) == "" {
		perffull = "%"
	}
	if strings.Replace(alarmarea, " ", "", -1) == "" {
		alarmarea = "%"
	}
	if strings.Replace(timein, " ", "", -1) == "" {
		timein = "346953600000"
	}
	if strings.Replace(timelast, " ", "", -1) == "" {
		timelast = strconv.Itoa(time.Now().Nanosecond() * 1e6)
	}
	sqlstring := "exec [dbo].[" + table + "timeselect] @tagname='" + tagname + "',@perffull='" + perffull + "',@alarmarea='" + alarmarea + "',@timein='" + timestamptotime(timein) + "',@timelast='" + timestamptotime(timelast) + "'"

	selectrows, err := db.Query(sqlstring)
	if err != nil {
		return finalresultmap, err
	}
	defer selectrows.Close()
	var selectrowsData []*ALM_SELECT

	for selectrows.Next() {
		var selectrow = new(ALM_SELECT)
		selectrows.Scan(&selectrow.ALM_TAGNAME, &selectrow.ALM_ALMSTATUS, &selectrow.ALM_NATIVETIMEIN, &selectrow.ALM_NATIVETIMELAST, &selectrow.ALM_VALUE, &selectrow.ALM_TAGDESC)
		selectrowsData = append(selectrowsData, selectrow)
	}
	for _, artselect := range selectrowsData {
		resultmap = make(map[string]string)
		if strings.Replace(artselect.ALM_TAGNAME, " ", "", -1) == "" {
			continue
		}
		resultmap["ALM_TAGNAME"] = strings.Replace(artselect.ALM_TAGNAME, " ", "", -1)
		resultmap["ALM_ALMSTATUS"] = strings.Replace(artselect.ALM_ALMSTATUS, " ", "", -1)
		resultmap["ALM_NATIVETIMEIN"] = timetotimestamp(artselect.ALM_NATIVETIMEIN)
		resultmap["ALM_NATIVETIMELAST"] = timetotimestamp(artselect.ALM_NATIVETIMELAST)
		resultmap["ALM_VALUE"] = strings.Replace(artselect.ALM_VALUE, " ", "", -1)
		resultmap["ALM_DESCR"] = strings.Replace(artselect.ALM_TAGDESC, " ", "", -1)
		resultinterface = append(resultinterface, resultmap)
	}
	finalresultmap["data"] = resultinterface
	return finalresultmap, err
}

// timestamptotime函数是指时间戳转换成时间
func timestamptotime(timestampinput string) string {
	// stringint, _ := strconv.Atoi(timestampinput)
	stringint, _ := strconv.ParseInt(timestampinput, 10, 64)
	tm := time.Unix(stringint/1000, 0)
	return tm.Format("2006-01-02 15:04:05")
}
