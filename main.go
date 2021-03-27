package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gookit/color"
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

func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

func main() {
	// style := color.New(color.FgWhite, color.BgBlack, color.OpBold)
	// style.Println("custom color style")
	cyan := color.FgLightCyan.Render
	blue := color.FgLightBlue.Render
	magenta := color.FgLightMagenta.Render
	var emailAddress = ""
	var err error
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s", magenta("Please Input your email address: "))
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
	// calling load session api
	loadSession := &Session{
		UID:           "1",
		HeroName:      "",
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
			},
		}
		fmt.Printf("%s\n", magenta("You can choose one Hero from below list:"))
		for index, hero := range heroList.Heros {
			fmt.Printf("%s%s %s \n", blue(fmt.Sprintf("%d", index)), blue(")"), cyan(hero.Name))
		}
		var heroNum int
		fmt.Scanf("%d", &heroNum)
		choosenHero := heroList.Heros[heroNum].Name
		// call get hero api
		fmt.Printf("%s\n", magenta("You have choosen Hero: ", cyan(choosenHero)))

	}
	for {
		fmt.Printf("%s\n %s %s \n %s %s\n %s %s\n",
			magenta("Please choose a game action:"), blue("1)"), cyan("Fight"), blue("2)"), cyan("Archive"), blue("3)"), cyan("Quit"))
		action, err := reader.ReadString('\n')
		if err != nil {
			color.Error.Println("An error occured while reading input. Please try again", err)
			return
		}
		actionNum1 := strings.TrimSuffix(action, "\n")
		actionNum, err := strconv.ParseInt(actionNum1, 10, 32)
		if err != nil {
			color.Error.Println("An error occured while parse the input.", err)
			return
		}
		// call fight api
		switch actionNum {
		case 1:
			fmt.Printf("Calling %s API\n", cyan("Fight"))
			fight := &Fight{
				GameOver:  false,
				NextLevel: true,
				Score:     0,
				HeroBlood: 20,
				BossBlood: 80,
			}
			if fight.BossBlood == 0 {
				fmt.Println("Enter Next Level")
			} else {
				if fight.HeroBlood == 0 {
					break
				}
			}
		case 2:
			fmt.Printf("Calling %s API\n", cyan("Archive"))
		case 3:
			fmt.Printf("%s\n", cyan("Quit"))
			return
		default:
			fmt.Printf("Calling %s API\n", cyan("Fight"))
		}
	}

}
