package activity

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/model"
	entity2 "mio/internal/pkg/model/entity"
	entity "mio/internal/pkg/model/entity/activity"
	"mio/internal/pkg/repository"
	repoactivity "mio/internal/pkg/repository/activity"
	"mio/internal/pkg/service"
)

var DefaultGDdbService = GDdbService{repo: repoactivity.DefaultGDDonationBookRepository}

type GDdbService struct {
	repo repoactivity.GDDonationBookRepository
}

//CreateUser 创建活动用户
func (srv GDdbService) CreateUser(userId, inviteId int64) (entity.GDDonationBookRecord, error) {
	//检查是否已经存在
	record := srv.repo.FindBy(repoactivity.FindRecordBy{
		UserId: userId,
	})
	if record.ID != 0 {
		return record, nil
	}
	if userId == inviteId {
		inviteId = 0
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
		return entity.GDDonationBookRecord{}, err
	}
	return record, nil
}

// HomePage 首页返回数据
func (srv GDdbService) HomePage(userId, inviteId int64) (GDDbHomePageResponse, error) {

	//返回用户信息
	userAnswerRes := repoactivity.GDDbHomePageUserInfo{
		UserInfo:    repoactivity.GDDbUserInfo{},
		InviteInfo:  repoactivity.GDDbUserInfo{},
		InvitedInfo: make([]repoactivity.GDDbUserInfo, 0),
	}
	schoolRes := make([]entity.GDDbSchoolRank, 0)
	record := GDDbHomePageResponse{
		User:   userAnswerRes,
		School: schoolRes,
	}
	if userId != 0 {
		userInfo, err := srv.CreateUser(userId, inviteId)
		if err != nil {
			return record, err
		}
		userAnswerRes, err = srv.GetUser(&userInfo)
		if err != nil {
			return record, err
		}
	}

	//返回学校捐赠排行
	schoolRes = repoactivity.DefaultGDDbSchoolRankRepository.GetRank()
	//组装数据
	record = GDDbHomePageResponse{
		User:   userAnswerRes,
		School: schoolRes,
	}
	return record, nil
}

// GetUser 用于返回首页数据
func (srv GDdbService) GetUser(userRecord *entity.GDDonationBookRecord) (repoactivity.GDDbHomePageUserInfo, error) {
	var userResult, inviteResult repoactivity.GDDbUserInfo
	invitedResult := make([]repoactivity.GDDbUserInfo, 0)
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
	if len(invitedRes) > 0 {
		invitedIds := make([]int64, 0)
		for _, invited := range invitedRes {
			invitedIds = append(invitedIds, invited.UserId)
		}
		invitedUsers := userRepo.GetUserListBy(repository.GetUserListBy{UserIds: invitedIds})
		invitedUsersMap := make(map[int64]entity2.User)
		for _, invited := range invitedUsers {
			invitedUsersMap[invited.ID] = invited
		}
		for _, invited := range invitedRes {
			invitedResult = append(invitedResult, repoactivity.GDDbUserInfo{
				GDDonationBookRecord: invited,
				AvatarUrl:            invitedUsersMap[invited.UserId].AvatarUrl,
				Nickname:             invitedUsersMap[invited.UserId].Nickname,
			})
		}
	}
	//组装数据
	userResult = repoactivity.GDDbUserInfo{
		GDDonationBookRecord: *userRecord,
		AvatarUrl:            user.AvatarUrl,
		Nickname:             user.Nickname,
	}

	//返回数据
	return repoactivity.GDDbHomePageUserInfo{
		UserInfo:    userResult,
		InviteInfo:  inviteResult,
		InvitedInfo: invitedResult,
	}, nil
}

func (srv GDdbService) GetUserSchool(userId int64) (repoactivity.GDDbUserSchool, error) {
	resp := repoactivity.GDDbUserSchool{}
	donationBookResult := srv.repo.GetUserBy(repoactivity.FindRecordBy{UserId: userId})
	user, _ := service.DefaultUserService.GetUserById(donationBookResult.UserId)
	userSchool := repoactivity.DefaultGDDbUserSchoolRepository.FindBy(repoactivity.FindRecordBy{UserId: userId})
	if userSchool.ID != 0 {
		school := repoactivity.DefaultGDDbSchoolRepository.FindById(userSchool.SchoolId)
		grade := repoactivity.DefaultGDDbGradeRepository.FindById(userSchool.GradeId)
		city := repoactivity.DefaultGDDbCityRepository.FindById(school.CityId)
		resp = repoactivity.GDDbUserSchool{
			GDDbUserInfo: repoactivity.GDDbUserInfo{
				GDDonationBookRecord: donationBookResult,
				AvatarUrl:            user.AvatarUrl,
				Nickname:             user.Nickname,
			},
			UserName:    userSchool.UserName,
			CityName:    city.CityName,
			GradeName:   grade.Grade,
			ClassNumber: userSchool.ClassNumber,
			SchoolName:  school.SchoolName,
		}
	}
	return resp, nil
}

// UpdateActivityUser 更新证书链接地址
func (srv GDdbService) UpdateActivityUser(userId int64, t int, url string) error {
	record := srv.repo.FindBy(repoactivity.FindRecordBy{UserId: userId})
	if record.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	if t == 1 {
		record.TitleUrl = url
	} else {
		record.CertificateUrl = url
	}
	return srv.repo.Save(&record)
}

// CheckActivityStatus 检测成团状态
func (srv GDdbService) CheckActivityStatus(userId, schoolId int64) error {
	var userInfo, inviteInfo entity.GDDonationBookRecord
	userInfo = srv.repo.FindBy(repoactivity.FindRecordBy{
		UserId: userId,
	})
	if userInfo.InviteId != 0 && userInfo.ID != 0 {
		inviteInfo = srv.repo.FindBy(repoactivity.FindRecordBy{
			UserId: userInfo.InviteId,
		})
		// 答题晚了
		if inviteInfo.IsSuccess == 1 && userInfo.IsSuccess == 0 {
			//更新用户状态
			userInfo.InviteType = 0
			_ = srv.repo.Save(&userInfo)
			return errors.New("慢了一步，好友已和他人完成共同捐赠")
		}
		//正常答题 更新状态
		ids := []int64{userInfo.UserId, userInfo.InviteId}
		if err := app.DB.Model(entity.GDDonationBookRecord{}).Where("user_id in ?", ids).Updates(entity.GDDonationBookRecord{IsSuccess: 1}).Error; err != nil {
			return errors.New("更新答题状态失败")
		}
		//被邀请者答题完成，更新学校排行榜
		err := srv.IncrRank(userInfo.InviteId)
		if err != nil {
			return errors.New("更新捐书人数失败")
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
		//更新
		record.UserName = userName
		record.SchoolId = schoolId
		record.GradeId = gradeId
		record.ClassNumber = classNumber
		err = repoactivity.DefaultGDDbUserSchoolRepository.Save(&record)
	} else {
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

// IncrRankBack IncrRank  学校捐赠书+1 废弃
func (srv GDdbService) IncrRankBack(userId int64) error {
	activityUser := repoactivity.DefaultGDDonationBookRepository.FindBy(repoactivity.FindRecordBy{UserId: userId})
	if activityUser.ID != 0 && activityUser.InviteType == 1 && activityUser.IsSuccess == 1 {
		//获取学校id
		var userSchoolList []entity.GDDbUserSchool
		schoolIds := make([]int64, 0)
		err := app.DB.Model(entity.GDDbUserSchool{}).Where("user_id = ? or user_id = ?", activityUser.UserId, activityUser.InviteId).Find(&userSchoolList).Error
		if err != nil {
			return err
		}
		for _, userSchool := range userSchoolList {
			schoolIds = append(schoolIds, userSchool.SchoolId)
		}
		//获取学校信息
		schoolList := repoactivity.DefaultGDDbSchoolRepository.FindAllBy(repoactivity.FindSchoolBy{SchoolIds: schoolIds})
		rankInfo := repoactivity.DefaultGDDbSchoolRankRepository.FindAllBy(repoactivity.FindSchoolBy{SchoolIds: schoolIds})
		if len(rankInfo) > 1 {
			updateReq := make([]entity.GDDbSchoolRank, 0)
			for _, rank := range rankInfo {
				updateReq = append(updateReq, entity.GDDbSchoolRank{DonateNumber: rank.DonateNumber + 1})
			}
			err = app.DB.Model(entity.GDDbSchoolRank{}).Where("school_id in ?", schoolIds).Updates(&rankInfo).Error
		} else {
			insertReq := make([]entity.GDDbSchoolRank, 0)
			for _, school := range schoolList {
				insertReq = append(insertReq, entity.GDDbSchoolRank{
					SchoolId:     school.ID,
					SchoolName:   school.SchoolName,
					DonateNumber: 1,
					CreatedAt:    model.Time{},
					UpdatedAt:    model.Time{},
				})
			}
			err = app.DB.Model(entity.GDDbSchoolRank{}).Create(&insertReq).Error
		}
		if err != nil {
			fmt.Printf("error:%e", err)
			return err
		}
	}

	return nil
}

// IncrRank  学校捐赠书+1
func (srv GDdbService) IncrRank(userId int64) error {
	activityUser := repoactivity.DefaultGDDonationBookRepository.FindBy(repoactivity.FindRecordBy{UserId: userId})
	var err error
	if activityUser.ID != 0 && activityUser.IsSuccess == 1 {
		//获取学校id
		USchool := repoactivity.DefaultGDDbUserSchoolRepository.FindBy(repoactivity.FindRecordBy{UserId: activityUser.UserId})
		//获取学校信息
		schoolInfo := repoactivity.DefaultGDDbSchoolRepository.FindBy(repoactivity.FindSchoolBy{SchoolId: USchool.SchoolId})
		rankInfo := repoactivity.DefaultGDDbSchoolRankRepository.FindBy(repoactivity.FindSchoolBy{SchoolId: USchool.SchoolId})
		if rankInfo.ID != 0 {
			rankInfo.DonateNumber++
			err = repoactivity.DefaultGDDbSchoolRankRepository.Save(&rankInfo)
		} else {
			insertReq := entity.GDDbSchoolRank{
				SchoolId:     schoolInfo.ID,
				SchoolName:   schoolInfo.SchoolName,
				DonateNumber: 1,
				CreatedAt:    model.Time{},
				UpdatedAt:    model.Time{},
			}
			err = repoactivity.DefaultGDDbSchoolRankRepository.Create(&insertReq)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (srv GDdbService) GetCityList() []entity.GDDbCity {
	list := make([]entity.GDDbCity, 0)
	list = repoactivity.DefaultGDDbCityRepository.FindAll()
	return list
}

// GetGradeList 获取年级列表
func (srv GDdbService) GetGradeList() []entity.GDDbGrade {
	list := make([]entity.GDDbGrade, 0)
	list = repoactivity.DefaultGDDbGradeRepository.FindAll()
	return list
}

// GetSchoolList 获取学校列表
func (srv GDdbService) GetSchoolList(schoolName string, cityId, gradeId int64) []entity.GDDbSchool {
	list := make([]entity.GDDbSchool, 0)
	gradeInfo := repoactivity.DefaultGDDbGradeRepository.FindById(gradeId)

	find := repoactivity.FindSchoolBy{
		CityId:     cityId,
		SchoolName: schoolName,
		GradeType:  gradeInfo.Type,
	}
	list = repoactivity.DefaultGDDbSchoolRepository.FindAllBy(find)

	return list
}

// CreateSchool  新建学校信息
func (srv GDdbService) CreateSchool(schoolName string, cityId int64, gradeType int) (int64, error) {
	//获取年级类型
	schoolRes := repoactivity.DefaultGDDbSchoolRepository.FindBy(repoactivity.FindSchoolBy{
		CityId:     cityId,
		GradeType:  gradeType,
		SchoolName: schoolName,
	})
	if schoolRes.ID != 0 {
		return 0, errors.New("学校已存在")
	}
	school := &entity.GDDbSchool{
		CityId:     cityId,
		Type:       gradeType,
		SchoolName: schoolName,
		CreatedAt:  model.Time{},
		UpdatedAt:  model.Time{},
	}
	err := repoactivity.DefaultGDDbSchoolRepository.Create(school)
	if err != nil {
		return 0, err
	}
	return school.ID, nil
}

// GetAchievement 获取我的成就
func (srv GDdbService) GetAchievement(userId int64) entity.GDDonationBookRecord {
	return srv.repo.FindBy(repoactivity.FindRecordBy{UserId: userId})
}

func (srv GDdbService) CloseLateTips(userId int64) error {
	record := srv.repo.FindBy(repoactivity.FindRecordBy{UserId: userId})
	if record.ID == 0 {
		return gorm.ErrRecordNotFound
	}
	if record.InviteId == 0 {
		return nil
	}
	record.InviteId = 0
	return srv.repo.Save(&record)
}