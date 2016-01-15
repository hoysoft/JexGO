package JexGO

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/hoysoft/JexGO/logger"
	_ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)
type Db struct {
    Cnf  DbConfig
	db gorm.DB
}

type DbConfig struct {
	//"sqlite3"
	DriverName string
	//"/tmp/post_db.bin"
	DataSourceName string
	UserName string
	Password string
	Host string
	Encoding string
}

//初始化数据接口
type IDB interface {
	InitData()
}


func NewDb(Cnf  DbConfig) *Db {
	var err error
	var db gorm.DB
	switch Cnf.DriverName {
	case "sqlite3":
		db, err = gorm.Open("sqlite3", Cnf.DataSourceName)
		logger.CheckFatal(err, "Got error when connect database")
	//this.DbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	case "mysql":
		//db, err := sql.Open("mysql", "user:password@tcp(localhost:5555)/dbname?charset=utf8&parseTime=True&loc=Local")
		db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			Cnf.UserName,
			Cnf.Password,
			Cnf.Host,
			Cnf.DataSourceName,
			Cnf.Encoding))
		logger.CheckFatal(err,"Got error when connect database")

	//this.DbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	default:
		logger.Fatal("The types of database does not support:"+Cnf.DriverName , err)
	}

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// Disable table name's pluralization
	db.SingularTable(true)

	// construct a gorp DbMap

	//

	return  &Db{Cnf:Cnf,db:db}
}




func (this *Db)AutoMigrate(tables ...interface{}) *Db{
	if this.Cnf.DriverName=="mysql" {
		this.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(tables...)
	}else{
		this.db.AutoMigrate(tables...)
	}

	//初始化数据表数据
	for _,value:=range  tables {
		if idb,ok:=value.(IDB);ok{
			idb.InitData()
		}
	}
	return this
}