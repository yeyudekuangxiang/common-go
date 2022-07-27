package service

import (
	"bytes"
	"fmt"
	"github.com/medivhzhan/weapp/v3"
	"github.com/medivhzhan/weapp/v3/security"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/pkg/wxapp"
	"mio/pkg/wxoa"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var DefaultTopicService = NewTopicService()

func NewTopicService() TopicService {
	return TopicService{
		topic: repository.DefaultTopicRepository,
		tag:   repository.DefaultTagRepository,
	}
}

type TopicService struct {
	topic       repository.ITopicRepository
	tag         repository.ITagRepository
	TokenServer *wxoa.AccessTokenServer
}

func (srv TopicService) GetTopicList(param repository.GetTopicPageListBy) ([]entity.Topic, error) {
	//获取置顶数据，按时间倒序排序
	//topList := make([]entity.Topic, 0)
	//app.DB.Model(entity.Topic{}).Preload("Tags").
	//	Where("status = ? and sort = ?", 3, 9999).
	//	Where("id > ?", param.ID).
	//	Where("tag_id = ?", param.TopicTagId).
	//	Order("updated_at desc").
	//	Limit(param.Limit).
	//	Find(&topList)
	//获取精华数据，按时间倒序排序
	//essenceList := make([]entity.Topic, 0)
	//剩下的数据，按时间倒叙排序
	//获取所有数据数量
	return nil, nil
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
		likeList := repository.TopicLikeRepository{DB: app.DB}.GetListBy(repository.GetTopicLikeListBy{
			TopicIds: topicIds,
			UserId:   userId,
		})
		for _, like := range likeList {
			topicLikeMap[int64(like.TopicId)] = like.Status == 1
		}
	}

	//整理数据
	detailList := make([]TopicDetail, 0)
	for _, topic := range topicList {
		user, err := DefaultUserService.GetUserById(topic.UserId)
		if err != nil {
			return nil, err
		}
		detailList = append(detailList, TopicDetail{
			Topic:         topic,
			IsLike:        topicLikeMap[topic.Id],
			UpdatedAtDate: topic.UpdatedAt.Format("01-02"),
			User:          user.ShortUser(),
		})
	}

	return detailList, nil
}

// GetTopicDetailPageList 通过topic表直接查询获取内容列表
func (srv TopicService) GetTopicDetailPageList(param repository.GetTopicPageListBy) ([]TopicDetail, int64, error) {
	list, total := srv.topic.GetTopicPageList(param)

	//更新曝光和查看次数
	srv.UpdateTopicFlowListShowCount(list, param.UserId)
	if param.ID != 0 && len(list) > 0 {
		app.Logger.Info("更新查看次数", list[0].Id, param.UserId)
		srv.UpdateTopicSeeCount(list[0].Id, param.UserId)
	}

	detailList, err := srv.fillTopicList(list, param.UserId)
	if err != nil {
		return nil, 0, err
	}
	return detailList, total, nil
}

// GetTopicDetailPageListByFlow 通过topic_flow内容流表获取内容列表 当topic_flow数据不存在时 会后台任务进行初始化并且调用 GetTopicDetailPageList 方法返回数据
func (srv TopicService) GetTopicDetailPageListByFlow(param repository.GetTopicPageListBy) ([]TopicDetail, int64, error) {

	topicList, total := srv.topic.GetFlowPageList(repository.GetTopicFlowPageListBy{
		Offset:     param.Offset,
		Limit:      param.Limit,
		UserId:     param.UserId,
		TopicId:    param.ID,
		TopicTagId: param.TopicTagId,
		Status:     entity.TopicStatusPublished,
	})
	if total == 0 {
		DefaultTopicFlowService.InitUserFlowByMq(param.UserId)
		return srv.GetTopicDetailPageList(param)
	}

	//更新曝光和查看次数
	srv.UpdateTopicFlowListShowCount(topicList, param.UserId)

	if param.ID != 0 && len(topicList) > 0 {
		app.Logger.Info("更新查看次数", param.UserId, topicList[0].Id)
		srv.UpdateTopicSeeCount(topicList[0].Id, param.UserId)
	}

	topicDetailList, err := srv.fillTopicList(topicList, param.UserId)

	if err != nil {
		return nil, 0, err
	}
	return topicDetailList, total, nil
}

// UpdateTopicSeeCount 更新内容的查看次数加1
func (srv TopicService) UpdateTopicSeeCount(topicId int64, userId int64) {
	err := initUserFlowPool.Submit(func() {
		topic := srv.topic.FindById(topicId)
		if topic.Id == 0 {
			return
		}

		seeCount := topic.SeeCount + 1

		if err := srv.topic.UpdateColumn(topic.Id, "see_count", seeCount); err != nil {
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

// GetShareWeappQrCode 获取小程序端内容详情页分享小程序码
func (srv TopicService) GetShareWeappQrCode(userId int, topicId int) ([]byte, string, error) {
	resp, err := wxapp.NewClient(app.Weapp).GetUnlimitedQRCodeResponse(&weapp.UnlimitedQRCode{
		Scene:     fmt.Sprintf("tid=%d&uid=%d&s=p", topicId, userId),
		Page:      "pages/cool-mio/mio-detail/index",
		Width:     100,
		IsHyaline: true,
	})
	if err != nil {
		return nil, "", err
	}
	if resp.ErrCode != 0 {
		return nil, "", errors.New(resp.ErrMsg)
	}
	return resp.Buffer, resp.ContentType, nil
}

// UpdateTopicSort 更新内容的排序权重
func (srv TopicService) UpdateTopicSort(topicId int64, sort int) error {
	topic := srv.topic.FindById(topicId)
	if topic.Id == 0 {
		return errors.New("未查询到此内容")
	}
	err := srv.topic.UpdateColumn(topicId, "sort", sort)
	if err != nil {
		return err
	}
	DefaultTopicFlowService.AfterUpdateTopicByMq(topicId)
	return nil
}

func (srv TopicService) ImportUser(filename string) error {

	file, err := excelize.OpenFile(filename)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if file.SheetCount == 0 {
		return errors.New("没有数据")
	}

	rows, err := file.GetRows(file.GetSheetList()[0])
	if err != nil {
		return errors.WithStack(err)
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}
		importId := row[0]
		nickname := row[1]
		//avatarImage := row[2]
		wechat := row[3]
		phone := row[4]
		user := entity.User{}
		if nickname == "星星充电" {
			app.DB.Where("openid = ?", wechat).First(&user)
		} else {
			app.DB.Where("nick_name = ?", nickname).First(&user)
		}

		if user.ID != 0 {
			if user.PhoneNumber != phone {
				return errors.Errorf("存在同名但手机号不同用户 %s", nickname)
			}
			log.Println("用户已存在", user)
			continue
		}

		avatar, err := srv.uploadImportUserAvatar(path.Join(strings.Split(filename, "_")[0], importId+".jpg"))
		if err != nil {
			return errors.WithMessage(err, "上传头像失败"+importId)
		}
		_, err = DefaultUserService.CreateUser(CreateUserParam{
			OpenId:      wechat,
			AvatarUrl:   avatar,
			Nickname:    nickname,
			PhoneNumber: phone,
			Source:      entity.UserSourceMio,
		})
		if err != nil {
			return errors.Errorf("创建用户失败 %s %v", nickname, err)
		}
	}
	return nil
}

func (srv TopicService) uploadImportUserAvatar(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}

	_, fileName := path.Split(filepath)

	name := fmt.Sprintf("images/topic/%s/%s", path.Base(path.Dir(filepath)), fileName)

	fmt.Println("上传头像", filepath, name)

	defer file.Close()

	return DefaultOssService.PutObject(name, file)
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
		u, err := DefaultOssService.PutObject(name, file)
		if err != nil {
			file.Close()
			return nil, err
		}
		images = append(images, u)
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
		return errors.New("没有数据")
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
			return errors.New("未查询到用户`" + nickname + "`,请先导入用户")
		}
		if len(users) > 1 {
			return errors.New("检测到有多个昵称为`" + nickname + "`的用户,请手动处理后再导入")
		}

		topicTag := ""
		topicTagId := ""
		var tag1, tag2 *entity.Tag
		if tag1Text != "" {
			err = app.DB.Where("name = ?", tag1Text).First(&tag1).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return errors.New("未查询到话题`" + tag1Text + "`,请先导入话题")
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
					return errors.New("为查询到话题`" + tag2Text + "`,请先导入话题")
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
func (srv TopicService) CreateTopic(userId int64, avatarUrl, nikeName, openid string, title, content string, tagIds []int64, images []string) error {
	//检查内容
	err := srv.checkMsgs(openid, content)
	if err != nil {
		return err
	}
	//处理images
	imageStr := strings.Join(images, ",")
	//tag
	tagModel := make([]entity.Tag, 0)
	for _, tagId := range tagIds {
		tagModel = append(tagModel, entity.Tag{
			Id: tagId,
		})
	}
	tag := srv.tag.GetById(tagIds[0])
	//topic
	topicModel := &entity.Topic{
		UserId:    userId,
		Title:     title,
		TopicTag:  tag.Name,
		Content:   content,
		ImageList: imageStr,
		Status:    entity.TopicStatusNeedVerify,
		Avatar:    avatarUrl,
		Nickname:  nikeName,
		Tags:      tagModel,
		CreatedAt: model.Time{},
		UpdatedAt: model.Time{},
	}
	if err := srv.topic.Save(topicModel); err != nil {
		return err
	}
	return nil
}

// UpdateTopic 更新帖子
func (srv TopicService) UpdateTopic(userId int64, avatarUrl, nikeName, openid string, topicId int64, title, content string, tagIds []int64, images []string) error {
	//检查内容
	err := srv.checkMsgs(openid, content)
	//查询记录是否存在
	topicModel := srv.topic.FindById(topicId)
	if topicModel.Id == 0 {
		return errors.New("该帖子不存在")
	}
	if topicModel.UserId != userId {
		return errors.New("无权限修改")
	}
	//处理images
	imageStr := strings.Join(images, ",")
	//tag
	tagModel := make([]entity.Tag, 0)
	for _, tagId := range tagIds {
		tagModel = append(tagModel, entity.Tag{
			Id: tagId,
		})
	}
	tag := srv.tag.GetById(tagIds[0])

	//更新帖子
	topicModel.Title = title
	topicModel.Avatar = avatarUrl
	topicModel.Nickname = nikeName
	topicModel.ImageList = imageStr
	topicModel.Content = content
	topicModel.TopicTag = tag.Name
	if err := app.DB.Model(&entity.Topic{}).Updates(topicModel).Error; err != nil {
		return err
	}
	err = app.DB.Model(&topicModel).Association("Tags").Replace(tagModel)
	if err != nil {
		return err
	}
	return nil
}

// DetailTopic 获取topic详情
func (srv TopicService) DetailTopic(topicId int64) (entity.Topic, error) {
	//查询数据是否存在
	topic := srv.topic.FindById(topicId)
	if topic.Id == 0 {
		return entity.Topic{}, errors.New("数据不存在")
	}
	//查找关联关系
	tagModels := make([]entity.Tag, 0)
	err := app.DB.Model(entity.Topic{}).Association("Tags").Find(&tagModels)
	if err != nil {
		return entity.Topic{}, err
	}
	topic.Tags = tagModels
	//更新查看次数
	return topic, nil
}

// DelTopic 软删除
func (srv TopicService) DelTopic(userId, topicId int64) error {
	topicModel := srv.topic.FindById(topicId)
	if topicModel.Id == 0 {
		return errors.New("该帖子不存在")
	}
	if topicModel.UserId != userId {
		return errors.New("无权限删除")
	}
	topicModel.Status = 4
	if err := app.DB.Model(&entity.Topic{}).Save(topicModel).Error; err != nil {
		return err
	}
	return nil
}

func (srv TopicService) checkMsgs(openid, content string) error {
	length := len(content)
	if length > 2500 {
		s := []rune(content)
		var buffer bytes.Buffer
		for i, str := range s {
			buffer.WriteString(string(str))
			if i > 0 && (i+1)%2500 == 0 {
				params := &security.MsgSecCheckRequest{
					Content: "",
					Version: 2,
					Scene:   3,
					Openid:  openid,
				}
				err := srv.checkMsg(params)
				buffer.Reset()
				if err != nil {
					return err
				}
			}
		}
		params := &security.MsgSecCheckRequest{
			Content: "",
			Version: 2,
			Scene:   3,
			Openid:  openid,
		}
		err := srv.checkMsg(params)
		buffer.Reset()
		if err != nil {
			return err
		}
		return nil
	}
	//处理内容
	params := &security.MsgSecCheckRequest{
		Content: content,
		Version: 2,
		Scene:   3,
		Openid:  openid,
	}
	err := srv.checkMsg(params)
	if err != nil {
		return err
	}
	return nil
}

func (srv TopicService) checkMsg(params *security.MsgSecCheckRequest) error {
	check, err := app.Weapp.NewSecurity().MsgSecCheck(params)
	if err != nil {
		return err
	}
	if check.ErrCode != 0 {
		return fmt.Errorf("check error: %s", check.ErrMSG)
	}
	return nil
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
