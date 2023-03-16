package model

type GoQueryData struct {
	DataModel
	Title   string `gorm:"column:title;comment:标题" json:"title"`     //标题
	Url     string `gorm:"column:url;comment:网址" json:"url"`         //网址
	Content string `gorm:"column:content;comment:内容" json:"content"` //内容
}
