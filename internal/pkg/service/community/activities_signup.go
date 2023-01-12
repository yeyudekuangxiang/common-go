package community

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/community"
	"mio/pkg/errno"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type (
	ActivitiesSignupService interface {
		GetPageList(params community.FindAllActivitiesSignupParams) ([]*entity.APIActivitiesSignup, int64, error)
		GetSignupInfo(id int64) (*entity.APIActivitiesSignup, error)
		FindAll(params community.FindAllActivitiesSignupParams) ([]*entity.CommunityActivitiesSignup, int64, error)
		FindSignupList(params community.FindAllActivitiesSignupParams) ([]*entity.APISignupList, int64, error)
		Signup(params SignupParams) error    //报名
		CancelSignup(Id, userId int64) error //取消报名
		Export(w http.ResponseWriter, r *http.Request, topicId int64)
	}

	defaultCommunityActivitiesSignupService struct {
		ctx         *mioContext.MioContext
		signupModel community.ActivitiesSignupModel
	}
)

func (srv defaultCommunityActivitiesSignupService) Export(w http.ResponseWriter, r *http.Request, topicId int64) {
	list, _, err := srv.signupModel.FindSignupList(community.FindAllActivitiesSignupParams{TopicId: topicId})
	f := excelize.NewFile()

	// 创建一个工作表
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		app.Logger.Errorf(fmt.Sprintf("活动报名数据Export Error:%s", err.Error()))
	}

	// 设置单元格的值
	f.SetCellValue("Sheet1", "A1", "昵称")
	f.SetCellValue("Sheet1", "B1", "真实姓名")
	f.SetCellValue("Sheet1", "C1", "联系电话")
	f.SetCellValue("Sheet1", "D1", "wechat")
	f.SetCellValue("Sheet1", "E1", "年龄")
	f.SetCellValue("Sheet1", "F1", "性别")
	f.SetCellValue("Sheet1", "G1", "居住城市")
	f.SetCellValue("Sheet1", "H1", "报名备注")
	//
	for i, item := range list {
		gender := "未知"
		if item.Gender == 1 {
			gender = "男"
		} else if item.Gender == 2 {
			gender = "女"
		}

		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), item.User.Nickname)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), item.RealName)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), item.Phone)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), item.Wechat)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), item.Age)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), gender)
		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", i+2), item.City)
		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", i+2), item.Remarks)
	}

	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	fileName := fmt.Sprintf("export_data_%d-%d.xlsx", time.Now().Unix(), topicId)

	file, err := os.OpenFile(filepath.Clean(fileName), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		app.Logger.Errorf(fmt.Sprintf("活动报名数据Export Error:%s", err.Error()))
	}
	_ = f.Write(file)

	defer func() {
		if err := file.Close(); err != nil {
			app.Logger.Errorf(fmt.Sprintf("活动报名数据Export Error:%s", err.Error()))
		}
		if err := f.Close(); err != nil {
			app.Logger.Errorf(fmt.Sprintf("活动报名数据Export Error:%s", err.Error()))
		}
	}()
	if err != nil {
		app.Logger.Errorf(fmt.Sprintf("活动报名数据Export Error:%s", err.Error()))
	}
	w.Header().Add("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	//var buffer bytes.Buffer
	buf, _ := f.WriteToBuffer()
	content := bytes.NewReader(buf.Bytes())
	http.ServeContent(w, r, fileName, time.Now(), content)
}

func (srv defaultCommunityActivitiesSignupService) FindSignupList(params community.FindAllActivitiesSignupParams) ([]*entity.APISignupList, int64, error) {
	list, total, err := srv.signupModel.FindSignupList(params)
	if err != nil {
		return nil, 0, errno.ErrInternalServer.WithMessage(err.Error())
	}
	return list, total, nil
}

func (srv defaultCommunityActivitiesSignupService) FindAll(params community.FindAllActivitiesSignupParams) ([]*entity.CommunityActivitiesSignup, int64, error) {
	list, total, err := srv.signupModel.FindAll(params)
	if err != nil {
		return nil, 0, errno.ErrInternalServer.WithMessage(err.Error())
	}
	return list, total, nil
}

func (srv defaultCommunityActivitiesSignupService) GetPageList(params community.FindAllActivitiesSignupParams) ([]*entity.APIActivitiesSignup, int64, error) {
	list, total, err := srv.signupModel.FindAllAPISignup(params)
	if err != nil {
		return nil, 0, errno.ErrInternalServer.WithMessage(err.Error())
	}
	for _, item := range list {
		item.Topic.Activity.Status = 1
		if item.Topic.Activity.SignupDeadline.Before(time.Now()) {
			item.Topic.Activity.Status = 2
		}
	}
	return list, total, nil
}

func (srv defaultCommunityActivitiesSignupService) GetSignupInfo(id int64) (*entity.APIActivitiesSignup, error) {
	signup, err := srv.signupModel.FindOneAPISignup(community.FindOneActivitiesSignupParams{Id: id})
	if err != nil {
		return &entity.APIActivitiesSignup{}, errno.ErrInternalServer.WithMessage(err.Error())
	}
	if signup.Id == 0 {
		return &entity.APIActivitiesSignup{}, errno.ErrCommon.WithMessage("未找到该标签")
	}
	return signup, nil
}

func (srv defaultCommunityActivitiesSignupService) Signup(params SignupParams) error {
	signup, err := srv.signupModel.FindOne(community.FindOneActivitiesSignupParams{
		TopicId:      params.TopicId,
		UserId:       params.UserId,
		SignupStatus: 1,
	})
	if err != nil {
		return err
	}
	if signup.Id != 0 {
		return errno.ErrCommon.WithMessage("不能重复报名哦")
	}

	signupModel := &entity.CommunityActivitiesSignup{}
	marshal, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(marshal, signupModel)
	if err != nil {
		return err
	}
	err = srv.signupModel.Create(signupModel)
	if err != nil {
		return err
	}
	return nil
}

func (srv defaultCommunityActivitiesSignupService) CancelSignup(id, userId int64) error {
	signup, err := srv.signupModel.FindOne(community.FindOneActivitiesSignupParams{Id: id, UserId: userId})
	if err != nil {
		return err
	}
	if signup.Id == 0 {
		return errno.ErrRecordNotFound
	}

	err = srv.signupModel.CancelSignup(signup)
	if err != nil {
		return err
	}
	return nil
}

func NewCommunityActivitiesSignupService(ctx *mioContext.MioContext) ActivitiesSignupService {
	return defaultCommunityActivitiesSignupService{
		ctx:         ctx,
		signupModel: community.NewCommunityActivitiesSignupModel(ctx),
	}
}
