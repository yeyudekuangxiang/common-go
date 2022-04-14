package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

var DefaultHttp = HttpClient{Timeout: 10 * time.Second}

type HttpClient struct {
	Timeout time.Duration
}

func (c HttpClient) PutJson(url string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return c.PostJsonBytes(url, body)
}
func (c HttpClient) PostJson(url string, data interface{}) ([]byte, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	client := http.Client{
		Timeout: c.Timeout,
	}

	res, err := client.Do(req)

	if res.StatusCode != 200 {
		return nil, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
func (c HttpClient) PostJsonBytes(url string, data []byte) ([]byte, error) {

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	client := http.Client{
		Timeout: c.Timeout,
	}

	res, err := client.Do(req)

	if res.StatusCode != 200 {
		return nil, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
func (c HttpClient) PostMapFrom(url string, data map[string]string) ([]byte, error) {
	body := c.encode(data)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: c.Timeout,
	}

	res, err := client.Do(req)

	if res.StatusCode != 200 {
		return nil, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
func (c HttpClient) PostFrom(url string, data url.Values) ([]byte, error) {
	body := data.Encode()

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: c.Timeout,
	}

	res, err := client.Do(req)

	if res.StatusCode != 200 {
		return nil, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
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
func (c HttpClient) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: c.Timeout,
	}

	res, err := client.Do(req)

	if res.StatusCode != 200 {
		return nil, errors.New("status:" + res.Status)
	}

	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}
