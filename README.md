# cmdb-utils
cmdb utils
#### 初始化http客户端
```
client := cmdbutils.InitClient("xxx.cmdb.com")
```
#### 关闭http客户端
```
cmdbutils.Close(client)
```

#### post
```
cmdbutils.CmdbPost(uri string, data map[string]interface{}, ak string, sk string, client *resty.Client, statusCode *int, body *[]byte) error
```
* uri: 调用路径
* data:post 请求体
* ak:cmdb ak
* sk:cmdb sk
* client: http客户端
* statusCode: http状态码
* body: http response body字节码

#### delete
```
cmdbutils.CmdbDelete(uri string, data map[string]string, ak string, sk string, client *resty.Client, statusCode *int, body *[]byte) error
```
* uri: 调用路径
* data:delete 参数
* ak:cmdb ak
* sk:cmdb sk
* client: http客户端
* statusCode: http状态码
* body: http response body字节码