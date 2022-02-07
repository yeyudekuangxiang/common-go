package util

import (
	"encoding/base64"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/pkg/errors"
	"mio/core/app"
	"strconv"
	"strings"
	"time"
)

func ParseToken(token string, v jwt.Claims) error {
	_, err := jwt.ParseWithClaims(token, v, func(token *jwt.Token) (interface{}, error) {
		tokenKey, err := GetAppConfig("app.TokenKey")
		if err != nil {
			app.Logger.Error(err)
			return "", errors.New("系统错误,请联系管理员")
		}
		return []byte(tokenKey), nil
	})
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

func IsOldToken(token string) bool {
	return strings.Contains(token, "newXToken")
}

type oldToken struct {
	In string `json:"in"`
}

func (oldToken) Valid() error {
	return nil
}
func ParseOldToken(token string, v interface{}) error {
	tokenKey, err := GetAppConfig("app.TokenKey")
	if err != nil {
		app.Logger.Error(err)
		return errors.New("系统错误,请联系管理员")
	}

	token = strings.ReplaceAll(token, "newXToken", "")

	claims := oldToken{}
	_, err = jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenKey), nil
	})
	if err != nil {
		return err
	}

	in := claims.In
	str, err := authCode(in, "DECODE", tokenKey, 0)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), v)
}
func authCode(str string, operation string, key string, expiry int) (string, error) {
	// 动态密匙长度，相同的明文会生成不同密文就是依靠动态密匙
	cKeyLength := 4
	// 密匙
	key = util.Md5(key)
	// 密匙a会参与加解密
	//substr($key, 0, 16)
	keyA := util.Md5(key[:16])
	// 密匙b会用来做数据完整性验证
	keyB := util.Md5(key[16:])
	// 密匙c用于变化生成的密文

	keyC := ""
	if cKeyLength > 0 {
		if operation == "DECODE" {
			keyC = str[0:cKeyLength]
		} else {
			//substr(md5(microtime()), -$cKeyLength)): ''
		}
	}

	// 参与运算的密匙
	cryptKey := keyA + util.Md5(keyA+keyC)
	keyLength := len(cryptKey)
	// 明文，前10位用来保存时间戳，解密时验证数据有效性，10到26位用来保存$keyB(密匙b)，
	// 解密时会通过这个密匙验证数据完整性
	// 如果是解码的话，会从第$ckey_length位开始，因为密文前$ckey_length位保存 动态密匙，以保证解密正确
	if operation == "DECODE" {
		decoded, err := base64.RawStdEncoding.DecodeString(str[cKeyLength:])
		if err != nil {
			return "", err
		}
		str = string(decoded)
	} else {
		//sprintf('%010d', $expiry ? $expiry + time() : 0).substr(md5($string. $keyB), 0, 16).$string
	}

	stringLength := len(str)
	result := ""
	box := make([]uint8, 0)
	for i := 0; i < 256; i++ {
		box = append(box, uint8(i))
	}
	rndKey := make([]uint8, 0)
	// 产生密匙簿
	for i := 0; i <= 255; i++ {
		rndKey = append(rndKey, cryptKey[i%keyLength])
	}

	// 用固定的算法，打乱密匙簿，增加随机性，好像很复杂，实际上对并不会增加密文的强度
	j := 0
	for i := 0; i < 256; i++ {
		j = int(uint8(j)+box[i]+rndKey[i]) % 256
		tmp := box[i]
		box[i] = box[j]
		box[j] = tmp
	}

	var (
		a, i int
	)
	j = 0
	for ; i < stringLength; i++ {
		a = (a + 1) % 256
		j = (j + int(box[a])) % 256
		tmp := box[a]
		box[a] = box[j]
		box[j] = tmp
		result += string(str[i] ^ (box[int(box[a]+box[j])%256]))
	}
	if operation == "DECODE" {
		// 验证数据有效性，请看未加密明文的格式
		rr, _ := strconv.Atoi(result[0:10])
		if (rr == 0 || int64(rr)-time.Now().Unix() > 0) &&
			result[10:26] == util.Md5(result[26:] + keyB)[0:16] {
			return result[26:], nil
		}
		return "", nil
	}
	return keyC + strings.ReplaceAll(base64.StdEncoding.EncodeToString([]byte(result)), "=", ""), nil
}

func CreateToken(v jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, v)
	tokenKey, err := GetAppConfig("app.TokenKey")
	if err != nil {
		app.Logger.Error(err)
		return "", errors.New("系统错误,请联系管理员")
	}
	return token.SignedString([]byte(tokenKey))
}
