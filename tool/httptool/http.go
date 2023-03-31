package httptool

import (
	"net/http"
	"net/url"
	"time"
)

var DefaultHttpClient = NewHttpClient(&http.Client{Timeout: time.Second * 10})

func SetLogger(logger Logger) {
	DefaultHttpClient.logger = logger
}
func PutJson(url string, data interface{}, options ...RequestOption) ([]byte, error) {
	return DefaultHttpClient.PutJson(url, data, options...)
}
func PostJson(url string, data interface{}, options ...RequestOption) ([]byte, error) {
	return DefaultHttpClient.PostJson(url, data, options...)
}
func PostJsonBytes(url string, data []byte, options ...RequestOption) ([]byte, error) {
	return DefaultHttpClient.PostJsonBytes(url, data, options...)
}
func PostMapFrom(url string, data map[string]string, options ...RequestOption) ([]byte, error) {
	return DefaultHttpClient.PostMapFrom(url, data, options...)
}
func PostFrom(url string, data url.Values, options ...RequestOption) ([]byte, error) {
	return DefaultHttpClient.PostFrom(url, data, options...)
}
func Get(url string, options ...RequestOption) ([]byte, error) {
	return DefaultHttpClient.Get(url, options...)
}
func OriginJson(url string, method string, data []byte, options ...RequestOption) (*HttpResult, error) {
	return DefaultHttpClient.OriginJson(url, method, data, options...)
}
func OriginForm(url string, method string, data []byte, options ...RequestOption) (*HttpResult, error) {
	return DefaultHttpClient.OriginForm(url, method, data, options...)
}
func OriginGet(url string, options ...RequestOption) (*HttpResult, error) {
	return DefaultHttpClient.OriginGet(url, options...)
}
