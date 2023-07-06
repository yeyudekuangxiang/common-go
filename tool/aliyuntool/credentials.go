package aliyuntool

import (
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

var ErrCredentialNotFound = errors.New("no credential found")

func NewCredential(config *credentials.Config) (credentials.Credential, error) {
	cred, err := credentials.NewCredential(config)
	if err != nil {
		if err.Error() == "No credential found" {
			return nil, ErrCredentialNotFound
		}
		return nil, err
	}
	return cred, nil
}
func UpdateAliYunStsTokenFunc(config *credentials.Config) (func() (accessKeyID, accessKeySecret, securityToken string, expireTime time.Time, err error), error) {
	cred, err := NewCredential(config)
	if err != nil {
		return nil, err
	}
	return func() (accessKeyID, accessKeySecret, securityToken string, expireTime time.Time, err error) {
		id, err := cred.GetAccessKeyId()
		if err != nil {
			return "", "", "", time.Time{}, err
		}
		secret, err := cred.GetAccessKeySecret()
		if err != nil {
			return "", "", "", time.Time{}, err
		}
		t, err := cred.GetSecurityToken()
		if err != nil {
			return "", "", "", time.Time{}, err
		}
		return *id, *secret, *t, time.Now().Add(time.Minute * 60), nil
	}, nil
}
func NewOssCredentialProvider(config *credentials.Config) (*OssCredentialProvider, error) {
	cred, err := NewOssCredential(config)
	if err != nil {
		return nil, err
	}
	return &OssCredentialProvider{
		cred: cred,
	}, nil
}

type OssCredentialProvider struct {
	cred oss.Credentials
}

func (c OssCredentialProvider) GetCredentials() oss.Credentials {
	return c.cred
}

func NewOssCredential(config *credentials.Config) (*OssCredential, error) {
	cred, err := NewCredential(config)
	if err != nil {
		return nil, err
	}
	return &OssCredential{
		cred: cred,
	}, nil
}

type OssCredential struct {
	cred credentials.Credential
}

func (k *OssCredential) GetAccessKeyID() string {
	id, err := k.cred.GetAccessKeyId()
	if err != nil {
		logx.Errorf("获取AccessKeyId异常", err)
		return ""
	}
	return *id
}
func (k *OssCredential) GetAccessKeySecret() string {
	secret, err := k.cred.GetAccessKeySecret()
	if err != nil {
		logx.Errorf("获取AccessKeySecret异常", err)
		return ""
	}
	return *secret
}
func (k *OssCredential) GetSecurityToken() string {
	t, err := k.cred.GetSecurityToken()
	if err != nil {
		logx.Errorf("获取SecurityToken异常", err)
		return ""
	}
	return *t
}
