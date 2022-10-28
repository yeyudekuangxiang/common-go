package service

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/service/track"
	"mio/internal/pkg/util"
	"mio/internal/pkg/util/validator"
	"mio/pkg/errno"
	"strconv"
	"time"
)

var DefaultCommentService = NewCommentService(repository.NewCommentRepository(context.NewMioContext()))

type (
	CommentService interface {
		FindOne(commentId int64) (*entity.CommentIndex, error)
		FindOneQuery(data *entity.CommentIndex) (*entity.CommentIndex, error)
		FindAll(data *entity.CommentIndex) ([]*entity.CommentIndex, int64, error)
		FindSubList(data *entity.CommentIndex, offset, limit int) ([]*entity.CommentIndex, int64, error)
		FindListAndChild(data *entity.CommentIndex, offset, limit int) ([]*entity.CommentIndex, int64, error)
		CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string, openId string) (entity.CommentIndex, int64, error)
		UpdateComment(userId, commentId int64, message string) error
		DelComment(userId, commentId int64) error
		DelCommentSoft(userId, commentId int64) error
		Like(userId, commentId int64, openId string) (*entity.CommentLike, int64, error)
		AddTopicLikeCount(commentId int64, num int) error
	}
)

type defaultCommentService struct {
	commentModel repository.CommentModel
}

func NewCommentService(model repository.CommentModel) CommentService {
	return &defaultCommentService{
		commentModel: model,
	}
}

func (srv *defaultCommentService) FindOne(commentId int64) (*entity.CommentIndex, error) {
	one, err := srv.commentModel.FindOne(commentId)
	if err != nil {
		return nil, err
	}
	return one, nil
}

func (srv *defaultCommentService) FindOneQuery(data *entity.CommentIndex) (*entity.CommentIndex, error) {
	builder := srv.commentModel.RowBuilder()
	if data.Id != 0 {
		builder.Where("id = ?", data.Id)
	}
	if data.MemberId != 0 {
		builder.Where("member_id = ?", data.MemberId)
	}
	if data.RootCommentId != 0 {
		builder.Where("RootCommentId = ?", data.RootCommentId)
	}
	if data.ToCommentId != 0 {
		builder.Where("ToCommentId = ?", data.ToCommentId)
	}
	if data.ObjId != 0 {
		builder.Where("obj_id = ?", data.ObjId)
	}
	if data.Attrs != 0 {
		builder.Where("attrs = ?", data.Attrs)
	}
	builder.Where("state = ?", 0)
	query, err := srv.commentModel.FindOneQuery(builder)
	if err != nil {
		return nil, err
	}
	return query, err
}

func (srv *defaultCommentService) FindAll(data *entity.CommentIndex) ([]*entity.CommentIndex, int64, error) {
	builder := srv.commentModel.RowBuilder()
	if data.MemberId != 0 {
		builder.Where("member_id = ?", data.MemberId)
	}
	if data.RootCommentId != 0 {
		builder.Where("RootCommentId = ?", data.RootCommentId)
	}
	if data.ToCommentId != 0 {
		builder.Where("ToCommentId = ?", data.ToCommentId)
	}
	if data.ObjId != 0 {
		builder.Where("obj_id = ?", data.ObjId)
	}
	if data.Attrs != 0 {
		builder.Where("attrs = ?", data.Attrs)
	}
	builder.Where("state = ?", 0)
	all, err := srv.commentModel.FindAll(builder, "like_count DESC")
	if err != nil {
		return nil, 0, err
	}
	count, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id"))
	if err != nil {
		return nil, 0, err
	}
	return all, count, nil
}

// FindListAndChild 获取顶级评论列表及各分属下的三条子评论
func (srv *defaultCommentService) FindListAndChild(params *entity.CommentIndex, offset, limit int) ([]*entity.CommentIndex, int64, error) {
	commentList := make([]*entity.CommentIndex, 0)
	var total int64
	err := app.DB.Model(&entity.CommentIndex{}).
		Preload("RootChild", func(db *gorm.DB) *gorm.DB {
			return db.Where("(select count(*) from comment_index index where index.root_comment_id = comment_index.root_comment_id and index.id <= comment_index.id) <= ?", 3).
				Order("comment_index.like_count desc").Preload("Member")
		}).
		Preload("Member").
		Where("to_comment_id = ?", 0).
		Where("obj_id = ?", params.ObjId).
		Where("state = ?", 0).
		Count(&total).
		Order("comment_index.like_count desc,comment_index.id asc").
		Limit(limit).
		Offset(offset).
		Find(&commentList).Error
	if err != nil {
		return nil, 0, err
	}
	return commentList, total, nil
}

//FindSubList 分页获取子评论列表
func (srv *defaultCommentService) FindSubList(data *entity.CommentIndex, offSize, limit int) ([]*entity.CommentIndex, int64, error) {
	commentList := make([]*entity.CommentIndex, 0)
	var total int64
	err := srv.commentModel.RowBuilder().
		Preload("Member").
		Where("root_comment_id = ?", data.RootCommentId).
		Where("state = ?", 0).
		Count(&total).
		Order("comment_index.like_count desc, comment_index.id asc").
		Find(&commentList).Error
	if err != nil {
		return nil, 0, err
	}
	return commentList, total, nil
}

func (srv *defaultCommentService) UpdateComment(userId, commentId int64, message string) error {

	req := entity.CommentIndex{
		Id:       commentId,
		Message:  message,
		MemberId: userId,
	}
	err := srv.commentModel.Update(&req)
	if err != nil {
		return err
	}
	return nil
}

func (srv *defaultCommentService) DelComment(userId, commentId int64) error {
	err := srv.commentModel.Delete(commentId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (srv *defaultCommentService) DelCommentSoft(userId, commentId int64) error {
	err := srv.commentModel.DeleteSoft(commentId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (srv *defaultCommentService) CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string, openId string) (entity.CommentIndex, int64, error) {
	topic, err := DefaultTopicService.DetailTopic(topicId)
	if err != nil {
		return entity.CommentIndex{}, 0, err
	}

	if message != "" {
		//检查内容
		if err = validator.CheckMsgWithOpenId(openId, message); err != nil {
			app.Logger.Error(fmt.Errorf("create Comment error:%s", err.Error()))
			zhuGeAttr := make(map[string]interface{}, 0)
			zhuGeAttr["场景"] = "发布评论"
			zhuGeAttr["失败原因"] = err.Error()
			track.DefaultZhuGeService().Track(config.ZhuGeEventName.MsgSecCheck, openId, zhuGeAttr)
			return entity.CommentIndex{}, 0, errno.ErrCommon.WithMessage(err.Error())
		}
	}

	comment := entity.CommentIndex{
		ObjId:         topicId,
		Message:       message,
		MemberId:      userId,
		RootCommentId: RootCommentId,
		ToCommentId:   ToCommentId,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
	if topic.UserId == userId {
		comment.IsAuthor = 1
	}
	_, err = srv.commentModel.Insert(&comment)
	if err != nil {
		return entity.CommentIndex{}, 0, err
	}
	//更新topic
	app.DB.Model(&topic).Update("updated_at", model.Time{Time: time.Now()})
	//更新count数据
	if ToCommentId != 0 {
		err = srv.commentModel.Trans(func(tx *gorm.DB) error {
			//回复的评论
			ToCommentIdRow, err := srv.commentModel.FindOne(ToCommentId)
			if err != nil {
				return err
			}
			ToCommentIdRow.RootCount++ //更新父评论的根评论数量
			ToCommentIdRow.Count++     //更新父评论的评论数量
			err = tx.Model(ToCommentIdRow).Select("root_count", "count").Updates(ToCommentIdRow).Error
			if err != nil {
				return errors.WithMessage(err, "update ToCommentIdRow:")
			}
			//本条评论
			ToCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("to_comment_id = ?", ToCommentId))
			if err != nil {
				return err
			}
			if ToCommentIdRow.RootCommentId != 0 {
				ToCommentId = ToCommentIdRow.RootCommentId
			}
			comment.RootCommentId = ToCommentId          //更新RootCommentId
			comment.Floor = int32(ToCommentIdChildCount) //更新楼层
			err = tx.Model(comment).Select("root_comment_id", "floor").Updates(comment).Error
			if err != nil {
				return errors.WithMessage(err, "update dataRow:")
			}

			//顶级评论
			if ToCommentIdRow.RootCommentId != 0 {
				RootCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("root_comment_id = ?", ToCommentIdRow.RootCommentId))
				if err != nil {
					return err
				}
				//更新顶级评论下的评论数量
				err = tx.Model(&entity.CommentIndex{}).Where("id = ?", ToCommentIdRow.RootCommentId).Update("count", RootCommentIdChildCount).Error
				if err != nil {
					return errors.WithMessage(err, "update RootCommentIdRow:")
				}
			}
			return nil
		})
		if err != nil {
			return entity.CommentIndex{}, 0, err
		}
	}
	//find comment
	one, err := srv.commentModel.FindOne(comment.Id)
	if err != nil {
		return entity.CommentIndex{}, 0, err
	}
	//更新积分
	messagerune := []rune(comment.Message)
	if len(messagerune) > 8 {
		message = string(messagerune[0:8])
	}

	point := int64(entity.PointCollectValueMap[entity.POINT_COMMENT])
	pointService := NewPointService(context.NewMioContext())
	_, err = pointService.IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:       openId,
		Type:         entity.POINT_COMMENT,
		BizId:        util.UUID(),
		ChangePoint:  point,
		AdminId:      0,
		Note:         "评论" + message + "..." + "成功",
		AdditionInfo: strconv.FormatInt(topic.Id, 10) + "#" + strconv.FormatInt(comment.Id, 10),
	})
	if err != nil {
		point = 0
	}
	return *one, point, nil
}

func (srv *defaultCommentService) Like(userId, commentId int64, openId string) (*entity.CommentLike, int64, error) {
	comment, err := srv.commentModel.FindOne(commentId)
	if err != nil {
		return &entity.CommentLike{}, 0, err
	}

	like, point, err := DefaultCommentLikeService.Like(userId, commentId, comment.Message, openId)
	if err != nil {
		return &entity.CommentLike{}, 0, err
	}
	return like, point, nil
}

func (srv *defaultCommentService) AddTopicLikeCount(commentId int64, num int) error {
	err := srv.commentModel.AddTopicLikeCount(commentId, num)
	if err != nil {
		return err
	}
	return nil
}
