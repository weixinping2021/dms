package redis

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hdt3213/rdb/bytefmt"
	"github.com/hdt3213/rdb/model"
	"github.com/hdt3213/rdb/parser"
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
type Redis struct {
	ctx context.Context
}

type RedisKey struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Expire string `json:"expire"`
	Size   string `json:"size"`
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

// NewApp creates a new App application struct
func NewRedis() *Redis {
	return &Redis{}
}

func (a *Redis) AnalyseRdb(rdbFileName string) string {

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

func (a *Redis) GetRedisMemory(formattedTime string) RedisMomery {
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

func (a *Redis) GetRedisKeys(filed string, formattedTime string) []RedisKey {

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

func (a *Redis) GetPrefixkeys(prefix string, formattedTime string) []RedisKey {
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

func (a *Redis) GetRedisTop500Prefix(formattedTime string) []RedisKey {

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

func (a *Redis) GetRdbResultTitle() []map[string]string {
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
