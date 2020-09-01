package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	Url "net/url"
	"strings"
	"time"
	"github.com/godblesshugh/form"
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

func Get(url string, timeout int64, token string) ResData {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return requestError(err)
	}
	return request(req, timeout, token)
}

func PostParams(url, token string, body map[string]string) ResData {
	v := Url.Values{}
	for key, value := range body {
		v.Set(key, value)
	}
	return Request(url, "POST", v.Encode(), "application/x-www-form-urlencoded", token)
}

func PostJson(url, token, body string) ResData {
	return Request(url, "POST", body, "application/json", token)
}


func PostForm(url, token string, data map[string]interface{}) ResData {
	formData, err := form.EncodeToString(data)
	if err != nil {
		return requestError(err)
	}
	return Request(url, "POST", formData, "application/x-www-form-urlencoded", token)
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
