package community

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/common-go/baidu"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/queue/producer/common"
	"mio/internal/pkg/queue/types/message/smsmsg"
	communityModel "mio/internal/pkg/repository/community"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/community"
	"mio/internal/pkg/service/message"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/apiutil"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultTopicController = TopicController{}

type TopicController struct {
}

//List 获取文章列表
func (ctr *TopicController) List(c *gin.Context) (gin.H, error) {
	form := GetTopicPageListRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	list, total, err := community.DefaultTopicService.GetRecommendList(community.TopicListParams{
		Offset: form.Page,
		Limit:  form.PageSize,
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

//func (ctr *TopicController) ListFlow(c *gin.Context) (gin.H, error) {
//	form := GetTopicPageListRequest{}
//	if err := apiutil.BindForm(c, &form); err != nil {
//		return nil, err
//	}
//
//	user := apiutil.GetAuthUser(c)
//
//	list, total, err := community.DefaultTopicService.GetTopicDetailPageListByFlow(repository.GetTopicPageListBy{
//		ID:         form.ID,
//		TopicTagId: form.TopicTagId,
//		Offset:     form.Offset(),
//		Limit:      form.Limit(),
//		UserId:     user.ID,
//	})
//	if err != nil {
//		return nil, err
//	}
//	ids := make([]int64, 0)
//	for _, item := range list {
//		ids = append(ids, item.Id)
//	}
//	app.Logger.Infof("user:%d form:%+v ids:%+v", user.ID, form, ids)
//
//	return gin.H{
//		"list":     list,
//		"total":    total,
//		"page":     form.Page,
//		"pageSize": form.PageSize,
//	}, nil
//}

//GetShareWeappQrCode 获取分享二维码
func (ctr *TopicController) GetShareWeappQrCode(c *gin.Context) (gin.H, error) {
	form := GetWeappQrCodeRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)

	page := "pages/community/details/index"
	scene := fmt.Sprintf("tid=%d&uid=%d&s=p&m=c", form.TopicId, user.ID)

	qr, err := service.NewQRCodeService().GetUnlimitedQRCodeRaw(page, scene, 100)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"qrcode": qr,
	}, nil
}

//ChangeTopicLike 点赞 / 取消点赞
func (ctr *TopicController) ChangeTopicLike(c *gin.Context) (gin.H, error) {
	form := ChangeTopicLikeRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	topicLikeService := community.NewTopicLikeService(ctx)
	messageService := message.NewWebMessageService(ctx)

	resp, err := topicLikeService.ChangeLikeStatus(form.TopicId, user.ID, user.OpenId)
	if err != nil {
		return nil, err
	}

	title := resp.TopicTitle
	if len([]rune(title)) > 8 {
		title = string([]rune(title)[0:8]) + "..."
	}

	var point int64
	if resp.LikeStatus == 1 && resp.IsFirst == true {
		pointService := service.NewPointService(ctx)
		_, err := pointService.IncUserPoint(srv_types.IncUserPointDTO{
			OpenId:       user.OpenId,
			Type:         entity.POINT_LIKE,
			BizId:        util.UUID(),
			ChangePoint:  int64(entity.PointCollectValueMap[entity.POINT_LIKE]),
			AdminId:      0,
			Note:         "为文章 \"" + title + "\" 点赞",
			AdditionInfo: strconv.FormatInt(resp.TopicId, 10),
		})

		if err == nil {
			point = int64(entity.PointCollectValueMap[entity.POINT_LIKE])
		}

		//发送消息
		err = messageService.SendMessage(message.SendWebMessage{
			SendId:       user.ID,
			RecId:        resp.TopicUserId,
			Key:          "like_topic",
			Type:         message.MsgTypeLike,
			TurnType:     message.MsgTurnTypeArticle,
			TurnId:       resp.TopicId,
			MessageNotes: title,
		})
		if err != nil {
			app.Logger.Errorf("文章点赞站内信发送失败:%s", err.Error())
		}

	}

	return gin.H{
		"point":  point,
		"status": resp.LikeStatus,
	}, nil
}

//ListTopic 帖子列表+顶级评论+顶级评论下子评论3条
func (ctr *TopicController) ListTopic(c *gin.Context) (gin.H, error) {
	form := GetTopicPageListRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))

	params := community.TopicListParams{
		TopicTagId: form.TopicTagId,
		Offset:     form.Offset(),
		Limit:      form.Limit(),
		Label:      form.Order,
	}

	if params.Label == "recommend" && params.TopicTagId == 0 {
		params.Limit = form.PageSize
		params.Offset = form.Page
	}

	list, total, err := community.DefaultTopicService.GetTopicList(params)
	if err != nil {
		return nil, err
	}
	resList := make([]*entity.Topic, 0)
	//点赞数据
	likeMap := make(map[int64]struct{}, 0)
	topicLikeService := community.NewTopicLikeService(ctx)
	likeList, _ := topicLikeService.GetLikeInfoByUser(user.ID)
	if len(likeList) > 0 {
		for _, item := range likeList {
			likeMap[item.TopicId] = struct{}{}
		}
	}
	//收藏数据
	collectionMap := make(map[int64]struct{}, 0)
	collectionService := community.NewCollectionService(ctx)
	collectionIds := collectionService.Collections(user.OpenId, 0, 0, 0)
	for _, collectionId := range collectionIds {
		collectionMap[collectionId] = struct{}{}
	}

	//评论数据
	ids := make([]int64, 0) //topicId
	for _, item := range list {
		ids = append(ids, item.Id)
	}

	rootCommentCount := community.DefaultTopicService.GetRootCommentCount(ids)
	//组装数据---帖子的顶级评论数量
	topic2comment := make(map[int64]int64, 0)
	for _, item := range rootCommentCount {
		topic2comment[item.TopicId] = item.Total
	}
	for _, item := range list {
		//res := item.TopicItemRes()
		item.CommentCount = topic2comment[item.Id]
		if _, ok := likeMap[item.Id]; ok {
			item.IsLike = 1
		}
		if _, ok := collectionMap[item.Id]; ok {
			item.IsCollection = 1
		}
		if item.Type == 1 {
			item.Activity.Status = 1
			if item.Activity.SignupDeadline.Before(time.Now()) {
				item.Activity.Status = 2
			}
		}
		resList = append(resList, item)
	}
	app.Logger.Infof("GetTopicDetailPageListByFlow user:%d form:%+v ids:%+v", user.ID, form, ids)
	return gin.H{
		"list":     resList,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, nil
}

//CreateTopic 创建帖子
func (ctr *TopicController) CreateTopic(c *gin.Context) (gin.H, error) {
	user := apiutil.GetAuthUser(c)
	if user.Auth != 1 {
		return nil, errno.ErrCommon.WithMessage("无权限")
	}
	form := CreateTopicRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	if form.Type == 0 && len(form.TagIds) > 2 {
		return nil, errno.ErrCommon.WithMessage("话题数量最多选2个哦")
	}

	//审核
	//title审核
	err := validator.CheckMsgWithOpenId(user.OpenId, form.Title)
	if err != nil {
		return nil, errno.ErrCommon.WithMessage("标题审核未通过")
	}

	// 文本内容审核
	if form.Content != "" {
		if err := validator.CheckMsgWithOpenId(user.OpenId, form.Content); err != nil {
			app.Logger.Error(fmt.Errorf("create Topic error:%s", err.Error()))
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "发帖-文本内容审核"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, user.OpenId, zhuGeAttr)
			return nil, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	// 图片内容审核
	if len(form.Images) >= 1 {
		//reviewSrv := service.DefaultReviewService()
		for i, imgUrl := range form.Images {
			if err := validator.CheckMediaWithOpenId(user.OpenId, imgUrl); err != nil {
				//if err := reviewSrv.ImageReview(baidu.ImageReviewParam{ImgUrl: imgUrl}); err != nil {
				app.Logger.Error(fmt.Errorf("create Topic error:%s", err.Error()))
				zhuGeAttr := make(map[string]interface{}, 0)
				zhuGeAttr["场景"] = "发帖-图片内容审核"
				zhuGeAttr["失败原因"] = err.Error()
				track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, user.OpenId, zhuGeAttr)
				return nil, errno.ErrCommon.WithMessage("图片: " + strconv.Itoa(i) + " " + err.Error())
			}
		}
	}

	//创建帖子
	marshal, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	var params community.CreateTopicParams
	err = json.Unmarshal(marshal, &params)
	if err != nil {
		return nil, err
	}

	topic, err := community.DefaultTopicService.CreateTopic(user.ID, params)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": topic,
		"point": 0,
	}, nil
}

func (ctr *TopicController) UpdateTopic(c *gin.Context) (gin.H, error) {
	//user
	user := apiutil.GetAuthUser(c)
	if user.Auth != 1 {
		return nil, errno.ErrCommon.WithMessage("无权限")
	}

	form := UpdateTopicRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	if form.Type == 0 && len(form.TagIds) > 2 {
		return nil, errno.ErrCommon.WithMessage("话题数量最多选2个哦")
	}

	//审核
	if form.Content != "" {
		//检查内容
		if err := validator.CheckMsgWithOpenId(user.OpenId, form.Content); err != nil {
			app.Logger.Error(fmt.Errorf("update Topic error:%s", err.Error()))
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "更新帖子"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, user.OpenId, zhuGeAttr)
			return nil, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	if len(form.Images) > 1 {
		reviewSrv := service.DefaultReviewService()
		for i, imgUrl := range form.Images {
			if err := reviewSrv.ReviewImage(baidu.ImageReviewParam{ImgUrl: imgUrl}); err != nil {
				app.Logger.Error(fmt.Errorf("create Topic error:%s", err.Error()))
				zhuGeAttr := make(map[string]interface{}, 0)
				zhuGeAttr["场景"] = "发帖-图片内容审核"
				zhuGeAttr["失败原因"] = err.Error()
				track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, user.OpenId, zhuGeAttr)
				return nil, errno.ErrCommon.WithMessage("图片: " + strconv.Itoa(i) + " " + err.Error())
			}
		}
	}

	//更新帖子
	marshal, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	var params community.UpdateTopicParams
	err = json.Unmarshal(marshal, &params)
	if err != nil {
		return nil, err
	}
	topic, err := community.DefaultTopicService.UpdateTopic(user.ID, params)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"topic": topic,
	}, nil
}

func (ctr *TopicController) DelTopic(c *gin.Context) (gin.H, error) {
	form := IdRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	//user
	user := apiutil.GetAuthUser(c)
	if user.Auth != 1 {
		return nil, errno.ErrCommon.WithMessage("无权限")
	}
	ctx := context.NewMioContext()
	ActivitiesSignupService := community.NewCommunityActivitiesSignupService(ctx)
	TopicService := community.NewTopicService(ctx)

	//更新帖子
	topic, err := TopicService.DelTopic(user.ID, form.ID)
	if err != nil {
		return nil, err
	}
	// 报名活动删除的话 通知报名者
	if topic.Type == 1 {
		signupList, count, err := ActivitiesSignupService.FindAll(communityModel.FindAllActivitiesSignupParams{
			TopicId: topic.Id,
		})
		if err != nil {
			app.Logger.Errorf("【短信发送失败】取消报名活动通知短信发送失败: %s", err.Error())
		}
		if count == 0 {
			return nil, nil
		}
		for _, item := range signupList {
			//发送短信
			err := common.SendSms(smsmsg.SmsMessage{
				Phone:       item.Phone,
				Args:        topic.Title,
				TemplateKey: message.SmsActivityCancel,
			})
			if err != nil {
				app.Logger.Error("【短信发送失败】取消报名活动通知短信发送失败: %s", err.Error())
				break
			}
		}
	}

	return nil, nil
}

func (ctr *TopicController) DetailTopic(c *gin.Context) (gin.H, error) {
	form := IdRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	topicService := community.NewTopicService(ctx)
	topicLikeService := community.NewTopicLikeService(ctx)
	collectService := community.NewCollectionService(ctx)
	signupService := community.NewCommunityActivitiesSignupService(ctx)

	//获取帖子
	topic, err := topicService.DetailTopic(community.FindTopicParams{
		TopicId: form.ID,
		UserId:  user.ID,
	})
	if err != nil {
		return nil, err
	}
	//获取评论数量

	CommentCount := topicService.GetCommentCount([]int64{topic.Id})
	// 组装数据
	// 评论
	if len(CommentCount) > 0 {
		topic.CommentCount = CommentCount[0].Total
	}

	// 点赞
	like, err := topicLikeService.GetOneByTopic(topic.Id, user.ID)
	if err == nil {
		topic.IsLike = int(like.Status)
	}
	// 收藏
	collection, err := collectService.FindOneByTopic(topic.Id, user.OpenId)

	if err == nil {
		topic.IsCollection = collection.Status
	}
	if topic.Type == communityModel.TopicTypeActivity {
		info, err := signupService.GetSignupInfo(communityModel.FindOneActivitiesSignupParams{
			TopicId:      topic.Id,
			UserId:       user.ID,
			SignupStatus: communityModel.SignupStatusTrue,
		})
		if err == nil {
			topic.Activity.SignupStatus = info.SignupStatus
		}
		topic.Activity.Status = 1
		if topic.Activity.SignupDeadline.Before(time.Now()) {
			topic.Activity.Status = 2
		}
	}

	return gin.H{
		"topic": topic,
	}, nil
}

func (ctr *TopicController) MyTopic(c *gin.Context) (gin.H, error) {
	form := MyTopicRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	user := apiutil.GetAuthUser(c)
	status := form.Status
	if form.UserId != 0 {
		u, b, err := service.DefaultUserService.GetUserByID(form.UserId)
		if err != nil {
			return nil, errno.ErrCommon
		}
		if !b {
			return nil, errno.ErrUserNotFound
		}
		user = *u
		status = 3
	}

	list, total, err := community.DefaultTopicService.GetMyTopicList(community.MyTopicListParams{
		UserId: user.ID,
		Status: status,
		Type:   form.Type,
		Limit:  form.Limit(),
		Offset: form.Offset(),
	})

	if err != nil {
		return nil, err
	}

	resList := make([]*entity.Topic, 0)

	//点赞数据
	likeMap := make(map[int64]struct{}, 0)
	topicLikeService := community.NewTopicLikeService(ctx)
	likeList, _ := topicLikeService.GetLikeInfoByUser(user.ID)
	if len(likeList) > 0 {
		for _, item := range likeList {
			likeMap[item.TopicId] = struct{}{}
		}
	}

	//评论数据
	ids := make([]int64, 0) //topicId
	for _, item := range list {
		ids = append(ids, item.Id)
	}

	//收藏数据
	collectionMap := make(map[int64]struct{}, 0)
	collectionService := community.NewCollectionService(ctx)
	collectionIds := collectionService.Collections(user.OpenId, 0, 0, 0)
	for _, collectionId := range collectionIds {
		collectionMap[collectionId] = struct{}{}
	}

	rootCommentCount := community.DefaultTopicService.GetRootCommentCount(ids)
	// 组装数据---帖子的顶级评论数量
	topic2comment := make(map[int64]int64, 0)
	for _, item := range rootCommentCount {
		topic2comment[item.TopicId] = item.Total
	}
	//组装数据---点赞数据 收藏数据
	for _, item := range list {
		item.CommentCount = topic2comment[item.Id]
		if _, ok := likeMap[item.Id]; ok {
			item.IsLike = 1
		}
		if _, ok := collectionMap[item.Id]; ok {
			item.IsCollection = 1
		}
		item.Activity.Status = 1
		if item.Activity.SignupDeadline.Before(time.Now()) {
			item.Activity.Status = 2
		}
		resList = append(resList, item)
	}

	return gin.H{
		"list":     resList,
		"total":    total,
		"page":     form.Page,
		"pageSize": form.PageSize,
	}, err
}

func (ctr *TopicController) SignupTopic(c *gin.Context) (gin.H, error) {
	form := SignupTopicRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	signupService := community.NewCommunityActivitiesSignupService(ctx)
	params := community.SignupParams{
		TopicId:      form.TopicId,
		UserId:       user.ID,
		OpenId:       user.OpenId,
		RealName:     form.RealName,
		Phone:        form.Phone,
		Gender:       form.Gender,
		Age:          form.Age,
		Wechat:       form.Wechat,
		City:         form.City,
		Remarks:      form.Remarks,
		SignupTime:   time.Now(),
		SignupStatus: 1,
	}
	err := signupService.Signup(params)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (ctr *TopicController) CancelSignupTopic(c *gin.Context) (gin.H, error) {
	form := IdRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}
	user := apiutil.GetAuthUser(c)
	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	signupService := community.NewCommunityActivitiesSignupService(ctx)
	err := signupService.CancelSignup(form.ID, user.ID)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (ctr *TopicController) MySignup(c *gin.Context) (gin.H, error) {
	form := MySignupRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	user := apiutil.GetAuthUser(c)

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	signupService := community.NewCommunityActivitiesSignupService(ctx)
	list, total, err := signupService.GetPageList(communityModel.FindAllActivitiesSignupParams{
		UserId: user.ID,
		Offset: form.Offset(),
		Limit:  form.Limit(),
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

func (ctr *TopicController) MySignupDetail(c *gin.Context) (gin.H, error) {
	form := IdRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	//user := apiutil.GetAuthUser(c)

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	signupService := community.NewCommunityActivitiesSignupService(ctx)
	signInfo, err := signupService.GetSignupInfo(communityModel.FindOneActivitiesSignupParams{Id: form.ID})
	if err != nil {
		return nil, err
	}
	return gin.H{
		"data": signInfo,
	}, nil
}

//报名数据
func (ctr *TopicController) SignupList(c *gin.Context) (gin.H, error) {
	form := IdRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		return nil, err
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	user := apiutil.GetAuthUser(c)
	signupService := community.NewCommunityActivitiesSignupService(ctx)
	topicService := community.NewTopicService(ctx)
	topic, err := topicService.FindTopic(community.FindTopicParams{
		TopicId: form.ID,
		UserId:  user.ID,
		Type:    communityModel.TopicTypeActivity,
	})
	//仅发起人可查看
	if err != nil {
		if err == entity.ErrNotFount {
			return nil, errno.ErrRecordNotFound
		}
		return nil, errno.ErrCommon
	}

	if topic.UserId != user.ID {
		return nil, nil
	}

	signupList, total, err := signupService.FindSignupList(communityModel.FindAllActivitiesSignupParams{
		TopicId: topic.Id,
	})

	if err != nil {
		return nil, err
	}

	return gin.H{
		"seeCount":    topic.SeeCount,
		"signupCount": total,
		"signupList":  signupList,
	}, nil
}

//导出报名数据excel文件路径
func (ctr *TopicController) ExportSignupList(c *gin.Context) {
	form := IdRequest{}
	if err := apiutil.BindForm(c, &form); err != nil {
		app.Logger.Errorf("参数错误")
	}

	ctx := context.NewMioContext(context.WithContext(c.Request.Context()))
	user := apiutil.GetAuthUser(c)
	signupService := community.NewCommunityActivitiesSignupService(ctx)
	topicService := community.NewTopicService(ctx)
	topic, err := topicService.FindTopic(community.FindTopicParams{
		TopicId: form.ID,
		UserId:  user.ID,
		Type:    communityModel.TopicTypeActivity,
	})
	//仅发起人可查看
	if err != nil {
		app.Logger.Errorf(err.Error())
	}

	if topic.UserId != user.ID {
		app.Logger.Errorf("非创建者本人查看")
	}

	signupService.Export(c.Writer, c.Request, topic.Id)
}
