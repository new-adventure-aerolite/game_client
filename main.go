// Copyright (c) 2021 Li Wen and others

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"bufio"
	"fmt"

	"game/auth"
	"os"
	"regexp"
	"strings"

	"github.com/go-ini/ini"
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/interact"
)

type SessionView struct {
	Hero    Hero
	Boss    Boss
	Seesion Session
}

type Boss struct {
	Name         string
	Details      string
	AttackPower  int32
	DefensePower int32
	Blood        int32
	Level        int32
}

type Session struct {
	UID           string
	HeroName      string
	LiveHeroBlood int32
	LiveBossBlood int32
	CurrentLevel  int32
	Score         int32
}

type HeroList struct {
	Heros []Hero
}

type Hero struct {
	Name         string
	Details      string
	AttackPower  int32
	DefensePower int32
	Blood        int32
}

type Fight struct {
	GameOver  bool
	NextLevel bool
	Score     int32
	HeroBlood int32
	BossBlood int32
}

type Level struct {
	Msg     string
	Session string
}

const ConfigFile = ".game/config"

func GetToken() string {
	homePath := os.Getenv("HOME")
	configFilePath := fmt.Sprintf("%s/%s", homePath, ConfigFile)
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		color.Error.Printf("Fail to read config file: %v", err)
		os.Exit(1)
	}
	return cfg.Section("").Key("id_token").String()
}

func SetToken(token string) {
	homePath := os.Getenv("HOME")
	configFilePath := fmt.Sprintf("%s/%s", homePath, ConfigFile)
	cfg, err := ini.Load(configFilePath)
	if err != nil {
		color.Error.Printf("Fail to read config file: %v", err)
		os.Exit(1)
	}

	cfg.Section("").Key("id_token").SetValue(token)
	cfg.SaveTo(configFilePath)

}

func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func VerifyPasscode(passcode string) bool {
	pattern := `\d{6}`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(passcode)
}

func ShowSessionInfo(session, hero map[string]interface{}) {
	color.Info.Printf("%30s %s\n", "  Hero Name: ", hero["name"])
	color.Info.Printf("%30s %s\n", "    Description: ", hero["details"])
	liveheroblood := session["live_hero_blood"].(float64)
	livebossblood := session["live_boss_blood"].(float64)
	currentlevel := session["current_level"].(float64)
	// score := session["score"].(float64)
	var liveheroblood1 int = int(liveheroblood)
	var livebossblood1 int = int(livebossblood)
	var currentlevel1 int = int(currentlevel)
	color.Info.Printf("%30s %d\n", "    Live Hero Blood: ", liveheroblood1)
	color.Info.Printf("%30s %d\n", "    Live Boss Blood: ", livebossblood1)
	color.Info.Printf("%30s %d\n", "    Current Level: ", currentlevel1)
	var score float64 = 0
	if session["score"] == nil {
		score = 0
	} else {
		score = session["score"].(float64)
	}
	var score1 int = int(score)
	color.Info.Printf("%30s %d\n", "    Score: ", score1)
}

func main() {

	var err error
	var statusCode int
	reader := bufio.NewReader(os.Stdin)
	var heroName string = ""
Reset:
	token := GetToken()
	if token == "" {
		color.Green.Printf("Please Open %s in your browser\n", auth.AuthLoginUrl)
		color.Green.Printf("Then copy passcode from your browser to here: ")
		var passcode string
		for {
			passcode, err = reader.ReadString('\n')
			if err != nil {
				color.Error.Println("An error occured while get passcode. Please try again", err)
				return
			}
			// remove the delimeter from the string
			passcode = strings.TrimSuffix(passcode, "\n")
			if passcode != "" && len(passcode) == 6 {
				if !VerifyPasscode(passcode) {
					color.Warn.Println("Your passcode is incorrect, Please try again", err)
				}
			}
			token, statusCode, err = auth.RequestToken(passcode)
			if err != nil {
				color.Error.Println("An error occured while request token !!!!", err)
				return
			}
			if statusCode == 401 {
				color.Error.Println("Unauthorized occured while request token !!!!", statusCode)
			} else {
				break
			}
		}
		SetToken(token)
		color.Green.Printf("Login Successfully")
	}

	method := interact.SelectOne(
		"You want to play game from beginning or continue from last time?",
		map[string]string{"1": "From beginning", "2": "Continue"},
		"1",
		false,
	)
	color.Info.Println("Your select is:", method)
	if method != "Continue" {
		err = auth.ClearSession(token)
		if err != nil {
			color.Error.Println("An error occured while clear session !!!!", err)
		}
	}
	sessionView, err := auth.LoadSession(token)
	if err != nil {
		color.Error.Println("An error occured while load session !!!!", err)
	}

	if len(sessionView.Hero) == 0 || sessionView.Hero["name"] == nil {
		// calling  get hero list api
		heros, err := auth.RequestHeros(token)
		if err != nil {
			color.Error.Println("An error occured while load heros!!!!", err)
		}
		fmt.Println("----------------------------------------------------------")
		color.Info.Println("You can choose one Hero from below list:")

		var heroNameList = make([]string, len(heros))
		for index, hero := range heros {
			mapHero := hero.(map[string]interface{})
			heroNameList[index] = mapHero["name"].(string)
		}

		choosenHero := interact.Choice(
			"Choose Hero?",
			heroNameList,
			"",
			false,
		)
		// call get hero api
		color.Info.Println("Your select is:", choosenHero)
		setHero, err := auth.SetHero(choosenHero, token)
		if err != nil {
			color.Error.Println("An error occured while set hero !!!!", err)
		}
		// color.Info.Println("Description: ", setHero.Hero["details"])
		ShowSessionInfo(setHero.Session, setHero.Hero)
		heroName = setHero.Hero["name"].(string)
	} else {
		ShowSessionInfo(sessionView.Session, sessionView.Hero)
		heroName = sessionView.Hero["name"].(string)
	}
	fmt.Println("----------------------------------------------------------")
	for {
		action := interact.SingleSelect(
			"Your action?",
			map[string]string{"1": "Fight", "2": "Archive", "3": "Reset", "4": "Quit"},
			"1",
			false,
		)

		switch action {
		case "Fight":
			fmt.Println("----------------------------------------------------------")
			color.Info.Printf("  %s %s \n", heroName, action)
			fightResp, err := auth.Fight(token)
			if err != nil {
				color.Error.Println("An error occured while fight !!!!", err)
			}
			color.Info.Printf("%30s %d\n", "      Hero Blood  ", fightResp.HeroBlood)
			color.Info.Printf("%30s %d\n", "      Boss Blood  ", fightResp.BossBlood)
			color.Info.Printf("%30s %v\n", "      Next Level  ", fightResp.NextLevel)
			color.Info.Printf("%30s %d\n", "      Score  ", fightResp.Score)

			if fightResp.GameOver || fightResp.HeroBlood == 0 {
				color.Info.Println("  Game Over")
				msg, err := auth.QuitSession(token)
				if err != nil {
					color.Error.Println("An error occured while quit session !!!!", err)
				}
				color.Info.Println(msg)
				return
			}
			if fightResp.NextLevel || fightResp.BossBlood == 0 {
				session1, err := auth.LoadSession(token)
				if err != nil {
					color.Error.Println("An error occured while load seesion !!!!", err)
				}
				if session1.Session["current_level"].(float64) >= 2 {
					// game over
					color.Info.Println("  Congratulations, You Win the Game, ByeBye")
					err := auth.ClearSession(token)
					if err != nil {
						color.Error.Println("An error occured while clear session !!!!", err)
					}
					return
				}
				nextLevelResp, err := auth.NextLevel(token)
				if err != nil {
					color.Error.Println("An error occured while goes into next level !!!!", err)
				}
				if nextLevelResp != nil && nextLevelResp.Passed {
					color.Info.Println("  Congratulations, You Win the Game, ByeBye")
					err := auth.ClearSession(token)
					if err != nil {
						color.Error.Println("An error occured while clear session !!!!", err)
					}
					return
				}
				color.Info.Println("  You have gone into next level !!!")
			}
		case "Archive":
			sessionid, err := auth.ArchiveSession(token)
			if err != nil {
				color.Error.Println("An error occured while archive session !!!!", err)
			}
			color.Info.Println(fmt.Sprintf("SessionID %s is archived", sessionid))
		case "Reset":
			goto Reset
		case "Quit":
			msg, err := auth.QuitSession(token)
			if err != nil {
				color.Error.Println("An error occured while quit session !!!!", err)
			}
			color.Info.Println(msg)
			return
		default:
			color.Info.Println(action)
		}
	}
}
