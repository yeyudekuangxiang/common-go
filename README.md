# 绿喵后端项目
> 项目描述
## 目录结构
mp2c-go
- build 项目部署相关文件
- controller 控制器目录
- core 项目需要的全局对象以及初始化方法
- internal 项目内部的用的包
- mock 用于测试时对数据库、Repository或者Service进行mock
- repository 仓库(数据库curd)
- server 路由和中间件
- service  逻辑
- tests 功能测试
- config.ini 配置文件
- .gitlab-ci.yml  gitlab ci&cd文件
- main.go 入口文件
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
2. 项目的部署 目前仅有develop和master分支有部署  代码合并到develop分支后会自动构建部署 代码合并到master后需要手动打tag才会构建部署
