package server

import "errors"

type ServerService struct {
	Busy      bool
	ServerDir string
}

func Start(serverService *ServerService) error {
	if serverService.Busy {
		return errors.New("ServerService is busy")
	}

	return nil
}
