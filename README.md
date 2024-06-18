# fastApi

fastApi: 基于Gin框架构建的web服务开发手脚架，统一规范开发，快速开发Api接口

https://github.com/kwinH/fastApi

## 目的

本项目采用了一系列Golang中比较流行的组件，可以以本项目为基础快速搭建Restful Web API

## 特色

本项目已经整合了许多开发API所必要的组件：

1. [cobra](github.com/spf13/cobra): cobra既是一个用于创建强大现代CLI应用程序的库，也是一个生成应用程序和命令文件的程序
2. [Gin](https://github.com/gin-gonic/gin): 轻量级Web框架，Go世界里最流行的Web框架
3. [Gin-Cors](https://github.com/gin-contrib/cors): Gin框架提供的跨域中间件
4. [GORM](https://gorm.io/index.html): ORM工具。本项目需要配合Mysql使用
5. [Go-Redis](https://github.com/go-redis/redis): Golang Redis客户端
6. [viper](github.com/spf13/viper): Viper是适用于Go应用程序的完整配置解决方案
7. [JWT](github.com/golang-jwt/jwt): 使用jwt-go这个库来实现我们生成JWT和解析JWT的功能
8. [Zap](go.uber.org/zap): Zap日志库,配置traceId,可进行链路追踪
9. [swag](github.com/swaggo/swag): 使用swag快速生成RESTFul API文档
10. [cron](github.com/robfig/cron/v3): cron是golang中广泛使用的一个定时任务库
11. [go-nsq](github.com/nsqio/go-nsq): nsq 是一款基于 go 语言开发实现的分布式消息队列组件

本项目已经预先实现了一些常用的代码方便参考和复用:

1. 创建了用户模型
2. 实现了```/api/v1/user/register```用户注册接口
3. 实现了```/api/v1/user/login```用户登录接口
4. 实现了```/api/v1/user/me```用户资料接口(需要登录后获取token)

本项目已经预先创建了一系列文件夹划分出下列模块:
```
.
├── Dockerfile
├── LICENSE.md
├── Makefile
├── README.md
├── app
│   ├── http
│   │   ├── controller  MVC框架的controller，负责协调各部件完成任务
│   │   ├── middleware  中间件
│   │   ├── request     请求参数结构体
│   │   ├── response    响应结构体
│   │   └── service     业务逻辑处理
│   └── model                 数据库模型
│       └── user.go
├── bin
├── cmd                             命令行
│   ├── main.go
│   ├── migrate.go
│   ├── server
│   │   ├── cron.go
│   │   ├── gin.go
│   │   └── mq.go
│   └── version.go
├── common                          公共的
│   └── services              公共的服务
├── config.yaml                     配置文件
├── core                            一些核心组件初始化
│   ├── global
│   │   └── core.go
│   ├── gorm.go
│   ├── init.go
│   ├── logger
│   │   ├── gorm_logger.go
│   │   └── zap.go
│   ├── middleware
│   │   └── zap.go
│   ├── nsq_producer.go
│   ├── redis.go
│   ├── trans.go
│   └── viper.go
├── crontab                         定时任务
│   ├── cron_init.go
│   ├── schedule.go
│   └── testJob.go
├── docker                          docker相关
│   └── nsq
│       └── docker-compose.yaml
├── docs                            swagger生成的文档
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── log.log
├── main.go
├── mq                              消息队列
│   ├── mq_base.go
│   └── send_registered_email.go
├── router
│   ├── api.go
│   ├── router.go
│   ├── swagger.go
│   └── warp.go
├── start.sh
└── util                            工具类
    ├── common.go
    ├── gin.go
    └── ip.go

```

## config.yaml

项目在启动的时候依赖config.yaml配置文件

## Go Mod

本项目使用[Go Mod](https://github.com/golang/go/wiki/Modules)管理依赖。

```shell
go mod init go-crud
export GOPROXY=http://mirrors.aliyun.com/goproxy/
go run main.go // 自动安装
```

## 运行HTTP

> 项目运行后启动在3000端口（可以修改，参考gin文档)

```shell
go run main.go server -c config.yaml
```

## 定时任务

```shell
go run main.go cron -c config.yaml
```

## 消息队列

### 先启动nsq服务

```shell
cd docker/nsq
docker-compose up -d
```

### 开启config.yaml配置中的`nsq`选项

```yaml
...
nsq:
  producer: 127.0.0.1:4150
  consumer: 127.0.0.1:4161
```

### 运行消费者

```shell
go run main.go nq -c config.yaml
```
