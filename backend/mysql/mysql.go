package mysql

import (
	"context"
	"database/sql"
	utils "dms/backend"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/percona/go-mysql/query"
)

type Mysql struct {
	ctx context.Context
}

type MysqlProcess struct {
	ID      int    `json:"id"`
	User    string `json:"user"`
	Host    string `json:"host"`
	Dbname  string `json:"dbname"`
	Command string `json:"command"`
	Time    int    `json:"time"`
	Status  string `json:"status"`
	Sql     string `json:"sql"`
	Key     string `json:"key"`
}

type MysqlProcessF struct {
	ID       int            `json:"id"`
	User     string         `json:"user"`
	Host     string         `json:"host"`
	Dbname   string         `json:"dbname"`
	Command  string         `json:"command"`
	Time     int            `json:"time"`
	Status   string         `json:"status"`
	Sql      string         `json:"sql"`
	Key      string         `json:"key"`
	Children []MysqlProcess `json:"children"`
}

type MysqlConsStatus struct {
	Name    string `json:name`
	Actives int    `json:actives`
	Total   int    `json:total`
}

func getDb(dbId string) *sql.DB {
	dataEncoded, _ := os.ReadFile("./cons/" + dbId)
	m := new(utils.Connection)
	json.Unmarshal(dataEncoded, &m)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema", m.User, m.Password, m.Host, m.Port) //替换为你自己的数据库信息
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func NewMysql() *Mysql {
	return &Mysql{}
}

// 获取数据库连接,可以分别获取活跃或者非活跃的
func (a *Mysql) GetMysqlProcesslist(dbId string, sleep string) []MysqlProcessF {
	db := getDb(dbId)
	defer db.Close()
	sql := "SELECT ID, USER, HOST, ifnull(DB,'') as DB, COMMAND,TIME, STATE, IFNULL(INFO, 'No Query') AS INFO FROM INFORMATION_SCHEMA.PROCESSLIST"
	if sleep == "alive" {
		sql = sql + " where command != 'sleep'"
	}

	// 执行SHOW PROCESSLIST查询获取连接信息
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var processesFs []MysqlProcessF
	sqlFingerPrintMap := make(map[string][]MysqlProcess)
	for rows.Next() {
		var processes MysqlProcess
		// 从查询结果中扫描数据到结构体
		err := rows.Scan(&processes.ID, &processes.User, &processes.Host, &processes.Dbname, &processes.Command, &processes.Time, &processes.Status, &processes.Sql)
		if err != nil {
			log.Fatal(err)
		}
		processes.Key = strconv.Itoa(processes.ID)
		if processes.Command == "Query" {
			if _, exists := sqlFingerPrintMap[query.Fingerprint(processes.Sql)]; !exists {
				sqlFingerPrintMap[query.Fingerprint(processes.Sql)] = []MysqlProcess{} // 初始化切片
			}

			// 添加元素
			sqlFingerPrintMap[query.Fingerprint(processes.Sql)] = append(sqlFingerPrintMap[query.Fingerprint(processes.Sql)], processes)
		} else {
			var processF = MysqlProcessF{processes.ID, processes.User, processes.Host, processes.Dbname, processes.Command, processes.Time, processes.Status, processes.Sql, strconv.Itoa(processes.ID), nil}
			processesFs = append(processesFs, processF)
		}
	}
	for _, items := range sqlFingerPrintMap {
		if len(items) == 1 {
			for _, item := range items {
				var processF = MysqlProcessF{item.ID, item.User, item.Host, item.Dbname, item.Command, item.Time, item.Status, item.Sql, strconv.Itoa(item.ID), nil}
				processesFs = append(processesFs, processF)
			}
		} else {
			var processF = MysqlProcessF{items[0].ID, items[0].User, items[0].Host, items[0].Dbname, items[0].Command, items[0].Time, items[0].Status, query.Fingerprint(items[0].Sql), strconv.Itoa(items[0].ID), []MysqlProcess{}}
			for _, item := range items {
				process := MysqlProcess{item.ID, item.User, item.Host, item.Dbname, item.Command, item.Time, item.Status, item.Sql, strconv.Itoa(items[0].ID) + "-" + strconv.Itoa(item.ID)}
				processF.Children = append(processF.Children, process)
			}
			processesFs = append(processesFs, processF)
		}
	}
	return processesFs
}

func (a *Mysql) KillMysqlProcesses(dbId string, processes []MysqlProcessF) string {
	//fmt.Println(len(processes))
	db := getDb(dbId)
	defer db.Close()
	for _, process := range processes {
		fmt.Println("---------------")
		//fmt.Println(process)
		if process.Children != nil {
			continue
		} else {
			db.Query("kill " + strconv.Itoa(process.ID))
		}
	}

	return "success"
}

func (a *Mysql) GetConsStatus(dbId string) map[string][]MysqlConsStatus {

	status := make(map[string][]MysqlConsStatus)
	db := getDb(dbId)
	defer db.Close()

	sqlU := "select ifnull(user,''),count(info),count(*) from information_schema.PROCESSLIST p group by user"
	sqlI := "select ifnull(substring_index(host, ':', 1),'') ,count(info),count(*) from information_schema.PROCESSLIST p group by ifnull(substring_index(host, ':', 1),'')"
	sqlD := "select ifnull(DB,'') ,count(info),count(*) from information_schema.PROCESSLIST p group by DB "

	status["User"] = []MysqlConsStatus{}
	status["Ip"] = []MysqlConsStatus{}
	status["Db"] = []MysqlConsStatus{}

	rowsU, _ := db.Query(sqlU)
	defer rowsU.Close()

	for rowsU.Next() {
		var conStatus MysqlConsStatus
		err := rowsU.Scan(&conStatus.Name, &conStatus.Actives, &conStatus.Total)
		if err != nil {
			fmt.Println(err)
		} else {
			status["User"] = append(status["User"], conStatus)
		}
		fmt.Println(conStatus)
	}
	rowI, err := db.Query(sqlI)
	if err != nil {
		fmt.Println(err)
	} else {
		for rowI.Next() {
			var conStatus MysqlConsStatus
			err := rowI.Scan(&conStatus.Name, &conStatus.Actives, &conStatus.Total)
			if err != nil {
				fmt.Println(err)
			} else {
				status["Ip"] = append(status["Ip"], conStatus)
			}
		}
	}

	rowD, err := db.Query(sqlD)
	if err != nil {
		fmt.Println(err)
	} else {
		for rowD.Next() {
			var conStatus MysqlConsStatus
			err := rowD.Scan(&conStatus.Name, &conStatus.Actives, &conStatus.Total)
			if err != nil {
				fmt.Println(err)
			} else {
				status["Db"] = append(status["Db"], conStatus)
			}
		}
	}

	return status
}

func (a *Mysql) GetConspercent(dbId string) map[string]string {

	db := getDb(dbId)
	defer db.Close()

	status := make(map[string]string)
	countT := 0
	countA := 0
	row := db.QueryRow("select count(*),count(info) from information_schema.processlist")

	row.Scan(&countT, &countA)
	status["total"] = strconv.Itoa(countT)
	status["active"] = strconv.Itoa(countA)
	fmt.Println(status)
	return status
}

func (a *Mysql) GetMysqlLock(dbId string) []MysqlProcessF {
	db := getDb(dbId)
	defer db.Close()

	//先获取阻塞的组合,然后根据组合去循环
	rows, err := db.Query("select requesting_trx_id,blocking_trx_id  from information_schema.INNODB_LOCK_WAITS ")
	defer rows.Close()
	var processesFs []MysqlProcessF

	if err == nil {
		for rows.Next() {
			fmt.Println("start process")
			var requesting_trx_id int
			var blocking_trx_id int
			rows.Scan(&requesting_trx_id, &blocking_trx_id)
			var requestingProcess MysqlProcess
			var blockingProcess MysqlProcess

			err := db.QueryRow("select p.ID, p.USER, p.HOST, ifnull(p.DB,'') as DB, p.COMMAND,p.TIME, t.trx_state, IFNULL(p.INFO, 'No Query') AS INFO FROM INFORMATION_SCHEMA.PROCESSLIST p, information_schema.INNODB_TRX t where p.id = t.trx_mysql_thread_id  and t.trx_id = ?", requesting_trx_id).Scan(&requestingProcess.ID, &requestingProcess.User, &requestingProcess.Host, &requestingProcess.Dbname, &requestingProcess.Command, &requestingProcess.Time, &requestingProcess.Status, &requestingProcess.Sql)

			if err != nil {
				fmt.Println(err)
			}

			err = db.QueryRow("select p.ID, p.USER, p.HOST, ifnull(p.DB,'') as DB, p.COMMAND,p.TIME, t.trx_state, IFNULL(p.INFO, 'No Query') AS INFO FROM INFORMATION_SCHEMA.PROCESSLIST p, information_schema.INNODB_TRX t where p.id = t.trx_mysql_thread_id  and t.trx_id = ?", blocking_trx_id).Scan(&blockingProcess.ID, &blockingProcess.User, &blockingProcess.Host, &blockingProcess.Dbname, &blockingProcess.Command, &blockingProcess.Time, &blockingProcess.Status, &blockingProcess.Sql)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(requestingProcess)
			fmt.Println(blockingProcess)
			var processF = MysqlProcessF{requestingProcess.ID, requestingProcess.User, requestingProcess.Host, requestingProcess.Dbname, requestingProcess.Command, requestingProcess.Time, requestingProcess.Status, requestingProcess.Sql, strconv.Itoa(requestingProcess.ID), []MysqlProcess{}}
			processF.Key = strconv.Itoa(requestingProcess.ID)
			requestingProcess.Key = strconv.Itoa(requestingProcess.ID) + "-" + strconv.Itoa(requestingProcess.ID)
			blockingProcess.Key = strconv.Itoa(requestingProcess.ID) + "-" + strconv.Itoa(blockingProcess.ID)
			processF.Children = append(processF.Children, requestingProcess, blockingProcess)
			processesFs = append(processesFs, processF)

		}
	} else {
		fmt.Println(err)
	}

	return processesFs
}
