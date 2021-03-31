package types

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

const (
	GameOver string = `
	░██████╗░░█████╗░███╗░░░███╗███████╗  ░█████╗░██╗░░░██╗███████╗██████╗░
	██╔════╝░██╔══██╗████╗░████║██╔════╝  ██╔══██╗██║░░░██║██╔════╝██╔══██╗
	██║░░██╗░███████║██╔████╔██║█████╗░░  ██║░░██║╚██╗░██╔╝█████╗░░██████╔╝
	██║░░╚██╗██╔══██║██║╚██╔╝██║██╔══╝░░  ██║░░██║░╚████╔╝░██╔══╝░░██╔══██╗
	╚██████╔╝██║░░██║██║░╚═╝░██║███████╗  ╚█████╔╝░░╚██╔╝░░███████╗██║░░██║
	░╚═════╝░╚═╝░░╚═╝╚═╝░░░░░╚═╝╚══════╝  ░╚════╝░░░░╚═╝░░░╚══════╝╚═╝░░╚═╝`
	WinGame string = `
	░██╗░░░░░░░██╗██╗███╗░░██╗  ████████╗██╗░░██╗███████╗  ░██████╗░░█████╗░███╗░░░███╗███████╗
	░██║░░██╗░░██║██║████╗░██║  ╚══██╔══╝██║░░██║██╔════╝  ██╔════╝░██╔══██╗████╗░████║██╔════╝
	░╚██╗████╗██╔╝██║██╔██╗██║  ░░░██║░░░███████║█████╗░░  ██║░░██╗░███████║██╔████╔██║█████╗░░
	░░████╔═████║░██║██║╚████║  ░░░██║░░░██╔══██║██╔══╝░░  ██║░░╚██╗██╔══██║██║╚██╔╝██║██╔══╝░░
	░░╚██╔╝░╚██╔╝░██║██║░╚███║  ░░░██║░░░██║░░██║███████╗  ╚██████╔╝██║░░██║██║░╚═╝░██║███████╗
	░░░╚═╝░░░╚═╝░░╚═╝╚═╝░░╚══╝  ░░░╚═╝░░░╚═╝░░╚═╝╚══════╝  ░╚═════╝░╚═╝░░╚═╝╚═╝░░░░░╚═╝╚══════╝`
)
