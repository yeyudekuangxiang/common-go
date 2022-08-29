package service

import (
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/repository/repotypes"
	"time"
)

func NewUserFriendService(ctx *context.MioContext) UserFriendService {
	return UserFriendService{ctx: ctx,
		repo: repository.NewUserFriendRepository(ctx)}
}

var DefaultUserFriendService = NewUserFriendService(context.NewMioContext())

type UserFriendService struct {
	ctx  *context.MioContext
	repo repository.UserFriendRepository
}

func (srv UserFriendService) GetUserFriendList(uid int64, userIds []int64) (map[int64]entity.UserFriend, error) {
	list := srv.repo.GetListBy(repotypes.GetUserFriendListBy{FUserIds: userIds, Uid: uid})
	listMap := make(map[int64]entity.UserFriend, 0)
	for _, j := range list {
		listMap[j.FUid] = entity.UserFriend{
			FUid: j.FUid,
		}
	}
	return listMap, nil
}

func (srv UserFriendService) Create(user *entity.User, invitedBy string) (*entity.UserFriend, error) {
	if user.ID == 0 {
		return &entity.UserFriend{}, nil
	}
	friend := entity.UserFriend{}
	InviteUser := entity.User{}
	errUser := app.DB.Where("openid = ?", invitedBy).First(&InviteUser).Error
	if errUser != nil && errUser != gorm.ErrRecordNotFound {
		panic(errUser)
	}
	fUid := InviteUser.ID
	if friend.ID != 0 && fUid != 0 {
		return &friend, nil
	}
	friend = entity.UserFriend{
		FUid:      fUid,
		Uid:       user.ID,
		Type:      1,
		Source:    1,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	friend = entity.UserFriend{
		FUid:      user.ID,
		Uid:       fUid,
		Type:      2,
		Source:    1,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
	return &friend, app.DB.Create(&friend).Error
}
