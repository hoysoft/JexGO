package JexGO

import (
	"github.com/go-martini/martini"
	"flag"
	"os"
	"runtime"
	"fmt"

	"github.com/hoysoft/JexGO/logger"

	"github.com/hoysoft/JexGO/utils"
)


var SpiderHttp *Spider=newSpider();

type Spider struct {
	Cnf     *AppConf
	Db      *DB
	Martini *martini.ClassicMartini
	controllers map[string]controllerInfo
	Tasks *TaskManager  //后台任务管理器
}


type controllerInfo struct {
	ic IController
	handlers []martini.Handler
}


//func  SpiderHttp()*Spider{
//	if !SpiderHttp==nil {
//		SpiderHttp=newSpider()
//	}
//	return SpiderHttp;
//}

func newSpider(cnfFileName ...string) *Spider {
	s := new(Spider)
	s.controllers=make(map[string]controllerInfo)

	s.Db = newDB()
	if len(cnfFileName)==0{
		cnfFileName=append(cnfFileName,"app.cnf")
	}
	s.Cnf = NewAppConf(cnfFileName[0])

	s.checkFlag();
	s.Martini = martini.Classic()
	s.Tasks=NewTasks()
	return s
}



/**
 声明控制器
 */
func (this *Spider)RegisterController(urlpath string,ic IController,h ...martini.Handler){
	ic.SetPath(urlpath)
   this.controllers[urlpath]=controllerInfo{ic:ic,handlers:[]martini.Handler(h)}
}

func (this *Spider)checkFlag(){
	var flagHelp bool
	flag.BoolVar(&flagHelp, "h", false, "view this help")
	flag.StringVar(&this.Cnf.Port, "p", this.Cnf.Port, "http listen port")
	flag.StringVar(&this.Cnf.Env, "e", this.Cnf.Env, "select environment <development|production|test>")

	flag.Parse()

	if flagHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func (this *Spider)GetTables()[]interface{}{
	return db_Tables
}



func (this *Spider)Run() {

	//初始化
	runtime.GOMAXPROCS(runtime.NumCPU())
	martini.Env=this.Cnf.Env
	this.Db.DataSourceName=this.Cnf.DB.DataSourceName
	this.Db.DriverName=this.Cnf.DB.DriverName
	if martini.Env=="development" {
		this.Db.Db.Debug()
		logger.SetConsole(true)
		logger.SetLevel(logger.DEBUG)
	}else{
		logger.SetConsole(false)
		logger.SetLevel(logger.WARN)
	}

	this.Cnf.LogDir=utils.GetWokingDirectory("logs")
	//logger.SetRollingDaily(this.Cnf.LogDir, martini.Env+".log")

	logger.SetRollingFile(this.Cnf.LogDir, martini.Env+".log", 15, 5, logger.MB)
	this.Db.initDb()
	defer this.Db.Db.Close()


	//控制器路由设置
	for key,value:=range this.controllers{
		this.Martini.Group(key,value.ic.Router,value.handlers... );
	}

	//开始服务
	fmt.Println("Version:", this.Cnf.Version)
	os.Setenv("PORT",this.Cnf.Port)

	logger.Info("start webservice")
	this.Martini.Run()
}

