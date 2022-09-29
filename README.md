# 绿喵微服务

## 微服务框架使用[go-zero](https://github.com/zeromicro/go-zero)
## 项目架构参考 [go-zero-looklook](https://github.com/Mikaelemmmm/go-zero-looklook)

## 开发准备
### 1. 克隆运行环境
`git clone -b runtime git@gitlab.miotech.com:miotech-application/backend/mp2c-micro.git mp2c-micro-runtime`

### 2. 打开命令行，并将目录切换到mp2c-micro-runtime的文件夹下
`cd mp2c-micro-runtime`

### 3. 运行 ```make init``` 使用自定义的goctl模版文件

### 4. 运行```docker-compose -f docker-compose-env.yml up -d``` 命令安装运行环境(非必需)
运行此命令将会在docker中安装 jaeger、prometheus、grafana、elasticsearch、kibana、go-stash、filebeat、zookeeper、kafka、asynqmon、mysql、redis

### 5. 运行```docker-compose up -d```(非必需)
运行此命令将会安装nginx用于网关、并且将所有的微服务应用在 lyumikael/gomodd 容器中运行

## 项目开发
### 1. 将本项目clone到本地
`git clone git@gitlab.miotech.com:miotech-application/backend/mp2c-micro.git`  
`cd mp2c-micro`
## 开发model模块
### 1. 运行  ```make model``` 生成数据库代码
运行此命令将要求输入模块名称和表名 用空格隔开 输入后回车 将会在 ```app/应用名/model``` 目录下生成数据库代码  
此命令默认的数据库配置为测试数据库
并且要求数据库中表结构已经存在才可以生成数据库代码  
同时要求将数据库sql文件放在 ```app/应用名/model/sql``` 目录下
## 开发rpc应用
### 1.运行 ```make rpc``` 生成rpc proto文件
运行此命令将要求输入要开发的模块名称 输入后回车 将会生成```app/应用名/cmd/pb/应用名.prpto```文件   
在此文件中编辑服务方法

### 2.运行 ```make rpcgo```  生成rpc应用代码
运行此命令将要求输入要开发的模块名称 输入后回车 将会根据```app/应用名/cmd/rpc/pb/应用名.prpto```文件生成rpc应用代码  
在```app/应用名/cmd/rpc/internal/config``` 下编辑配置文件模版  
在```app/应用名/cmd/rpc/etc``` 下编辑配置文件  
在```app/应用名/cmd/rpc/internal/logic``` 下完成业务逻辑  
在```app/应用名/cmd/rpc/internal/svc``` 初始化第三方rpc服务、数据库model、redis等操作

## 开发api应用

### 1.运行 ```make api```
运行此命令将要求输入要开发的模块名称 输入后回车 将会生成```app/应用名/cmd/api/desc/应用名.api```文件   
在此文件中编辑api接口信息

### 2.运行 ```make apigo```
运行此命令将要求输入要开发的模块名称 输入后回车 将会根据```app/应用名/cmd/api/desc/应用名.api```文件生成api应用代码  
在```app/应用名/cmd/api/internal/config``` 下编辑配置文件模版  
在```app/应用名/cmd/api/etc``` 下编辑配置文件  
在```app/应用名/cmd/api/internal/logic``` 下完成业务逻辑  
在```app/应用名/cmd/api/internal/svc``` 初始化第三方rpc服务等操作

## 项目部署
部署之前如有配置文件的变更或者新增 联系阿金编辑  
在develop分支上打tag develop-应用名-应用类型(api|rpc)-v1.0.0 会自动将此服务部署到测试环境  
在master分支上打tag production-应用名-应用类型(api|rpc)-v1.0.0 会自动将此服务部署到正式环境  
每次在master分支上打tag后要在master分支上再打一个类型 v1.0.0 的标签 方便mp2c-go项目引用此项目