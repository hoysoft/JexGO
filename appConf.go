package JexGO
import (

)

/**
  读写项目配置文件
  */
type AppConf struct {
	AppName string
	Env string
	Version string
	Port string
    DB
	LogDir string
	config *CnfConfig
}

func NewAppConf(cnfFileName string) *AppConf {
	c:=new(AppConf)

	c.config=SetConfig("app.cnf")
	c.AppName=c.config.GetValue("Global","appName")
    c.Port=c.config.GetValue("Global","port","8888")
	c.Version=c.config.GetValue("Global","version","0.0.1")

	c.Env=c.config.GetValue("Global","Env","development")
	c.DriverName=c.config.GetValue(c.Env,"adapter")
	c.DataSourceName=c.config.GetValue(c.Env,"database")
    c.UserName=c.config.GetValue(c.Env,"username")
	c.Password=c.config.GetValue(c.Env,"password")
	c.Host=c.config.GetValue(c.Env,"host")
	c.Encoding=c.config.GetValue(c.Env,"encoding")
	return  c
}