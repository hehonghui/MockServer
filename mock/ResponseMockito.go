package mock

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type RespMapping struct {
	Path string	`json:"path"`			// 请求的path
	Method string	`json:"method"`		// 请求方式
	RespFile string	`json:"resp_file"`		// 存储返回值的文件路径
}

type MappingGroup struct {
	GETS []RespMapping
	POSTS []RespMapping
}

func (data *RespMapping) String() string  {
	return fmt.Sprintf("path : %s, method : %s, file : %s .", data.Path, data.Method, data.RespFile)
}

func ReadConfig(configFile string) *MappingGroup  {
	var slice = make([]RespMapping, 0)
	configData, err := ioutil.ReadFile(configFile)
	//fmt.Println("read data from ", configFile, ", data is : \n ", string(configData))
	if  err == nil {
		json.Unmarshal(configData, &slice)
		return groupRequest(slice)
	} else {
		fmt.Println("read resp mapping file error !!!")
		return nil
	}
}

func groupRequest(mappings []RespMapping) *MappingGroup {
	group := new(MappingGroup)
	for _, item := range mappings  {
		fmt.Println("resp mock item :", item.String() )

		if item.Method == "GET" {
			group.GETS = append(group.GETS, item)
		}
		if item.Method == "POST" {
			group.POSTS = append(group.POSTS, item)
		}
	}
	return group
}



