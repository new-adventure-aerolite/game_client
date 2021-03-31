package client

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/interact"
	"github.com/mitchellh/go-homedir"
	"github.com/new-adventure-aerolite/game-client/pkg/auth"
	"github.com/pkg/browser"
	"gopkg.in/ini.v1"
)

const ConfigFile = ".game/config"

func GetToken() string {
	// homePath := os.Getenv("HOME")
	homePath, err := homedir.Dir()
	if err != nil {
		color.Error.Printf("Fail to read home dir: %v", err)
		os.Exit(1)
	}
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

func ShowFightInfo(heroblood, bossblood, currentlevel, score int32) {
	color.Info.Printf("%30s %d\n", "      Hero Blood  ", heroblood)
	color.FgLightMagenta.Printf("%30s %d\n", "      Boss Blood  ", bossblood)
	color.Info.Printf("%30s %v\n", "      Current Level  ", currentlevel)
	color.Info.Printf("%30s %d\n", "      Score  ", score)
}

func ShowSessionInfo(session, hero, boss map[string]interface{}) {
	if hero != nil {
		color.Info.Printf("%30s %s\n", "  Hero: ", hero["name"])
		color.Info.Printf("%30s %s\n", "  Description: ", hero["details"])

		attachpower := hero["attack_power"].(float64)
		defensepower := hero["defense_power"].(float64)
		blood := hero["blood"].(float64)

		var attachpower1 int = int(attachpower)
		var defensepower1 int = int(defensepower)
		var blood1 int = int(blood)

		color.Info.Printf("%30s %d\n", fmt.Sprintf("%s Attach Power: ", hero["name"]), attachpower1)
		color.Info.Printf("%30s %d\n", fmt.Sprintf("%s Defense Power: ", hero["name"]), defensepower1)
		color.Info.Printf("%30s %d\n", fmt.Sprintf("%s Blood: ", hero["name"]), blood1)

		liveheroblood := session["live_hero_blood"].(float64)
		livebossblood := session["live_boss_blood"].(float64)
		currentlevel := session["current_level"].(float64)
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
	fmt.Println("----------------------------------------------------------")
	if boss != nil {
		color.FgLightMagenta.Printf("%30s %s\n", "  Boss: ", boss["name"])
		color.FgLightMagenta.Printf("%30s %s\n", "  Description: ", boss["details"])
		attachpower := boss["attack_power"].(float64)
		defensepower := boss["defense_power"].(float64)
		blood := boss["blood"].(float64)
		var attachpower1 int = int(attachpower)
		var defensepower1 int = int(defensepower)
		var blood1 int = int(blood)
		color.FgLightMagenta.Printf("%30s %d\n", fmt.Sprintf("%s Attach Power: ", boss["name"]), attachpower1)
		color.FgLightMagenta.Printf("%30s %d\n", fmt.Sprintf("%s Defense Power: ", boss["name"]), defensepower1)
		color.FgLightMagenta.Printf("%30s %d\n", fmt.Sprintf("%s Blood: ", boss["name"]), blood1)

	}
}

func Start() {
	var err error
	var statusCode int
	var heroName string = ""
	var bossName string = ""
	var currentLevel float64 = 0
	var sessionView *auth.SessionViewResponse = nil
	reader := bufio.NewReader(os.Stdin)
Reset:
	token := GetToken()
	if token == "" {
		color.Green.Printf("Please Open %s in your browser\n", auth.Url+"/login")
		browser.OpenURL(auth.Url + "/login")
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
		color.Green.Println("Login Successfully")
	}

	method := interact.SelectOne(
		"You want to play game from beginning or continue from last time?",
		map[string]string{"1": "New Start", "2": "Continue"},
		"1",
		false,
	)
	color.Info.Println("Your select is:", method)
	if method != "Continue" {
		err = auth.ClearSession(token)
		if err != nil {
			color.Error.Println("An error occured while clear session !!!!", err)
			return
		}
	}
	sessionView, err = auth.LoadSession(token)
	if sessionView == nil || err != nil {
		color.Error.Println("An error occured while load session !!!!", err)
		return
	}

	if len(sessionView.Hero) == 0 || sessionView.Hero["name"] == nil {
		// calling  get hero list api
		heros, err := auth.RequestHeros(token)
		if err != nil {
			color.Error.Println("An error occured while load heros!!!!", err)
			return
		}
		fmt.Println("----------------------------------------------------------")
		color.Info.Println("You can choose one Hero from below list:")

		var heroNameList = make([]string, len(heros))
		for index, hero := range heros {
			heroNameList[index] = hero.Name
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
		if setHero != nil && err != nil {
			color.Error.Println("An error occured while set hero !!!!", err)
			return
		}
		ShowSessionInfo(setHero.Session, setHero.Hero, sessionView.Boss)
		heroName = setHero.Hero["name"].(string)
	} else {
		ShowSessionInfo(sessionView.Session, sessionView.Hero, sessionView.Boss)
		heroName = sessionView.Hero["name"].(string)
	}
	bossName = sessionView.Boss["name"].(string)
	currentLevel = sessionView.Session["current_level"].(float64)
	var currentlevel1 int32 = int32(currentLevel)
	fmt.Println("----------------------------------------------------------")
	for {
		// var prefightResp *auth.FightResponse = nil
		action := interact.SingleSelect(
			"Your action?",
			map[string]string{"1": "Fight", "2": "Save", "3": "Reset", "4": "Quit"},
			"1",
			false,
		)
		switch action {
		case "Fight":
			fmt.Println("----------------------------------------------------------")
			color.Info.Printf("  %s(Hero) vs %s(Boss) \n", heroName, bossName)
			fightResp, err := auth.DoFight(token)
			if err != nil {
				color.Error.Println("An error occured while fight !!!!", err)
				return
			}

			if fightResp.GameOver || fightResp.HeroBlood == 0 {
				ShowFightInfo(fightResp.HeroBlood, fightResp.BossBlood, currentlevel1, fightResp.Score)
				color.Info.Println("  Game Over")
				msg, err := auth.QuitSession(token)
				if err != nil {
					color.Error.Println("An error occured while quit session !!!!", err)
				}
				color.Info.Println(msg)
				return
			}
			if fightResp.NextLevel || fightResp.BossBlood == 0 {
				if currentlevel1 >= 4 {
					ShowFightInfo(fightResp.HeroBlood, fightResp.BossBlood, currentlevel1, fightResp.Score)
					// game over
					color.Info.Printf("  Congratulations, %s Win the Game, ByeBye\n", heroName)
					err := auth.ClearSession(token)
					if err != nil {
						color.Error.Println("An error occured while clear session !!!!", err)
					}
					return
				}
				nextLevelResp, err := auth.NextLevel(token)
				if err != nil {
					color.Error.Println("An error occured while goes into next level !!!!", err)
					return
				}

				if nextLevelResp != nil && nextLevelResp.Passed {
					color.Info.Printf("  Congratulations, %s Win the Game, ByeBye\n", heroName)
					err := auth.ClearSession(token)
					if err != nil {
						color.Error.Println("An error occured while clear session !!!!", err)
					}
					return
				}
				currentlevel1 += 1
				// here LoadSession in order to get boss info when hero goes into next level
				session1, err := auth.LoadSession(token)
				if session1 == nil || err != nil {
					color.Error.Println("An error occured while load seesion !!!!", err)
					return
				}
				// update bossName
				bossName = session1.Boss["name"].(string)
			}
			ShowFightInfo(fightResp.HeroBlood, fightResp.BossBlood, currentlevel1, fightResp.Score)
		case "Save":
			sessionid, err := auth.ArchiveSession(token)
			if err != nil {
				color.Error.Println("An error occured while save session !!!!", err)
			}
			color.Info.Println(fmt.Sprintf("SessionID %s is saved", sessionid))
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
