package main

import (
	"context"
	utils "dms/backend"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx     context.Context
	workDir string
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

func (a *App) SetWorkDir(workdir string) string {
	// 获取当前用户的主目录
	homeDir, _ := os.UserHomeDir()
	configFile := filepath.Join(homeDir, "dms.txt")

	if a.workDir == "" {
		err := os.WriteFile(configFile, []byte(workdir), 0644)
		if err != nil {
			return err.Error()
		}
	}
	return "success"
}

func (a *App) GetWorkDir() (string, error) {
	//需要查看是否存储目录,如果不存在就报错
	_, err := os.Stat(a.workDir)

	if err != nil {
		return "", err // 返回错误
	}
	return a.workDir, nil
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
	conDir := utils.GetCondir()

	os.MkdirAll(conDir, os.ModePerm)

	entries, _ := os.ReadDir(conDir)

	var cons []utils.Connection
	// 遍历目录中的条目
	for _, entry := range entries {
		// 如果是文件（不是目录）
		if !entry.IsDir() {
			//fmt.Println("文件:", entry.Name())
			// 在这里进行你的处理逻辑
			filepath := path.Join(conDir, entry.Name())
			dataEncoded, _ := os.ReadFile(filepath)
			m := new(utils.Connection)
			json.Unmarshal(dataEncoded, &m)
			//fmt.Println(m)
			cons = append(cons, *m)
		}
	}
	//fmt.Print(cons)
	return cons[:]
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

func (a *App) OpenDir() string {
	// 返回选择的文件

	options := runtime.OpenDialogOptions{
		// Filters: ,
	}

	dirpath, err := runtime.OpenDirectoryDialog(a.ctx, options)

	if err != nil {
		return "选择文件失败"
	}

	return dirpath
}
