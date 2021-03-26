package client

import (
"bufio"
"fmt"
"os"
"strconv"
"strings"
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

func Start() {
	var emailAddress = ""
	var err error
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Please Input your email address:")
		// ReadString will block until the delimiter is entered
		emailAddress, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		// remove the delimeter from the string
		emailAddress = strings.TrimSuffix(emailAddress, "\n")
		if emailAddress != "" && len(emailAddress) != 0 {
			fmt.Println(emailAddress)
			break
		} else {
			fmt.Println("Your email address is empty, Please try again", err)
		}
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
		fmt.Println("HeroName is empty")
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
		fmt.Println("You can choose one Hero from below list:")
		for index, hero := range heroList.Heros {
			fmt.Printf("%d: %s\n", index, hero.Name)
		}
		var heroNum int
		fmt.Scanf("%d", &heroNum)
		choosenHero := heroList.Heros[heroNum].Name
		// call get hero api
		fmt.Println("You have choosen Hero:", choosenHero)

	}
	for {
		fmt.Print("Please choose a game action:\n 1: Fight \n 2: Archive\n 3: Quit\n")
		action, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}
		actionNum1 := strings.TrimSuffix(action, "\n")
		actionNum, err := strconv.ParseInt(actionNum1, 10, 32)
		fmt.Println(actionNum)
		// call fight api
		switch actionNum {
		case 1:
			fmt.Println("Calling Fight API: ")
			fight := &Fight{
				GameOver:  false,
				NextLevel: true,
				Score:     0,
				HeroBlood: 20,
				BossBlood: 80,
			}
			if fight.BossBlood == 0 {
				fmt.Println("Enter next level:")
			} else {
				if fight.HeroBlood == 0 {
					break
				}
			}
		case 2:
			fmt.Println("Calling Archive API:")
		case 3:
			fmt.Println("Quit:")
			return
		default:
			fmt.Println("Calling Fight API:")
		}
	}

}