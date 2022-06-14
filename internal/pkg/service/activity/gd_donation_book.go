package activity

import (
	"errors"
	"gorm.io/gorm"
	"mio/internal/pkg/model"
	entity2 "mio/internal/pkg/model/entity"
	entity "mio/internal/pkg/model/entity/activity"
	"mio/internal/pkg/repository"
	repoactivity "mio/internal/pkg/repository/activity"
)

var DefaultGDdbService = GDdbService{repo: repoactivity.DefaultGDDonationBookRepository}

type GDdbService struct {
	repo repoactivity.GDDonationBookRepository
}

//CreateUser 创建活动用户
func (srv GDdbService) CreateUser(userId, inviteId int64) (*entity.GDDonationBookRecord, error) {
	//检查是否已经存在
	record := srv.repo.FindBy(repoactivity.FindRecordBy{
		UserId: userId,
	})
	if record.ID != 0 {
		return &record, nil
	}
	record = entity.GDDonationBookRecord{
		UserId:       userId,
		AnswerStatus: 0,
		IsSuccess:    0,
		InviteId:     inviteId,
		CreatedAt:    model.Time{},
		UpdatedAt:    model.Time{},
	}
	if inviteId != 0 {
		record.InviteType = 1
	}
	err := srv.repo.Create(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// HomePage 首页返回数据
func (srv GDdbService) HomePage(userId, inviteId int64) (*GDDbHomePageResponse, error) {
	//返回用户信息
	var userAnswerRes repoactivity.GDDbHomePageUserInfo

	if userId != 0 {
		userInfo, err := srv.CreateUser(userId, inviteId)
		if err != nil {
			return nil, err
		}
		userAnswerRes, err = srv.GetUser(userInfo)
		if err != nil {
			return nil, err
		}
	}

	//返回学校捐赠排行
	schoolRes := repoactivity.DefaultGDDbSchoolRankRepository.GetRank()
	//组装数据
	record := &GDDbHomePageResponse{
		User:   userAnswerRes,
		School: schoolRes,
	}
	return record, nil
}

// GetUser 用于返回首页数据
func (srv GDdbService) GetUser(userRecord *entity.GDDonationBookRecord) (repoactivity.GDDbHomePageUserInfo, error) {
	var userResult, inviteResult repoactivity.GDDbUserInfo
	var invitedResult []repoactivity.GDDbUserInfo
	userRepo := repository.NewUserRepository()
	//用户是被邀请人
	if userRecord.InviteType != 0 {
		//获取邀请人信息
		invite := srv.repo.GetUserBy(repoactivity.FindRecordBy{UserId: userRecord.InviteId})
		inviteUser := userRepo.GetUserById(userRecord.InviteId)
		inviteResult = repoactivity.GDDbUserInfo{
			GDDonationBookRecord: invite,
			AvatarUrl:            inviteUser.AvatarUrl,
			Nickname:             inviteUser.Nickname,
		}
		//判断邀请人已经成团，但自己没有成团
		if invite.IsSuccess == 1 && userRecord.IsSuccess == 0 {
			//更新成团长
			userRecord.InviteType = 0
			userRecord.InviteId = 0
			inviteResult = repoactivity.GDDbUserInfo{}
			err := srv.repo.Save(userRecord)
			if err != nil {
				return repoactivity.GDDbHomePageUserInfo{}, err
			}
		}
	}

	//当前用户活动记录
	user := userRepo.GetUserById(userRecord.UserId)
	//受邀者用户活动记录
	invitedRes := srv.repo.GetInvitedBy(repoactivity.FindRecordBy{UserId: userRecord.UserId})
	invitedIds := make([]int64, 0)
	for _, invited := range invitedRes {
		invitedIds = append(invitedIds, invited.UserId)
	}
	invitedUsers := userRepo.GetUserListBy(repository.GetUserListBy{UserIds: invitedIds})
	invitedUsersMap := make(map[int64]entity2.User)
	for _, invited := range invitedUsers {
		invitedUsersMap[invited.ID] = invited
	}
	//组装数据
	userResult = repoactivity.GDDbUserInfo{
		GDDonationBookRecord: *userRecord,
		AvatarUrl:            user.AvatarUrl,
		Nickname:             user.Nickname,
	}
	for _, invited := range invitedRes {
		invitedResult = append(invitedResult, repoactivity.GDDbUserInfo{
			GDDonationBookRecord: invited,
			AvatarUrl:            invitedUsersMap[invited.UserId].AvatarUrl,
			Nickname:             invitedUsersMap[invited.UserId].Nickname,
		})
	}
	//返回数据
	res := repoactivity.GDDbHomePageUserInfo{
		UserInfo:    userResult,
		InviteInfo:  inviteResult,
		InvitedInfo: invitedResult,
	}

	return res, nil
}

// CheckActivityStatus 检测成团状态
func (srv GDdbService) CheckActivityStatus(userId, schoolId int64) error {
	var userInfo, inviteInfo entity.GDDonationBookRecord
	userInfo = srv.repo.FindBy(repoactivity.FindRecordBy{
		UserId: userId,
	})
	if userInfo.InviteId != 0 {
		inviteInfo = srv.repo.FindBy(repoactivity.FindRecordBy{
			UserId: userInfo.InviteId,
		})
		if inviteInfo.IsSuccess == 1 && userInfo.IsSuccess == 0 {
			//更新用户状态
			userInfo.InviteId = 0
			userInfo.InviteType = 0
			err := srv.repo.Save(&userInfo)
			if err != nil {
				return err
			}
			return errors.New("慢了一步，好友已和他人完成共同捐赠")
		}
	}

	return nil
}

func (srv GDdbService) SaveSchoolInfo(userName string, schoolId, gradeId, userId int64, classNumber uint32) error {
	//更新答题状态
	err := srv.UpdateAnswerStatus(userId, 2)
	if err != nil {
		return err
	}
	record := repoactivity.DefaultGDDbUserSchoolRepository.FindBy(repoactivity.FindRecordBy{UserId: userId})
	if record.ID != 0 {
		//创建
		req := &entity.GDDbUserSchool{
			UserId:      userId,
			UserName:    userName,
			SchoolId:    schoolId,
			GradeId:     gradeId,     //年级
			ClassNumber: classNumber, //班级
			CreatedAt:   model.Time{},
			UpdatedAt:   model.Time{},
		}
		err = repoactivity.DefaultGDDbUserSchoolRepository.Create(req)
	} else {
		//更新
		record.SchoolId = schoolId
		record.ClassNumber = classNumber
		err = repoactivity.DefaultGDDbUserSchoolRepository.Save(&record)
	}
	if err != nil {
		return err
	}
	return nil
}

// UpdateAnswerStatus 更新答题状态
func (srv GDdbService) UpdateAnswerStatus(userId int64, status int) error {
	record := srv.repo.FindBy(repoactivity.FindRecordBy{UserId: userId})
	if record.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	if record.AnswerStatus == 2 {
		return nil
	}
	record.AnswerStatus = status
	return srv.repo.Save(&record)
}

// IncrRank  学校捐赠书+1
func (srv GDdbService) IncrRank(schoolId int64) error {
	rankInfo := repoactivity.DefaultGDDbSchoolRankRepository.FindBy(repoactivity.FindRecordBy{UserId: schoolId})
	var err error
	if rankInfo.ID == 0 {
		schoolInfo := repoactivity.DefaultGDDbSchoolRepository.FindById(schoolId)
		//新增
		req := &entity.GDDbSchoolRank{
			SchoolId:     schoolInfo.ID,
			SchoolName:   schoolInfo.SchoolName,
			DonateNumber: 1,
			CreatedAt:    model.Time{},
			UpdatedAt:    model.Time{},
		}
		err = repoactivity.DefaultGDDbSchoolRankRepository.Create(req)
	} else {
		//更新
		rankInfo.DonateNumber++
		err = repoactivity.DefaultGDDbSchoolRankRepository.Save(&rankInfo)
	}
	if err != nil {
		return err
	}
	return nil
}

func (srv GDdbService) GetCityList() []entity.GDDbCity {
	return repoactivity.DefaultGDDbCityRepository.FindAll()
}

// GetGradeList 获取年级列表
func (srv GDdbService) GetGradeList() []entity.GDDbGrade {
	return repoactivity.DefaultGDDbGradeRepository.FindAll()
}

// GetSchoolList 获取学校列表
func (srv GDdbService) GetSchoolList(schoolName string, cityId, gradeId int64) []entity.GDDbSchool {
	//获取年级类型
	gradeInfo := repoactivity.DefaultGDDbGradeRepository.FindById(gradeId)
	find := repoactivity.FindSchoolBy{
		CityId:     cityId,
		SchoolName: schoolName,
		GradeType:  gradeInfo.Type,
	}

	return repoactivity.DefaultGDDbSchoolRepository.FindAllBy(find)
}

// CreateSchool  新建学校信息
func (srv GDdbService) CreateSchool(schoolName string, cityId int64, gradeType int) error {
	//获取年级类型
	schoolRes := repoactivity.DefaultGDDbSchoolRepository.FindBy(repoactivity.FindSchoolBy{
		CityId:     cityId,
		GradeType:  gradeType,
		SchoolName: schoolName,
	})
	if schoolRes.ID != 0 {
		return errors.New("学校已存在")
	}
	err := repoactivity.DefaultGDDbSchoolRepository.Create(&entity.GDDbSchool{
		CityId:     cityId,
		Type:       gradeType,
		SchoolName: schoolName,
		CreatedAt:  model.Time{},
		UpdatedAt:  model.Time{},
	})
	if err != nil {
		return err
	}
	return nil
}

// GetAchievement 获取我的成就
func (srv GDdbService) GetAchievement(userId int64) entity.GDDonationBookRecord {
	return srv.repo.FindBy(repoactivity.FindRecordBy{UserId: userId})
}
