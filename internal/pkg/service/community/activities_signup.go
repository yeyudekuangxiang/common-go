package community

import (
	"encoding/json"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/community"
	"mio/pkg/errno"
	"time"
)

type (
	ActivitiesSignupService interface {
		GetPageList(params community.FindAllActivitiesSignupParams) ([]*entity.APIActivitiesSignup, int64, error)
		GetOne(id int64) (*entity.CommunityActivitiesSignup, error)
		FindAll(params community.FindAllActivitiesSignupParams) ([]*entity.CommunityActivitiesSignup, int64, error)
		FindSignupList(params community.FindAllActivitiesSignupParams) ([]*entity.APISignupList, int64, error)
		Signup(params SignupParams) error    //报名
		CancelSignup(Id, userId int64) error //取消报名
	}

	defaultCommunityActivitiesSignupService struct {
		ctx         *mioContext.MioContext
		signupModel community.ActivitiesSignupModel
	}
)

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
		if item.Topic.Activity.Status != 3 {
			item.Topic.Activity.Status = 3
		}
	}
	return list, total, nil
}

func (srv defaultCommunityActivitiesSignupService) GetOne(id int64) (*entity.CommunityActivitiesSignup, error) {
	signup, err := srv.signupModel.FindOne(community.FindOneActivitiesSignupParams{Id: id})
	if err != nil {
		return &entity.CommunityActivitiesSignup{}, errno.ErrInternalServer.WithMessage(err.Error())
	}
	if signup.Id == 0 {
		return &entity.CommunityActivitiesSignup{}, errno.ErrCommon.WithMessage("未找到该标签")
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
		return nil
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
