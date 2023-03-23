package model

type ExcelData struct {
	DataModel
	Name   string  `gorm:"column:name;comment:姓名" json:"name"`                      //姓名
	Sex    string  `gorm:"column:sex;comment:性别(1:男 2:女)" json:"sex"`               //性别
	Age    int     `gorm:"column:age;comment:年龄" json:"age"`                        //年龄
	IDCard string  `gorm:"column:id_card;comment:身份证号" json:"id_card"`              //身份证号
	Money  float64 `gorm:"column:money;type:decimal(10,2);comment:余额" json:"money"` //余额,保留2位小数
}

func (o *ExcelData) List(cur_page int, page_size int) (data map[string]interface{}, err error) {
	var excel_data []ExcelData
	var count int64
	offset := page_size * (cur_page - 1)
	tx := db.Model(DBExcelDataModel)
	if ret := tx.Count(&count); ret.Error != nil {
		err = ret.Error
		return
	}
	if ret := tx.Offset(offset).Limit(page_size).Find(&excel_data); ret.Error != nil {
		err = ret.Error
		return
	}
	data = map[string]interface{}{
		"total": count,
		"data":  excel_data,
	}
	return
}

func (o *ExcelData) Create() (err error) {
	if ret := db.Create(o); ret.Error != nil {
		err = ret.Error
		return
	}
	return
}
