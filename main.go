package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/new-adventure-aerolite/game_client/auth/auth"

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
	Seesion string
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
			color.Error.Println("An error occured while request token!!!!", err)
		}
		color.Info.Println(token)
		SetToken(token)
	}
	session, err := auth.LoadSession(token)
	if err != nil {
		color.Error.Println("An error occured while load session!!!!", err)
	}
	color.Info.Println(session)
	// calling load session api
	loadSession := &Session{
		UID:           "1",
		HeroName:      "tq",
		LiveHeroBlood: 3,
		LiveBossBlood: 4,
		CurrentLevel:  5,
		Score:         6,
	}

	if loadSession.HeroName == "" && len(loadSession.HeroName) == 0 {
		// calling  get hero list api
		heroList := &HeroList{
			Heros: []Hero{
				{
					Name:         "hero1",
					Details:      "1",
					AttackPower:  2,
					DefensePower: 3,
					Blood:        4,
				},
				{
					Name:         "hero2",
					Details:      "2",
					AttackPower:  2,
					DefensePower: 3,
					Blood:        4,
				},
				{
					Name:         "hero3",
					Details:      "3",
					AttackPower:  2,
					DefensePower: 3,
					Blood:        4,
				},
			},
		}
		fmt.Println("----------------------------------------------------------")
		color.Info.Println("You can choose one Hero from below list:")
		var heroNameList = make([]string, len(heroList.Heros))
		for index, hero := range heroList.Heros {
			heroNameList[index] = hero.Name
		}
		choosenHero := interact.Choice(
			"Choose Hero(use string slice/array)?",
			heroNameList,
			"",
			false,
		)
		// call get hero api
		color.Info.Println("Your select is:", choosenHero)

	}
	fmt.Println("----------------------------------------------------------")
	for {
		action := interact.SingleSelect(
			"Your action(use map)?",
			map[string]string{"1": "Fight", "2": "Archive"},
			"1",
		)
		// call fight api
		switch action {
		case "Fight":
			color.Info.Println("Your select is:", action)
			color.Info.Println("Calling Fight API")
			fight := &Fight{
				GameOver:  false,
				NextLevel: true,
				Score:     0,
				HeroBlood: 20,
				BossBlood: 80,
			}
			if fight.BossBlood == 0 {
				color.Info.Println("Enter Next Level")
			} else {
				if fight.HeroBlood == 0 {
					break
				}
			}
		case "Archive":
			color.Info.Println("Your select is:", action)
			color.Info.Println("Calling Archive API")
		case "quit":
			color.Info.Println("Your select is:", action)
			color.Info.Println("quit", "quit")
			return
		default:
			color.Info.Println("Your select is:", action)
			color.Info.Println("Calling Fight API")
		}
	}

}
