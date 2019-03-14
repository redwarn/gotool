package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ResData struct {
	StatusCode int
	Body       []byte
}

func request(req *http.Request, timeout int64, token string) ResData {
	res := ResData{StatusCode: 0, Body: nil}
	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = time.Duration(timeout) * time.Second
	}
	setRequestHeader(req, token)
	resp, err := client.Do(req)
	if err != nil {
		res.Body = []byte(fmt.Sprintf("execute request error: %s", err.Error()))
		return res
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Body = []byte(fmt.Sprintf("read response body error: %s", err.Error()))
		return res
	}
	res.StatusCode = resp.StatusCode
	res.Body = body
	return res
}

func setRequestHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", token)
}

func requestError(err error) ResData {
	errorMessage := fmt.Sprintf("new requester error: %s", err.Error())
	return ResData{0, []byte(errorMessage)}
}

func Get(url string, timeout int64) ResData {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return requestError(err)
	}
	return request(req, timeout, "")
}

func Request(url, method, body, contentType, token string) ResData {

	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest(strings.ToUpper(method), url, buf)
	if err != nil {
		return requestError(err)
	}
	req.Header.Set("Content-type", contentType)
	return request(req, 60, token)
}
