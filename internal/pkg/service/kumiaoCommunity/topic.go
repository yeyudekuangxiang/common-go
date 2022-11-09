package kumiaoCommunity

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"io/ioutil"
	"math"
	"math/rand"
	"mio/config"
	"mio/internal/pkg/core/app"
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/oss"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"mio/pkg/wxoa"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var DefaultTopicService = NewTopicService(mioContext.NewMioContext())

func NewTopicService(ctx *mioContext.MioContext) TopicService {
	return TopicService{
		ctx:            ctx,
		topicModel:     repository.NewTopicModel(ctx),
		topicLikeModel: repository.NewTopicLikeRepository(ctx),
		tagModel:       repository.NewTagModel(ctx),
		userModel:      repository.NewUserRepository(),
	}
}

type TopicService struct {
	ctx            *mioContext.MioContext
	topicModel     repository.TopicModel
	topicLikeModel repository.TopicLikeModel
	tagModel       repository.TagModel
	tokenServer    *wxoa.AccessTokenServer
	userModel      repository.UserRepository
}

//将 entity.Topic 列表填充为 TopicDetail 列表
func (srv TopicService) fillTopicList(topicList []entity.Topic, userId int64) ([]TopicDetail, error) {
	//查询点赞信息
	topicIds := make([]int64, 0)
	for _, topic := range topicList {
		topicIds = append(topicIds, topic.Id)
	}
	topicLikeMap := make(map[int64]bool)
	if userId > 0 {
		likeList := srv.topicLikeModel.GetListBy(repository.GetTopicLikeListBy{
			TopicIds: topicIds,
			UserId:   userId,
		})
		for _, like := range likeList {
			topicLikeMap[int64(like.TopicId)] = like.Status == 1
		}
	}

	//整理数据
	detailList := make([]TopicDetail, 0)
	uidsMap := make(map[int64]struct{}, len(topicList))
	uids := make([]int64, 0)
	uinfoMap := make(map[int64]entity.ShortUser, 0)
	for _, topic := range topicList {
		if _, ok := uidsMap[topic.UserId]; ok {
			continue
		}
		uidsMap[topic.UserId] = struct{}{}
		uids = append(uids, topic.UserId)
	}
	uList := srv.userModel.GetUserListBy(repository.GetUserListBy{
		UserIds: uids,
	})

	for _, user := range uList {
		if _, ok := uidsMap[user.ID]; ok {
			uinfoMap[user.ID] = user.ShortUser()
		}
	}

	for _, topic := range topicList {
		detailList = append(detailList, TopicDetail{
			Topic:         topic,
			IsLike:        topicLikeMap[topic.Id],
			UpdatedAtDate: topic.UpdatedAt.Format("01-02"),
			User:          uinfoMap[topic.UserId],
		})
	}

	return detailList, nil
}

// GetTopicDetailPageList 通过topic表直接查询获取内容列表
func (srv TopicService) GetTopicDetailPageList(param repository.GetTopicPageListBy) ([]*entity.Topic, int64, error) {
	//list, total := srv.topicModel.GetTopicPageList(param)
	srv.zAddTopic()

	total := app.Redis.ZCard(srv.ctx.Context, "topic:rank").Val()

	ids, err := app.Redis.ZRevRange(srv.ctx.Context, "topic:rank", int64(param.Offset), int64(param.Limit)).Result()

	if err != nil {
		app.Logger.Errorf("Topic 获取topicId错误:%s", err.Error())
	}

	list, err := srv.topicModel.GetTopicListV2(repository.GetTopicPageListBy{
		Rids: ids,
	})

	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (srv TopicService) zAddTopic() {
	n := app.Redis.Exists(srv.ctx.Context, "topic:rank").Val()
	if n != 0 {
		return
	}
	//从redis获取topicIds一小时更新一次
	var results []entity.Topic

	app.DB.Model(&entity.Topic{}).
		Where("status = ?", 3).
		Where("is_top = ?", 0).
		Where("is_essence = ?", 0).
		FindInBatches(&results, 1000, func(tx *gorm.DB, batch int) error {
			var members []redis.Z
			for _, topic := range results {
				hot := util.NewHot()
				score := hot.Hot(int64(topic.SeeCount), topic.LikeCount, topic.CommentCount, int64(topic.IsEssence),
					topic.CreatedAt.Time)
				members = append(members, redis.Z{
					Score:  score,
					Member: topic.Id,
				})
			}
			app.Redis.ZAddArgs(srv.ctx.Context, "topic:rank", redis.ZAddArgs{
				NX:      false,
				XX:      false,
				LT:      false,
				GT:      false,
				Ch:      true,
				Members: members,
			})
			return nil
		})
	app.Redis.Expire(srv.ctx.Context, "topic:rank", time.Hour*24)
}

// GetTopicList 分页获取帖子，且分页获取顶级评论，且获取顶级评论下3条子评论。
func (srv TopicService) GetTopicList(param repository.GetTopicPageListBy) ([]*entity.Topic, int64, error) {
	list, i, err := srv.topicModel.GetTopicList(param)
	if err != nil {
		return nil, 0, err
	}
	return list, i, nil
}

func (srv TopicService) GetMyTopicList(param repository.GetTopicPageListBy) ([]*entity.Topic, int64, error) {
	topic, i, err := srv.topicModel.GetMyTopic(param)
	if err != nil {
		return nil, 0, err
	}
	return topic, i, nil
}

// GetTopicDetailPageListByFlow 通过topic_flow内容流表获取内容列表 当topic_flow数据不存在时 会后台任务进行初始化并且调用 GetTopicDetailPageList 方法返回数据
//func (srv TopicService) GetTopicDetailPageListByFlow(param repository.GetTopicPageListBy) ([]TopicDetail, int64, error) {
//
//	topicList, total := srv.topicModel.GetFlowPageList(repository.GetTopicFlowPageListBy{
//		Offset:     param.Offset,
//		Limit:      param.Limit,
//		UserId:     param.UserId,
//		TopicId:    param.ID,
//		TopicTagId: param.TopicTagId,
//		Status:     entity.TopicStatusPublished,
//	})
//	if total == 0 {
//		DefaultTopicFlowService.InitUserFlowByMq(param.UserId)
//		return srv.GetTopicDetailPageList(param)
//	}
//
//	//更新曝光和查看次数
//	srv.UpdateTopicFlowListShowCount(topicList, param.UserId)
//
//	if param.ID != 0 && len(topicList) > 0 {
//		app.Logger.Info("更新查看次数", param.UserId, topicList[0].Id)
//		srv.UpdateTopicSeeCount(topicList[0].Id, param.UserId)
//	}
//
//	topicDetailList, err := srv.fillTopicList(topicList, param.UserId)
//
//	if err != nil {
//		return nil, 0, err
//	}
//	return topicDetailList, total, nil
//}

// UpdateTopicSeeCount 更新内容的查看次数加1
func (srv TopicService) UpdateTopicSeeCount(topicId int64, userId int64) {
	err := initUserFlowPool.Submit(func() {
		topic := srv.topicModel.FindById(topicId)
		if topic.Id == 0 {
			return
		}

		seeCount := topic.SeeCount + 1

		if err := srv.topicModel.UpdateColumn(topic.Id, "see_count", seeCount); err != nil {
			app.Logger.Error("更新topic查看次数失败", topicId, userId)
			return
		}
		DefaultTopicFlowService.AddUserFlowSeeCount(userId, topicId)
		DefaultTopicFlowService.AfterUpdateTopic(topicId)
	})
	if err != nil {
		app.Logger.Error("提交更新topic查看次数任务失败", userId, topicId, err)
	}
}

// UpdateTopicFlowListShowCount 更新内容流的曝光次数加1
func (srv TopicService) UpdateTopicFlowListShowCount(list []entity.Topic, userId int64) {
	err := initUserFlowPool.Submit(func() {
		for _, topic := range list {
			DefaultTopicFlowService.AddUserFlowShowCount(userId, topic.Id)
		}
	})
	if err != nil {
		app.Logger.Error("提交更新topic曝光次数任务失败", userId, err)
	}
}

//根据id列表对 entity.Topic 列表排序
func (srv TopicService) sortTopicListByIds(list []entity.Topic, ids []int64) []entity.Topic {
	topicMap := make(map[int64]entity.Topic)
	for _, topic := range list {
		topicMap[topic.Id] = topic
	}

	newList := make([]entity.Topic, 0)
	for _, id := range ids {
		newList = append(newList, topicMap[id])
	}
	return newList
}

// FindById 根据id查询 entity.Topic
func (srv TopicService) FindById(topicId int64) *entity.Topic {
	return srv.topicModel.FindById(topicId)
}

// UpdateTopicSort 更新内容的排序权重
func (srv TopicService) UpdateTopicSort(topicId int64, sort int) error {
	topic := srv.topicModel.FindById(topicId)
	if topic.Id == 0 {
		return errno.ErrCommon.WithMessage("未查询到此内容")
	}
	err := srv.topicModel.UpdateColumn(topicId, "sort", sort)
	if err != nil {
		return err
	}
	DefaultTopicFlowService.AfterUpdateTopicByMq(topicId)
	return nil
}

//func (srv TopicService) ImportUser(filename string) error {
//
//	file, err := excelize.OpenFile(filename)
//	if err != nil {
//		return errors.WithStack(err)
//	}
//	defer file.Close()
//
//	if file.SheetCount == 0 {
//		return errno.ErrCommon.WithMessage("没有数据")
//	}
//
//	rows, err := file.GetRows(file.GetSheetList()[0])
//	if err != nil {
//		return errors.WithStack(err)
//	}
//
//	for i, row := range rows {
//		if i == 0 {
//			continue
//		}
//		importId := row[0]
//		nickname := row[1]
//		//avatarImage := row[2]
//		wechat := row[3]
//		phone := row[4]
//		user := entity.User{}
//		if nickname == "星星充电" {
//			app.DB.Where("openid = ?", wechat).First(&user)
//		} else {
//			app.DB.Where("nick_name = ?", nickname).First(&user)
//		}
//
//		if user.ID != 0 {
//			if user.PhoneNumber != phone {
//				return errors.Errorf("存在同名但手机号不同用户 %s", nickname)
//			}
//			log.Println("用户已存在", user)
//			continue
//		}
//
//		avatar, err := srv.uploadImportUserAvatar(path.Join(strings.Split(filename, "_")[0], importId+".jpg"))
//		if err != nil {
//			return errors.WithMessage(err, "上传头像失败"+importId)
//		}
//		_, err = service.DefaultUserService.CreateUser(service.CreateUserParam{
//			OpenId:      wechat,
//			AvatarUrl:   avatar,
//			Nickname:    nickname,
//			PhoneNumber: phone,
//			Source:      entity.UserSourceMio,
//		})
//		if err != nil {
//			return errors.Errorf("创建用户失败 %s %v", nickname, err)
//		}
//	}
//	return nil
//}

func (srv TopicService) uploadImportUserAvatar(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	_, fileName := path.Split(filepath)

	name := fmt.Sprintf("images/topic/%s/%s", path.Base(path.Dir(filepath)), fileName)

	fmt.Println("上传头像", filepath, name)

	defer file.Close()

	avatarPath, err := oss.DefaultOssService.PutObject(name, file)
	if err != nil {
		return "", err
	}
	return oss.DefaultOssService.FullUrl(avatarPath), nil
}

func (srv TopicService) UploadImportTopicImage(dirPath string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	images := make([]string, 0)
	for _, fileInfo := range fileInfos {

		if strings.Contains(fileInfo.Name(), "DS_Store") {
			continue
		}

		file, err := os.Open(path.Join(dirPath, fileInfo.Name()))
		if err != nil {
			return nil, err
		}

		name := fmt.Sprintf("images/topic/%s/%s/%s", path.Base(path.Dir(dirPath)), path.Base(dirPath), fileInfo.Name())

		fmt.Println("上传内容图片", path.Join(dirPath, fileInfo.Name()), name)
		topicPath, err := oss.DefaultOssService.PutObject(name, file)
		if err != nil {
			file.Close()
			return nil, err
		}
		images = append(images, oss.DefaultOssService.FullUrl(topicPath))
	}
	return images, nil
}

// ImportTopic 从xlsx中导入内容
func (srv TopicService) ImportTopic(filename string, baseImportId int) error {
	file, err := excelize.OpenFile(filename)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if file.SheetCount == 0 {
		return errno.ErrCommon.WithMessage("没有数据")
	}

	rows, err := file.GetRows(file.GetSheetList()[0])
	if err != nil {
		return errors.WithStack(err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		importId, err := strconv.Atoi(row[0])
		if err != nil {
			return errors.WithStack(err)
		}

		nickname := row[1]
		title := row[2]
		content := strings.Trim(row[3], " ")
		tag1Text := strings.Trim(row[6], " ")
		tag2Text := ""
		if len(row) >= 8 {
			tag2Text = strings.Trim(row[7], " ")
		}

		users := make([]entity.User, 0)

		if nickname == "咚的一声不响" {
			err = app.DB.Where("phone_number = ?", "19988433595").Find(&users).Error
			if err != nil {
				return errors.WithStack(err)
			}
		} else if nickname == "星星充电" {
			err = app.DB.Where("openid = ?", "XXCD_2022").Find(&users).Error
			if err != nil {
				return errors.WithStack(err)
			}
		} else {
			if nickname == "阿斌啊" {
				nickname = "阿斌啊哦"
			}
			err = app.DB.Where("nick_name = ?", nickname).Find(&users).Error
			if err != nil {
				return errors.WithStack(err)
			}
		}

		if len(users) == 0 {
			return errno.ErrCommon.WithMessage("未查询到用户`" + nickname + "`,请先导入用户")
		}
		if len(users) > 1 {
			return errno.ErrCommon.WithMessage("检测到有多个昵称为`" + nickname + "`的用户,请手动处理后再导入")
		}

		topicTag := ""
		topicTagId := ""
		var tag1, tag2 *entity.Tag
		if tag1Text != "" {
			err = app.DB.Where("name = ?", tag1Text).First(&tag1).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errno.ErrCommon.WithMessage("未查询到话题`" + tag1Text + "`,请先导入话题")
				}
				return errors.WithStack(err)
			}
			topicTag += tag1Text + ","
			topicTagId += strconv.Itoa(int(tag1.Id)) + ","
		}

		if tag2Text != "" {
			err = app.DB.Where("name = ?", tag2Text).First(&tag2).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errno.ErrCommon.WithMessage("为查询到话题`" + tag2Text + "`,请先导入话题")
				}
				return errors.WithStack(err)
			}
			topicTag += tag2Text + ","
			topicTagId += strconv.Itoa(int(tag2.Id)) + ","
		}
		topicTag = strings.TrimRight(topicTag, ",")
		topicTagId = strings.TrimRight(topicTagId, ",")

		imageList, err := srv.UploadImportTopicImage(path.Join(strings.Split(filename, "_")[0], row[0]))
		if err != nil {
			return errors.WithMessagef(err, "获取topic%d图片失败", importId)
		}

		topic := entity.Topic{
			TopicTag:   topicTag,
			UserId:     users[0].ID,
			Title:      title,
			Content:    content,
			Avatar:     users[0].AvatarUrl,
			Nickname:   users[0].Nickname,
			TopicTagId: topicTagId,
			ImportId:   importId + baseImportId,
			ImageList:  strings.Join(imageList, ","),
			Status:     3,
		}
		err = app.DB.Where("import_id = ?", importId+baseImportId).First(&topic).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.WithStack(err)
		}
		if topic.Id == 0 {
			err = app.DB.Create(&topic).Error
			if err != nil {
				return errors.WithMessagef(err, "创建topic%d失败", importId)
			}
			app.DB.Model(entity.Topic{}).Where("id = ?", topic.Id).Updates(map[string]interface{}{
				"created_at": time.Now().Add(time.Duration(int64(math.Ceil(rand.Float64()*-1200))) * time.Hour),
				"updated_at": time.Now().Add(time.Duration(int64(math.Ceil(rand.Float64()*-1200))) * time.Hour),
			})
		}

		if tag1 != nil {
			topicTag := entity.TopicTag{
				TopicId: topic.Id,
				TagId:   tag1.Id,
			}
			err = app.DB.Where("topic_id = ? and tag_id = ?", topic.Id, tag1.Id).First(&topicTag).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return errors.WithStack(err)
			}
			if topicTag.Id == 0 {
				err = app.DB.Create(&topicTag).Error
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
		if tag2 != nil {
			topicTag := entity.TopicTag{
				TopicId: topic.Id,
				TagId:   tag2.Id,
			}
			err = app.DB.Where("topic_id = ? and tag_id = ?", topic.Id, tag2.Id).First(&topicTag).Error
			if err != nil && err != gorm.ErrRecordNotFound {
				return errors.WithStack(err)
			}
			if topicTag.Id == 0 {
				err = app.DB.Create(&topicTag).Error
				if err != nil {
					return errors.WithStack(err)
				}
			}
		}
	}

	return nil
}

//CreateTopic 创建文章
func (srv TopicService) CreateTopic(userId int64, avatarUrl, nikeName, openid string, title, content string, tagIds []int64, images []string) (*entity.Topic, error) {
	topicModel := &entity.Topic{}

	// 文本内容审核
	if content != "" {
		if err := validator.CheckMsgWithOpenId(openid, content); err != nil {
			app.Logger.Error(fmt.Errorf("create Topic error:%s", err.Error()))
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "发帖-文本内容审核"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, openid, zhuGeAttr)
			return topicModel, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	// 图片内容审核
	if len(images) > 1 {
		for i, imgUrl := range images {
			if err := validator.CheckMediaWithOpenId(openid, imgUrl); err != nil {
				app.Logger.Error(fmt.Errorf("create Topic error:%s", err.Error()))
				zhuGeAttr := make(map[string]interface{}, 0)
				zhuGeAttr["场景"] = "发帖-图片内容审核"
				zhuGeAttr["失败原因"] = err.Error()
				track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, openid, zhuGeAttr)
				return topicModel, errno.ErrCommon.WithMessage("图片: " + strconv.Itoa(i) + " " + err.Error())
			}
		}
	}

	//处理images
	imageStr := strings.Join(images, ",")

	//topic
	topicModel = &entity.Topic{
		UserId:    userId,
		Title:     title,
		Content:   content,
		ImageList: imageStr,
		Status:    entity.TopicStatusNeedVerify,
		Avatar:    avatarUrl,
		Nickname:  nikeName,
		CreatedAt: model.Time{Time: time.Now()},
		UpdatedAt: model.Time{Time: time.Now()},
	}

	if len(tagIds) > 0 {
		//tag
		tagModel := make([]entity.Tag, 0)
		for _, tagId := range tagIds {
			tagModel = append(tagModel, entity.Tag{
				Id: tagId,
			})
		}

		tag := srv.tagModel.GetById(tagIds[0])
		topicModel.TopicTag = tag.Name
		topicModel.TopicTagId = strconv.FormatInt(tag.Id, 10)
		topicModel.Tags = tagModel
	}

	if err := srv.topicModel.Save(topicModel); err != nil {
		return topicModel, errno.ErrCommon.WithMessage("帖子保存失败")
	}

	return topicModel, nil
}

// UpdateTopic 更新帖子
func (srv TopicService) UpdateTopic(userId int64, avatarUrl, nikeName, openid string, topicId int64, title, content string, tagIds []int64, images []string) (*entity.Topic, error) {

	//查询记录是否存在
	topicModel := srv.topicModel.FindById(topicId)

	if topicModel.Id == 0 {
		return topicModel, errno.ErrCommon.WithMessage("该帖子不存在")
	}
	if topicModel.UserId != userId {
		return topicModel, errno.ErrCommon.WithMessage("无权限修改")
	}
	if content != "" {
		//检查内容
		if err := validator.CheckMsgWithOpenId(openid, content); err != nil {
			app.Logger.Error(fmt.Errorf("update Topic error:%s", err.Error()))
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "更新帖子"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, openid, zhuGeAttr)
			return topicModel, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	if len(images) > 1 {
		for i, imgUrl := range images {
			if err := validator.CheckMediaWithOpenId(openid, imgUrl); err != nil {
				app.Logger.Error(fmt.Errorf("create Topic error:%s", err.Error()))
				zhuGeAttr := make(map[string]interface{}, 0)
				zhuGeAttr["场景"] = "发帖-图片内容审核"
				zhuGeAttr["失败原因"] = err.Error()
				track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, openid, zhuGeAttr)
				return topicModel, errno.ErrCommon.WithMessage("图片: " + strconv.Itoa(i) + " " + err.Error())
			}
		}
	}

	//处理images
	imageStr := strings.Join(images, ",")

	//更新帖子
	topicModel.Title = title
	topicModel.Avatar = avatarUrl
	topicModel.Nickname = nikeName
	topicModel.ImageList = imageStr
	topicModel.Content = content

	if topicModel.Status != 3 {
		topicModel.Status = 1
	}
	//tag
	if len(tagIds) > 0 {
		tagModel := make([]entity.Tag, 0)
		for _, tagId := range tagIds {
			tagModel = append(tagModel, entity.Tag{
				Id: tagId,
			})
		}
		tag := srv.tagModel.GetById(tagIds[0])
		topicModel.TopicTag = tag.Name
		topicModel.TopicTagId = strconv.FormatInt(tag.Id, 10)
		if err := app.DB.Model(&topicModel).Association("Tags").Replace(tagModel); err != nil {
			return topicModel, errno.ErrCommon.WithMessage("Tag更新失败")
		}

	} else {
		topicModel.TopicTag = ""
		topicModel.TopicTagId = ""
		err := app.DB.Model(&topicModel).Association("Tags").Clear()
		if err != nil {
			return topicModel, errno.ErrCommon.WithMessage("Tag更新失败")
		}
	}

	if err := app.DB.Model(&topicModel).Updates(&topicModel).Error; err != nil {
		return topicModel, errno.ErrCommon.WithMessage("帖子更新失败")
	}
	return topicModel, nil
}

// DetailTopic 获取topic详情
func (srv TopicService) DetailTopic(topicId int64) (*entity.Topic, error) {
	//查询数据是否存在
	topic := srv.topicModel.FindById(topicId)
	if topic.Id == 0 {
		return topic, errno.ErrCommon.WithMessage("数据不存在")
	}
	//views+1
	if err := srv.topicModel.AddTopicSeeCount(topic.Id, 1); err != nil {
		app.Logger.Errorf("更新topic查看次数失败 : %s", err.Error())
	}
	return topic, nil
}

// DelTopic 软删除
func (srv TopicService) DelTopic(userId, topicId int64) error {
	topicModel := srv.topicModel.FindById(topicId)
	if topicModel.Id == 0 {
		return errno.ErrCommon.WithMessage("该帖子不存在")
	}
	if topicModel.UserId != userId {
		return errno.ErrCommon.WithMessage("无权限删除")
	}

	if err := app.DB.Delete(&topicModel).Error; err != nil {
		return errno.ErrInternalServer
	}
	return nil
}

func (srv TopicService) GetSubCommentCount(ids []int64) (result []CommentCount) {
	app.DB.Model(&entity.CommentIndex{}).
		Select("root_comment_id as total_id, count(*) as total").
		Where("state = ?", 0).
		Where("root_comment_id in ?", ids).
		Group("root_comment_id").
		Find(&result)
	return result
}

func (srv TopicService) GetRootCommentCount(ids []int64) (result []CommentCount) {
	app.DB.Model(&entity.CommentIndex{}).Select("obj_id as topic_id, count(*) as total").
		Where("obj_id in ?", ids).
		Where("to_comment_id = 0").
		Where("state = ?", 0).
		Group("obj_id").
		Find(&result)
	return result
}

func (srv TopicService) GetCommentCount(ids []int64) (result []CommentCount) {
	app.DB.Model(&entity.CommentIndex{}).Select("obj_id as topic_id, count(*) as total").
		Where("obj_id in ?", ids).
		Where("state = ?", 0).
		Group("obj_id").
		Find(&result)
	return result
}

func (srv TopicService) getTopicImage(importId int, p string) ([]string, error) {
	files, err := ioutil.ReadDir(path.Join(p, strconv.Itoa(importId)))
	if err != nil {
		return nil, err
	}
	list := make([]string, 0)
	for _, f := range files {
		list = append(list, fmt.Sprintf("https://miotech-resource.oss-cn-hongkong.aliyuncs.com/static/mp2c/images/topic/info/%d/%s", importId, f.Name()))
	}
	return list, nil
}

func (srv TopicService) CountTopic(param repository.GetTopicCountBy) (int64, error) {
	var total int64
	query := app.DB.Model(&entity.Topic{})
	if param.UserId != 0 {
		query.Where("topic.user_id = ?", param.UserId)
	}
	err := query.Count(&total).Group("topic.id").Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

func (srv TopicService) UpdateAuthor(userId, passiveUserId int64) error {
	return app.DB.Model(&entity.Topic{}).Where("topic.user_id = ?", passiveUserId).Update("user_id", userId).Error
}
