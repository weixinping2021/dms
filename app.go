package main

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hdt3213/rdb/bytefmt"
	"github.com/hdt3213/rdb/model"
	"github.com/hdt3213/rdb/parser"
	"github.com/percona/go-mysql/query"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// MemoryUsage represents the memory usage of a prefix
type MemoryUsage struct {
	Prefix string
	Memory int64
}
type ByMemory []MemoryUsage

func (a ByMemory) Len() int           { return len(a) }
func (a ByMemory) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByMemory) Less(i, j int) bool { return a[i].Memory > a[j].Memory } // Descending order

type Record struct {
	Row    []string
	SortBy float64
}

// App struct
type App struct {
	ctx context.Context
}

type RedisKey struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Expire string `json:"expire"`
	Size   string `json:"size"`
}

type Connection struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type RedisMomery struct {
	MomeryForever int `json:"momeryForever" default:"0"` //不过期的内存
	CountForever  int `json:"countForever" default:"0"`  //不过期的key的数量
	Momerys3      int `json:"momerys3" default:"0"`      //过期时间小于3天的内存占用
	Counts3       int `json:"counts3" default:"0"`       //过期时间小于3天的key的数量
	Momeryb3s7    int `json:"momeryb3s7" default:"0"`    //过期时间大于3小于7的内存占用
	Countb3s7     int `json:"countb3s7" default:"0"`     //过期时间大于3小于7的key的数量
	Momeryb7      int `json:"momeryb7" default:"0"`      //过期时间大于7天的内存占用
	Countb7       int `json:"countb7" default:"0"`       //过期时间大于7天的key的数量
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

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// return all the connetions user stored
// 字段包括 连接名称 name, 主机ip,端口,账号,密码
func (a *App) GetCons() []map[string]string {

	//判断存储连接信息的目录是否存在
	_, err := os.Stat("./cons")
	if os.IsNotExist(err) {
		// 如果不存在，则创建目录
		err := os.MkdirAll("./cons", os.ModePerm)
		if err != nil {
			fmt.Errorf("创建目录失败: %w", err)
		}
		fmt.Println("目录已创建:", "./cons")
		return make([]map[string]string, 0)
	} else if err != nil {
		fmt.Errorf("检查目录失败: %w", err)
	} else {
		fmt.Println("目录已存在:", "./cons")
		entries, err := os.ReadDir("./cons")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(len(entries))
		cons := make([]map[string]string, len(entries))
		// 遍历目录中的条目
		i := 0
		for _, entry := range entries {
			// 如果是文件（不是目录）
			if !entry.IsDir() {
				fmt.Println("文件:", entry.Name())
				// 在这里进行你的处理逻辑
				cons[i] = make(map[string]string, 1)
				cons[i]["value"] = entry.Name()
				cons[i]["label"] = entry.Name()
				i = i + 1
			}
		}
		return cons[:]
	}
	return nil
}

// 把单个连接存储成单个文件,并且按照连接的名字命名
func (a *App) AddCon(con Connection) string {

	//需要进行判断,如果存在同名的连接,直接返回错误
	fmt.Println("打印返回的连接信息")
	fmt.Println(con)
	_, err := os.Stat("./cons/" + con.Name)
	// 如果返回的错误是不存在，则文件不存在
	if os.IsNotExist(err) {
		new_con := Connection(con)

		data, _ := json.Marshal(new_con)

		err := os.WriteFile("./cons/"+con.Name, data, 0644)

		if err != nil {
			panic(err)
		}
		return "succsess"
	} else {
		return "duplicate"
	}

}

// 把单个连接存储成单个文件,并且按照连接的名字命名
func (a *App) DeleteCon(con Connection) string {

	fmt.Println(con.Name)
	//需要进行判断,如果存在同名的连接,直接返回错误
	err := os.Remove("./cons/" + con.Name)
	// 如果返回的错误是不存在，则文件不存在
	if err != nil {
		fmt.Println(err)
		return "error"
	} else {
		return "succsess"
	}

}

func (a *App) GetFullCons() []Connection {

	//判断存储连接信息的目录是否存在
	_, err := os.Stat("./cons")
	if os.IsNotExist(err) {
		// 如果不存在，则创建目录
		err := os.MkdirAll("./cons", os.ModePerm)
		if err != nil {
			fmt.Errorf("创建目录失败: %w", err)
		}
		fmt.Println("目录已创建:", "./cons")
		return make([]Connection, 0)
	} else if err != nil {
		// 其他错误
		fmt.Errorf("检查目录失败: %w", err)
	} else {
		fmt.Println("目录已存在:", "./cons")
		entries, err := os.ReadDir("./cons")
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(len(entries))
		cons := make([]Connection, len(entries))
		// 遍历目录中的条目
		i := 0
		for _, entry := range entries {
			// 如果是文件（不是目录）
			if !entry.IsDir() {
				fmt.Println("文件:", entry.Name())
				// 在这里进行你的处理逻辑
				dataEncoded, _ := os.ReadFile("./cons/" + entry.Name())
				m := new(Connection)
				json.Unmarshal(dataEncoded, &m)
				cons[i] = *m
				i = i + 1
			}
		}
		return cons[:]
	}
	return nil
}

func (a *App) GetPeople() []Person {
	return []Person{
		{Name: "Alice", Age: 25},
		{Name: "Bob", Age: 30},
		{Name: "Charlie", Age: 35},
	}
}

// 打开选择文件对话框
func (a *App) OpenDialog() string {

	// 返回选择的文件

	options := runtime.OpenDialogOptions{
		// Filters: ,
	}

	filepath, err := runtime.OpenFileDialog(a.ctx, options)

	if err != nil {
		return "选择文件失败"
	}

	return filepath

}

func (a *App) AnalyseRdb(rdbFileName string) string {

	rdbFile, err := os.Open(rdbFileName)
	defer rdbFile.Close()

	currentTime := time.Now()

	// 自定义格式化: "2022-10-10:10:10:10"
	formattedTime := currentTime.Format("2006-01-02:15:04:05")

	csvKyes := fmt.Sprintf("./redis/%s_keys.csv", formattedTime)

	csvFile, err := os.Create(csvKyes)

	if err != nil {
		fmt.Println("create json %s failed, %v", csvKyes, err)
	}

	defer csvFile.Close()

	var dec = parser.NewDecoder(rdbFile)

	_, err = csvFile.WriteString("database,key,type,size,size_readable,element_count,encoding,expiration,days\n")

	csvWriter := csv.NewWriter(csvFile)

	defer csvWriter.Flush()

	formatExpiration := func(o model.RedisObject) string {
		expiration := o.GetExpiration()
		if expiration == nil {
			return ""
		}
		return expiration.Format(time.RFC3339)
	}

	memory := new(RedisMomery)

	dec.Parse(func(object model.RedisObject) bool {

		//计算内存占用
		days := 0
		redisTime := object.GetExpiration()
		if redisTime != nil {
			// 计算两个时间的差异
			duration := redisTime.Sub(time.Now())
			days = int(duration.Hours() / 24) // 将小时转换为天数
		}

		if days == 0 {
			memory.MomeryForever = memory.MomeryForever + object.GetSize()
			memory.CountForever = memory.CountForever + 1
		} else if days <= 3 && days > 0 {
			memory.Momerys3 = memory.Momerys3 + object.GetSize()
			memory.Counts3 = memory.Counts3 + 1
		} else if days > 3 && days <= 7 {
			memory.Momeryb3s7 = memory.Momeryb3s7 + object.GetSize()
			memory.Countb3s7 = memory.Countb3s7 + 1
		} else if days > 7 {
			memory.Momeryb7 = memory.Momeryb7 + object.GetSize()
			memory.Countb7 = memory.Countb7 + 1
		}

		err = csvWriter.Write([]string{
			strconv.Itoa(object.GetDBIndex()),
			object.GetKey(),
			object.GetType(),
			strconv.Itoa(object.GetSize()),
			bytefmt.FormatSize(uint64(object.GetSize())),
			strconv.Itoa(object.GetElemCount()),
			object.GetEncoding(),
			formatExpiration(object),
			fmt.Sprintf("%d", days),
		})
		if err != nil {
			fmt.Printf("csv write failed: %v", err)
			return false
		}
		return true
	})
	fmt.Print(memory)
	// 将结构体的格式，转为json字符串的格式。这里用的到库包是"github.com/pquerna/ffjson/ffjson"
	data, err := json.Marshal(memory)
	if err != nil {
		fmt.Println(err)
	}
	// 将json格式的数据写入文件
	err = os.WriteFile(fmt.Sprintf("./redis/%s_memory.csv", formattedTime), data, 0777)
	if err != nil {
		fmt.Println(err)
	}
	return formattedTime
}

func (a *App) GetRedisMemory(formattedTime string) RedisMomery {
	// 打开 JSON 文件
	file, err := os.Open(fmt.Sprintf("./redis/%s_memory.csv", formattedTime))
	if err != nil {
		fmt.Println("无法打开文件:", err)
	}
	defer file.Close()

	// 创建解码器并解析 JSON 文件
	var memory RedisMomery
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&memory)
	if err != nil {
		fmt.Println("JSON 解析错误:", err)
	}
	return memory
}

func (a *App) GetRedisKeys(filed string, formattedTime string) []RedisKey {

	file, err := os.Open(fmt.Sprintf("./redis/%s_keys.csv", formattedTime))
	if err != nil {
		fmt.Println("无法打开文件:", err)
	}
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 读取表头
	headers, err := reader.Read()
	if err != nil {
		fmt.Println("读取表头失败:", err)
	}

	// 确定需要排序的字段（假设字段为 "score"）
	sortField := filed
	sortFieldIndex := -1

	for i, header := range headers {
		if header == sortField {
			sortFieldIndex = i
			break
		}
	}

	if sortFieldIndex == -1 {
		fmt.Printf("字段 %s 未找到\n", sortField)
	}

	// 定义用于存储排序数据的切片
	var records []Record

	// 逐行读取数据
	for {
		row, err := reader.Read()
		if err != nil {
			// 如果读取到文件末尾，则退出循环
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("读取行失败:", err)
		}

		// 将指定字段转换为数值
		sortValue, err := strconv.ParseFloat(row[sortFieldIndex], 64)
		if err != nil {
			// 如果转换失败，设为默认值 0
			sortValue = 0
		}

		// 将行数据添加到切片中
		records = append(records, Record{
			Row:    row,
			SortBy: sortValue,
		})
	}

	// 按排序字段降序排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].SortBy > records[j].SortBy
	})

	// 输出表头
	fmt.Println(headers)

	// 输出前 500 行
	limit := 500
	if len(records) < limit {
		limit = len(records)
	}

	rediskeys := make([]RedisKey, limit)

	for i := 0; i < limit; i++ {
		number := strconv.Itoa(i)
		size, _ := strconv.ParseInt(records[i].Row[3], 10, 64)
		return_size := formatBytes(size)
		key := RedisKey{number, records[i].Row[1], records[i].Row[2], records[i].Row[7], return_size}
		rediskeys[i] = key
		fmt.Println(records[i].Row)
	}
	return rediskeys[:]

}

func (a *App) GetPrefixkeys(prefix string, formattedTime string) []RedisKey {
	// 打开 CSV 文件
	file, err := os.Open(fmt.Sprintf("./redis/%s_keys.csv", formattedTime))
	if err != nil {
		fmt.Println("无法打开文件:", err)
	}
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 读取表头
	headers, err := reader.Read()
	if err != nil {
		fmt.Println("读取表头失败:", err)
	}

	// 确定需要排序的字段（假设字段为 "score"）
	sortField := "size"
	sortFieldIndex := -1

	for i, header := range headers {
		if header == sortField {
			sortFieldIndex = i
			break
		}
	}

	if sortFieldIndex == -1 {
		fmt.Printf("字段 %s 未找到\n", sortField)
	}

	// 定义用于存储排序数据的切片
	var records []Record

	// 逐行读取数据
	for {
		row, err := reader.Read()
		if err != nil {
			// 如果读取到文件末尾，则退出循环
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("读取行失败:", err)
		}
		if strings.HasPrefix(row[1], prefix) {
			// 将指定字段转换为数值
			sortValue, err := strconv.ParseFloat(row[sortFieldIndex], 64)
			if err != nil {
				// 如果转换失败，设为默认值 0
				sortValue = 0
			}

			// 将行数据添加到切片中

			records = append(records, Record{
				Row:    row,
				SortBy: sortValue,
			})
		} else {
			fmt.Println("Goodbye")
		}

	}

	// 按排序字段降序排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].SortBy > records[j].SortBy
	})

	// 输出表头
	fmt.Println(headers)

	// 输出前 500 行
	limit := 500
	if len(records) < limit {
		limit = len(records)
	}

	rediskeys := make([]RedisKey, limit)

	for i := 0; i < limit; i++ {
		number := strconv.Itoa(i)
		size, _ := strconv.ParseInt(records[i].Row[3], 10, 64)
		return_size := formatBytes(size)
		key := RedisKey{number, records[i].Row[1], records[i].Row[2], records[i].Row[7], return_size}
		rediskeys[i] = key
		fmt.Println(records[i].Row)
	}
	return rediskeys[:]

}

func (a *App) GetRedisTop500Prefix(formattedTime string) []RedisKey {

	file, err := os.Open(fmt.Sprintf("./redis/%s_keys.csv", formattedTime))
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	// Create a map to store memory sizes by key prefix
	prefixMemory := make(map[string]int64)
	prefixCount := make(map[string]int)

	// Create a CSV reader to read the file in chunks
	reader := csv.NewReader(bufio.NewReader(file))

	// Read the header line (if any)
	_, err = reader.Read() // Skip the header if there is one
	if err != nil {
		fmt.Println("Error reading header:", err)
	}

	// Process the CSV file line by line
	for {
		record, err := reader.Read()
		if err != nil {
			// End of file or error
			if err.Error() != "EOF" {
				fmt.Println("Error reading record:", err)
			}
			break
		}

		// Check if the row has the necessary data
		if len(record) < 2 {
			continue // Skip rows that don"t have enough data
		}

		// Extract the key and memory size
		key := record[1]
		memSize, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			fmt.Println("Invalid memory size:", record[1])
			continue
		}

		// Split the key and extract the prefix
		split := splitKey(key)
		if len(split) > 0 {
			prefix := split[0] // The first part is considered the prefix
			// Add memory size to the corresponding prefix
			prefixMemory[prefix] += memSize
			prefixCount[prefix] += 1
		}
	}

	// Convert the map to a slice for sorting
	var memoryUsage []MemoryUsage
	for prefix, totalMem := range prefixMemory {
		memoryUsage = append(memoryUsage, MemoryUsage{Prefix: prefix, Memory: totalMem})
	}
	// Sort the memory usage by memory size (descending order)
	sort.Sort(ByMemory(memoryUsage))
	limit := 500
	if len(memoryUsage) < limit {
		limit = len(memoryUsage)
	}
	i := 0
	rediskeys := make([]RedisKey, limit)
	// Print the memory usage by prefix, sorted by memory
	fmt.Println("Memory usage by prefix (sorted):")

	for _, usage := range memoryUsage {
		number := strconv.Itoa(i)
		size := formatBytes(usage.Memory)
		count := strconv.Itoa(prefixCount[usage.Prefix])
		key := RedisKey{number, usage.Prefix, count, "", size}
		rediskeys[i] = key
		i = i + 1
		if i == 500 {
			break
		}
		fmt.Printf("Prefix: %s, Total Memory: %d bytes\n", usage.Prefix, usage.Memory)
	}
	return rediskeys[:]
}

func (a *App) GetRdbResultTitle() []map[string]string {
	dir := "./redis" // 当前目录，你可以根据需要修改

	// 读取目录下的文件
	files, err := os.ReadDir(dir)
	fmt.Println(len(files) / 2)
	filesPrefix := make([]map[string]string, len(files)/2)
	if err != nil {
		fmt.Println(err)
	}

	i := 0
	// 循环输出文件名
	for _, file := range files {
		if strings.Contains(file.Name(), "keys.csv") {
			// 只输出文件名
			filesPrefix[i] = map[string]string{"value": strings.Split(file.Name(), "_")[0], "label": strings.Split(file.Name(), "_")[0]}
			i = i + 1
		}
		fmt.Println(file.Name())
	}

	return filesPrefix
}

func formatBytes(size int64) string {
	const (
		_  = iota // ignore first value by assigning to underscore
		KB = 1 << (10 * iota)
		MB
		GB
		TB
	)

	switch {
	case size >= TB:
		return fmt.Sprintf("%.2f TB", float64(size)/TB)
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

func splitKey(key string) []string {
	// Define the delimiters
	delimiters := ":;,_-+@=|#"
	return strings.FieldsFunc(key, func(r rune) bool {
		return strings.ContainsRune(delimiters, r)
	})
}

// 获取数据库连接,可以分别获取活跃或者非活跃的
func (a *App) GetMysqlProcesslist(dbId string, sleep string) []MysqlProcessF {
	dataEncoded, _ := os.ReadFile("./cons/" + dbId)
	m := new(Connection)
	json.Unmarshal(dataEncoded, &m)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema", m.User, m.Password, m.Host, m.Port) //替换为你自己的数据库信息
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

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

	// 遍历查询结果,这里就需要区分了,如果Command是query的,然后在用go-mysql里面的指纹插件搞成一组
	/*
		q := "SELECT c FROM t WHERE a=1 AND b='foo'\n" + ORDER BY c ASC LIMIT 1"
		f := query.Fingerprint(q)
		fmt.Println(f)
	*/
	//var processes []MysqlProcess
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

func (a *App) KillMysqlProcesses(dbId string, processes []MysqlProcessF) string {
	//fmt.Println(len(processes))
	dataEncoded, _ := os.ReadFile("./cons/" + dbId)
	m := new(Connection)
	json.Unmarshal(dataEncoded, &m)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema", m.User, m.Password, m.Host, m.Port) //替换为你自己的数据库信息
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

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

func (a *App) GetConsStatus(dbId string) map[string][]MysqlConsStatus {

	status := make(map[string][]MysqlConsStatus)
	dataEncoded, _ := os.ReadFile("./cons/" + dbId)
	m := new(Connection)
	json.Unmarshal(dataEncoded, &m)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema", m.User, m.Password, m.Host, m.Port) //替换为你自己的数据库信息
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

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

func (a *App) GetConspercent(dbId string) map[string]string {

	dataEncoded, _ := os.ReadFile("./cons/" + dbId)
	m := new(Connection)
	json.Unmarshal(dataEncoded, &m)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema", m.User, m.Password, m.Host, m.Port) //替换为你自己的数据库信息
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

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

func (a *App) GetMysqlLock(dbId string) []MysqlProcessF {
	dataEncoded, _ := os.ReadFile("./cons/" + dbId)
	m := new(Connection)
	json.Unmarshal(dataEncoded, &m)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema", m.User, m.Password, m.Host, m.Port) //替换为你自己的数据库信息
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

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
