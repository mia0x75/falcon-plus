package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	httpDefaultTimeout = 5 // 5秒超时
)

var tr = &http.Transport{
	MaxIdleConnsPerHost: 256,
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
	DisableCompression:  false,
	DisableKeepAlives:   false,
	TLSHandshakeTimeout: httpDefaultTimeout * time.Second,
	Dial: func(netw, addr string) (net.Conn, error) {
		deadline := time.Now().Add(httpDefaultTimeout * time.Second)
		dial := net.Dialer{
			Timeout:   httpDefaultTimeout * time.Second,
			KeepAlive: 86400 * time.Second,
		}
		conn, err := dial.Dial(netw, addr)
		if err != nil {
			return conn, err
		}
		conn.SetDeadline(deadline)
		return conn, nil
	},
}

// timeout
var defaultHttpClient = http.Client{
	Transport: tr,
}

type Http struct {
	url       string
	userAgent string
	headers   map[string]string
}

func NewHttp(url string) *Http {
	request := &Http{
		url: url,
	}
	return request
}

func Post(addr string, postData []byte) ([]byte, error) {
	c := NewHttp(addr)
	return c.Post(postData)
}

func Get(addr string) ([]byte, error) {
	c := NewHttp(addr)
	return c.Get()
}

func (req *Http) request(method string, url string, data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	r, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	if req.headers != nil {
		for k, v := range req.headers {
			r.Header.Set(k, v)
		}
	} else {
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Connection", "keep-alive")
	}
	resp, err := defaultHttpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("error http status：%d", resp.StatusCode))
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (req *Http) Put(data []byte) ([]byte, error) {
	return req.request("PUT", req.url, data)
}

func (req *Http) Post(data []byte) ([]byte, error) {
	return req.request("POST", req.url, data)
}

func (req *Http) Get() ([]byte, error) {
	return req.request("GET", req.url, nil)
}

func (req *Http) Delete() ([]byte, error) {
	return req.request("DELETE", req.url, nil)
}

func (http *Http) SetUserAgent(userAgent string) {
	http.userAgent = userAgent
}

func (http *Http) SetHeaders(headers map[string]string) {
	http.headers = headers
}
