package question

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	qnrEntity "mio/internal/pkg/model/entity/question"
	repoQnr "mio/internal/pkg/repository/question"
	"mio/internal/pkg/repository/repotypes"
	"mio/internal/pkg/service/srv_types"
)

var DefaultSubjectService = SubjectService{ctx: context.NewMioContext()}

func NewSubjectService(ctx *context.MioContext) *SubjectService {
	return &SubjectService{
		ctx:        ctx,
		repo:       repoQnr.NewSubjectRepository(ctx),
		repoOption: repoQnr.NewOptionRepository(ctx),
		repoUser:   repoQnr.NewUserRepository(ctx),
	}
}

type SubjectService struct {
	ctx        *context.MioContext
	repo       *repoQnr.SubjectRepository
	repoOption *repoQnr.OptionRepository
	repoUser   *repoQnr.UserRepository
}

func (srv SubjectService) GetPageList(dto srv_types.GetQnrSubjectDTO) ([]qnrEntity.Subject, error) {
	list, err := srv.repo.List(repotypes.GetQuestionSubjectGetListBy{QuestionId: dto.QnrId})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (srv SubjectService) CreateInBatches(dto []qnrEntity.Subject) error {
	err := srv.repo.CreateInBatches(dto)
	return err
}

func (srv SubjectService) GetList(openid string) (gin.H, error) {
	//查询用户是否入库，入库并回答过问题
	info := srv.repoUser.FindBy(repotypes.GetQuestionUserGetById{OpenId: openid})
	isSubmit := 0
	if info.UserId != 0 {
		isSubmit = 1
	}
	//所有的题目
	subjectList, subjectErr := srv.repo.List(repotypes.GetQuestionSubjectGetListBy{
		QuestionId: 1, //金融调查问卷
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
	subjectMap := make(map[int64][]api_types.QuestionVo, 0)
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
	//题目和分类组装
	typeMap := []api_types.QuestionCategory{
		{Id: 1, Title: "衣", Desc: "1111"},
		{Id: 2, Title: "食", Desc: "222"},
		{Id: 3, Title: "住", Desc: "333"},
		{Id: 4, Title: "用", Desc: "444"},
		{Id: 5, Title: "行", Desc: "555"},
	}
	list := make([]api_types.QuestionListVo, 0)
	for _, v := range typeMap {
		l, err := subjectMap[v.Id]
		if err {
			list = append(list, api_types.QuestionListVo{Title: v.Title, List: l, Desc: v.Desc})
		}
	}
	return gin.H{"subject": list, "isSubmit": isSubmit, "subjectCount": len(list)}, nil
}
