package tests

import (
	"github.com/gin-gonic/gin"
	"log"
	"mio/core/initialize"
	"mio/internal/util"
	mock_repository "mio/mock/repository"
	"mio/model/auth"
	"mio/service"
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
}
func SetupMock() {
	//real 真实环境 mock mock环境测试
	if TestEnv != "real" {
		onceMock.Do(func() {
			service.DefaultAdminService = service.NewAdminService(mock_repository.NewAdminMockRepository())
			service.DefaultUserService = service.NewUserService(mock_repository.NewUserMockRepository())
		})
	}
}
func SetupServer() *gin.Engine {
	once.Do(func() {
		initialize.InitIni("../../config.ini")
		if TestEnv == "real" {
			initialize.InitDB()
		}
		initialize.InitValidator()

		router = initialize.InitServer().Handler.(*gin.Engine)
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
func createUserToken() string {
	token, err := util.CreateToken(auth.User{
		Guid: "27c1c190-4c5f-42a1-aa2d-77298b97fe90",
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
