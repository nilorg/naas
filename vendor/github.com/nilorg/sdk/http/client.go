package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// 默认超时
	defaultTimeout = time.Second * 3
)

// Client Http客户端
type Client struct {
	Header *Header
	client *http.Client
}

// NewClient 创建HttpClient
func NewClient() *Client {
	return &Client{
		Header: new(Header),
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// NewHeaderClient 创建HttpClient
func NewHeaderClient(header *Header) *Client {
	return &Client{
		Header: header,
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// SetHeader 设置 http request header
func (h *Client) SetHeader(header *Header) {
	h.Header = header
}

// SetTimeout 设置超时
func (h *Client) SetTimeout(timeout time.Duration) {
	h.client.Timeout = timeout
}

// SetTransport 设置Transport
func (h *Client) SetTransport(transport http.RoundTripper) {
	h.client.Transport = transport
}

// Get send get request.
func (h *Client) Get(url string, values url.Values) (body []byte, statusCode int, err error) {
	if values != nil {
		url += "?" + values.Encode()
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	return h.do(req)
}

// Post send post request.
func (h *Client) Post(url string, values url.Values) (body []byte, statusCode int, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(values.Encode()))
	if err != nil {
		return
	}
	return h.do(req)
}

// PostData send post request.
func (h *Client) PostData(url string, values []byte) (body []byte, statusCode int, err error) {
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(values))
	if err != nil {
		return
	}
	return h.do(req)
}

// PostJSON send post request.
func (h *Client) PostJSON(url string, jsonObject interface{}) (body []byte, statusCode int, err error) {
	json, err := json.Marshal(jsonObject)
	if err != nil {
		return
	}
	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(json))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	return h.do(req)
}

func (h *Client) do(req *http.Request) (body []byte, statusCode int, err error) {
	response, err := h.Do(req)
	if err != nil {
		return
	}
	statusCode = response.StatusCode
	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)
	return
}

// Do 发送
func (h *Client) Do(req *http.Request) (response *http.Response, err error) {
	// 设置Header
	for k, v := range *h.Header {
		req.Header.Set(k, v)
	}
	response, err = h.client.Do(req)
	return
}
