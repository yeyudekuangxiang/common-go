package community

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository/community"
	"mio/pkg/errno"
)

type (
	ActivitiesSignupService interface {
		GetPageList(params community.FindAllActivitiesSignupParams) ([]entity.APIActivitiesSignup, int64, error)
		GetOne(id int64) (entity.CommunityActivitiesSignup, error)
		Signup(params SignupParams) error    //报名
		CancelSignup(Id, userId int64) error //取消报名
	}

	defaultCommunityActivitiesSignupService struct {
		ctx         *mioContext.MioContext
		signupModel community.ActivitiesSignupModel
	}
)

func NewCommunityActivitiesSignupService(ctx *mioContext.MioContext) ActivitiesSignupService {
	return defaultCommunityActivitiesSignupService{
		ctx:         ctx,
		signupModel: community.NewCommunityActivitiesSignupModel(ctx),
	}
}

func (srv defaultCommunityActivitiesSignupService) GetPageList(params community.FindAllActivitiesSignupParams) ([]entity.APIActivitiesSignup, int64, error) {
	list, total, err := srv.signupModel.FindAllAPISignup(params)
	if err != nil {
		return nil, 0, errno.ErrInternalServer.WithMessage(err.Error())
	}
	return list, total, nil
}

func (srv defaultCommunityActivitiesSignupService) GetOne(id int64) (entity.CommunityActivitiesSignup, error) {
	signup, err := srv.signupModel.FindOne(community.FindOneActivitiesSignupParams{Id: id})
	if err != nil {
		return entity.CommunityActivitiesSignup{}, errno.ErrInternalServer.WithMessage(err.Error())
	}
	if signup.Id == 0 {
		return entity.CommunityActivitiesSignup{}, errno.ErrCommon.WithMessage("未找到该标签")
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

	signupModel := &entity.CommunityActivitiesSignup{
		TopicId:      params.TopicId,
		UserId:       params.UserId,
		RealName:     params.RealName,
		Phone:        params.Phone,
		Gender:       params.Gender,
		Age:          params.Age,
		Wechat:       params.Wechat,
		City:         params.City,
		Remarks:      params.Remarks,
		SignupTime:   params.SignupTime,
		SignupStatus: params.SignupStatus,
	}
	//
	//marshal, err := json.Marshal(params)
	//if err != nil {
	//	return err
	//}
	//
	//err = json.Unmarshal(marshal, signupModel)
	//if err != nil {
	//	return err
	//}
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

	err = srv.signupModel.CancelSignup(&signup)
	if err != nil {
		return err
	}
	return nil
}
