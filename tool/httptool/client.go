package httptool

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type HttpClient struct {
	client *http.Client
	logger Logger
}

func NewHttpClient(client *http.Client, opts ...Option) *HttpClient {
	c := &HttpClient{client: client, logger: NewConsoleLogger()}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type HttpResult struct {
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
	start := time.Now()

	req, err := c.newRequest(method, url, data, start)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	return c.do(start, data, req, options...)
}
func (c HttpClient) OriginForm(url string, method string, data []byte, options ...RequestOption) (*HttpResult, error) {
	start := time.Now()
	req, err := c.newRequest(method, url, data, start)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	return c.do(start, data, req, options...)
}
func (c HttpClient) OriginGet(url string, options ...RequestOption) (*HttpResult, error) {
	start := time.Now()
	req, err := c.newRequest("GET", url, nil, start)
	if err != nil {
		return nil, err
	}

	result, err := c.do(start, nil, req, options...)

	return result, err
}
func (c HttpClient) do(start time.Time, data []byte, req *http.Request, options ...RequestOption) (*HttpResult, error) {
	for _, op := range options {
		op(req)
	}

	res, err := c.client.Do(req)
	if err != nil {
		err = errors.WithStack(err)
		c.log(LogData{
			Url:          req.URL.String(),
			Header:       req.Header,
			Start:        start,
			Duration:     time.Now().Sub(start).Milliseconds(),
			Method:       req.Method,
			StatusCode:   nil,
			RequestBody:  data,
			ResponseBody: nil,
		}, err)
		return nil, err
	}

	if res.StatusCode != 200 {
		c.log(LogData{
			Url:            req.URL.String(),
			Header:         req.Header,
			ResponseHeader: res.Header,
			Start:          start,
			Duration:       time.Now().Sub(start).Milliseconds(),
			Method:         req.Method,
			StatusCode:     &res.StatusCode,
			RequestBody:    data,
			ResponseBody:   nil,
		}, nil)
		return &HttpResult{Response: res}, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		err = errors.WithStack(err)
		c.log(LogData{
			Url:            req.URL.String(),
			Header:         req.Header,
			ResponseHeader: res.Header,
			Start:          start,
			Duration:       time.Now().Sub(start).Milliseconds(),
			Method:         req.Method,
			StatusCode:     &res.StatusCode,
			RequestBody:    data,
			ResponseBody:   nil,
		}, err)
		return &HttpResult{Response: res}, err
	}
	c.log(LogData{
		Url:            req.URL.String(),
		Header:         req.Header,
		ResponseHeader: res.Header,
		Start:          start,
		Duration:       time.Now().Sub(start).Milliseconds(),
		Method:         req.Method,
		StatusCode:     &res.StatusCode,
		RequestBody:    data,
		ResponseBody:   body,
	}, nil)
	return &HttpResult{Response: res, Body: body}, nil
}
func (c HttpClient) encode(data map[string]string) string {
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
func (c HttpClient) newRequest(method string, url string, data []byte, start time.Time) (*http.Request, error) {
	var reader io.Reader
	if data != nil {
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		err = errors.WithStack(err)
		c.log(LogData{
			Url:         url,
			Start:       start,
			Duration:    time.Now().Sub(start).Milliseconds(),
			Method:      method,
			RequestBody: data,
		}, err)
		return nil, err
	}

	return req, nil
}
func (c HttpClient) log(data LogData, err error) {
	if c.logger == nil {
		return
	}
	c.logger.Log(data, err)
}
