package JexGO

/**
  控制器基类
 */
import (

	"github.com/go-martini/martini"
	"strings"


	"net/http"
	"strconv"
	"github.com/jinzhu/gorm"
)

type IController interface {
	Init(r martini.Router)
	Router(martini.Router)
	SetPath(path string)
    Urls(parms ...string) string
	Paths(parms ...string) string
	//Prepare(...interface{})
}

type Controller struct {
	templatePath string
	DB  gorm.DB
	Data map[string]interface{}
}

func (this *Controller) Init(r martini.Router){
	this.DB=DB
	this.Data=make(map[string]interface{})
    this.Data["title"]=JexHttp.Cnf.AppName
}

func (this *Controller)SetPath(path string){
	this.templatePath=path
}

func (this *Controller)Urls(parms ...string)string{
	 url :=this.templatePath

	for _, parm := range parms {
		if len(url)>0{
			url += "/"+parm
		}else {
			url = parm
		}
	}
	url=strings.Replace(url,"//","/",1)
//	fmt.Println(url)
	return url
}

func (this *Controller)Paths(parms ...string)string{
	url :=this.templatePath

	for _, parm := range parms {
		if len(url)>0{
			url += "/"+parm
		}else {
			url = parm
		}
	}
	//fmt.Println(url)
    url=strings.TrimLeft(url,"/")

	//fmt.Println(url)
	return url
}


type TablePaginatorConfig struct {
	Fields []string
    Sortby []string
	Order []string
}

//分页读取数据
func (this *Controller)GetTablePaginatorData(tableDB *gorm.DB,req *http.Request,limit int64)   {

	p := req.URL.Query().Get("p")
	pageNo, _ := strconv.Atoi(p)
	if pageNo == 0 {
		pageNo = 1
	}
	//var limit int64 = 10 //每页10行显示
	var offset int64 = (int64(pageNo) - 1) * limit //起始位置

    tbDB:=tableDB.Limit(limit).Offset(offset)
	var count int
    tbDB.Count(count)


	this.Data["paginator"] = NewPaginator(req, int(limit), count)

}

/**
m.Get("/", func() {
  // 显示
})

m.Patch("/", func() {
  // 更新
})

m.Post("/", func() {
  // 创建
})

m.Put("/", func() {
  // 替换
})

m.Delete("/", func() {
  // 删除
})

m.Options("/", func() {
  // http 选项
})

m.NotFound(func() {
  // 处理 404
})
 */