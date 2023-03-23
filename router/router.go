package router

import (
	"github.com/HountryLiu/go-study-tool/controller/excel"
	"github.com/HountryLiu/go-study-tool/controller/ffmpeg"
	"github.com/HountryLiu/go-study-tool/controller/g_json"
	"github.com/HountryLiu/go-study-tool/controller/go_query"
	"github.com/HountryLiu/go-study-tool/controller/validate"
	_ "github.com/HountryLiu/go-study-tool/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(router *gin.Engine) {
	//swagger集成
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := router.Group("/api")

	//excel操作学习
	api_excel := api.Group("/excel")
	excel.Router(api_excel)

	//validate学习
	api_validate := api.Group("/validate")
	validate.Router(api_validate)

	//gjson学习
	api_gjson := api.Group("/gjson")
	g_json.Router(api_gjson)

	//goquery学习
	api_goquery := api.Group("/goquery")
	go_query.Router(api_goquery)

	//ffmpeg学习
	api_ffmpeg := api.Group("/ffmpeg")
	ffmpeg.Router(api_ffmpeg)
}
