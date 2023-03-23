package utils

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
)

var (
	Validate      *validator.Validate
	validateTrans ut.Translator
)

// 表单验证
func InitValidate() (err error) {
	zh := zh.New() //中文翻译器

	// 第一个参数是必填，如果没有其他的语言设置，就用这第一个
	// 后面的参数是支持多语言环境（
	// uni := ut.New(en, en) 也是可以的
	// uni := ut.New(en, zh, tw)
	uni := ut.New(zh)
	validateTrans, _ = uni.GetTranslator("zh") //获取需要的语言
	Validate = validator.New()
	err = zhtrans.RegisterDefaultTranslations(Validate, validateTrans)
	return
}
