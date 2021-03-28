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

func main() {
	var emailAddress = ""
	var err error
	reader := bufio.NewReader(os.Stdin)
	for {
		color.Green.Printf("Please Input your email address: ")
		// ReadString will block until the delimiter is entered
		emailAddress, err = reader.ReadString('\n')
		if err != nil {
			color.Error.Println("An error occured while reading input. Please try again", err)
			return
		}
		// remove the delimeter from the string
		emailAddress = strings.TrimSuffix(emailAddress, "\n")
		if emailAddress != "" && len(emailAddress) != 0 {
			// fmt.Println(emailAddress)
			if VerifyEmailFormat(emailAddress) {
				break
			}

		}
		color.Warn.Println("Your email address is empty or incorrect, Please try again", err)
	}

	token := GetToken()
	if token == "" {
		token, err = auth.RequestToken("116407")
		if err != nil {
			color.Error.Println("An error occured while request token !!!!", err)
		}
		color.Info.Println(token)
		SetToken(token)
	}
	sessionView, err := auth.LoadSession(token)
	if err != nil {
		color.Error.Println("An error occured while load session !!!!", err)
	}
	color.Info.Println(sessionView.Hero)

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
		color.Info.Println(setHero)

	}
	fmt.Println("----------------------------------------------------------")
	for {
		action := interact.SingleSelect(
			"Your action?",
			map[string]string{"1": "Fight", "2": "Archive", "3": "Quit"},
			"1",
			false,
		)
		// call fight api
		switch action {
		case "Fight":
			color.Info.Println("Your select is:", action)
			fightResp, err := auth.Fight(token)
			if err != nil {
				color.Error.Println("An error occured while fight !!!!", err)
			}
			fmt.Println("GameOver", fightResp.GameOver)
			fmt.Println("NextLevel", fightResp.NextLevel)
			fmt.Println("HeroBlood", fightResp.HeroBlood)
			fmt.Println("BossBlood", fightResp.BossBlood)
			if fightResp.GameOver || fightResp.HeroBlood == 0 {
				color.Info.Println("Game Over")
				msg, err := auth.QuitSession(token)
				if err != nil {
					color.Error.Println("An error occured while quit session !!!!", err)
				}
				color.Info.Println(msg)
				return
			}
			if fightResp.NextLevel || fightResp.BossBlood == 0 {
				nextLevelResp, err := auth.NextLevel(token)
				if err != nil {
					color.Error.Println("An error occured while goes into next level !!!!", err)
				}
				if nextLevelResp != nil && nextLevelResp.Session != nil {
					curLevel := nextLevelResp.Session["current_level"]
					if curLevel != nil {
						if curLevel.(float64) > 2 {
							return
						}
					}
				}
			}
		case "Archive":
			color.Info.Println("Your select is:", action)
			sessionid, err := auth.ArchiveSession(token)
			if err != nil {
				color.Error.Println("An error occured while archive session !!!!", err)
			}
			color.Info.Println(fmt.Sprintf("SessionID %s is archived", sessionid))
		case "Quit":
			color.Info.Println("Your select is:", action)
			msg, err := auth.QuitSession(token)
			if err != nil {
				color.Error.Println("An error occured while quit session !!!!", err)
			}
			color.Info.Println(msg)
			return
		default:
			color.Info.Println("Your select is:", action)
			color.Info.Println("Calling Fight API")
		}
	}
}
