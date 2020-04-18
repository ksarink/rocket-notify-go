package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/user"
	"strings"
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

func GetPipeInput() string {
	info, err := os.Stdin.Stat()
	if err != nil {
		// no pipe input
		return ""
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		// no pipe input
		return ""
	}

	scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
	var output []string

	for scanner.Scan() {
		output = append(output, scanner.Text())
	}

	return strings.Join(output, "\\n")
}

func SendMessage(restyclient *resty.Client, cfg *Config, rid string) {
	hostname, _ := os.Hostname()
	alias := flag.String("sender", hostname, "The name of the sender (if omitted: hostname)")
	emoji := flag.String("emoji", ":robot:", "The emoji used as avatar (e.g. :robot:)")
	flag.Parse()
	msg := strings.Join(flag.Args(), " ")

	body := `{"message": { "rid": "` + rid +
		`", "alias": "` + *alias +
		`", "emoji": "` + *emoji +
		`", "msg": "` + msg + `" }}`

	_, _ = restyclient.R().
		EnableTrace().
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Auth-Token", cfg.Token).
		SetHeader("X-User-Id", cfg.Userid).
		SetBody(body).
		Post(cfg.Baseurl + "/api/v1/chat.sendMessage")
}

func main() {
	var cfg Config
	LoadConfig(&cfg)

	client := resty.New()

	username := GetTargetUser()
	rid := GetRoomId(client, &cfg, username)
	SendMessage(client, &cfg, rid)
}
