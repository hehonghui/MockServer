package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"mockserver/common"
)

// request & response mapping struct
type RespMapping struct {
	Path         string `json:"path"`      // 请求的path
	Method       string `json:"method"`    // 请求方式
	RespFilePath string `json:"resp_file"` // 存储返回值的文件路径
}


func (data *RespMapping) String() string  {
	return fmt.Sprintf("path : %s, method : %s, file : %s .", data.Path, data.Method, data.RespFilePath)
}

// mappings slice
var mappings []*RespMapping


// 1. 读取mock配置
// 2. 添加router 处理
// 3. 返回结果

func main() {
	mappingFilePath := "/Users/mrsimple/GoProjects/src/mockserver/config/mapping.json"
	if len(os.Args) == 2 {
		mappingFilePath = os.Args[1]
	}
	fmt.Println("====> Mapping file is ", mappingFilePath)

	err := readConfig(mappingFilePath)
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
	fmt.Println("=================================================> init router START ")
	router := http.NewServeMux()
	for _, item := range mappings {
		fmt.Println("route : " , item)
		router.HandleFunc(item.Path, processRequest)
	}
	fmt.Println("================================================> init router END ! ")
	fmt.Println()

	err := http.ListenAndServe(":10086", router)
	log.Fatal(err)
}

func processRequest(w http.ResponseWriter, r *http.Request)  {
	defer func() {
		// error 异常处理.
		if err := recover(); err != nil {
			resultMap := common.StringMap {"err_msg":err.(string), "path": r.URL.Path, "method": r.Method }
			result,_ := json.Marshal(resultMap)
			writeResponse(w, -1, result)
		}
	}()
	fmt.Println("==> Request path : ", r.URL.Path, " , method : ", r.Method)
	item, notFound := findRespMapping(mappings, r)
	if notFound != nil {
		panic("Not found RespMapping!")
	}
	result, err := ioutil.ReadFile(item.RespFilePath)
	if err == nil {
		writeResponse(w, 200, result)
	} else {
		panic("Not found response data for : " + r.URL.Path)
	}
}

func writeResponse(writer http.ResponseWriter, code int, body []byte)  {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(body)
}

func findRespMapping(mappings []*RespMapping, r *http.Request) (*RespMapping, error) {
	for _, item := range mappings {
		if item.Path == r.URL.Path && item.Method == r.Method {
			fmt.Println(">>> find resp file : " , item.RespFilePath, ", path : ", r.URL.Path, "method : ", r.Method)
			return item, nil
		}
	}
	return nil, errors.New("Not found RespMapping!")
}