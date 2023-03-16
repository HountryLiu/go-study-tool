package model

type ValidateData struct {
	ID    uint   `gorm:"column:id;primary_key;comment:id" json:"id"`                //id
	Name  string `gorm:"column:name;comment:姓名" json:"name" validate:"required"`    //姓名
	IP    string `gorm:"column:ip;comment:ip" json:"ip"  validate:"ipv4"`           //ip
	Email string `gorm:"column:email;comment:email" json:"email"  validate:"email"` //email
	Data  string `gorm:"column:data;comment:data" json:"data"`                      //data
}

func (o *ValidateData) Create() (err error) {
	if ret := db.Create(o); ret.Error != nil {
		err = ret.Error
		return
	}
	return
}
