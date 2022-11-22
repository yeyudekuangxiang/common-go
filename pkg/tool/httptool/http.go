package httptool

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) *HttpClient {
	return &HttpClient{client: client}
}

type HttpResult struct {
	Err      error
	Response *http.Response
	Body     []byte
}
type RequestOption func(req *http.Request)

func HttpWithHeader(key, value string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

func (c HttpClient) PutJson(url string, data interface{}, options ...RequestOption) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	result, err := c.OriginJson(url, "PUT", body, options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}
func (c HttpClient) PostJson(url string, data interface{}, options ...RequestOption) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.PostJsonBytes(url, body, options...)
}
func (c HttpClient) PostJsonBytes(url string, data []byte, options ...RequestOption) ([]byte, error) {
	result, err := c.OriginJson(url, "POST", data, options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}
func (c HttpClient) PostMapFrom(url string, data map[string]string, options ...RequestOption) ([]byte, error) {
	body := c.encode(data)
	result, err := c.OriginForm(url, "POST", []byte(body), options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}
func (c HttpClient) PostFrom(url string, data url.Values, options ...RequestOption) ([]byte, error) {
	body := data.Encode()

	result, err := c.OriginForm(url, "POST", []byte(body), options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}
func (c HttpClient) Get(url string, options ...RequestOption) ([]byte, error) {
	result, err := c.OriginGet(url, options...)
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func (c HttpClient) OriginJson(url string, method string, data []byte, options ...RequestOption) (*HttpResult, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("content-type", "application/json")

	for _, op := range options {
		op(req)
	}

	res, err := c.client.Do(req)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode != 200 {
		return &HttpResult{Response: res, Err: errors.New("status:" + res.Status)}, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &HttpResult{Response: res, Err: err}, errors.WithStack(err)
	}

	return &HttpResult{
		Response: res,
		Body:     body,
	}, nil
}
func (c HttpClient) OriginForm(url string, method string, data []byte, options ...RequestOption) (*HttpResult, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	for _, op := range options {
		op(req)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode != 200 {
		return &HttpResult{Response: res, Err: errors.New("status:" + res.Status)}, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &HttpResult{Response: res, Err: err}, errors.WithStack(err)
	}
	return &HttpResult{Response: res, Body: body}, nil
}
func (c HttpClient) OriginGet(url string, options ...RequestOption) (*HttpResult, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, op := range options {
		op(req)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.StatusCode != 200 {
		return &HttpResult{Response: res, Err: errors.New("status:" + res.Status)}, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &HttpResult{Response: res, Err: err}, errors.New("status:" + res.Status)
	}
	return &HttpResult{Response: res, Body: body}, nil
}
func (c HttpClient) encode(data map[string]string, options ...RequestOption) string {
	var buf strings.Builder
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(data[k]))
	}
	return buf.String()
}

var DefaultHttpClient = NewHttpClient(&http.Client{Timeout: time.Second * 10})

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
