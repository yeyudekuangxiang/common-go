package service

import (
	"fmt"
	"io"
	"mio/config"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/service/oss"
	"mio/internal/pkg/service/srv_types"
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
	ocrPath, err := oss.DefaultOssService.PutObject(key, reader)
	if err != nil {
		return "", err
	}
	return oss.DefaultOssService.FullUrl(ocrPath), nil
}

func (srv UploadService) UploadImage(openid string, reader io.Reader, filename string, scene string) (string, error) {
	key := fmt.Sprintf("%s%s/%s%s", strings.TrimLeft(scene, oss.DefaultOssService.BasePath), openid, util.UUID(), path.Ext(filename))
	ocrPath, err := oss.DefaultOssService.PutObject(key, reader)
	if err != nil {
		return "", err
	}
	return oss.DefaultOssService.FullUrl(ocrPath), nil
}

//CreateUploadToken operatorId 上传者id operatorType上传者类型 1用户 2管理员 3企业版用户 scene上传场景
func (srv UploadService) CreateUploadToken(operatorId int64, operatorType int8, scene string) (*srv_types.UploadTokenInfo, error) {
	uploadScene, err := DefaultUploadSceneService.FindUploadScene(srv_types.FindSceneParam{
		Scene: scene,
	})
	if err != nil {
		return nil, err
	}
	if uploadScene.ID == 0 {
		return nil, errno.ErrRecordNotFound.With(err)
	}
	if uploadScene.MustLogin && operatorId == 0 {
		return nil, errno.ErrValidation.WithCaller()
	}

	if operatorId != 0 {
		lockKey := fmt.Sprintf("UploadToken%d%d", operatorType, operatorId)
		if !util.DefaultLock.LockNum(lockKey, uploadScene.MaxCount, time.Hour*24) {
			return nil, errno.ErrLimit.WithCaller()
		}
	}

	log, err := DefaultUploadLogService.Create(srv_types.CreateUploadLogParam{
		OssPath:      uploadScene.OssDir,
		OperatorId:   operatorId,
		OperatorType: operatorType,
		SceneId:      uploadScene.ID,
	})
	if err != nil {
		return nil, err
	}

	tokenInfo, err := oss.DefaultOssService.GetPolicyToken(srv_types.GetOssPolicyTokenParam{
		ExpireTime:  time.Minute * 5,
		MaxSize:     uploadScene.MaxSize,
		UploadDir:   uploadScene.OssDir,
		CallbackUrl: util.LinkJoin(config.Config.App.Domain, "/api/mp2c/upload/callback?logId="+log.LogId),
		MimeTypes:   uploadScene.MimeTypes,
		MaxAge:      uploadScene.MaxAge,
	})

	if err != nil {
		return nil, err
	}

	return &srv_types.UploadTokenInfo{
		OssPolicyToken: *tokenInfo,
		MimeTypes:      uploadScene.MimeTypes,
		MaxSize:        uploadScene.MaxSize,
		UploadId:       log.LogId,
		Domain:         config.Config.OSS.CdnDomain,
		MaxAge:         uploadScene.MaxAge,
	}, nil
}
func (srv UploadService) UploadCallback(param srv_types.UploadCallbackParam) error {
	_, err := DefaultUploadLogService.UpdateLog(param.LogId, param.Filename, param.Size)
	if err != nil {
		app.Logger.Error("上传文件回调失败", param, err)
	}
	return err
}
func (srv UploadService) GetUrlByLogId(logId string) (string, error) {
	log, err := DefaultUploadLogService.FindUploadLog(logId)
	if err != nil {
		return "", err
	}
	return log.Url, nil
}

func (srv UploadService) MultipartUploadOcrImage(openid string, reader io.Reader, filename string) (string, error) {
	key := fmt.Sprintf("multpart/%s/%s/%s", openid, util.UUID(), path.Ext(filename))
	ocrPath, err := oss.DefaultOssService.MultipartPutObject(key, reader, filename)
	if err != nil {
		return "", err
	}
	return ocrPath, nil
}

//CreateStsToken operatorId 上传者id operatorType上传者类型 1用户 2管理员 3企业版用户 scene上传场景
func (srv UploadService) CreateStsToken(operatorId int64, operatorType int8, scene string) (*srv_types.OssStsInfo, error) {
	uploadScene, err := DefaultUploadSceneService.FindUploadScene(srv_types.FindSceneParam{
		Scene: scene,
	})
	if err != nil {
		return nil, err
	}
	if uploadScene.ID == 0 {
		return nil, errno.ErrRecordNotFound.With(err)
	}
	if uploadScene.MustLogin && operatorId == 0 {
		return nil, errno.ErrValidation.WithCaller()
	}

	if operatorId != 0 {
		lockKey := fmt.Sprintf("UploadToken%d%d", operatorType, operatorId)
		if !util.DefaultLock.LockNum(lockKey, uploadScene.MaxCount, time.Hour*24) {
			return nil, errno.ErrLimit.WithCaller()
		}
	}

	log, err := DefaultUploadLogService.Create(srv_types.CreateUploadLogParam{
		OssPath:      uploadScene.OssDir,
		OperatorId:   operatorId,
		OperatorType: operatorType,
		SceneId:      uploadScene.ID,
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("路径", path.Join("acs:oss:*:*:miotech-resource", uploadScene.OssDir, "*"))
	cert, err := oss.DefaultOssService.GetSTSToken(srv_types.AssumeRoleParam{
		Scheme:          "https",
		Method:          "POST",
		RoleArn:         "acs:ram::1742387841614768:role/corplinkosscrurole-miotech-resource",
		RoleSessionName: "OssMultipartUpload",
		DurationSeconds: time.Minute * 15,
		Policy: srv_types.StsPolicy{
			Version: "1",
			Statement: []srv_types.Statement{
				{
					Effect: "Allow",
					Action: []string{
						"oss:PutObject",
						"oss:InitiateMultipartUpload",
						"oss:UploadPart",
						"oss:UploadPartCopy",
						"oss:CompleteMultipartUpload",
						"oss:AbortMultipartUpload",
						"oss:ListParts",
						"oss:ListMultipartUploads",
					},
					Resource: []string{
						path.Join("acs:oss:*:*:miotech-resource", uploadScene.OssDir, "*"),
					},
				},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	return &srv_types.OssStsInfo{
		Credentials: *cert,
		Region:      "oss-cn-hongkong",
		/*CallbackBodyUrl:  util.LinkJoin(config.Config.App.Domain, "/api/mp2c/upload/callback?logId="+log.LogId),
		CallbackBody:     "filename=${object}&size=${size}&mimeType=${mimeType}&height=${imageInfo.height}&width=${imageInfo.width}",
		CallbackBodyType: "application/x-www-form-urlencoded",*/
		Bucket:    "miotech-resource",
		MimeTypes: uploadScene.MimeTypes,
		MaxSize:   uploadScene.MaxSize,
		UploadId:  log.LogId,
		Path:      path.Join(uploadScene.OssDir, "/"),
		MaxAge:    uploadScene.MaxAge,
	}, nil
}
