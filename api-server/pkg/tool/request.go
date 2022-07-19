package tool

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: time.Second * 3,
}

func HttpGet(url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}
