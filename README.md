# A simple mock server with Go

用Go语言实现的一个简单的mock server (只支持完整的path匹配,不支持参数化), 在配置文件 mapping.json 文件中配置每个请求的path、对应的请求方式、保存Response数据的文件路径. 

示例如下:

```json
[
  {
    "path": "/v4/users/",
    "method": "GET",
    "resp_file": "/Users/yourname/config/responses/users_resp.json"
  },
  {
    "path": "/v5/categories/rec/articles/next/",
    "method": "POST",
    "resp_file": "/Users/yourname/config/responses/news_list_data.json"
  }
]
```

## 启动mock server

进入到项目所在的目录,执行启动命令

```
go run mockserver.go ~/mapping.json
```

其中最后一个参数为 mapping.json 所在的路径, 修改为你的mapping路径即可.


## 验证mock server 是否正常启动

终端中输入如下Log 即表示成功解析mapping文件,并启动mock server.


```
Mapping file is  /Users/yourname/mapping.json
=================================================> init router START
route :  path : /v4/users/, method : GET, file : /Users/yourname/config/responses/users_resp.json .
route :  path : /v5/categories/rec/articles/next/, method : GET, file : /Users/yourname/config/responses/news_list_data.json .
route :  path : /v9/categories/sorted/, method : GET, file : /Users/yourname/config/responses/channels.json .
================================================> init router END !
```