package auth

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/new-adventure-aerolite/game-client/pkg/types"
)

const (
	AuthLoginUrl string = "https://authz.eastus.cloudapp.azure.com:5555"
	Url          string = "https://rpg-game.eastus.cloudapp.azure.com"
)

type SessionResp struct {
	Body       []byte
	StatusCode int // e.g. 200
}

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
	Passed  bool
}

type TokenAPIResponse struct {
	Token string `json:"id_token"`
}

type SessionViewResponse struct {
	Hero    map[string]interface{} `json:"hero"`
	Boss    map[string]interface{} `json:"boss"`
	Session map[string]interface{} `json:"session"`
}

func SendRequest(method, url, token, payload string) (*SessionResp, error) {
	var req *http.Request
	var err error

	sessionResp := &SessionResp{}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

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
	sessionResp.StatusCode = res.StatusCode
	sessionResp.Body = resBody
	// fmt.Println(string(resBody))
	return sessionResp, err
}

func RequestToken(passcode string) (string, int, error) {
	url := fmt.Sprintf("%s/passcode?passcode=%s", Url, passcode)
	token, err := SendRequest("GET", url, "", "")
	if err != nil {
		return "", token.StatusCode, err
	}
	tokenJson := TokenAPIResponse{}
	err = json.Unmarshal(token.Body, &tokenJson)
	if err != nil {
		return "", token.StatusCode, err
	}
	return tokenJson.Token, token.StatusCode, err
}

func LoadSession(token string) (*SessionViewResponse, int, error) {
	url := fmt.Sprintf("%s/session", Url)
	resBody, err := SendRequest("GET", url, token, "")

	if err != nil {
		return nil, resBody.StatusCode, err
	}

	object := SessionViewResponse{}
	err = json.Unmarshal(resBody.Body, &object)
	if err != nil {
		return nil, resBody.StatusCode, err
	}

	return &object, resBody.StatusCode, err
}

func RequestHeros(token string) ([]types.Hero, error) {
	url := fmt.Sprintf("%s/heros", Url)
	resBody, err := SendRequest("GET", url, token, "")

	if err != nil {
		return nil, err
	}
	var objects []types.Hero
	err = json.Unmarshal(resBody.Body, &objects)
	if err != nil {
		return nil, err
	}
	return objects, err
}

func SetHero(heroname, token string) (*SessionViewResponse, error) {
	urlstring := fmt.Sprintf("%s/session?hero=%s", Url, url.PathEscape(heroname))
	resBody, err := SendRequest("PUT", urlstring, token, "")

	if err != nil {
		return nil, err
	}

	object := SessionViewResponse{}
	err = json.Unmarshal(resBody.Body, &object)
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
	err = json.Unmarshal(resBody.Body, &object)
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
	err = json.Unmarshal(resBody.Body, &object)
	if err != nil {
		return "", err
	}

	return object.SessionID, err
}

func DoFight(token string) (*FightResponse, error) {
	url := fmt.Sprintf("%s/session/fight", Url)
	resBody, err := SendRequest("PUT", url, token, "")

	if err != nil {
		return nil, err
	}
	object := FightResponse{}
	err = json.Unmarshal(resBody.Body, &object)
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

	if resBody.StatusCode != 200 {
		if strings.Contains(string(resBody.Body), "404") {
			return &NextLevelResponse{
				Passed: true,
			}, nil
		}
	}

	object := NextLevelResponse{}
	err = json.Unmarshal(resBody.Body, &object)
	if err != nil {
		return nil, err
	}

	return &object, err
}

func ClearSession(token string) error {
	url := fmt.Sprintf("%s/session/clear", Url)
	resBody, err := SendRequest("POST", url, token, "")

	if err != nil {
		return err
	}
	object := QuitResponse{}
	err = json.Unmarshal(resBody.Body, &object)
	if err != nil {
		return err
	}

	return err
}
