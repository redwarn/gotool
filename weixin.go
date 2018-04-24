package gotool

import (
	"encoding/json"
	"fmt"

	"github.com/levigross/grequests"
	"github.com/tidwall/gjson"
)

var (
	WeiXinErr = func(code int64, msg string) error {
		return fmt.Errorf("wx return error,code: %d, msg: %s", code, msg)
	}
)

const (
	TEXT     = "text"
	TEXTCARD = "textcard"
)

type (
	Text struct {
		Content string `json:"content"`
	}

	TextCard struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		BtnTxt      string `json:"btntxt"`
	}

	Message struct {
		MsgType  string    `json:"msgtype"`
		ToUser   string    `json:"touser"`
		ToTag    string    `json:"totag"`
		ToParty  string    `json:"toparty"`
		AgentId  int64     `json:"agentid"`
		Safe     int64     `json:"safe"`
		Text     *Text     `json:"text"`
		TextCard *TextCard `json:"textcard"`
	}
	Client struct {
		TokenAPIURL string
		ApiURL      string
		CorpID      string
		CorpSecret  string
		Message     *Message
	}
)

func (c *Client) GetToken() (string, error) {
	params := &grequests.RequestOptions{
		Params: map[string]string{
			"corpid":     c.CorpID,
			"corpsecret": c.CorpSecret,
		},
	}

	data, err := grequests.Get(c.TokenAPIURL, params)
	if err != nil {
		return "", err
	}
	res := data.String()
	code := gjson.Get(res, "errcode")
	token := gjson.Get(res, "access_token")
	if code.Int() == 0 {
		return token.String(), nil
	}
	return "", WeiXinErr(code.Int(), gjson.Get(res, "errmsg").String())
}

func (c *Client) SendMessage() (bool, error) {
	data, err := json.Marshal(c.Message)
	if err != nil {
		return false, err
	}
	token, err := c.GetToken()

	if err != nil {
		return false, err
	}
	params := &grequests.RequestOptions{
		Params: map[string]string{
			"access_token": token,
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		JSON: data,
	}
	res, err := grequests.Post(c.ApiURL, params)
	if err != nil {
		return false, err
	}
	resJson := res.String()
	code := gjson.Get(resJson, "errcode")
	if code.Int() == 0 {
		return true, nil
	}
	return false, WeiXinErr(code.Int(), gjson.Get(resJson, "errmsg").String())
}
