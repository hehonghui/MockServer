package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
)


var mappings []*RespMapping


type RespMapping struct {
	Path         string `json:"path"`      // 请求的path
	Method       string `json:"method"`    // 请求方式
	RespFilePath string `json:"resp_file"` // 存储返回值的文件路径
}

type StringMap map[string]interface{}


func (data *RespMapping) String() string  {
	return fmt.Sprintf("path : %s, method : %s, file : %s .", data.Path, data.Method, data.RespFilePath)
}


/**
 *	1. 读取mock配置
 *	2. 添加router 处理
 *	3. 返回结果
 */
func main() {
	fmt.Println("init mock server ...")
	err := readConfig("/Users/mrsimple/GoProjects/src/mockserver/config/mapping.json")
	if err == nil {
		createRouter()
	} else {
		panic("init server failed !!!")
	}
}

func readConfig(configFile string) error  {
	configData, err := ioutil.ReadFile(configFile)
	//fmt.Println("read data from ", configFile, ", data is : \n ", string(configData))
	if  err == nil {
		return json.Unmarshal(configData, &mappings)
	} else {
		fmt.Println("read resp mapping file error !!!")
		return errors.New("read resp mapping file error !!!")
	}
}


func createRouter() {
	router := gin.Default()
	for _, item := range mappings {
		fmt.Println("route path , " , item.Path, ", resp : ", item.RespFilePath)
		router.Handle(item.Method, item.Path, processRequest)
	}
	router.Run(":10086")
}


// 根据请求的path, method 找到对应的映射, 然后解析为json, 如果解析为json失败则尝试解析为json array.
func processRequest(c *gin.Context) {
	defer func() {
		// error 异常处理.
		if err := recover(); err != nil {
			c.JSON(-1, StringMap{"-1": -1, "err_msg": err.(string)})
		}
	}()
	item, notFound := findRespMapping(c)
	if notFound != nil {
		panic("Not found RespMapping!")
	}
	result, err := ioutil.ReadFile(item.RespFilePath)
	//fmt.Println("resp file , " , item.RespFilePath, ", path : ", c.Request.URL.Path)
	if err == nil {
		// string to map, and then map to json
		var resp StringMap
		jsonErr := json.Unmarshal(result, &resp)
		if jsonErr == nil {
			c.JSON(200, resp )
		} else {
			fmt.Println("json parse error : ", jsonErr)
			// json array 的返回结果
			var arrayResult []StringMap
			jsonArrayErr := json.Unmarshal(result, &arrayResult)
			if jsonArrayErr == nil {
				c.JSON(200, arrayResult )
			} else {
				panic("Unmarshal data to json error")
			}
		}
	}
}

func findRespMapping(c *gin.Context) (*RespMapping, error) {
	for _, item := range mappings {
		//fmt.Println("find resp file , " , item.RespFilePath, ", path : ", c.Request.URL.Path)
		if item.Path == c.Request.URL.Path && item.Method == c.Request.Method {
			return item, nil
		}
	}
	return nil, errors.New("Not found RespMapping!")
}