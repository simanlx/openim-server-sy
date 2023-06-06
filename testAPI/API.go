package testAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PostAPI(url string, construct func() []byte) (int, int, error) {
	// 创建一个http请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(construct()))
	if err != nil {
		return 0, 0, err
	}
	// 读取请求返回
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	// 打印http返回内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}
	httpcode := resp.StatusCode

	type resptem struct {
		ErrCode int `json:"errCode"`
	}
	res := resptem{}
	err = json.Unmarshal(body, &res)
	fmt.Println("收到消息", string(body))
	if err != nil {
		return 0, 0, err
	}
	return httpcode, res.ErrCode, nil
}
