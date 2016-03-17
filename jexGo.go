package JexGO

import (
	"github.com/go-martini/martini"
	"flag"
	"os"
	"runtime"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/hoysoft/JexGO/logger"
	"github.com/hoysoft/JexGO/utils"
	"github.com/kardianos/service"
)

var (
	JexHttp *JexGo = newJexHttp()
	DB *gorm.DB
	Tasks *TaskManager  //后台任务管理器
	db_Tables  []interface{}

	//flag
	fHelp = flag.Bool("h", false, "view this help")
    //service
//	fInstall = flag.Bool("install", false, "Path to file")
//	fUninstall = flag.Bool("uninstall", false, "Path to file")
//	fStart = flag.Bool("start", false, "Path to file")
//	fStop = flag.Bool("stop", false, "Path to file")
//	fRestart = flag.Bool("restart", false, "Path to file")
)


type Config struct {
	AppName  string
	Version  string
	Port     string
	Env      string
	LogDir   string
	DBConfig DbConfig
}




type JexGo struct {
	Cnf         *Config
	Martini     *martini.ClassicMartini
	controllers map[string]controllerInfo
}


type controllerInfo struct {
	ic       IController
	handlers []martini.Handler
}


func (c *Config)loadfromFile(cnfile string) {
	config := SetConfig(cnfile)
	//全局配置
	c.AppName = config.GetValue("Global", "appName")
	c.Port = config.GetValue("Global", "port", "8888")
	c.Version = config.GetValue("Global", "version", "0.0.1")
	c.LogDir = config.GetValue("Global", "logdir", utils.GetWokingDirectory("logs"))
	c.Env = config.GetValue("Global", "Env", "development")
	//dbconfig
	c.DBConfig.DriverName = config.GetValue(c.Env, "adapter")
	c.DBConfig.DataSourceName = config.GetValue(c.Env, "database")
	c.DBConfig.UserName = config.GetValue(c.Env, "username")
	c.DBConfig.Password = config.GetValue(c.Env, "password")
	c.DBConfig.Host = config.GetValue(c.Env, "host")
	c.DBConfig.Encoding = config.GetValue(c.Env, "encoding")
}


//func  SpiderHttp()*Spider{
//	if !SpiderHttp==nil {
//		SpiderHttp=newSpider()
//	}
//	return SpiderHttp;
//}

func newJexHttp(cnfFileName ...string) *JexGo {

	//	if SpiderHttp==nil{
	n := &JexGo{}
	n.Cnf = &Config{}
	n.Cnf.loadfromFile("app.cnf")
	n.controllers = make(map[string]controllerInfo)
	n.Martini = martini.Classic()
	Tasks = NewTasks()
	n.checkFlag();
	return n
	//	}
	//   return SpiderHttp
}

func (this *JexGo)SetConfig(cnf *Config) {

}

/**
 声明控制器
 */
func (this *JexGo)RegisterController(urlpath string, ic IController, h ...martini.Handler) {
	ic.SetPath(urlpath)
	this.controllers[urlpath] = controllerInfo{ic:ic, handlers:[]martini.Handler(h)}
}

func (this *JexGo)checkFlag() {
	flag.StringVar(&this.Cnf.Port, "p", this.Cnf.Port, "http listen port")
	flag.StringVar(&this.Cnf.Env, "e", this.Cnf.Env, "select environment <development|production|test>")

	flag.Parse()


//	p := &program{}
//	s, err := service.New(p, &service.Config{
//		Name : "redhare",
//		DisplayName : "redhare",
//		Description : "redhare web service",
//	})
//	if err != nil {
//		log.Fatal(err)
//	}

	switch  {
	case *fHelp:
		flag.Usage()
		os.Exit(0)
//	case *fInstall:
//		err = s.Install()
//		if err != nil {
//			log.Fatal("install", err)
//		}
//		os.Exit(0)
//	case *fUninstall:
//		err = s.Uninstall()
//		if err != nil {
//			log.Fatal("Uninstall", err)
//		}
//		os.Exit(0)
//	case *fStart:
//		err = s.Start()
//		if err != nil {
//			log.Fatal("Start", err)
//		}
//	case *fStop:
//		err = s.Stop()
//		if err != nil {
//			log.Fatal("Stop", err)
//		}
//	case *fRestart:
//		err = s.Restart()
//		if err != nil {
//			log.Fatal("Restart", err)
//		}
//		os.Exit(0)
//	default:
//		err = s.Run()
//		if err != nil {
//			log.Fatal("Run", err)
//		}
	}
}

func (this *JexGo)GetTables() []interface{} {
	return db_Tables
}



func (this *JexGo)Run() {
	//初始化
	runtime.GOMAXPROCS(runtime.NumCPU())
	d := NewDb(this.Cnf.DBConfig)
	DB = d.db
	d.AutoMigrate(db_Tables...)
	defer DB.Close()

	martini.Env = this.Cnf.Env
	if martini.Env == "development" {
		DB.Debug()
		// DB.LogMode(true)
		logger.SetConsole(true)
		logger.SetLevel(logger.DEBUG)
	}else {
		logger.SetConsole(false)
		logger.SetLevel(logger.WARN)
	}
	//logger.SetRollingDaily(this.Cnf.LogDir, martini.Env+".log")

	logger.SetRollingFile(this.Cnf.LogDir, martini.Env + ".log", 15, 5, logger.MB)



	//控制器路由设置
	for key, value := range this.controllers {
		this.Martini.Group(key, value.ic.Router, value.handlers...);
	}

	//开始服务
	fmt.Println("Version:", this.Cnf.Version)
	os.Setenv("PORT", this.Cnf.Port)

	logger.Info("start webservice")
	this.Martini.Run()
}



/**
数据库部分
 */

func AddTable(table interface{}) {
	db_Tables = append(db_Tables, table)
}


type program struct{mainfun func()}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
	p.mainfun()
}
func (p *program) Stop(s service.Service) error {
	return nil
}