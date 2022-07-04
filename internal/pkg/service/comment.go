package service

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"time"
)

var DefaultCommentService = NewCommentService(repository.DefaultCommentRepository)

type (
	CommentService interface {
		FindOne(commentId int64) (*entity.CommentIndex, error)
		FindOneQuery(data *entity.CommentIndex) (*entity.CommentIndex, error)
		FindAll(data *entity.CommentIndex) ([]*entity.CommentIndex, int64, error)
		FindPageListByPage(data *entity.CommentIndex, page, pageSize int64, orderBy string) ([]*entity.CommentIndex, int64, error)
		FindPageListByIdDESC(data *entity.CommentIndex, preCommentMinId, pageSize int64) ([]*entity.CommentIndex, int64, error)
		FindPageListByIdASC(data *entity.CommentIndex, preCommentMinId, pageSize int64) ([]*entity.CommentIndex, int64, error)
		CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string) error
		UpdateComment(userId, commentId int64, message string) error
		DelComment(userId, commentId int64) error
		DelCommentSoft(userId, commentId int64) error
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
	if data.ID != 0 {
		builder.Where("id = ?", data.ID)
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

// FindPageListByPage 除去置顶评论外，按时间分页获取评论
func (srv *defaultCommentService) FindPageListByPage(data *entity.CommentIndex, page, pageSize int64, orderBy string) ([]*entity.CommentIndex, int64, error) {
	//1. 分页获取根评论列表及各自子评论总数
	rootComments := make([]map[string]interface{}, 0)
	err := app.DB.Table("comment_index as index").
		Select("index.*,(select count(*) from comment_index where RootCommentId_comment_id = index.id) as sub_count").
		Where("index.RootCommentId_comment_id = 0 and index.to_comment_id = 0").
		Where("index.obj_id = ?", data.ObjId).
		Find(&rootComments).Error
	if err != nil {
		return nil, 0, err
	}
	//2.根据根评论id获取各自前三条子评论,以like_count desc和id asc排序
	rootCommentIds := make([]int64, 0)
	memberIds := make([]int64, 0)
	for _, item := range rootComments {
		rootCommentIds = append(rootCommentIds, item["id"].(int64))
		memberIds = append(memberIds, item["member_id"].(int64))
	}
	subComments := make([]entity.CommentIndex, 0)
	err = app.DB.Table("comment_index as index").
		Where("(select count(id) from comment_index where RootCommentId_comment_id = index.RootCommentId_comment_id and like_count > index.like_count) < 3").
		Where("index.RootCommentId_comment_id in ?", rootCommentIds).
		Order("RootCommentId_comment_id").
		Find(subComments).Error
	if err != nil {
		return nil, 0, err
	}
	//3.获取用户昵称及头像

	usersInfo := make([]map[string]interface{}, 0)
	err = app.DB.Model(&entity.User{}).Select("id, nick_name, avatar_url").Where("id in = ?", memberIds).Find(&usersInfo).Error
	if err != nil {
		return nil, 0, err
	}
	//4.组合数据
	for _, rootCommentId := range rootComments {
		for _, subComment := range subComments {

		}
	}
	return nil, 0, nil
}

func (srv *defaultCommentService) FindPageListByIdDESC(data *entity.CommentIndex, preCommentMinId, pageSize int64) ([]*entity.CommentIndex, int64, error) {
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
	byPage, err := srv.commentModel.FindPageListByIdDESC(builder, preCommentMinId, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := srv.commentModel.FindCount(builder)
	if err != nil {
		return nil, 0, err
	}
	return byPage, count, err
}

func (srv *defaultCommentService) FindPageListByIdASC(data *entity.CommentIndex, preCommentMinId, pageSize int64) ([]*entity.CommentIndex, int64, error) {
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
	byPage, err := srv.commentModel.FindPageListByIdASC(builder, preCommentMinId, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := srv.commentModel.FindCount(builder)
	if err != nil {
		return nil, 0, err
	}
	return byPage, count, nil
}

func (srv *defaultCommentService) UpdateComment(userId, commentId int64, message string) error {
	req := entity.CommentIndex{
		ID:       commentId,
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

func (srv *defaultCommentService) CreateComment(userId, topicId, RootCommentId, ToCommentId int64, message string) error {
	data := &entity.CommentIndex{
		ObjId:         topicId,
		Message:       message,
		MemberId:      userId,
		RootCommentId: RootCommentId,
		ToCommentId:   ToCommentId,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
	_, err := srv.commentModel.Insert(data)
	if err != nil {
		return err
	}
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
			err = app.DB.Model(ToCommentIdRow).Select("RootCommentId_count", "count").Updates(ToCommentIdRow).Error
			if err != nil {
				return errors.WithMessage(err, "update ToCommentIdRow:")
			}
			//本条评论
			ToCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("ToCommentId = ?", ToCommentId))
			if err != nil {
				return err
			}
			if ToCommentIdRow.RootCommentId != 0 {
				ToCommentId = ToCommentIdRow.RootCommentId
			}
			data.RootCommentId = ToCommentId          //更新RootCommentId
			data.Floor = int32(ToCommentIdChildCount) //更新楼层
			err = app.DB.Model(data).Select("RootCommentId", "floor").Updates(data).Error
			if err != nil {
				return errors.WithMessage(err, "update dataRow:")
			}

			//顶级评论
			if ToCommentIdRow.RootCommentId != 0 {
				RootCommentIdChildCount, err := srv.commentModel.FindCount(srv.commentModel.CountBuilder("id").Where("RootCommentId = ?", ToCommentIdRow.RootCommentId))
				if err != nil {
					return err
				}
				//更新顶级评论下的评论数量
				err = app.DB.Model(&entity.CommentIndex{}).Where("id = ?", ToCommentIdRow.RootCommentId).Update("count", RootCommentIdChildCount).Error
				if err != nil {
					return errors.WithMessage(err, "update RootCommentIdRow:")
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
