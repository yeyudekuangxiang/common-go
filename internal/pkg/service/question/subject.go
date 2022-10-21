package question

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	qnrEntity "mio/internal/pkg/model/entity/question"
	repoQnr "mio/internal/pkg/repository/question"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
)

var DefaultSubjectService = SubjectService{ctx: context.NewMioContext()}

func NewSubjectService(ctx *context.MioContext) *SubjectService {
	return &SubjectService{
		ctx:        ctx,
		repo:       repoQnr.NewSubjectRepository(ctx),
		repoOption: repoQnr.NewOptionRepository(ctx),
		repoUser:   repoQnr.NewUserRepository(ctx),
		repoAnswer: repoQnr.NewAnswerRepository(ctx),
	}
}

type SubjectService struct {
	ctx        *context.MioContext
	repo       *repoQnr.SubjectRepository
	repoOption *repoQnr.OptionRepository
	repoUser   *repoQnr.UserRepository
	repoAnswer *repoQnr.AnswerRepository
}

func (srv SubjectService) GetPageList(dto srv_types.GetQuestionSubjectDTO) ([]qnrEntity.Subject, error) {
	list, err := srv.repo.List(repotypes.GetQuestionSubjectGetListBy{QuestionId: dto.QuestionId})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (srv SubjectService) CreateInBatches(dto []qnrEntity.Subject) error {
	err := srv.repo.CreateInBatches(dto)
	return err
}

func (srv SubjectService) GetList(openid string, questionId int64) (gin.H, error) {
	//查询用户是否入库，入库并回答过问题
	info := srv.repoUser.FindBy(repotypes.GetQuestionUserGetById{OpenId: openid})
	isSubmit := 0
	if info.UserId != 0 {
		isSubmit = 1
	}
	//所有的题目
	subjectList, subjectErr := srv.repo.List(repotypes.GetQuestionSubjectGetListBy{
		QuestionId: questionId, //金融调查问卷
	})
	if subjectErr != nil {
		return gin.H{}, nil
	}
	var subjectIds []model.LongID //获取所有的题目id
	for _, val := range subjectList {
		subjectIds = append(subjectIds, val.SubjectId)
	}

	//所有的答案
	optionList, optionErr := srv.repoOption.GetListBy(repotypes.GetQuestOptionGetListBy{
		SubjectIds: subjectIds,
	})
	if optionErr != nil {
		return gin.H{}, nil
	}

	optionMap := make(map[model.LongID][]api_types.QuestionOptionVO)
	for _, val := range optionList {
		optionMap[val.SubjectId] = append(optionMap[val.SubjectId], api_types.QuestionOptionVO{
			ID:             val.ID,
			Title:          val.Title,
			Remind:         val.Remind,
			JumpSubject:    val.JumpSubject,
			RelatedSubject: val.RelatedSubject,
			Carbon:         val.Carbon,
			SubjectId:      val.SubjectId,
		})
	}

	//答案和题目组装
	subjectMap := make(map[qnrEntity.QuestionCategoryType][]api_types.QuestionVo, 0)
	for _, val := range subjectList {
		option, ok := optionMap[val.SubjectId]
		if !ok {
			option = []api_types.QuestionOptionVO{}
		}
		subjectMap[val.CategoryId] = append(subjectMap[val.CategoryId], api_types.QuestionVo{
			ID:        val.ID,
			Title:     val.Title,
			Type:      val.Type,
			Remind:    val.Remind,
			IsHide:    val.IsHide,
			Option:    option,
			SubjectId: val.SubjectId,
		})
	}
	list := make([]api_types.QuestionListVo, 0)
	for _, v := range qnrEntity.QuestionCategoryTypeMap {
		l, err := subjectMap[v]
		if err {
			list = append(list, api_types.QuestionListVo{Title: v.Text(), List: l, Desc: v.DescText()})
		}
	}
	return gin.H{"subject": list, "isSubmit": isSubmit, "subjectCount": len(list)}, nil
}

func (srv SubjectService) GetUserQuestion(dto srv_types.GetQuestionUserDTO) (srv_types.AddUserCarbonInfoDTO, error) {
	//查询用户是否入库，入库并回答过问题
	info := srv.repoUser.FindBy(repotypes.GetQuestionUserGetById{OpenId: dto.OpenId})
	if info.UserId == 0 {
		return srv_types.AddUserCarbonInfoDTO{}, errno.ErrCommon.WithMessage("未提交年度排放问卷")
	}
	userId := info.UserId
	//总碳量
	carbon := srv.repoAnswer.GetUserCarbon(repotypes.GetQuestionUserCarbon{Uid: userId, QuestionId: dto.QuestionId})

	//用户碳量分类汇总
	carbonClassify := srv.repoAnswer.GetUserAnswer(repotypes.GetQuestionUserCarbon{Uid: userId, QuestionId: dto.QuestionId})

	var userCarbonClassify []srv_types.UserCarbonClassify
	for _, answerStruct := range carbonClassify {
		userCarbonClassify = append(userCarbonClassify, srv_types.UserCarbonClassify{
			CategoryId:   qnrEntity.QuestionCategoryType(answerStruct.CategoryId),
			Carbon:       util.CarbonToRate(answerStruct.Carbon),
			CategoryName: qnrEntity.QuestionCategoryType(answerStruct.CategoryId).Text(),
			CarbonValue:  answerStruct.Carbon,
		})
	}

	//今日碳量
	carbonToday := service.NewCarbonTransactionService(context.NewMioContext()).GetTodayCarbon(dto.UserId) //今日碳量

	//日均排放
	carbonDes := decimal.NewFromFloat(carbon)
	yesDes := decimal.NewFromFloat(365)
	carbonDay := carbonDes.Div(yesDes).Round(2)
	carbonDayFloat, _ := carbonDay.Float64()

	//碳中和完成度
	completion := "0"
	carbonTodayDes := decimal.NewFromFloat(carbonToday)
	if !carbonDay.IsZero() {
		completion = carbonTodayDes.Div(carbonDay).Mul(decimal.NewFromInt(100)).Round(2).String()
	}

	//属于用户群里
	userGroup := ""
	userGroupTips := ""
	level := int8(0)
	switch {
	case carbon > 0 && carbon <= 2000000:
		{
			userGroup = "减排王者"
			userGroupTips = "你在低碳环保这条路上一骑绝尘！"
			level = 0
			break
		}
	case carbon > 2000000 && carbon <= 6000000:
		{
			userGroup = "低碳先锋"
			userGroupTips = "优秀的低碳选手，要继续保持哦！"
			level = 1
			break
		}
	case carbon > 6000000 && carbon <= 12000000:
		{
			userGroup = "减碳新人"
			userGroupTips = "环保不难，从小事做起节能减排！"
			level = 2
			break
		}
	case carbon > 12000000 && carbon <= 24000000:
		{
			userGroup = "减碳小白"
			userGroupTips = "别躺平～美丽地球，没你真不行！"
			level = 3
			break
		}
	default:
		{
			break
		}
	}

	compareWithCountry := "低于"
	compareWithGlobal := "低于"
	if carbon > 6800000 {
		compareWithCountry = "高于"
	}
	if carbon > 4400000 {
		compareWithGlobal = "高于"
	}
	return srv_types.AddUserCarbonInfoDTO{
		UserGroup:          userGroup, //属于用户群里
		UserGroupTips:      userGroupTips,
		CarbonYear:         util.CarbonToRate(carbon),         //总碳量
		CarbonToday:        util.CarbonToRate(carbonToday),    //今日碳量
		CarbonClassify:     userCarbonClassify,                //用户碳量分类汇总
		CarbonDay:          util.CarbonToRate(carbonDayFloat), //日均排放
		CarbonCompletion:   completion + "%",                  //碳中和完成度
		CompareWithCountry: compareWithCountry,
		CompareWithGlobal:  compareWithGlobal,
		Level:              level,
	}, nil
}
