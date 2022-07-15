package service

import (
	"fmt"
	"io"
	"mio/config"
	"mio/internal/pkg/service/service_types"
	"mio/internal/pkg/util"
	"mio/pkg/errno"
	"path"
	"strings"
	"time"
)

var DefaultUploadService = UploadService{}

type UploadService struct {
}

func (srv UploadService) UploadOcrImage(openid string, reader io.Reader, filename string, collectType PointCollectType) (string, error) {
	key := fmt.Sprintf("ocr/%s/%s/%s%s", strings.ToLower(string(collectType)), openid, util.UUID(), path.Ext(filename))
	imgUrl, err := DefaultOssService.PutObject(key, reader)
	return imgUrl, err
}
func (srv UploadService) CreateUploadToken(userId int64, scene string) (*service_types.UploadTokenInfo, error) {
	uploadScene, err := DefaultUploadSceneService.FindUploadScene(service_types.FindSceneParam{
		Scene: scene,
	})
	if err != nil {
		return nil, err
	}
	if uploadScene.ID == 0 {
		return nil, errno.ErrRecordNotFound.With(err)
	}
	if uploadScene.MustLogin && userId == 0 {
		return nil, errno.ErrValidation.WithCaller()
	}

	lockKey := fmt.Sprintf("UploadToken%d", userId)
	if !util.DefaultLock.LockNum(lockKey, uploadScene.MaxCount, time.Hour*24) {
		return nil, errno.ErrLimit.WithCaller()
	}

	log, err := DefaultUploadLogService.Create(service_types.CreateUploadLogParam{
		OssPath: uploadScene.OssDir,
		UserId:  userId,
		SceneId: uploadScene.ID,
	})
	if err != nil {
		return nil, err
	}

	tokenInfo, err := DefaultOssService.GetPolicyToken(service_types.GetOssPolicyTokenParam{
		ExpireTime:  time.Minute * 5,
		MaxSize:     uploadScene.MaxSize,
		UploadDir:   uploadScene.OssDir,
		CallbackUrl: config.Config.App.Domain + "/api/mp2c/upload/callback",
		MimeTypes:   uploadScene.MimeTypes,
		MaxAge:      uploadScene.MaxAge,
	})

	if err != nil {
		return nil, err
	}

	return &service_types.UploadTokenInfo{
		OssPolicyToken: *tokenInfo,
		MimeTypes:      uploadScene.MimeTypes,
		MaxSize:        uploadScene.MaxSize,
		UploadId:       log.LogId,
		Domain:         util.LinkJoin(config.Config.OSS.CdnDomain, uploadScene.OssDir),
		MaxAge:         uploadScene.MaxAge,
	}, nil
}
func (srv UploadService) UploadCallback(param service_types.UploadCallbackParam) error {
	_, err := DefaultUploadLogService.UpdateLog(param.LogId, param.Filename, param.Size)
	return err
}
func (srv UploadService) GetUrlByLogId(logId string) (string, error) {
	log, err := DefaultUploadLogService.FindUploadLog(logId)
	if err != nil {
		return "", err
	}
	return log.Url, nil
}
