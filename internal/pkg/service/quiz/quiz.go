package quiz

import (
	"fmt"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/timetool"
	"gorm.io/gorm"
	"math/rand"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/context"
	"mio/internal/pkg/model/entity"
	carbonProducer "mio/internal/pkg/queue/producer/carbon"
	"mio/internal/pkg/queue/producer/quizpdr"
	carbonmsg "mio/internal/pkg/queue/types/message/carbon"
	"mio/internal/pkg/queue/types/message/quizmsg"
	"mio/internal/pkg/service"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"time"
)

const OneDayAnswerNum = 4
const questionToPointRatio = 25

var DefaultQuizService = QuizService{}

type QuizService struct {
}

func (srv QuizService) DailyQuestions(openid string, ActivityChannel string) ([]entity.QuizQuestionV2, error) {
	DefaultQuizSingleRecordService.ClearTodayRecord(openid)
	if ActivityChannel != "" {
		switch ActivityChannel {
		case "dove-low-carbon-chellenge":
			list, err := DefaultQuizQuestionService.GetDailyQuestions(3)
			if err != nil {
				return nil, err
			}

			// 创建一个切片
			Ids := []int64{5678672527225073676, 5678672527225073675, 5678672527225073674, 5678672527225073673}
			// 设置随机种子
			rand.Seed(time.Now().UnixNano())
			// 生成一个随机索引
			randomIndex := rand.Intn(len(Ids))
			// 从切片中获取随机元素
			randomElement := Ids[randomIndex]
			channelList, err := DefaultQuizQuestionService.GetDailyQuestionsById(1, randomElement)
			if err != nil {
				return nil, err
			}
			for _, v2 := range channelList {
				list = append(list, v2)
			}
			// 使用Fisher-Yates算法随机排序切片
			for i := len(list) - 1; i > 0; i-- {
				j := rand.Intn(i + 1)
				list[i], list[j] = list[j], list[i]
			}
			return list, err
		}
	}
	return DefaultQuizQuestionService.GetDailyQuestions(OneDayAnswerNum)
}

// Availability 是否可以答题 true可以 false不可以
func (srv QuizService) Availability(openid string) (bool, error) {
	isAnsweredToday, err := DefaultQuizDailyResultService.IsAnsweredToday(openid, timetool.NowDate())
	if err != nil {
		return false, err
	}
	return !isAnsweredToday, nil
}

func (srv QuizService) AnswerQuestion(openid, questionId, answer string) (*AnswerQuestionResult, error) {
	if !util.DefaultLock.Lock("QuizAnswerQuestion"+openid, time.Second*5) {
		return nil, errno.ErrLimit
	}
	defer util.DefaultLock.UnLock("QuizAnswerQuestion" + openid)

	if err := DefaultQuizSingleRecordService.IsAnswered(openid, questionId, timetool.NowDate()); err != nil {
		return nil, err
	}

	todayAnsweredNum := DefaultQuizSingleRecordService.GetTodayAnswerNum(openid)
	if todayAnsweredNum >= OneDayAnswerNum {
		return nil, errno.ErrCommon.WithMessage("答题数量超出限制")
	}

	question := entity.QuizQuestionV2{}

	err := app.DB.Where("id = ?", questionId).Take(&question).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		panic(err)
	}

	if question.ID == 0 {
		return nil, errno.ErrCommon.WithMessage("题目不存在")
	}

	isRight := question.AnswerStatement == answer
	_, err = DefaultQuizSingleRecordService.CreateSingleRecord(CreateSingleRecordParam{
		OpenId:     openid,
		QuestionId: questionId,
		Correct:    isRight,
	})
	if err != nil {
		return nil, err
	}
	return &AnswerQuestionResult{
		IsMatched:           isRight,
		DetailedDescription: question.DetailedDescription,
		CurrentIndex:        DefaultQuizSingleRecordService.GetTodayAnswerNum(openid),
	}, nil
}
func (srv QuizService) Submit(openId string, uid int64) (int, error) {
	if !util.DefaultLock.Lock(fmt.Sprintf("QUIZ_Ssubmit%s", openId), time.Second*10) {
		return 0, errno.ErrLimit
	}
	defer util.DefaultLock.UnLock(fmt.Sprintf("QUIZ_Ssubmit%s", openId))

	todayResult, err := DefaultQuizDailyResultService.CompleteTodayQuiz(openId, timetool.Now())
	if err != nil {
		return 0, err
	}
	err = DefaultQuizSummaryService.UpdateTodaySummary(UpdateSummaryParam{
		OpenId:           openId,
		TodayCorrectNum:  todayResult.CorrectNum,
		TodayAnsweredNum: todayResult.IncorrectNum + todayResult.CorrectNum,
	})
	if err != nil {
		return 0, err
	}
	answerPoint, err := srv.SendAnswerPoint(openId, todayResult.CorrectNum)
	if err != nil {
		return 0, err
	}
	bizId := util.UUID()
	quizpdr.SendMessage(quizmsg.QuizMessage{
		Uid:              uid,
		OpenId:           openId,
		TodayCorrectNum:  todayResult.CorrectNum,
		TodayAnsweredNum: todayResult.IncorrectNum,
		QuizTime:         time.Now().Unix(),
		BizId:            bizId,
	})

	//投递mq
	if err = carbonProducer.ChangeSuccessToQueue(carbonmsg.CarbonChangeSuccess{
		Openid:        openId,
		UserId:        uid,
		TransactionId: bizId,
		Type:          string(entity.POINT_QUIZ),
		City:          "",
		Value:         0,
		Info:          fmt.Sprintf("%+v", todayResult),
	}); err != nil {
		app.Logger.Errorf("ChangeSuccessToQueue 投递失败:%v", err)
	}

	return answerPoint, nil
}
func (srv QuizService) SendAnswerPoint(openId string, correctNum int) (int, error) {
	if correctNum > OneDayAnswerNum {
		correctNum = OneDayAnswerNum
	}

	point := correctNum * questionToPointRatio
	_, err := service.NewPointService(context.NewMioContext()).IncUserPoint(srv_types.IncUserPointDTO{
		OpenId:      openId,
		Type:        entity.POINT_QUIZ,
		BizId:       util.UUID(),
		ChangePoint: int64(point),
	})
	return point, err
}
func (srv QuizService) DailyResult(openId string) (*srv_types.QuizDailyResult, error) {
	result, err := DefaultQuizDailyResultService.FindTodayResult(openId, timetool.NowDate())
	if err != nil {
		return nil, err
	}
	qdResult := srv_types.QuizDailyResult{
		QuizDailyResult: *result,
		Point:           result.CorrectNum * questionToPointRatio,
	}
	switch result.CorrectNum {
	case 0:
		qdResult.Text = []string{"妥妥的一只零碳萌新~", "还需努力哟！"}
	case 1:
		qdResult.Text = []string{"零碳见习生你好呀~", "要再接再厉嗷~"}
	case 2:
		qdResult.Text = []string{"你已达到零碳高中生水平！", "加油吧少年~"}
	case 3:
		qdResult.Text = []string{"你已成功晋级零碳大学生！", "距离巅峰只有一步之遥~"}
	case 4:
		qdResult.Text = []string{"你已达到零碳博士水准~", "再创辉煌吧！"}
	}
	return &qdResult, nil
}
func (srv QuizService) GetSummary(openId string) (*entity.QuizSummary, error) {
	return DefaultQuizSummaryService.FindOrCreateSummary(openId)
}
