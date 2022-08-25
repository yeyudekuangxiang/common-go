package qnr

import (
	"github.com/gin-gonic/gin"
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/core/context"
	qnrService "mio/internal/pkg/service/qnr"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util/apiutil"
)

var DefaultSubjectController = SubjectController{}

type SubjectController struct {
}

type Ans struct {
	Id     int64  `json:"id"`
	Answer string `json:"answer"`
}

func (SubjectController) Create(ctx *gin.Context) (gin.H, error) {
	answerServer := qnrService.NewAnswerService(context.NewMioContext(context.WithContext(ctx)))
	form := api_types.GetQnrSubjectCreateDTO{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(ctx)
	err := answerServer.Add(srv_types.AddQnrAnswerDTO{
		OpenId: user.OpenId,
		UserId: user.ID,
		Answer: form.Answer})
	if err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

func (SubjectController) GetList(ctx *gin.Context) (gin.H, error) {
	subjectServer := qnrService.NewSubjectService(context.NewMioContext(context.WithContext(ctx)))
	optionServer := qnrService.NewOptionService(context.NewMioContext(context.WithContext(ctx)))

	//所有的题目
	subjectList, subjectErr := subjectServer.GetPageList(srv_types.GetQnrSubjectDTO{
		QnrId: 1, //金融调查问卷
	})
	if subjectErr != nil {
		return gin.H{}, nil
	}
	var subjectIds []int64 //获取所有的题目id
	for _, val := range subjectList {
		subjectIds = append(subjectIds, val.ID)
	}

	//所有的答案
	optionList, optionErr := optionServer.GetPageList(srv_types.GetQnrOptionDTO{
		SubjectIds: subjectIds,
	})
	if optionErr != nil {
		return gin.H{}, nil
	}

	optionMap := make(map[int64][]api_types.OptionVO)
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
		option, ok := optionMap[val.ID]
		if !ok {
			option = []api_types.OptionVO{}
		}
		subjectMap[val.CategoryId] = append(subjectMap[val.CategoryId], api_types.QnrVo{
			ID:     val.ID,
			Title:  val.Title,
			Type:   val.Type,
			Remind: val.Remind,
			Option: option,
		})
	}

	//题目和分类组装
	list := make([]api_types.QnrListVo, 0)
	BannerTypeMap := map[int64]string{
		1: "一、 个人信息",
		2: "二、 绿色金融市场建设",
		3: "三、 绿色金融工具",
		4: "四、 配套保障与政府支持",
		5: "五、 企业活动",
		6: "六、 生态空间和城市基建",
		7: "总评分",
	}
	for k, v := range BannerTypeMap {
		l, err := subjectMap[k]
		if err {
			list = append(list, api_types.QnrListVo{Title: v, List: l})
		}
	}
	return gin.H{"subject": list, "isSubmit": 0}, nil
}
