package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"os"
)

func main() {
	client := resty.New()

	token := ""
	userid := ""
	baseurl := "https://abc.de"
	user := ""

	body := `{ "username": "` + user + `" }`
	resp, _ := client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Auth-Token", token).
		SetHeader("X-User-Id", userid).
		SetBody(body).
		Post(baseurl + "/api/v1/im.create")

	var result map[string]interface{}
	_ = json.Unmarshal([]byte(resp.String()), &result)
	var rid = result["room"].(map[string]interface{})["_id"].(string)

	msg := ""
	argsWithoutProg := os.Args[1:]

	for i, arg := range argsWithoutProg {
		if i > 0 {
			msg += " "
		}
		msg += arg
	}
	body = `{"message": { "rid": "` + rid + `", "alias": "GoLang", "emoji": ":robot:", "msg": "` + msg + `" }}`

	resp, _ = client.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Auth-Token", token).
		SetHeader("X-User-Id", userid).
		SetBody(body).
		Post(baseurl + "/api/v1/chat.sendMessage")
}
