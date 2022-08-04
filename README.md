# 绿喵后端项目
> 项目描述
## 目录结构
mp2c-go
- [build](build) 各个项目部署相关文件
- [cmd](cmd) 各个项目的启动目录
- [config](config) 应用程序的全局配置目录 
    - [config.go](config/config.go) 全局配置
    - [redis.go](config/redis.go) redis键配置列表
- [internal](internal) 各个项目的主要逻辑目录
    - [app](internal/app) 各个项目的实际应用代码
    - [pkg](internal/pkg) 各个项目公用的逻辑代码
        - [core](internal/pkg/core) 应用程序通用的客户端和初始化方法
        - [model](internal/pkg/model) 各种实体文件
        - [repository](internal/pkg/repository) 数据库增删改查逻辑
        - [service](internal/pkg/service) 项目的逻辑处理
        - [util](internal/pkg/util) 工具目录
- [mock](mock) 用于测试时对数据库、Repository或者Service进行mock
- [pkg](pkg) 外部应用程序可以使用的库代码
- [tests](tests) 功能测试
- [config.ini.example](config.ini.example) 配置文件模版
- [.gitlab-ci.yml](.gitlab-ci.yml)  gitlab ci&cd文件

### 数据流动 controller->service->repository
## 编码风格
1. 接口返回数据为application/json  返回数据结构为`{"code":200,data:{},messsage:"操作成功"}`
所有返回字段键名均为小驼峰形式
2. 所有用于获取参数、返回数据、方法参数的结构体都定义在同级目录下的types.go文件里 若结构体太多也可考虑放在同级types目录下面
3. controller代码获取参数均使用 [gin binding](https://gin-gonic.com/zh-cn/docs/examples/multipart-urlencoded-binding/) 和 [go-playground/validator](https://github.com/go-playground/validator) 进行绑定和验证
4. service和repository中的方法必须写在结构体中 controller不做必须要求但尽量也写成结构体的形式
5. 所有写的接口尽量在tests目录下写功能测试、有空余时间也可多写单元测试
6. 自己写的模块下要写readme文件以帮助其他同事使用
## 项目开发流程
> 分支管理规范参考 [gitlab flow](https://docs.gitlab.com/ee/topics/gitlab_flow.html)  
> 代码提交规范参考 <https://blog.csdn.net/github_39506988/article/details/90298780>
1. 新功能的开发 - 从远程develop分支切出 /^feature-.+$/ 分支,开发完成推送远程后手动运行通道会将此分支部署到测试环境 部署成功后地址为 dev-domain/feature-(.+)/ 具体地址可在 [环境](https://gitlab.miotech.com/miotech-application/backend/mp2c-go/-/environments) 中查看 此分支删除后会自动删除已部署的环境
2. 测试发布 - 新功能开发完成后请求合并到develop分支 , 合并完成后会自动将develop分支部署到测试环境 部署成功后地址为 dev-domain 具体地址可在 [环境](https://gitlab.miotech.com/miotech-application/backend/mp2c-go/-/environments) 中查看
3. 预发布 - develop分支请求合并到master分支,合并完成后会自动将master分支部署到正式环境 部署成功后地址为 prod-domain/pre-prod/ 具体地址可在 [环境](https://gitlab.miotech.com/miotech-application/backend/mp2c-go/-/environments) 中查看
4. 发布正式版本 - 在master分支上打tag后会将此tag自动部署到正式环境 部署成功后地址为 prod-domain 具体地址可在 [环境](https://gitlab.miotech.com/miotech-application/backend/mp2c-go/-/environments) 中查看
5. bug的修复 紧急bug可直接从master分支切出 /^hotfix-.+$/分支,推送远程后手动运行通道会将此分支自动部署的正式环境  部署成功后地址为 prod-domain/hotfix-(.+)/ 具体地址可在 [环境](https://gitlab.miotech.com/miotech-application/backend/mp2c-go/-/environments) 中查看 此分支删除后会自动删除已部署的环境。修复完成后手动将此分支分别请求合并到master和develop分支
## 项目运行
1. `git clone https://gitlab.miotech.com/miotech-application/backend/mp2c-go.git`
2. 下载所需库文件 `go mod download`
3. 复制配置文件`copy config.ini.example config.ini` 并且完善配置文件config.ini
4. `go run cmd/mp2c/main.go`

## [开发规范](DEVELOPMENT_GUIDE.md)
## [避坑指南](NOTES.md)