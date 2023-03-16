package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	config Config
)

type Config struct {
	AppName  string `ini:"app_name"`
	LogLevel string `ini:"log_level"`
	IP       string `ini:"ip"`
	Port     string `ini:"port"`
	Database `ini:"database"`
	App      `ini:"app"`
}

type Database struct {
	Driver   string `ini:"driver"`
	IP       string `ini:"ip"`
	Port     string `ini:"port"`
	Username string `ini:"username"`
	Password string `ini:"password"`
	Database string `ini:"database"`
	Sslmode  string `ini:"sslmode"`
}

type App struct {
	ExportSavePath string `ini:"export_save_path"`
}

func Init() (err error) {

	//加载系统配置
	cfg, err := ini.Load("./config/app.ini")
	if err != nil {
		return
	}
	err = cfg.MapTo(&config)

	if err != nil {
		fmt.Println("cfg.MapTo err :", err)
	}
	return
}

func GetConfigObj() Config {
	return config
}

func GetHostConfig() string {
	return config.IP + ":" + config.Port
}

func GetDbConfig() string {
	db := config.Database
	db_config := db.Driver + "://" + db.Username + ":" + db.Password + "@" + db.IP + ":" + db.Port + "/" + db.Database + "?sslmode=" + db.Sslmode
	return db_config
}

func GetExcelPath() string {
	return config.App.ExportSavePath
}
