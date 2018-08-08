# go-file-store

> beego 实现restful风格文件存储系统

# 部署方式

1、下载源码

```shell
https://github.com/jiangyx3915/go-file-store.git
```

2、配置数据库信息和文件上传路径

修改conf -> app.conf

```text
dbHost =    连接地址
dbPort =    端口
dbUser =    连接用户
dbPassword = 密码
dbDatabase = 库名

uploadFilesDir = 文件上传路径
```

3、运行启动

在项目根目录运行

```shell
go run main.go
```

# 使用指南

1、创建client(用户)

使用POST方法

```text
http://127.0.0.1:8080/v1/client
```

将会返回Client_Id 和 client_secret 

2、创建token

使用POST方法

```text
http://127.0.0.1:8080/v1/token

所带参数
client_id
client_secret
```

返回  token

3、上传文件

POST方法

```text
http://127.0.0.1:8080/v1/file/upload

token       token
file        上传的文件

```

上传成功会返回slug

4、下载文件

GET方法

```text
http://127.0.0.1:8080/v1/file/download

token       token
slug        文件上传后返回的slug
```