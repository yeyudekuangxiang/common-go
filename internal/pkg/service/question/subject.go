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

func (srv SubjectService) GetUserQuestion(dto srv_types.GetQuestionUserDTO) srv_types.AddUserCarbonInfoDTO {
	//总碳量
	carbon := srv.repoAnswer.GetUserCarbon(repotypes.GetQuestionUserCarbon{Uid: dto.UserId, QuestionId: dto.QuestionId})

	//用户碳量分类汇总
	carbonClassify := srv.repoAnswer.GetUserAnswer(repotypes.GetQuestionUserCarbon{Uid: dto.UserId, QuestionId: dto.QuestionId})

	var userCarbonClassify []srv_types.UserCarbonClassify
	for _, answerStruct := range carbonClassify {
		userCarbonClassify = append(userCarbonClassify, srv_types.UserCarbonClassify{
			CategoryId:   answerStruct.CategoryId,
			Carbon:       util.CarbonToRate(answerStruct.Carbon),
			CategoryName: answerStruct.CategoryId.Text(),
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
	carbonTodayDes := decimal.NewFromFloat(carbonToday)
	completion := carbonTodayDes.Div(carbonDay).Round(2).String()

	//属于用户群里
	personName := ""
	switch {
	case carbon > 1500:
		{
			personName = "低碳环保人群"
		}
	case carbon > 1500 && carbon <= 4000:
		{
			personName = "低碳环保人群1"
		}
	case carbon > 4000 && carbon <= 16000:
		{
			personName = "低碳环保人群2"
		}
	default:
		{
			personName = "低碳环保人群3"
		}
	}
	compareWithCountry := "高于"
	compareWithGlobal := "低于"
	return srv_types.AddUserCarbonInfoDTO{
		PersonName:         personName,                        //属于用户群里
		CarbonYear:         util.CarbonToRate(carbon),         //总碳量
		CarbonToday:        util.CarbonToRate(carbonToday),    //今日碳量
		CarbonClassify:     userCarbonClassify,                //用户碳量分类汇总
		CarbonDay:          util.CarbonToRate(carbonDayFloat), //日均排放
		CarbonCompletion:   completion,                        //碳中和完成度
		CompareWithCountry: compareWithCountry,
		CompareWithGlobal:  compareWithGlobal,
	}
}
