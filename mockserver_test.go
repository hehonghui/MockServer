package main

import (
	"testing"
	"net/http"
	"net/url"
)

// =======================================> 单元测试


// 读取 mapping 配置文件测试用例
func TestReadConfigFromFile(t *testing.T) {
	if len(mappings) == 0 {
		t.Log("mapping is empty.")
	} else {
		t.Error("mapping is not empty !")
	}
	// 读取参数
	readConfig("/Users/mrsimple/GoProjects/src/mockserver/config/mapping.json")
	// 再次验证
	if len(mappings) == 0 {
		t.Error("mapping is still empty after readConfig !")
	}
}

// 为请求找到对应的返回结果测试用例
func TestFindMappingForRequest(t *testing.T) {
	mockMappings := make([]*RespMapping, 0)

	map1 := RespMapping {
		"/v4/users/", "GET",
		"/Users/mrsimple/GoProjects/src/mockserver/config/responses/users_resp.json",
	}
	map2 := RespMapping {
		"/v9/categories/sorted/", "GET",
		"/Users/mrsimple/GoProjects/src/mockserver/config/responses/channels.json",
	}

	mockMappings = append(mockMappings, &map1, &map2)
	t.Logf("mock mapping data : %v ", mockMappings)

	req := new(http.Request)
	req.URL = &url.URL{
		Path:"/v4/users/",
	}
	req.Method = "GET"

	// /v4/users/
	item, err := findRespMapping(mockMappings, req)
	if err != nil {
		t.Error("Not found response. ", item)
	} else {
		t.Log("found response ", item.String())
	}

	// sorted request
	sortedReq := new(http.Request)
	sortedReq.URL = &url.URL{
		Path:"/v9/categories/sorted/",
	}
	sortedReq.Method = "GET"
	sItem, sErr := findRespMapping(mockMappings, req)
	if sErr != nil {
		t.Error("Not found response. ", sItem)
	} else {
		t.Log("found response ", sItem.String())
	}

	// 找不到匹配的 request
	notFoundReq := new(http.Request)
	notFoundReq.URL = &url.URL{
		Path:"/v1/notfound/",
	}
	notFoundReq.Method = "GET"
	notItem, notErr := findRespMapping(mockMappings, req)
	if notErr == nil {
		t.Log("I have not config for ", notItem.Path)
	} else {
		t.Error(notItem.String() + " should not be found !")
	}
}

