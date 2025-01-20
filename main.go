package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"dms/backend/mysql"
	"dms/backend/redis"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	//先创建dms的配置文件,如果有就不创建了
	var workdir string
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting user home directory:", err)
	}

	configFile := filepath.Join(homeDir, "dms.txt")

	// 使用 os.Stat 获取文件信息
	_, err = os.Stat(configFile)

	if err != nil {
		workdir = ""
	} else {
		// 读取文件内容
		data, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatal(err)
		}
		// 打印文件内容
		fmt.Println(string(data))
		workdir = string(data)
	}
	// Create an instance of the app structure
	app := NewApp()
	app.workDir = workdir
	mysql := mysql.NewMysql()
	redis := redis.NewRedis()
	// Create application with options
	err = wails.Run(&options.App{
		Title:  "dms",
		Width:  1500,
		Height: 1000,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			mysql,
			redis,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
