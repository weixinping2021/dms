package utils

import (
	"log"
	"os"
	"path/filepath"
)

type Connection struct {
	Key      string `json:"key"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
}

func GetWorkDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error getting user home directory:", err)
	}

	configFile := filepath.Join(homeDir, "dms.txt")

	data, err := os.ReadFile(configFile)
	return string(data)
}

func GetCondir() string {
	dir := GetWorkDir()
	return filepath.Join(dir, "cons")
}

func GetReisdir() string {
	dir := GetWorkDir()
	return filepath.Join(dir, "redis")
}
