package admin

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service"
	"mio/internal/pkg/util/apiutil"
)

var DefaultTopicController = TopicController{}

type TopicController struct {
}

func (ctr TopicController) List(c *gin.Context) (gin.H, error) {
	form := TopicListRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//get topic by params
	list, total, err := service.DefaultTopicAdminService.GetTopicList(repository.TopicListRequest{
		ID:        form.ID,
		Title:     form.Title,
		UserId:    form.ID,
		UserName:  form.UserName,
		TagId:     form.TagId,
		Status:    form.Status,
		IsTop:     form.IsTop,
		IsEssence: form.IsEssence,
		Offset:    form.Offset(),
		Limit:     form.Limit(),
	})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"list":     list,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

func (ctr TopicController) Detail(c *gin.Context) (gin.H, error) {
	form := TopicDetailRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	topic, err := service.DefaultTopicAdminService.DetailTopic(form.ID)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": topic,
	}, nil
}

func (ctr TopicController) Create(c *gin.Context) (gin.H, error) {
	form := CreateTopicRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//user
	user := apiutil.GetAuthUser(c)
	//更新帖子
	err := service.DefaultTopicAdminService.CreateTopic(user.ID, user.AvatarUrl, user.Nickname, user.OpenId, form.Title, form.Content, form.TagIds, form.Images)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

func (ctr TopicController) Update(c *gin.Context) (gin.H, error) {
	form := UpdateTopicRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	//更新帖子
	err := service.DefaultTopicAdminService.UpdateTopic(user.OpenId, form.ID, form.Title, form.Content, form.TagIds, form.Images)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Review 审核
func (ctr TopicController) Review(c *gin.Context) (gin.H, error) {
	form := ChangeTopicStatus{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.DefaultTopicAdminService.Review(form.ID, form.Status)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Top 置顶
func (ctr TopicController) Top(c *gin.Context) (gin.H, error) {
	form := ChangeTopicStatus{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.DefaultTopicAdminService.Top(form.ID, form.IsTop)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Essence 精华
func (ctr TopicController) Essence(c *gin.Context) (gin.H, error) {
	form := ChangeTopicStatus{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	err := service.DefaultTopicAdminService.Essence(form.ID, form.IsEssence)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
