package JexGO

/**
数据库管理接口
 */

import (
    "github.com/jinzhu/gorm"
    _ "github.com/lib/pq"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/mattn/go-sqlite3"


	//"database/sql"
	"fmt"
"log"

	"github.com/hoysoft/JexGO/utils"
)

var db_Tables  []interface{};

//初始化数据接口
type IDB interface {
	InitData()
}



type DB struct {
	//"sqlite3"
	DriverName string
	//"/tmp/post_db.bin"
	DataSourceName string
	UserName string
	Password string
	Host string
	Encoding string
	Db gorm.DB
}


//存在数据时update，否则插入
//func (this *DB) MergeInto(value  interface{}) error {
//
//		s := reflect.ValueOf(value).Elem()
//
//	fmt.Println("s.Type():",s.Type())
//		tb,err:=  this.DbMap.TableFor(s.Type(),true)
//	CheckErr(err,"tttt")
//		unique_fields:=""
//	fmt.Println("tb.Columns:",tb)
//		for _,field :=range tb.Columns {
//			if field.Unique {
//				if len(unique_fields)>0{
//					unique_fields=unique_fields+" and "
//				}
//				unique_fields =unique_fields + field.ColumnName+" = "+s.FieldByName(field.ColumnName).Interface().(string)
//			}
//		}
//	fmt.Println("select count(*) from "+tb.TableName+" where "+unique_fields)
// 		count,_:=this.DbMap.SelectInt("select count(*) from "+tb.TableName+" where "+unique_fields)
//
//		if count>0{
//			//update
//			 _,err:=this.DbMap.Update(value)
//			return err
//		}else{
//			//insert
//			return this.DbMap.Insert(value)
//		}
//}

func init(){

}

func newDB() *DB{
	db:=new(DB)
  return db
}

func AddTable(table interface{}){
	db_Tables=append(db_Tables,table)
}

func (this *DB)initDb() {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish

	var err error
	switch this.DriverName {
	case "sqlite3":
		this.Db, err = gorm.Open("sqlite3", this.DataSourceName)
		utils.CheckErr(err, "Got error when connect database")
		//this.DbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	case "mysql":
		//db, err := sql.Open("mysql", "user:password@tcp(localhost:5555)/dbname?charset=utf8&parseTime=True&loc=Local")
		this.Db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			this.UserName,
			this.Password,
		    this.Host,
			this.DataSourceName,
			this.Encoding))
		utils.CheckErr(err, "Got error when connect database")
		//this.DbMap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	default:
		log.Fatalln("The types of database does not support:"+this.DriverName , err)
	}

	this.Db.DB()
	this.Db.DB().Ping()
	this.Db.DB().SetMaxIdleConns(10)
	this.Db.DB().SetMaxOpenConns(100)

	// Disable table name's pluralization
	this.Db.SingularTable(true)

	// construct a gorp DbMap

	//


		if this.DriverName=="mysql" {
			this.Db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(db_Tables...)
		}else{
			this.Db.AutoMigrate(db_Tables...)
		}

	//初始化数据表数据
	for _,value:=range  db_Tables {
		if idb,ok:=value.(IDB);ok{
			idb.InitData()
		}

	}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	//dbmap.AddTableWithName(Post{}, "posts").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	//err = this.DbMap.CreateTablesIfNotExists()
	// CheckErr(err, "Create tables failed")
}



