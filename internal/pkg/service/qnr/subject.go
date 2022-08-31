package question

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	qnrEntity "mio/internal/pkg/model/entity/qnr"
	repoQnr "mio/internal/pkg/repository/qnr"
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
	list, err := srv.repo.List(repotypes.GetQuestSubjectGetListBy{QnrId: dto.QnrId})
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
	info := srv.repoUser.FindBy(repotypes.GetQuestUserGetById{OpenId: openid})
	isSubmit := 0
	if info.UserId != 0 {
		isSubmit = 1
	}

	//所有的题目
	subjectList, subjectErr := srv.repo.List(repotypes.GetQuestSubjectGetListBy{
		QnrId: 1, //金融调查问卷
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

	optionMap := make(map[model.LongID][]api_types.OptionVO)
	for _, val := range optionList {
		optionMap[val.SubjectId] = append(optionMap[val.SubjectId], api_types.OptionVO{
			ID:             val.ID,
			Title:          val.Title,
			Remind:         val.Remind,
			JumpSubject:    val.JumpSubject,
			RelatedSubject: val.RelatedSubject,
		})
	}

	//答案和题目组装
	subjectMap := make(map[int64][]api_types.QnrVo, 0)
	for _, val := range subjectList {
		option, ok := optionMap[val.SubjectId]
		if !ok {
			option = []api_types.OptionVO{}
		}
		subjectMap[val.CategoryId] = append(subjectMap[val.CategoryId], api_types.QnrVo{
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
	typeMap := []api_types.QnrCategory{
		{Id: 1, Title: "一、 个人信息"},
		{Id: 2, Title: "二、 绿色金融市场建设"},
		{Id: 3, Title: "三、 绿色金融工具"},
		{Id: 4, Title: "四、 配套保障与政府支持"},
		{Id: 5, Title: "五、 企业活动"},
		{Id: 6, Title: "六、 生态空间和城市基建"},
		{Id: 7, Title: "七、 总评分"},
	}
	list := make([]api_types.QnrListVo, 0)
	for _, v := range typeMap {
		l, err := subjectMap[v.Id]
		if err {
			list = append(list, api_types.QnrListVo{Title: v.Title, List: l})
		}
	}
	return gin.H{"subject": list, "isSubmit": isSubmit, "subjectCount": 31}, nil
}
