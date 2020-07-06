package client

import (
	"testing"

	json "github.com/json-iterator/go"
)

func TestCallEndApi(t *testing.T) {
	InitClient(nil)

	params := make(map[string]string)
	params["name"] = "jack"

	data, err := CallEndApi("127.0.0.1:1234", "hello", "world", params)
	if err != nil {
		t.Error(err)
	}

	type Resp struct {
		Success bool        `json:"success"`
		PayLoad interface{} `json:"payload"`
		Err     struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	resp := &Resp{}
	if err = json.Unmarshal(data, resp); err != nil {
		t.Error(err)
	}

	if !resp.Success {
		t.Error(resp.Err.Code, resp.Err.Message)
	}

	t.Log(resp)
}
