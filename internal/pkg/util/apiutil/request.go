package apiutil

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/model"
	entity "mio/internal/pkg/model/entity"
	"mio/internal/pkg/model/entity/business"
	"mio/pkg/errno"
	"mio/pkg/validator"
)

func BindForm(c *gin.Context, data interface{}) error {
	if err := c.ShouldBind(data); err != nil {
		err = validator.TranslateError(err)
		return errno.ErrBind.With(err)
	}
	return nil
}
func GetAuthAdmin(c *gin.Context) entity.SystemAdmin {
	if admin, ok := c.Get("AuthAdmin"); ok {
		return admin.(entity.SystemAdmin)
	}
	return entity.SystemAdmin{}
}
func GetAuthUser(c *gin.Context) entity.User {
	if user, ok := c.Get("AuthUser"); ok {
		return user.(entity.User)
	}
	return entity.User{}
}
func GetAuthBusinessUser(c *gin.Context) business.User {
	if user, ok := c.Get("BusinessUser"); ok {
		return user.(business.User)
	}
	return business.User{
		ID:            3,
		Uid:           "mock-uid-3",
		BCompanyId:    1,
		BDepartmentId: 2,
		Nickname:      "greencat",
		Mobile:        "13000000000",
		TelephoneCode: "86",
		Realname:      "绿喵",
		Avatar:        "https://miotech-resource.oss-cn-hongkong.aliyuncs.com/static/mp2c/images/topic/mio-kol/mio-avatar.jpg",
		Status:        1,
		CreatedAt:     model.NewTime(),
		UpdatedAt:     model.NewTime(),
	}
}
