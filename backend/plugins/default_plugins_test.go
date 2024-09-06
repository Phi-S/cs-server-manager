package plugins

import (
	"encoding/json"
	"os"
	"testing"
)

func Benchmark_createDefaultPluginsJson(b *testing.B) {
	cs2PracticeMode := getDefaultPlugins()
	jsonStr, err := json.MarshalIndent(cs2PracticeMode, "", "  ")
	if err != nil {
		b.Log(err)
		b.FailNow()
	}
	if err := os.WriteFile("default-plugins.json", jsonStr, os.ModePerm); err != nil {
		b.Log(err)
		b.FailNow()
	}
}
