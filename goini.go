package JexGO

import (
	"bufio"
	"io"
	"os"
	"strings"
	"github.com/hoysoft/JexGO/utils"
)

type CnfConfig struct {
	filepath string                         //your ini file path directory+file
	conflist []map[string]map[string]string //configuration information slice
}

//Create an empty configuration file
func SetConfig(filepath string) *CnfConfig {
	c := new(CnfConfig)
	c.filepath = filepath

	return c
}

//To obtain corresponding value of the key values
func (c *CnfConfig) GetValue(section, name string, defval ...interface{}) string {
	c.ReadList()
	conf := c.ReadList()
	for _, v := range conf {
		for key, value := range v {
			if key == section {
				if len(value[name])==0 && len(defval)>0{
					return defval[0].(string)
                }else{
					return value[name]
				}
			}
		}
	}
	if len(defval)>0{
		return defval[0].(string)
	}else{
		return ""
	}
}

//Set the corresponding value of the key value, if not add, if there is a key change
func (c *CnfConfig) SetValue(section, key, value string) bool {
	c.ReadList()
	data := c.conflist
	var ok bool
	var index = make(map[int]bool)
	var conf = make(map[string]map[string]string)
	for i, v := range data {
		_, ok = v[section]
		index[i] = ok
	}

	i, ok := func(m map[int]bool) (i int, v bool) {
		for i, v := range m {
			if v == true {
				return i, true
			}
		}
		return 0, false
	}(index)

	if ok {
		c.conflist[i][section][key] = value
		return true
	} else {
		conf[section] = make(map[string]string)
		conf[section][key] = value
		c.conflist = append(c.conflist, conf)
		return true
	}

	return false
}

//Delete the corresponding key values
func (c *CnfConfig) DeleteValue(section, name string) bool {
	c.ReadList()
	data := c.conflist
	for i, v := range data {
		for key, _ := range v {
			if key == section {
				delete(c.conflist[i][key], name)
				return true
			}
		}
	}
	return false
}

//List all the configuration file
func (c *CnfConfig) ReadList() []map[string]map[string]string {

	file, err := os.Open(c.filepath)
	if err != nil {
		utils.CheckErr(err,"open inifile failure.")
	}
	defer file.Close()
	var data map[string]map[string]string
	var section string
	buf := bufio.NewReader(file)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				utils.CheckErr(err,"read inifile failure.")
			}
			if len(line) == 0 {
				break
			}
		}
		switch {
		case len(line) == 0:
		case strings.HasPrefix(line, "#"):
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
			data = make(map[string]map[string]string)
			data[section] = make(map[string]string)
		default:
			i := strings.IndexAny(line, "=")
			value := strings.TrimSpace(line[i+1 : len(line)])
			data[section][strings.TrimSpace(line[0:i])] = value
			if c.uniquappend(section) == true {
				c.conflist = append(c.conflist, data)
			}
		}

	}

	return c.conflist
}



//Ban repeated appended to the slice method
func (c *CnfConfig) uniquappend(conf string) bool {
	for _, v := range c.conflist {
		for k, _ := range v {
			if k == conf {
				return false
			}
		}
	}
	return true
}