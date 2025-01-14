package main

import (
	"embed"

	"dms/backend/mysql"
	"dms/backend/redis"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	mysql := mysql.NewMysql()
	redis := redis.NewRedis()
	// Create application with options
	err := wails.Run(&options.App{
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
