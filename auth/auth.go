package auth

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	AuthLoginUrl string = "https://authz.eastus.cloudapp.azure.com:5555"
	Url          string = "https://app.eastus.cloudapp.azure.com:8000"
)

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
	fmt.Println(string(resBody))
	return resBody, err
}

func RequestToken(passcode string) (string, error) {
	url := fmt.Sprintf("%s/passcode?passcode=%s", AuthLoginUrl, passcode)
	token, err := SendRequest("GET", url, "", "")
	return string(token), err
}
