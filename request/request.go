package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpClient struct {
	*http.Client
}

func NewHttpClient(dur time.Duration) *HttpClient {
	return &HttpClient{
		Client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   dur,
		},
	}
}

var httpDefaultClient = NewHttpClient(time.Second * 5)

func (c *HttpClient) response(r *http.Request) (string, error) {
	rsp, err := c.Do(r)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	strBody := strings.Trim(string(body), "\r\n")
	if rsp.StatusCode != http.StatusOK {
		return strBody, fmt.Errorf("status err:%v", rsp.StatusCode)
	}
	return strBody, nil
}
func (c *HttpClient) httpGet(rawUrl string, data url.Values) (string, error) {
	if len(data) > 0 {
		rawUrl += "?" + data.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, rawUrl, nil)

	if err != nil {
		return "", err
	}
	return c.response(req)
}
func (c *HttpClient) httpPost(rawUrl string, params map[string]interface{}) (string, error) {
	bs, err := json.Marshal(params)
	if err != nil {
		return "", nil
	}
	reader := bytes.NewBuffer(bs)
	req, err := http.NewRequest(http.MethodPost, rawUrl, reader)
	if err != nil {
		return "", err
	}
	req.Header["Content-Type"] = []string{"application/json"}
	return c.response(req)
}
func HttpGet(rawUrl string, data url.Values) (string, error) {
	return httpDefaultClient.httpGet(rawUrl, data)
}
func HttpPost(rawUrl string, params map[string]interface{}) (string, error) {
	return httpDefaultClient.httpPost(rawUrl, params)
}
