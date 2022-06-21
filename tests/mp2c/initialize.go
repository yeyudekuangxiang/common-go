package mp2c

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mio/internal/app/mp2c/server"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/auth"
	service2 "mio/internal/pkg/service"
	"mio/internal/pkg/util"
	mock_repository "mio/mock/repository"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var router *gin.Engine
var once sync.Once
var onceMock sync.Once

var TestEnv string

func init() {
	TestEnv = strings.Trim(os.Getenv("TEST_ENV"), " ")
	if TestEnv == "" {
		TestEnv = "mock"
		_ = os.Setenv("TEST_ENV", TestEnv)
	}
	setupMock()
}
func setupMock() {
	//real 真实环境 mock mock环境测试
	if TestEnv != "real" {
		onceMock.Do(func() {
			service2.DefaultSystemAdminService = service2.NewSystemAdminService(mock_repository.NewAdminMockRepository())
			service2.DefaultUserService = service2.NewUserService(mock_repository.NewUserMockRepository())
		})
	}
}
func SetupServer() *gin.Engine {
	once.Do(func() {
		if TestEnv == "real" {
			fmt.Println(os.Getenv("Config"))
			initialize.Initialize(os.Getenv("Config"))
		}
		initialize.InitValidator()

		router = server.InitServer().Handler.(*gin.Engine)
	})
	return router
}
func AddAuthToken(req *http.Request) {
	req.Header.Set("Token", createUserToken())
}
func AddUserToken(req *http.Request) {
	req.Header.Set("Token", createUserToken())
}
func AddAdminToken(req *http.Request) {
	req.Header.Set("Token", createAdminToken())
}
func AddBusinessToken(req *http.Request) {
	req.Header.Set("Token", createBusinessUserToken())
}
func createBusinessUserToken() string {
	token, err := util.CreateToken(auth.BusinessUser{
		Uid:       "test",
		CreatedAt: model.NewTime(),
	})
	if err != nil {
		log.Fatal("create token err:", err)
	}
	return token
}
func createUserToken() string {
	token, err := util.CreateToken(auth.User{
		ID: 1,
	})
	if err != nil {
		log.Fatal("create token err:", err)
	}
	return token
}
func createAdminToken() string {
	token, err := util.CreateToken(auth.Admin{
		ID: 1,
	})
	if err != nil {
		log.Fatal("create token err:", err)
	}
	return token
}
