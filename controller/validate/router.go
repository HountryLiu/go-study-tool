package validate

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup) {
	router.POST("/create", Create)
}
