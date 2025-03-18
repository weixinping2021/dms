package redis

import (
	"bufio"
	"context"
	utils "dms/backend"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
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
	Memory int
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
	Key string `json:"key"`
	//	Name         string `json:"name"`
	Type         string `json:"type"`
	Expire       string `json:"expire"`
	SizeReadable string `json:"sizereadable"`
	Size         int    `json:"size"`
	ElementCount string `json:"elementcount"`
	Db           string `json:"db"`
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type RedisMomery struct {
	MomeryForever int `json:"momeryForever" default:"0"` //不过期的内存
	CountForever  int `json:"countForever" default:"0"`  //不过期的key的数量
	MomeryExpire  int `json:"momeryExpire" default:"0"`  //过期时间小于3天的内存占用
	CountExpire   int `json:"countExpire" default:"0"`   //过期时间小于3天的key的数量
}

// NewApp creates a new App application struct
func NewRedis() *Redis {
	return &Redis{}
}

// 分析redis rdb文件
func (a *Redis) AnalyseRdb(rdbFileName string) string {

	rdbFile, _ := os.Open(rdbFileName)
	defer rdbFile.Close()

	currentTime := time.Now()

	// 自定义格式化: "2022-10-10:10:10:10"
	formattedTime := currentTime.Format("2006-01-02:15:04:05")

	redisPath := utils.GetReisdir()

	redisFileName := path.Join(redisPath, fmt.Sprintf("%s_keys.csv", formattedTime))

	csvFile, _ := os.Create(redisFileName)
	defer csvFile.Close()

	var dec = parser.NewDecoder(rdbFile)

	_, err := csvFile.WriteString("database,key,type,size,size_readable,element_count,encoding,expiration,seconds\n")

	csvWriter := csv.NewWriter(csvFile)

	defer csvWriter.Flush()

	formatExpiration := func(o model.RedisObject) string {
		expiration := o.GetExpiration()
		if expiration == nil {
			return ""
		}
		return expiration.Format("2006-01-02 15:04:05")
	}

	memory := new(RedisMomery)

	dec.Parse(func(object model.RedisObject) bool {

		//计算内存占用
		seconds := 0
		redisTime := object.GetExpiration()
		if redisTime == nil {
			memory.MomeryForever = memory.MomeryForever + int(object.GetSize())
			memory.CountForever = memory.CountForever + 1
		} else {
			memory.MomeryExpire = memory.MomeryExpire + int(object.GetSize())
			memory.CountExpire = memory.CountExpire + 1
			seconds = redisTime.Second()
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
			strconv.Itoa(seconds),
		})

		if err != nil {
			log.Printf("csv write failed: %v", err)
			return false
		}
		return true
	})

	data, err := json.Marshal(memory)
	if err != nil {
		fmt.Println(err)
	}
	// 将json格式的数据写入文件
	redisMemoryFile := path.Join(redisPath, fmt.Sprintf("%s_memory.csv", formattedTime))
	err = os.WriteFile(redisMemoryFile, data, 0777)
	if err != nil {
		log.Println(err)
	}
	return formattedTime
}

func (a *Redis) GetForverMemoryPic(formattedTime string) [][]string {
	redisPath := utils.GetReisdir()

	redisFileName := path.Join(redisPath, fmt.Sprintf("%s_keys.csv", formattedTime))

	file, err := os.Open(redisFileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	// Create a map to store memory sizes by key prefix
	prefixMemory := make(map[string]int)
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
		memSize, err := strconv.Atoi(record[3])
		if err != nil {
			fmt.Println("Invalid memory size:", record[1])
			continue
		}

		if record[7] == "" {
			split := splitKey(key)
			if len(split) > 1 {
				prefix := split[0] // The first part is considered the prefix
				// Add memory size to the corresponding prefix
				prefixMemory[prefix] += memSize
				prefixCount[prefix] += 1
			} else {
				prefixMemory["无分隔符"] += memSize
				prefixCount["无分隔符"] += 1
			}
		}
		// Split the key and extract the prefix

	}

	// Convert the map to a slice for sorting
	var memoryUsage []MemoryUsage
	for prefix, totalMem := range prefixMemory {
		memoryUsage = append(memoryUsage, MemoryUsage{Prefix: prefix, Memory: totalMem})
	}
	// Sort the memory usage by memory size (descending order)
	sort.Sort(ByMemory(memoryUsage))
	limit := 30
	if len(memoryUsage) < limit {
		limit = len(memoryUsage)
	}
	i := 0

	// Print the memory usage by prefix, sorted by memory
	fmt.Println("Memory usage by prefix (sorted):")
	forverMemorys := make([][]string, 3)
	for _, usage := range memoryUsage {
		//number := strconv.Itoa(i)
		//size := formatBytes(usage.Memory)
		count := strconv.Itoa(prefixCount[usage.Prefix])
		//key := RedisKey{number, usage.Prefix, count, "", size, keyDate}
		//rediskeys[i] = key
		forverMemorys[0] = append(forverMemorys[0], usage.Prefix)
		forverMemorys[1] = append(forverMemorys[1], count)
		forverMemorys[2] = append(forverMemorys[2], fmt.Sprintf("%.2f", float64(usage.Memory)/1024/1024))

		i = i + 1
		if i == 30 {
			break
		}
		fmt.Printf("Prefix: %s, Total Memory: %d bytes\n", usage.Prefix, usage.Memory)
	}
	return forverMemorys
}

func (a *Redis) GetRedisMemory(formattedTime string) []map[string]string {
	// 打开 JSON 文件

	redisPath := utils.GetReisdir()

	redisFileName := path.Join(redisPath, fmt.Sprintf("%s_memory.csv", formattedTime))

	file, _ := os.Open(redisFileName)
	defer file.Close()

	// 创建解码器并解析 JSON 文件
	var memoryBase RedisMomery
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&memoryBase)

	if err != nil {
		log.Println("JSON 解析错误:", err)
	}

	memory := make([]map[string]string, 2)
	memory[0] = map[string]string{"key": "1", "tjfl": "内存统计", "forever": formatBytes(memoryBase.MomeryForever), "expire": formatBytes(memoryBase.MomeryExpire), "total": formatBytes(memoryBase.MomeryForever + memoryBase.MomeryExpire)}
	memory[1] = map[string]string{"key": "2", "tjfl": "数量统计", "forever": strconv.Itoa(memoryBase.CountForever), "expire": strconv.Itoa(memoryBase.CountExpire), "total": strconv.Itoa(memoryBase.CountForever + memoryBase.CountExpire)}
	return memory
}

func (a *Redis) GetExpirePrefixs(formattedTime string, keyDate string) []RedisKey {

	redisPath := utils.GetReisdir()
	redisFileName := path.Join(redisPath, fmt.Sprintf("%s_keys.csv", formattedTime))
	file, err := os.Open(redisFileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
	}
	defer file.Close()

	// Create a map to store memory sizes by key prefix
	prefixMemory := make(map[string]int)
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

		// Extract the key and memory size
		key := record[1]
		memSize, err := strconv.Atoi(record[3])
		if err != nil {
			fmt.Println("Invalid memory size:", record[1])
			continue
		}
		if keyDate != "" && record[7] != "" && record[7][0:10] == keyDate {
			split := splitKey(key)
			if len(split) > 1 {
				prefix := split[0] // The first part is considered the prefix
				// Add memory size to the corresponding prefix
				prefixMemory[prefix] += memSize
				prefixCount[prefix] += 1
			} else {
				prefixMemory["无分隔符"] += memSize
				prefixCount["无分隔符"] += 1
			}
		}
		// Split the key and extract the prefix

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
		size := formatBytes(usage.Memory)
		count := strconv.Itoa(prefixCount[usage.Prefix])
		key := RedisKey{usage.Prefix, "", keyDate, size, usage.Memory, count, ""}
		rediskeys[i] = key
		i = i + 1
		if i == 500 {
			break
		}
		fmt.Printf("Prefix: %s, Total Memory: %d bytes\n", usage.Prefix, usage.Memory)
	}
	return rediskeys[:]
}

func (a *Redis) GetExpireMemoryPic(formattedTime string) [][]string {
	redisPath := utils.GetReisdir()
	redisFileName := path.Join(redisPath, fmt.Sprintf("%s_keys.csv", formattedTime))
	// 打开 CSV 文件
	file, _ := os.Open(redisFileName)
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	keyMap := make(map[string]int)
	keyMemoryMap := make(map[string]int)
	// 逐行读取数据,并且根据日期进行统计
	for {
		row, err := reader.Read()
		if err != nil {
			// 如果读取到文件末尾，则退出循环
			if err.Error() == "EOF" {
				break
			}
			log.Println("读取行失败:", err)
		}
		if row[7] == "expiration" {
			continue
		}
		if row[7] != "" {
			value, exists := keyMap[row[7][:10]]
			if exists {
				keyMap[row[7][:10]] = value + 1
				number, _ := strconv.Atoi(row[3])
				keyMemoryMap[row[7][:10]] = keyMemoryMap[row[7][:10]] + number
			} else {
				keyMap[row[7][:10]] = 1
				number, _ := strconv.Atoi(row[3])
				keyMemoryMap[row[7][:10]] = number
			}
		}

	}
	// 提取 map 的键
	keys := make([]string, 0, len(keyMap))
	for key := range keyMap {
		keys = append(keys, key)
	}

	// 按照日期排序（字符串的格式是 YYYY-MM-DD，符合时间排序的规则）
	sort.Slice(keys, func(i, j int) bool {
		date1, _ := time.Parse("2006-01-02", keys[i])
		date2, _ := time.Parse("2006-01-02", keys[j])
		return date1.Before(date2)
	})

	// 打印排序后的结果
	rediskeys := make([][]string, 3)
	rediskeys[0] = make([]string, 0)
	rediskeys[1] = make([]string, 0)
	rediskeys[2] = make([]string, 0)
	for _, key := range keys {
		rediskeys[0] = append(rediskeys[0], key)
		rediskeys[1] = append(rediskeys[1], strconv.Itoa(keyMap[key]))
		rediskeys[2] = append(rediskeys[2], fmt.Sprintf("%.2f", float64(keyMemoryMap[key])/1024/1024))
		//fmt.Println(key)
	}

	//fmt.Println(rediskeys)
	return rediskeys
}

/*
func (a *Redis) GetRedisKeys(filed string, formattedTime string) []RedisKey {

	redisPath := utils.GetReisdir()

	redisFileName := path.Join(redisPath, fmt.Sprintf("%s_keys.csv", formattedTime))

	file, _ := os.Open(redisFileName)
	defer file.Close()

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 读取表头
	headers, _ := reader.Read()

	// 确定需要排序的字段
	sortField := filed
	sortFieldIndex := -1

	for i, header := range headers {
		if header == sortField {
			sortFieldIndex = i
			break
		}
	}

	if sortFieldIndex == -1 {
		log.Printf("字段 %s 未找到\n", sortField)
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
			log.Println("读取行失败:", err)
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
		key := RedisKey{number, records[i].Row[1], records[i].Row[2], records[i].Row[7], return_size, ""}
		rediskeys[i] = key
		//fmt.Println(records[i].Row)
	}
	return rediskeys[:]

}
*/
/*
func (a *Redis) GetPrefixkeys(prefix string, formattedTime string, keyDate string) []RedisKey {

	redisPath := utils.GetReisdir()

	redisFileName := path.Join(redisPath, fmt.Sprintf("%s_keys.csv", formattedTime))
	// 打开 CSV 文件
	file, _ := os.Open(redisFileName)
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
	if prefix == "无分隔符" {
		for {
			row, err := reader.Read()
			if err != nil {
				// 如果读取到文件末尾，则退出循环
				if err.Error() == "EOF" {
					break
				}
				fmt.Println("读取行失败:", err)
			}
			if len(splitKey(row[1])) == 1 && row[7] != "" && row[7][0:10] == keyDate {
				// 将指定字段转换为数值
				sortValue, err := strconv.ParseFloat(row[sortFieldIndex], 64)
				if err != nil {
					// 如果转换失败，设为默认值 0
					sortValue = 0
				}
				records = append(records, Record{
					Row:    row,
					SortBy: sortValue,
				})
			}

		}
	} else {
		for {
			row, err := reader.Read()
			if err != nil {
				// 如果读取到文件末尾，则退出循环
				if err.Error() == "EOF" {
					break
				}
				fmt.Println("读取行失败:", err)
			}
			if len(splitKey(row[1])) > 0 && splitKey(row[1])[0] == prefix && row[7][0:10] == keyDate {
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

		}
	}

	// 按排序字段降序排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].SortBy > records[j].SortBy
	})

	// 输出前 1000 行
	limit := 1000
	if len(records) < limit {
		limit = len(records)
	}

	rediskeys := make([]RedisKey, limit)

	for i := 0; i < limit; i++ {
		number := strconv.Itoa(i)
		size, _ := strconv.ParseInt(records[i].Row[3], 10, 64)
		return_size := formatBytes(size)
		key := RedisKey{number, records[i].Row[1], records[i].Row[2], records[i].Row[7], return_size, ""}
		rediskeys[i] = key
		//fmt.Println(records[i].Row)
	}
	return rediskeys[:]

}
*/
func (a *Redis) GetPrefixDetail(formattedTime string, prefix string, keyDate string) []RedisKey {

	//获取csv读取器
	redisPath := utils.GetReisdir()
	redisFileName := path.Join(redisPath, fmt.Sprintf("%s_keys.csv", formattedTime))
	file, _ := os.Open(redisFileName)
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Read() //过滤第一行,不进行处理,因为是字段名
	rediskeys := make([]RedisKey, 0)

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
		/*  Key          string `json:"key"`
		//	Name         string `json:"name"`
			Type         string `json:"type"`
			Expire       string `json:"expire"`
			SizeReadable string `json:"sizereadable"`
			Size         int    `json:"size"`
		database,key,type,size,size_readable,element_count,encoding,expiration,seconds*/

		if keyDate == "" && row[7] == "" && splitKey(row[1])[0] == prefix {
			size, _ := strconv.Atoi(row[3])
			key := RedisKey{row[1], row[2], row[7], formatBytes(size), size, row[5], row[0]}
			rediskeys = append(rediskeys, key)
		} else if keyDate != "" && row[7] != "" && keyDate == row[7][0:10] && splitKey(row[1])[0] == prefix {
			size, _ := strconv.Atoi(row[3])
			key := RedisKey{row[1], row[2], row[7], formatBytes(size), size, row[5], row[0]}
			rediskeys = append(rediskeys, key)
		} else if prefix == "无分隔符" && keyDate != "" && row[7] != "" && keyDate == row[7][0:10] && len(splitKey(row[1])) == 1 {
			size, _ := strconv.Atoi(row[3])
			key := RedisKey{row[1], row[2], row[7], formatBytes(size), size, row[5], row[0]}
			rediskeys = append(rediskeys, key)
		}
	}

	//按size排序
	sort.Slice(rediskeys, func(i, j int) bool {
		return rediskeys[i].Size > rediskeys[j].Size
	})

	if len(rediskeys) > 500 {
		return rediskeys[:500]
	} else {
		return rediskeys[:]
	}

}

/*
func (a *Redis) GetRedisTop500Prefix(formattedTime string, keyDate string) []RedisKey {

		redisPath := utils.GetReisdir()

		redisFileName := path.Join(redisPath, fmt.Sprintf("%s_keys.csv", formattedTime))

		file, err := os.Open(redisFileName)
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
			if keyDate != "" && record[7] != "" && record[7][0:10] == keyDate {
				split := splitKey(key)
				if len(split) > 1 {
					prefix := split[0] // The first part is considered the prefix
					// Add memory size to the corresponding prefix
					prefixMemory[prefix] += memSize
					prefixCount[prefix] += 1
				} else {
					prefixMemory["无分隔符"] += memSize
					prefixCount["无分隔符"] += 1
				}
			}
			// Split the key and extract the prefix

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
			key := RedisKey{number, usage.Prefix, count, "", size, keyDate}
			rediskeys[i] = key
			i = i + 1
			if i == 500 {
				break
			}
			fmt.Printf("Prefix: %s, Total Memory: %d bytes\n", usage.Prefix, usage.Memory)
		}
		return rediskeys[:]
	}
*/
func (a *Redis) GetRdbResultTitle() []map[string]string {
	dir := utils.GetReisdir()
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

func formatBytes(size int) string {
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
