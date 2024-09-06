package steamcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func Test_downloadSteamCmd(t *testing.T) {
	tempDirPath := filepath.Join(os.TempDir(), fmt.Sprintf("temp_test_%v", uuid.New()))
	if err := os.Mkdir(tempDirPath, os.ModePerm); err != nil {
		t.Fatal("os.Mkdir temp dir", err)
	}

	if err := downloadSteamCmd(tempDirPath); err != nil {
		t.Fatal(err)
	}

	if IsSteamCmdInstalled(tempDirPath) == false {
		t.Fatal("download successful but still not installed")
	}
}
