package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestHttpGet(t *testing.T) {
	url := "http://httpbin.org/get"
	params := make(map[string][]string)
	params["test"] = []string{"test1", "test2"}
	body, err := HttpGet(url, params)
	if err != nil {
		t.Error(err)
	}
	t.Log(body)
}
func TestHttpPost(t *testing.T) {
	url := "http://httpbin.org/post"
	body := fmt.Sprintf(`{"params1": "params1", "params2": "params2" }`)
	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &params); err != nil {
		t.Error(err)
	}
	body, err := HttpPost(url, params)
	if err != nil {
		t.Error(err)
	}
	t.Log(body)

}
