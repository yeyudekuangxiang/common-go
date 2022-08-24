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

func (SubjectController) Create(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetQnrSubjectListDTO{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
func (SubjectController) GetList(ctx *gin.Context) (gin.H, error) {
	form := api_types.GetQnrSubjectListDTO{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	subjectServer := qnrService.NewSubjectService(context.NewMioContext(context.WithContext(ctx)))

	//所有的题目
	subjectList, _ := subjectServer.GetPageList(srv_types.GetQnrSubjectDTO{
		QnrId: 1,
	})
	var subjectIds []int64 //获取所有的题目id
	for _, val := range subjectList {
		subjectIds = append(subjectIds, val.ID)
	}

	//所有的答案
	optionServer := qnrService.NewOptionService(context.NewMioContext(context.WithContext(ctx)))
	optionList, _ := optionServer.GetPageList(srv_types.GetQnrOptionDTO{
		SubjectIds: subjectIds,
	})
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

	list := make([]api_types.QnrListVo, 0)

	a1, err := subjectMap[1]
	if err {
		list = append(list, api_types.QnrListVo{Title: "一、 个人信息", List: a1})
	}
	a2, err2 := subjectMap[2]
	if err2 {
		list = append(list, api_types.QnrListVo{Title: "二、 绿色金融市场建设", List: a2})
	}
	a3, err3 := subjectMap[3]
	if err3 {
		list = append(list, api_types.QnrListVo{Title: "三、 绿色金融工具", List: a3})
	}
	a4, err4 := subjectMap[4]
	if err4 {
		list = append(list, api_types.QnrListVo{Title: "四、 配套保障与政府支持", List: a4})
	}
	a5, err5 := subjectMap[5]
	if err5 {
		list = append(list, api_types.QnrListVo{Title: "五、 企业活动", List: a5})
	}
	a6, err6 := subjectMap[6]
	if err6 {
		list = append(list, api_types.QnrListVo{Title: "六、 生态空间和城市基建", List: a6})
	}
	a7, err7 := subjectMap[7]
	if err7 {
		list = append(list, api_types.QnrListVo{Title: "总评分", List: a7})
	}
	return gin.H{"subject": list, "isSubmit": 0}, nil
}
