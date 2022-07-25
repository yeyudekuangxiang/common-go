package service

import (
	"mio/config"
	"mio/internal/pkg/model/entity"
	"mio/internal/pkg/repository"
	"mio/internal/pkg/service/srv_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
)

var DefaultUploadLogService = UploadLogService{repo: repository.DefaultUploadLogRepository}

type UploadLogService struct {
	repo repository.UploadLogRepository
}

func (srv UploadLogService) Create(param srv_types.CreateUploadLogParam) (*entity.UploadLog, error) {
	log := entity.UploadLog{
		LogId:        util.UUID(),
		OssPath:      param.OssPath,
		OperatorId:   param.OperatorId,
		OperatorType: param.OperatorType,
		SceneId:      param.SceneId,
	}
	return &log, srv.repo.Create(&log)
}
func (srv UploadLogService) UpdateLog(logId, filename string, size int64) (*entity.UploadLog, error) {
	log, err := srv.repo.FindLog(repository.FindLogBy{
		LogId: logId,
	})
	if err != nil {
		return nil, err
	}
	if log.ID == 0 {
		return nil, errno.ErrRecordNotFound.WithCaller()
	}
	log.Size = size
	log.Url = util.LinkJoin(config.Config.OSS.CdnDomain, log.OssPath, filename)
	return log, srv.repo.Save(log)
}
func (srv UploadLogService) FindUploadLog(logId string) (*entity.UploadLog, error) {
	return srv.repo.FindLog(repository.FindLogBy{
		LogId: logId,
	})
}
