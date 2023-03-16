package excel

import (
	"github.com/gin-gonic/gin"
)

const (
	CSV  = "csv"
	XLSX = "xlsx"
	XLS  = "xls"
)

func Router(router *gin.RouterGroup) {
	router.GET("/export", Export)
	router.POST("/import", Import)
}
