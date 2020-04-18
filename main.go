package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/user"
)

type Config struct {
	Baseurl string `yaml:"baseurl"`
	Userid  string `yaml:"userid"`
	Token   string `yaml:"token"`
}

func GetTargetUser() string {
	currentuser, _ := user.Current()
	sudouser := os.Getenv("SUDO_USER")
	if currentuser.Uid == "0" && len(sudouser) > 0 {
		return sudouser
	}
	username := currentuser.Username
	return username
}

func LoadConfig(cfg *Config) {
	f, err := os.Open("/etc/rocket-notify/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	//var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func GetRoomId(restyclient *resty.Client, cfg *Config, username string) string {
	body := `{ "username": "` + username + `" }`
	resp, _ := restyclient.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Auth-Token", cfg.Token).
		SetHeader("X-User-Id", cfg.Userid).
		SetBody(body).
		Post(cfg.Baseurl + "/api/v1/im.create")

	var result map[string]interface{}
	_ = json.Unmarshal([]byte(resp.String()), &result)
	var rid = result["room"].(map[string]interface{})["_id"].(string)
	return rid
}

func SendMessage(restyclient *resty.Client, cfg *Config, rid string, msg string) {
	body := `{"message": { "rid": "` + rid + `", "alias": "GoLang", "emoji": ":robot:", "msg": "` + msg + `" }}`

	_, _ = restyclient.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Auth-Token", cfg.Token).
		SetHeader("X-User-Id", cfg.Userid).
		SetBody(body).
		Post(cfg.Baseurl + "/api/v1/chat.sendMessage")
}

func CreateMessage() string {
	msg := ""
	argsWithoutProg := os.Args[1:]

	for i, arg := range argsWithoutProg {
		if i > 0 {
			msg += " "
		}
		msg += arg
	}
	return msg
}

func main() {
	var cfg Config
	LoadConfig(&cfg)

	client := resty.New()

	username := GetTargetUser()
	rid := GetRoomId(client, &cfg, username)
	msg := CreateMessage()
	SendMessage(client, &cfg, rid, msg)
}
