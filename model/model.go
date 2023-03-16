package model

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ...
var (
	db *gorm.DB

	DBExcelDataModel    = new(ExcelData)
	DBValidateDataModel = new(ValidateData)
	DBGoQueryDataModel  = new(GoQueryData)
)

// DataModel 数据库存储的基类
type DataModel struct {
	ID        uint       `gorm:"column:id;primary_key" json:"id"`     // 自增的id
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"` // 更新时间
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`              // 删除时间
}

var (
	models = []interface{}{
		// 所有需要存入数据库的模型放到这里
		DBExcelDataModel,
		DBValidateDataModel,
		DBGoQueryDataModel,
	}
)

// Init 初始化数据库
func Init(dbPath string) (err error) {
	dsn := dbPath
	uri, _ := url.Parse(dsn)
	pass, _ := uri.User.Password()
	port := "5432"
	if len(strings.Split(uri.Host, ":")) > 1 {
		port = strings.Split(uri.Host, ":")[1]
	}
	sslmode := uri.Query().Get("sslmode")
	if sslmode == "" {
		sslmode = "disable"
	}
	timezone := uri.Query().Get("TimeZone")
	timeout := uri.Query().Get("lock_timeout")
	v2path := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		strings.Split(uri.Host, ":")[0], uri.User.Username(), pass, uri.Path[1:], port, sslmode)
	if timezone != "" {
		v2path += " TimeZone=" + timezone
	}
	if timeout != "" {
		v2path += " lock_timeout=" + timeout
	}
	if err := InitDB(uri.Scheme, v2path, models...); err != nil {
		panic(err)
	}

	return
}

// InitDB 打开对应的数据库链接并创建相应数据模型名(单数形式)表结构
func InitDB(driver string, dbpath string, models ...interface{}) (err error) {
	var d gorm.Dialector
	driver = strings.ToLower(driver)
	switch driver {
	default:
		return errors.New("unsupport db driver")
	case "postgres":
		d = &postgres.Dialector{Config: &postgres.Config{DSN: dbpath}}
	case "mysql":
		d = &mysql.Dialector{Config: &mysql.Config{DSN: dbpath}}
	case "sqlite3":
		d = &sqlite.Dialector{DSN: dbpath}
		os.MkdirAll(filepath.Dir(dbpath), 0755)
	}
	db, err = gorm.Open(d, &gorm.Config{NamingStrategy: &schema.NamingStrategy{SingularTable: true}, DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		return err
	}

	if os.Getenv("DB_DEBUG") == "true" {
		db = db.Debug()
	}
	tx, _ := db.DB()
	if tx != nil {
		tx.SetMaxIdleConns(5)
		tx.SetMaxOpenConns(20)
	}
	if err != nil {
		return err
	}

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
	// 	logger.Config{
	// 		SlowThreshold:             time.Second, // 慢 SQL 阈值
	// 		LogLevel:                  logger.Info, // 日志级别
	// 		IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
	// 		Colorful:                  true,        // 彩色打印
	// 	},
	// )
	// db.Config.Logger = newLogger
	return db.AutoMigrate(models...)
}

// DB  ...
func DB() *gorm.DB {
	return db
}
