package server

type StartParameters struct {
	Hostname        string   `json:"hostname"`
	Password        string   `json:"password"`
	StartMap        string   `json:"start_map"`
	MaxPlayers      uint8    `json:"max_players"`
	SteamLoginToken string   `json:"steam_login_token"`
	Additional      []string `json:"additional"`
}

func DefaultStartParameters() *StartParameters {
	return &StartParameters{
		Hostname:        "cs server",
		Password:        "",
		StartMap:        "de_mirage",
		MaxPlayers:      10,
		SteamLoginToken: "",
		Additional:      []string{},
	}
}
