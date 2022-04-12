# 绿喵后端项目
> 项目描述
## 目录结构
mp2c-go
- build 各个项目部署相关文件
- cmd 各个项目的启动目录
- config 应用程序的全局配置目录 
    - config.go 全局配置
    - redis.go redis键配置列表
- internal 各个项目的主要逻辑目录
    - app 各个项目的实际应用代码
    - pkg 各个项目公用的逻辑代码
        - core 应用程序通用的客户端和初始化方法
        - model 各种实体文件
        - repository 数据库增删改查逻辑
        - service 项目的逻辑处理
        - util 工具目录
- mock 用于测试时对数据库、Repository或者Service进行mock
- pkg 外部应用程序可以使用的库代码
- tests 功能测试
- config.ini.example 配置文件模版
- .gitlab-ci.yml  gitlab ci&cd文件

### 数据流动 controller->service->repository
## 编码风格
1. 接口返回数据为application/json  返回数据结构为`{"code":200,data:{},messsage:"操作成功"}`
所有返回字段键名均为小驼峰形式
2. 所有用于获取参数、返回数据、方法参数的结构体都定义在同级目录下的types.go文件里 若结构体太多也可考虑放在同级types目录下面
3. controller代码获取参数均使用 [gin binding](https://gin-gonic.com/zh-cn/docs/examples/multipart-urlencoded-binding/) 和 [go-playground/validator](https://github.com/go-playground/validator) 进行绑定和验证
4. service和repository中的方法必须写在结构体中 controller不做必须要求但尽量也写成结构体的形式
5. 所有写的接口尽量在tests目录下写功能测试、有空余时间也可多写单元测试
## 项目开发流程
1. 分支管理规范参考 <https://www.jianshu.com/p/7ae40a051cb8>
2. 代码提交规范参考 <https://blog.csdn.net/github_39506988/article/details/90298780>
3. 项目的部署 目前仅有develop和master分支有部署  代码合并到develop分支后会自动构建部署 代码合并到master后需要手动打tag才会构建部署

## 项目运行
1. `git clone https://gitlab.miotech.com/miotech-application/backend/mp2c-go.git`
2. 下载所需库文件 `go mod download`
3. 复制配置文件`copy config.ini.example config.ini` 并且完善配置文件config.ini
4. `go run cmd/mp2c/main.go`