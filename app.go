package main

import (
	"context"
	utils "dms/backend"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
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
func (a *App) AddCon(con utils.Connection) string {

	//需要进行判断,如果存在同名的连接,直接返回错误
	fmt.Println("打印返回的连接信息")
	fmt.Println(con)
	_, err := os.Stat("./cons/" + con.Name)
	// 如果返回的错误是不存在，则文件不存在
	if os.IsNotExist(err) {
		new_con := utils.Connection(con)

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
func (a *App) DeleteCon(con utils.Connection) string {

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

func (a *App) GetFullCons() []utils.Connection {

	//判断存储连接信息的目录是否存在
	_, err := os.Stat("./cons")
	if os.IsNotExist(err) {
		// 如果不存在，则创建目录
		err := os.MkdirAll("./cons", os.ModePerm)
		if err != nil {
			fmt.Errorf("创建目录失败: %w", err)
		}
		fmt.Println("目录已创建:", "./cons")
		return make([]utils.Connection, 0)
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
		cons := make([]utils.Connection, len(entries))
		// 遍历目录中的条目
		i := 0
		for _, entry := range entries {
			// 如果是文件（不是目录）
			if !entry.IsDir() {
				fmt.Println("文件:", entry.Name())
				// 在这里进行你的处理逻辑
				dataEncoded, _ := os.ReadFile("./cons/" + entry.Name())
				m := new(utils.Connection)
				json.Unmarshal(dataEncoded, &m)
				cons[i] = *m
				i = i + 1
			}
		}
		return cons[:]
	}
	return nil
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
