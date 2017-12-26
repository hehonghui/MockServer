package main

import (
	"mockserver/mock"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"encoding/json"
)

/**
	1. 读取mock配置
	2. 添加router 处理
	3. 返回结果
 */
func main() {
	fmt.Println("init mock server ...")
	group := mock.ReadConfig("/Users/mrsimple/GoProjects/src/mockserver/config/mapping.json")
	if group != nil {
		createRouter(group)
	}
}

func createRouter(group *mock.MappingGroup) {
	router := gin.Default()
	for _, item := range group.GETS {
		fmt.Println("path , " , item.Path, ", resp : ", item.RespFile)
		router.Handle(item.Method, item.Path, func(c *gin.Context) {
			defer func() {
				// error 异常处理.
				if err := recover(); err != nil {
					c.JSON(1001, map[string]interface{}{"err_code": 1001, "err_msg": err.(string)})
				}
			}()
			result, err := ioutil.ReadFile(item.RespFile)
			fmt.Println("resp file , " , item.RespFile, ", path : ", c.Request.URL.Path)
			if err == nil {
				// string to map, and then map to json
				var resp map[string]interface{}
				json.Unmarshal(result, &resp)
				c.JSON(200, resp )
			}
		})

		// beta
		//router.GET(item.Path, func(c *gin.Context) {
		//	defer func() {
		//		// error 异常处理.
		//		if err := recover(); err != nil {
		//			c.JSON(1001, map[string]interface{}{"err_code": 1001, "err_msg": err.(string)})
		//		}
		//	}()
		//	result, err := ioutil.ReadFile(item.RespFile)
		//	fmt.Println("resp file , " , item.RespFile, ", path : ", c.Request.URL.Path)
		//	if err == nil {
		//		// string to map, and then map to json
		//		var resp map[string]interface{}
		//		json.Unmarshal(result, &resp)
		//		c.JSON(200, resp )
		//	}
		//})
	}
	router.Run(":10086")
}