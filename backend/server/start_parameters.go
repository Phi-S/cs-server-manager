package server

type StartParameters struct {
    Hostname        string
    Password        string
    StartMap        string
    MaxPlayers      uint8
    SteamLoginToken string
    Additional      []string
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
