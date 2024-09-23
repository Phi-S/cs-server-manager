package plugins

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Phi-S/cs-server-manager/gvalidator"
)

func executeCustomInstallAction(csgoDir, pluginName string) error {
	if pluginName == "metamod_source" {
		if err := metamodInstall(filepath.Join(csgoDir, "gameinfo.gi")); err != nil {
			return fmt.Errorf("metamodInstall: %w", err)
		}
	}

	return nil
}

func executeCustomUninstallAction(csgoDir, pluginName string) error {
	if pluginName == "metamod_source" {
		if err := metamodUninstall(filepath.Join(csgoDir, "gameinfo.gi")); err != nil {
			return fmt.Errorf("metamodUninstall: %w", err)
		}
	}

	return nil
}

var metamodLine = "\t\t\tGame csgo/addons/metamod"

func metamodInstall(gameinfoPath string) error {
	if err := gvalidator.Instance().Var(gameinfoPath, "required,file"); err != nil {
		return fmt.Errorf("gameinfo.gi path '%v' is not valid %w", gameinfoPath, err)
	}

	f, err := os.Open(gameinfoPath)
	if err != nil {
		return fmt.Errorf("failed to open gameinfo.gi '%v' %w", gameinfoPath, err)
	}

	lineAdded := false
	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		if line == "\t\t\tGame_LowViolence\tcsgo_lv // Perfect World content override" {
			lines = append(lines, metamodLine)
			lineAdded = true
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read gameinfo.gi '%v' %w", gameinfoPath, err)
	}

	if !lineAdded {
		return fmt.Errorf("failed to add required line to gaminfo.gi")
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close filestream %w", err)
	}

	if err := os.Remove(gameinfoPath); err != nil {
		return fmt.Errorf("failed to remove old gaminfo.io '%v' %w", gameinfoPath, err)
	}

	f, err = os.Create(gameinfoPath)
	if err != nil {
		return fmt.Errorf("failed to create new gameinfo.gi '%v' %w", gameinfoPath, err)
	}

	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write line '%v' to new gameinfo.gi '%v' %w", line, gameinfoPath, err)
		}
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close new gaminfo.gi file %w", err)
	}

	newGameinfoContent, err := os.ReadFile(gameinfoPath)
	if err != nil {
		return fmt.Errorf("failed to validate new gameinfo.gi %w", err)
	}

	if !strings.Contains(string(newGameinfoContent), metamodLine) {
		return fmt.Errorf("new gameinfo.gi is missing metamod line")
	}

	return nil
}

func metamodUninstall(gameinfoPath string) error {
	if err := gvalidator.Instance().Var(gameinfoPath, "required,file"); err != nil {
		return fmt.Errorf("gameinfo.gi path '%v' is not valid %w", gameinfoPath, err)
	}

	f, err := os.Open(gameinfoPath)
	if err != nil {
		return fmt.Errorf("failed to open gameinfo.gi '%v' %w", gameinfoPath, err)
	}

	lineRemoved := false
	lines := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line == metamodLine {
			lineRemoved = true
		} else {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read gameinfo.gi '%v' %w", gameinfoPath, err)
	}

	if !lineRemoved {
		return fmt.Errorf("failed to remove metamod line from gaminfo.gi")
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close filestream %w", err)
	}

	if err := os.Remove(gameinfoPath); err != nil {
		return fmt.Errorf("failed to remove old gaminfo.io '%v' %w", gameinfoPath, err)
	}

	f, err = os.Create(gameinfoPath)
	if err != nil {
		return fmt.Errorf("failed to create new gameinfo.gi '%v' %w", gameinfoPath, err)
	}

	for _, line := range lines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write line '%v' to new gameinfo.gi '%v' %w", line, gameinfoPath, err)
		}
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("failed to close new gaminfo.gi file %w", err)
	}

	newGameinfoContent, err := os.ReadFile(gameinfoPath)
	if err != nil {
		return fmt.Errorf("failed to validate new gameinfo.gi %w", err)
	}

	if !strings.Contains(string(newGameinfoContent), metamodLine) {
		return fmt.Errorf("new gameinfo.gi is still containing metamod line after uninstall")
	}

	return nil
}
