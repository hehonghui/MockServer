package main

import (
	"mockserver/mock"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
)


var mappingGroup *mock.MappingGroup

/**
	1. 读取mock配置
	2. 添加router 处理
	3. 返回结果
 */
func main() {
	fmt.Println("init mock server ...")
	mappingGroup, err := mock.ReadConfig("/Users/mrsimple/GoProjects/src/mockserver/config/mapping.json")
	if err == nil {
		createRouter(mappingGroup)
	}
}

func createRouter(group *mock.MappingGroup) {
	router := gin.Default()
	for _, item := range group.Maps {
		fmt.Println("path , " , item.Path, ", resp : ", item.RespFile)
		router.Handle(item.Method, item.Path, processRequest)
	}
	router.Run(":10086")
}


func processRequest(c *gin.Context) {
	defer func() {
		// error 异常处理.
		if err := recover(); err != nil {
			c.JSON(-1, map[string]interface{}{"-1": -1, "err_msg": err.(string)})
		}
	}()
	item, notFound := findRespMapping(c)
	fmt.Println(item, " not found , ", notFound)

	if notFound != nil {
		panic("Not found RespMapping!")
	}
	result, err := ioutil.ReadFile(item.RespFile)
	//fmt.Println("resp file , " , item.RespFile, ", path : ", c.Request.URL.Path)
	if err == nil {
		// string to map, and then map to json
		var resp map[string]interface{}
		json.Unmarshal(result, &resp)
		c.JSON(200, resp )
	}
}

func findRespMapping(c *gin.Context) (*mock.RespMapping, error) {
	fmt.Println("findRespMapping , " , mappingGroup, ", maps ", mappingGroup.Maps)

	for _, item := range mappingGroup.Maps {
		fmt.Println("find resp file , " , item.RespFile, ", path : ", c.Request.URL.Path)
		if item.Path == c.Request.URL.Path && item.Method == c.Request.Method {
			return item, nil
		}
	}
	return nil, errors.New("Not found RespMapping!")
}