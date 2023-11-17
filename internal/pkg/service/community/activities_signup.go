package community

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/converttool"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/sorttool"
	"gorm.io/gorm"
	"math"
	"mio/config"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/community"
	"mio/internal/pkg/service/track"
	"mio/pkg/errno"
	"strconv"
	"strings"
	"time"
)

var DefaultActivitiesSignupService = NewCommunityActivitiesSignupService(mioContext.NewMioContext())

type (
	ActivitiesSignupService interface {
		GetPageList(params community.FindAllActivitiesSignupParams) ([]*entity.APIActivitiesSignup, int64, error)
		GetSignupInfo(params community.FindOneActivitiesSignupParams) (*entity.APIActivitiesSignup, bool, error)
		FindAll(params community.FindAllActivitiesSignupParams) ([]*entity.CommunityActivitiesSignup, int64, error)
		FindSignupList(params community.FindAllActivitiesSignupParams) ([]*entity.APISignupList, int64, error)
		Signup(params SignupInfosParams) error //报名
		CancelSignup(Id, userId int64) error   //取消报名
		Export(topicId int64) (*bytes.Buffer, error)
		FindListCount(params FindListCountReq) ([]*entity.APIListCount, error)
		SignupV2(params SignupParams) error
		GetPageListV2(params community.FindAllActivitiesSignupParams) ([]*entity.APIActivitiesSignupV2, int64, error)
		GetSignupInfoV2(params community.FindOneActivitiesSignupParams) (*entity.APIActivitiesSignupV2, bool, error)
	}

	defaultCommunityActivitiesSignupService struct {
		ctx         *mioContext.MioContext
		signupModel community.ActivitiesSignupModel

		topicModel community.TopicModel
	}
)

func ToExcelColumn(column int) string {
	if column <= 0 {
		return ""
	}

	columnName := ""
	for column > 0 {
		// 由于Excel列从1开始，所以减去1以便从0开始计数
		column--
		// 计算当前位置的字母
		letter := rune('A' + (column % 26))
		// 将字母添加到列名的开头
		columnName = string(letter) + columnName
		// 移动到下一个字母位置
		column = column / 26
	}
	return columnName
}
func (srv defaultCommunityActivitiesSignupService) Export(topicId int64) (*bytes.Buffer, error) {

	//获取报名信息
	activityInfo := entity.CommunityActivities{}
	err := app.DB.Where("id = ?", topicId).Find(&activityInfo).Error
	if err != nil {
		app.Logger.Errorf(fmt.Sprintf("报名信息不存在 %d Error:%s", topicId, err.Error()))
		return nil, err
	}

	//获取列信息
	colInfoList := make([]SignupInfo, 0)
	err = json.Unmarshal([]byte(activityInfo.SATag), &colInfoList)
	if err != nil {
		app.Logger.Errorf(fmt.Sprintf("解析报名信息异常 %d Error:%s", topicId, err.Error()))
		return nil, err
	}

	//获取报名列表
	list, _, err := srv.signupModel.FindSignupList(community.FindAllActivitiesSignupParams{TopicId: topicId})
	if err != nil {
		return nil, err
	}
	// 用于存储没一列的数据
	colsListMap := make(map[string][]string)
	// 用户存储列信息
	colMap := make(map[string]string)

	// 处理数据
	for i, item := range list {
		signupInfos := make([]SignupInfo, 0)
		err = json.Unmarshal([]byte(item.SignupInfo), &signupInfos)
		if err != nil {
			return nil, err
		}
		for _, signupInfo := range signupInfos {
			value := signupInfo.Value
			if signupInfo.Code == "gender" {
				gender := srv.toString(signupInfo.Value)
				if gender == "1" {
					value = "女"
				} else if gender == "2" {
					value = "男"
				}
			}

			colName := signupInfo.Code
			rows, ok := colsListMap[colName]
			colMap[colName] = signupInfo.Title
			if !ok {
				rows = make([]string, len(list))
			}
			rows[i] = srv.toString(value)
			colsListMap[colName] = rows
		}
	}

	// 导出到excel
	f := excelize.NewFile()
	defer f.Close()
	colI := 1
	// 先按照发起人最后一次编辑的顺序设置列数据
	for _, col := range colInfoList {
		err = f.SetCellStr("Sheet1", ToExcelColumn(colI)+"1", col.Title)
		if err != nil {
			return nil, err
		}
		list := colsListMap[col.Code]
		for i, v := range list {
			err = f.SetCellStr("Sheet1", ToExcelColumn(colI)+strconv.Itoa(i+2), v)
			if err != nil {
				return nil, err
			}
		}
		delete(colMap, col.Code)
		colI++
	}

	var setErr error
	//剩余的列数据按照ascii码循序设置列数据
	sorttool.Map(colMap, func(key interface{}) {
		if setErr != nil {
			return
		}
		k := key.(string)
		title := colMap[k]
		setErr = f.SetCellStr("Sheet1", ToExcelColumn(colI)+"1", title)
		if setErr != nil {
			return
		}
		list := colsListMap[k]
		for i, v := range list {
			setErr = f.SetCellStr("Sheet1", ToExcelColumn(colI)+strconv.Itoa(i+2), v)
			if setErr != nil {
				return
			}
		}
		colI++
	})
	if setErr != nil {
		return nil, setErr
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf, nil
}
func (srv defaultCommunityActivitiesSignupService) toString(v interface{}) string {
	switch v.(type) {
	case bool:
		return strconv.FormatBool(v.(bool))
	case float64:
		v2 := v.(float64)
		if v2 == math.Trunc(v2) {
			return strconv.FormatInt(int64(v2), 10)
		} else {
			return fmt.Sprintf("%.2f", v2)
		}
	case string:
		return v.(string)
	case []interface{}:
		list := v.([]interface{})
		b := strings.Builder{}
		for _, item := range list {
			b.WriteString(srv.toString(item))
			b.WriteString("/")
		}
		return strings.TrimRight(b.String(), "/")
	case map[string]interface{}:
		return fmt.Sprintf("%+v", v)
	case nil:
		return ""
	}
	return fmt.Sprintf("%+v", v)
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

func (srv defaultCommunityActivitiesSignupService) GetPageListV2(params community.FindAllActivitiesSignupParams) ([]*entity.APIActivitiesSignupV2, int64, error) {
	list, total, err := srv.signupModel.FindAllAPISignupV2(params)
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

func (srv defaultCommunityActivitiesSignupService) GetSignupInfoV2(params community.FindOneActivitiesSignupParams) (*entity.APIActivitiesSignupV2, bool, error) {
	signup, err := srv.signupModel.FindOneAPISignupV2(community.FindOneActivitiesSignupParams{
		Id:           params.Id,
		TopicId:      params.TopicId,
		UserId:       params.UserId,
		SignupStatus: params.SignupStatus,
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &entity.APIActivitiesSignupV2{}, false, nil
		}
		return &entity.APIActivitiesSignupV2{}, false, errno.ErrInternalServer.WithMessage(err.Error())
	}

	return signup, true, nil
}

func (srv defaultCommunityActivitiesSignupService) Signup(params SignupInfosParams) error {
	topic, err := srv.findTopic(params.TopicId, 1, 3)
	if err != nil {
		return err
	}
	var phone string
	for _, item := range params.SignupInfos {
		if item.Code == "phone" {
			phone = item.Value.(string)
		}
	}
	//查看是否已经报名 用户id
	_, err = srv.findSignupRecord(params.TopicId, params.UserId, 1, phone)
	if err != nil {
		return err
	}
	//查看是否超过上限
	err = srv.checkSignupNum(topic.Id, int64(topic.Activity.SignupNumber))
	if err != nil {
		return err
	}

	marshal, err := json.Marshal(params.SignupInfos)
	if err != nil {
		return err
	}
	signupModel := &entity.CommunityActivitiesSignup{
		TopicId:      params.TopicId,
		UserId:       params.UserId,
		Phone:        phone,
		SignupInfo:   string(marshal),
		SignupTime:   params.SignupTime,
		SignupStatus: params.SignupStatus,
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

func (srv defaultCommunityActivitiesSignupService) SignupV2(params SignupParams) error {
	topic, err := srv.topicModel.FindOneTopic(repository.FindTopicParams{
		TopicId: params.TopicId,
		Type:    converttool.PointerInt(1),
		Status:  3,
	})
	if err != nil {
		if err == entity.ErrNotFount {
			return errno.ErrCommon.WithMessage("活动不存在")
		}
		return errno.ErrCommon.WithMessage(err.Error())
	}

	signup, err := srv.signupModel.FindOneV2(community.FindOneActivitiesSignupParams{
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

	signupModel := &entity.CommunityActivitiesSignupV2{}
	marshal, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(marshal, signupModel)
	if err != nil {
		return err
	}
	err = srv.signupModel.CreateV2(signupModel)
	if err != nil {
		return err
	}
	//诸葛打点
	/*	zhuGeAttr := make(map[string]interface{}, 0)
		zhuGeAttr["活动id"] = params.TopicId
		zhuGeAttr["活动名称"] = topic.Title
		zhuGeAttr["作者名称"] = topic.User.Nickname
		zhuGeAttr["报名者id"] = params.UserId
		zhuGeAttr["报名时间"] = params.SignupTime
		track.DefaultZhuGeService().Track(config.ZhuGeEventName.PostSignUp, params.OpenId, zhuGeAttr)
	*/
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

func (srv defaultCommunityActivitiesSignupService) findSignupRecord(id, uid int64, signupStatus int, phone string) (*entity.CommunityActivitiesSignup, error) {
	signup, err := srv.signupModel.FindOne(community.FindOneActivitiesSignupParams{
		TopicId:      id,
		UserId:       uid,
		Phone:        phone,
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

	if num == 0 {
		return errno.ErrCommon.WithMessage("报名人数已满")
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
