package utils

import (
"log"
"strings"
	"path/filepath"
	"os"
	"net/http"
	"io/ioutil"
	"regexp"
	"fmt"


)

func CheckErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//判断是否路径
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("not reached")
}

func HttpGetString(url string)([]byte, error){
	res, err := http.Get(url)
	if err != nil {
		return nil,err
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil,err
	}
	return   result,nil
}

//正则表达式map输出
func FindStringSubmatchMap(r *regexp.Regexp,s string) map[string]string{
	captures:=make(map[string]string)

	match:=r.FindStringSubmatch(s)
	if match==nil{
		return captures
	}

	for i,name:=range r.SubexpNames(){
		//Ignore the whole regexp match and unnamed groups
		if i==0||name==""{
			continue
		}

		captures[name]=match[i]
	}
	return captures
}

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}
func ByteUnitStr_uint(n uint64) string {
	return ByteUnitStr(int64(n))
}
func ByteUnitStr(n int64) string {
	var unit string
    size := float64(n)
	for i := 1; i < len(byteUnits); i++ {
		if size < 1000 {
			unit = byteUnits[i-1]
			break
		}

		size = size / 1000
	}

	return fmt.Sprintf("%.3g %s", size, unit)
}

