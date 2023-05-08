package qnr

import (
	"mio/config"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	qnrEntity "mio/internal/pkg/model/entity/qnr"
	repo "mio/internal/pkg/repository"
	repoQnr "mio/internal/pkg/repository/qnr"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
)

var DefaultAnswerService = OptionService{ctx: context.NewMioContext()}

func NewAnswerService(ctx *context.MioContext) *AnswerService {
	return &AnswerService{
		ctx:         ctx,
		answerRepo:  repoQnr.NewAnswerRepository(ctx),
		qrnUserRepo: repoQnr.NewUserRepository(ctx),
		Invite:      repo.NewInviteRepository(ctx),
		channel:     repo.DefaultUserChannelRepository,
		user:        repo.NewUserRepository(),
	}
}

type AnswerService struct {
	ctx         *context.MioContext
	answerRepo  *repoQnr.AnswerRepository
	qrnUserRepo *repoQnr.UserRepository
	Invite      *repo.InviteRepository
	user        repo.UserRepository
	channel     repo.UserChannelRepository
}

func (srv AnswerService) CreateInBatches(dto []srv_types.CreateQnrAnswerDTO) error {
	list := make([]qnrEntity.Answer, 0)
	for _, answerDTO := range dto {
		list = append(list, qnrEntity.Answer{
			QnrId:     answerDTO.QnrId,
			SubjectId: answerDTO.SubjectId,
			UserId:    answerDTO.UserId,
			Answer:    answerDTO.Answer,
		})
	}
	err := srv.answerRepo.CreateInBatches(list)
	return err
}

type Ans struct {
	Id     int64  `json:"id"`
	Answer string `json:"answer"`
}

func (srv AnswerService) Add(dto srv_types.AddQnrAnswerDTO) error {
	//查询用户是否入库，入库并回答过问题
	info := srv.qrnUserRepo.FindBy(repotypes.GetQuestUserGetById{OpenId: dto.OpenId})
	if info.UserId != 0 {
		return errno.ErrCommon.WithMessage("您已经提交了")
	}
	//获取用户信息
	userInfo := srv.user.GetUserById(dto.UserId)
	id, err2 := util.SnowflakeID()
	if err2 != nil {
		return errno.ErrCommon.With(err2)
	}
	//获取渠道信息
	channelName := ""
	channel := srv.channel.FindByCid(repo.FindUserChannelBy{Cid: userInfo.ChannelId})
	if channel.ID != 0 {
		channelName = channel.Name
	}

	//获取邀请关系
	InvitedByOpenId := ""
	inviteInfo := srv.Invite.GetInvite(userInfo.OpenId)
	if inviteInfo.InvitedByOpenId != "" {
		InvitedByOpenId = inviteInfo.InvitedByOpenId
	}

	//保存用户信息
	errUser := srv.qrnUserRepo.Create(&qnrEntity.User{
		UserId:      id.Int64(),
		ThirdId:     userInfo.OpenId,
		InvitedById: InvitedByOpenId,
		Phone:       userInfo.PhoneNumber,
		Channel:     channelName,
		Ip:          userInfo.Ip,
		City:        userInfo.CityCode,
	})
	if errUser != nil {
		return errno.ErrCommon.WithMessage("用户信息保存失败")
	}
	//保存答案
	createList := make([]srv_types.CreateQnrAnswerDTO, 0)
	for _, l := range dto.Answer {
		createList = append(createList, srv_types.CreateQnrAnswerDTO{
			Answer:    l.Answer,
			QnrId:     1,
			SubjectId: l.Id,
			UserId:    model.LongID(id.Int64()),
		})
	}
	err := srv.CreateInBatches(createList)
	if err != nil {
		return errno.ErrCommon.WithMessage("保存答案失败")
	}
	//诸葛上报
	/*attr := make(map[string]interface{}, 8)
	attr["昵称"] = userInfo.Nickname
	attr["手机号"] = userInfo.PhoneNumber
	attr["渠道名称"] = channelName
	attr["邀请人的openid"] = InvitedByOpenId
	srv.ToZhuGe(userInfo.OpenId, attr) */

	track.DefaultSensorsService().Track(true, config.SensorsEventName.Qnr, userInfo.GUID, map[string]interface{}{
		"channel_name": channelName,
		"invited_guid": InvitedByOpenId,
	})

	return nil
}

func (srv AnswerService) ToZhuGe(openId string, attr map[string]interface{}) {
	eventName := config.ZhuGeEventName.Qnr
	track.DefaultZhuGeService().Track(eventName, openId, attr)
}
