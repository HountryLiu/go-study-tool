//go:generate swag init -p pascalcase --parseDependency --output ./docs
package main

import (
	"flag"

	"github.com/HountryLiu/go-study-tool/config"
	"github.com/HountryLiu/go-study-tool/model"
	global_router "github.com/HountryLiu/go-study-tool/router"
	"github.com/HountryLiu/go-study-tool/utils"
	"github.com/gin-gonic/gin"
)

var (
	addr   *string
	dbPath *string
)

func main() {
	// 初始化系统配置
	if err := config.Init(); err != nil {
		panic(err)
	}
	addr = flag.String("s", config.GetHostConfig(), "server address")
	dbPath = flag.String("db", config.GetDbConfig(), "database path")

	//校验器初始化
	if err := utils.InitValidate(); err != nil {
		panic(err)
	}

	// 初始化数据库
	if err := model.Init(*dbPath); err != nil {
		panic(err)
	}

	router := gin.Default()
	global_router.InitRouter(router)

	router.Run(*addr)
}
