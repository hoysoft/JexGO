package JexGO
import (
	"html/template"
	"github.com/hoysoft/JexGO/utils"
)

var (
	AppFuncMaps    []template.FuncMap
	tplFuncMaps  template.FuncMap
)

func init() {

	tplFuncMaps = make(template.FuncMap)
	tplFuncMaps["title"]= GetAppTitle
	tplFuncMaps["getCnfValue"]=  GetCnfValue
	tplFuncMaps["substr"]= utils.Substr
	tplFuncMaps["byteUnitStr"]= utils.ByteUnitStr
	tplFuncMaps["byteUnitStr_uint"]= utils.ByteUnitStr_uint
	tplFuncMaps["Html2str"]=  Html2str
	AppFuncMaps=append(AppFuncMaps,tplFuncMaps)
}

func AddFuncMap(key string, funname interface{})   {
	tplFuncMaps[key]=funname
}

func GetAppTitle() string{
	return SpiderHttp.Cnf.AppName
}

