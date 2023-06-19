package community

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"mio/config"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/community"
	"mio/internal/pkg/service/track"
	"mio/pkg/errno"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var DefaultActivitiesSignupService = NewCommunityActivitiesSignupService(mioContext.NewMioContext())

type (
	ActivitiesSignupService interface {
		GetPageList(params community.FindAllActivitiesSignupParams) ([]*entity.APIActivitiesSignup, int64, error)
		GetSignupInfo(params community.FindOneActivitiesSignupParams) (*entity.APIActivitiesSignup, bool, error)
		FindAll(params community.FindAllActivitiesSignupParams) ([]*entity.CommunityActivitiesSignup, int64, error)
		FindSignupList(params community.FindAllActivitiesSignupParams) ([]*entity.APISignupList, int64, error)
		Signup(params SignupParams) error    //报名
		CancelSignup(Id, userId int64) error //取消报名
		Export(w http.ResponseWriter, r *http.Request, topicId int64)
		FindListCount(params FindListCountReq) ([]*entity.APIListCount, error)
	}

	defaultCommunityActivitiesSignupService struct {
		ctx         *mioContext.MioContext
		signupModel community.ActivitiesSignupModel
		topicModel  community.TopicModel
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

func (srv defaultCommunityActivitiesSignupService) GetSignupInfo(params community.FindOneActivitiesSignupParams) (*entity.APIActivitiesSignup, bool, error) {
	signup, err := srv.signupModel.FindOneAPISignup(community.FindOneActivitiesSignupParams{
		Id:           params.Id,
		TopicId:      params.TopicId,
		UserId:       params.UserId,
		SignupStatus: params.SignupStatus,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &entity.APIActivitiesSignup{}, false, nil
		}
		return &entity.APIActivitiesSignup{}, false, errno.ErrInternalServer.WithMessage(err.Error())
	}

	return signup, true, nil
}

func (srv defaultCommunityActivitiesSignupService) Signup(params SignupParams) error {
	topic, err := srv.findTopic(params.TopicId, 1, 3)
	if err != nil {
		return err
	}
	//查看是否已经报名
	_, err = srv.findSignupRecord(params.TopicId, params.UserId, 1)
	if err != nil {
		return err
	}
	//查看是否超过上限
	err = srv.checkSignupNum(topic.Id, int64(topic.Activity.SignupNumber))
	if err != nil {
		return err
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
	track.DefaultSensorsService().Track(false, config.SensorsEventName.ActivityApply, params.OpenId, map[string]interface{}{
		"title":      topic.Title,
		"topic_id":   int(params.TopicId),
		"apply_time": params.SignupTime,
	})
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

func (srv defaultCommunityActivitiesSignupService) FindListCount(params FindListCountReq) ([]*entity.APIListCount, error) {
	list, err := srv.signupModel.FindListCount(community.FindListCountParams{TopicIds: params.TopicIds})
	if err != nil {
		return nil, errno.ErrInternalServer.WithMessage(err.Error())
	}
	return list, nil
}

func NewCommunityActivitiesSignupService(ctx *mioContext.MioContext) ActivitiesSignupService {
	return defaultCommunityActivitiesSignupService{
		ctx:         ctx,
		signupModel: community.NewCommunityActivitiesSignupModel(ctx),
		topicModel:  community.NewTopicModel(ctx),
	}
}

func (srv defaultCommunityActivitiesSignupService) findTopic(id int64, tp, status int) (*entity.Topic, error) {
	topic := srv.topicModel.FindById(id)
	if topic.Id == 0 {
		return nil, errno.ErrCommon.WithMessage("活动不存在")
	}
	return topic, nil
}

func (srv defaultCommunityActivitiesSignupService) findSignupRecord(id, uid int64, signupStatus int) (*entity.CommunityActivitiesSignup, error) {
	signup, err := srv.signupModel.FindOne(community.FindOneActivitiesSignupParams{
		TopicId:      id,
		UserId:       uid,
		SignupStatus: signupStatus,
	})
	if err != nil {
		return nil, err
	}
	if signup.Id != 0 {
		return nil, errno.ErrCommon.WithMessage("不能重复报名哦")
	}
	return signup, nil
}

func (srv defaultCommunityActivitiesSignupService) checkSignupNum(id, num int64) error {
	count, err := srv.signupModel.FindListCount(community.FindListCountParams{TopicIds: []int64{id}})
	if err != nil {
		return nil
	}

	if len(count) > 0 && count[0].NumOfSignup >= num {
		return errno.ErrCommon.WithMessage("报名人数已满")
	}
	return nil
}

//func (srv defaultCommunityActivitiesSignupService) checkSignupInfo(params SignupParams) error {
//	info := "[{\"type\":1,\"title\":\"姓名\",\"sort\":1},{\"type\":2,\"title\":\"性别\",\"sort\":2,\"options\":{\"option1\":\"男\",\"option2\":\"女\"}},{\"type\":4,\"title\":\"爱好\",\"sort\":3,\"options\":{\"option1\":\"唱\",\"option2\":\"跳\",\"option3\":\"rap\",\"option4\":\"篮球\"}},{\"type\":3,\"title\":\"备注\",\"sort\":4}]"
//	infos := make([]interface{}, 0)
//	err := json.Unmarshal([]byte(info), &infos)
//	if err != nil {
//		return err
//	}
//	if len(infos) == 0 {
//		return errno.ErrCommon
//	}
//	//for _, item := range infos {
//		item是map[string]interface / map[string]map[string]interface
//
//	}
//	fmt.Println(infos)
//	return nil
//}
