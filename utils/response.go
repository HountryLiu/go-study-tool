package utils

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	Normal = 0
)

// 系统错误码
const (
	Unauthorized = -1 - iota
	InvalidAuthorization
	MethodNotAllowed
	InternalServerError
	ForwardError
	FileSystemError
	// ...
)

// 10001起，输入参数相关的错误
const (
	ParamsRequired = 10001 + iota
	UnmarshalError
	UnSupportFileType
	// ...
)

// 20001起，数据库操作相关的错误
const (
	RecordDuplicated = 20001 + iota
	DatabaseSelectError
	DatabaseDeleteError
	DatabaseUpdateError
	DatabaseInsertError
	// ...
)

// 30001起，逻辑处理过程的错误
const (
	CannotModify = 30001 + iota
	CannotDelete
	// ...
)

var statusText = map[int]string{
	Normal: "成功",

	Unauthorized:         "未授权/登陆",
	InvalidAuthorization: "无效的验证信息",
	MethodNotAllowed:     "不支持请求的方法",
	InternalServerError:  "内部错误",
	ForwardError:         "前置出错, 前置动作未完成/请求上级出错/必要的初始化未完成等",
	FileSystemError:      "读写文件出错",

	ParamsRequired:    "缺少必要的参数",
	UnmarshalError:    "数据反序列化失败",
	UnSupportFileType: "不支持该文件类型",

	RecordDuplicated:    "记录重复/违反唯一约束",
	DatabaseSelectError: "数据库查询错误",
	DatabaseDeleteError: "数据库删除错误",
	DatabaseUpdateError: "数据库更新错误",
	DatabaseInsertError: "数据库插入错误",

	CannotModify: "不可修改",
	CannotDelete: "不可删除",
}

func Text(code int) string {
	return statusText[code]
}

func Success(ctx *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"no":  Normal,
			"msg": Text(Normal),
		})
	} else if len(data) == 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"no":   Normal,
			"data": data[0],
			"msg":  Text(Normal),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"no":   Normal,
			"data": data,
			"msg":  Text(Normal),
		})
	}

}

func Error(ctx *gin.Context, no int, err error) {
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"no":  no,
			"msg": Text(no),
		})
	} else {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ctx.JSON(http.StatusOK, gin.H{
				"no":     no,
				"errors": err.Error(),
				"msg":    Text(no),
			})
		} else {
			// validator.ValidationErrors类型错误则进行翻译
			removeStructName := func(fields map[string]string) string {
				result := []string{}

				for _, err := range fields {
					result = append(result, err)
				}
				return strings.Join(result, ", ")
			}
			ctx.JSON(http.StatusOK, gin.H{
				"no":     no,
				"errors": removeStructName(errs.Translate(validateTrans)),
				"msg":    Text(no),
			})
		}
	}
}
