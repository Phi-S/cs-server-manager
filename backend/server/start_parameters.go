package server

type StartParameters struct {
	Hostname        string   `json:"hostname" validate:"required,lte=128"`
	Password        string   `json:"password" validate:"omitempty,alphanum,lte=32"`
	StartMap        string   `json:"start_map" validate:"required,printascii,lte=32"`
	MaxPlayers      uint8    `json:"max_players" validate:"required,number,lte=128"`
	SteamLoginToken string   `json:"steam_login_token" validate:"omitempty,alphanum,len=32"`
	Additional      []string `json:"additional" validate:"omitempty,dive"`
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
