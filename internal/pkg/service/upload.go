package service

import (
	"fmt"
	"io"
	"mio/internal/pkg/util"
	"path"
	"strings"
)

var DefaultUploadService = UploadService{}

type UploadService struct {
}

func (srv UploadService) UploadOcrImage(openid string, reader io.Reader, filename string, collectType PointCollectType) (string, error) {
	key := fmt.Sprintf("ocr/%s/%s/%s%s", strings.ToLower(string(collectType)), openid, util.UUID(), path.Ext(filename))
	imgUrl, err := DefaultOssService.PutObject(key, reader)
	return imgUrl, err
}
