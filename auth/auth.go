package auth

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	AuthLoginUrl string = "https://authz.eastus.cloudapp.azure.com:5555"
	Url          string = "https://app.eastus.cloudapp.azure.com:8000"
)

type QuitResponse struct {
	Msg string `json:"msg"`
}

type ArchiveResponse struct {
	Msg       string `json:"msg"`
	SessionID string `json:"session_id"`
}

type FightResponse struct {
	GameOver  bool  `json:"game_over"`
	NextLevel bool  `json:"next_level"`
	Score     int32 `json:"score"`
	HeroBlood int32 `json:"hero_blood"`
	BossBlood int32 `json:"boss_blood"`
}

type NextLevelResponse struct {
	Msg     string                 `json:"msg"`
	Session map[string]interface{} `json:"session"`
}

type TokenAPIResponse struct {
	Token string `json:"id_token"`
}

type SessionViewResponse struct {
	Hero    map[string]interface{} `json:"hero"`
	Boss    map[string]interface{} `json:"boss"`
	Session map[string]interface{} `json:"session"`
}

type Hero interface {
}

type Boss interface {
}

type Session interface {
}

func SendRequest(method, url, token, payload string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	var req *http.Request
	var err error

	if payload != "" {
		reqBody := strings.NewReader(payload)
		req, err = http.NewRequest(method, url, reqBody)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", "bearer "+token)
		req.Header.Add("Content-Type", "application/json")
	}
	// Send HTTP request
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// fmt.Println(string(resBody))
	return resBody, err
}

func RequestToken(passcode string) (string, error) {
	url := fmt.Sprintf("%s/passcode?passcode=%s", AuthLoginUrl, passcode)
	token, err := SendRequest("GET", url, "", "")
	if err != nil {
		return "", err
	}
	tokenJson := TokenAPIResponse{}
	err = json.Unmarshal(token, &tokenJson)
	if err != nil {
		return "", err
	}
	return tokenJson.Token, err
}

func LoadSession(token string) (*SessionViewResponse, error) {
	url := fmt.Sprintf("%s/session", Url)
	resBody, err := SendRequest("GET", url, token, "")

	if err != nil {
		return nil, err
	}

	object := SessionViewResponse{}
	err = json.Unmarshal(resBody, &object)
	if err != nil {
		return nil, err
	}

	return &object, err
}

func RequestHeros(token string) ([]Hero, error) {
	url := fmt.Sprintf("%s/heros", Url)
	resBody, err := SendRequest("GET", url, token, "")
	fmt.Println(string(resBody))
	if err != nil {
		return nil, err
	}
	var objects []Hero
	err = json.Unmarshal(resBody, &objects)
	if err != nil {
		return nil, err
	}
	return objects, err
}

func SetHero(heroname, token string) (*SessionViewResponse, error) {
	url := fmt.Sprintf("%s/session?hero=%s", Url, heroname)
	resBody, err := SendRequest("PUT", url, token, "")

	if err != nil {
		return nil, err
	}

	object := SessionViewResponse{}
	err = json.Unmarshal(resBody, &object)
	if err != nil {
		return nil, err
	}

	return &object, err
}

func QuitSession(token string) (string, error) {
	url := fmt.Sprintf("%s/session/quit", Url)
	resBody, err := SendRequest("POST", url, token, "")

	if err != nil {
		return "", err
	}
	object := QuitResponse{}
	err = json.Unmarshal(resBody, &object)
	if err != nil {
		return "", err
	}

	return object.Msg, err
}

func ArchiveSession(token string) (string, error) {
	url := fmt.Sprintf("%s/session/archive", Url)
	resBody, err := SendRequest("POST", url, token, "")

	if err != nil {
		return "", err
	}
	object := ArchiveResponse{}
	err = json.Unmarshal(resBody, &object)
	if err != nil {
		return "", err
	}

	return object.SessionID, err
}

func Fight(token string) (*FightResponse, error) {
	url := fmt.Sprintf("%s/session/fight", Url)
	resBody, err := SendRequest("PUT", url, token, "")

	if err != nil {
		return nil, err
	}
	object := FightResponse{}
	err = json.Unmarshal(resBody, &object)
	if err != nil {
		return nil, err
	}

	return &object, err
}

func NextLevel(token string) (*NextLevelResponse, error) {
	url := fmt.Sprintf("%s/session/level", Url)
	resBody, err := SendRequest("POST", url, token, "")

	if err != nil {
		return nil, err
	}
	object := NextLevelResponse{}
	err = json.Unmarshal(resBody, &object)
	if err != nil {
		return nil, err
	}

	return &object, err
}
