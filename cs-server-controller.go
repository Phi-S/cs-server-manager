package main

import (
	"cs-server-controller/event"
	steamcmd_service "cs-server-controller/steam_cmd"
	"fmt"
	"log"
	"log/slog"
	"path/filepath"
	"time"
)

func main() {
	var dataFolder = "/home/desk/programming/code/go/data/cs-server-controller"
	var steamCmdFolder = filepath.Join(dataFolder, "steamcmd")
	var serverFolder = filepath.Join(dataFolder, "server")

	var steamCMdService = steamcmd_service.New(steamCmdFolder, serverFolder)

	steamCMdService.EventOnUpdateOrInstallStarted().Register(func(dp event.DefaultPayload) {
		fmt.Println("EVENT: OnUpdateOrInstallStarted")
	})

	steamCMdService.EventOnUpdateOrInstallDone().Register(func(dp event.DefaultPayload) {
		fmt.Println("EVENT: OnUpdateOrInstallDone")
	})

	steamCMdService.EventOnUpdateOrInstallFailed().Register(func(pwd event.PayloadWithData[error]) {
		fmt.Println("EVENT: EventOnUpdateOrInstallFailed DATA: " + pwd.Data.Error())
	})

	steamCMdService.EventOnOutput().Register(func(dp event.PayloadWithData[string]) {
		fmt.Println("EVENT: OnOutput: out: " + dp.Data)
	})

	slog.Info("asdf")

	go PrintStatus(&steamCMdService)
	time.Sleep(time.Millisecond * 100)

	var err = steamCMdService.Start()
	if err != nil {
		fmt.Println("FAILED TO START: " + err.Error())
	}
	log.Println("started......")

	time.Sleep(time.Second * 10)
	fmt.Println("go cancel")
	steamCMdService.Cancel()
	log.Println("after cancel...")
	time.Sleep(time.Second * 5)

	err = steamCMdService.Start()
	if err != nil {
		fmt.Println("FAILED TO START: " + err.Error())
	}
	log.Println("started......")

	time.Sleep(time.Second * 5000)

}

func PrintStatus(steamCMdService *steamcmd_service.SteamcmdService) {
	for {
		fmt.Printf("==== is busy: %t \n", steamCMdService.IsBusy())
		time.Sleep(time.Millisecond * 1000)
	}
}
