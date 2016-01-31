package plugins
import (
"encoding/json"

	"html/template"

)



type JQueryDataTable struct {
	Ajax *jqdtb_ajax `json:"ajax"`
	Columns []*JQDBTable_Column `json:"columns"`
	Language   *jqdtb_language  `json:"language,omitempty"`
}

type jqdtb_ajax struct {
   Url string  `json:"url"`
   DataSrc string `json:"dataSrc"`
}

type JQDBTable_Column struct {
	DataField string `json:"data,omitempty"`
	Tile string `json:"title"`
	Sortable bool `json:"bSortable,omitempty"`
	SortDataType bool `json:"sSortDataType,omitempty"` //定义数据源类型的排序(主要是input类型)，dom-text,dom-select,dom-checkbox
	Class string `json:"sClass,omitempty"`
	Width string `json:"width,omitempty"`
	Searchable bool `json:"searchable,omitempty"`
	Type string  `json:"sType,omitempty"` //属性sType：定义这列数据类型，包括（string,numeric,date,html)。
	Render string `json:"render,omitempty"`
	Visible bool `json:"visible,omitempty"`
}
type jqdtb_language struct {
	Url string  `json:"url"`
}

func NewJQueryDataTable()*JQueryDataTable{
	 j:=&JQueryDataTable{
		 Ajax:&jqdtb_ajax{},
		 Language:&jqdtb_language{},
	 }
	 return j
}

func (j *JQueryDataTable)AddColumns(col ...*JQDBTable_Column)*JQueryDataTable{
    j.Columns=append(j.Columns,col...)
	return j
}

func (j *JQueryDataTable)SetAjaxSource(url string)*JQueryDataTable{
	j.Ajax.Url=url

	return j
}

func (j *JQueryDataTable)SetLanguage(url string)*JQueryDataTable{
	j.Language.Url=url
	return j
}

func  (j *JQueryDataTable)JS() template.JS{
	bytes,_:= json.Marshal(j)
//	fmt.Println("string:",string(bytes))
//	fmt.Println("html:", template.HTML( string(bytes)))
//	fmt.Println("HTMLEscapeString:", template.HTMLEscapeString( string(bytes)))

	return   template.JS(bytes)
}

var JQueryDataTable_LangeChina=`
{
"sProcessing": "处理中...",
"sLengthMenu": "显示 _MENU_ 项结果",
"sZeroRecords": "没有匹配结果",
"sInfo": "显示第 _START_ 至 _END_ 项结果，共 _TOTAL_ 项",
"sInfoEmpty": "显示第 0 至 0 项结果，共 0 项",
"sInfoFiltered": "(由 _MAX_ 项结果过滤)",
"sInfoPostFix": "",
"sSearch": "搜索:",
"sUrl": "",
"sEmptyTable": "表中数据为空",
"sLoadingRecords": "载入中...",
"sInfoThousands": ",",
"oPaginate": {
"sFirst": "首页",
"sPrevious": "上页",
"sNext": "下页",
"sLast": "末页"
},
"oAria": {
"sSortAscending": ": 以升序排列此列",
"sSortDescending": ": 以降序排列此列"
}
}`


