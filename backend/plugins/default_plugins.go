package plugins

func getDefaultPlugins() []Plugin {
	return []Plugin{getCounterStrikeSharp(), getCs2PracticeMode()}
}

func getCounterStrikeSharp() Plugin {
	return Plugin{
		Name:        "CounterStrikeSharp",
		Description: "CounterStrikeSharp allows you to write server plugins in C# for Counter-Strike 2/Source2/CS2",
		URL:         "https://github.com/roflmuffin/CounterStrikeSharp",
		InstallDir:  "/",
		Versions: []Version{
			{
				Name:        "v264",
				DownloadURL: "https://github.com/roflmuffin/CounterStrikeSharp/releases/download/v264/counterstrikesharp-with-runtime-build-264-linux-8f59fd5.zip",
				Dependencies: []PluginDependency{
					{
						Name:         "metamod_source",
						InstallDir:   "/",
						Version:      "2.0.0-git1313",
						DownloadURL:  "https://mms.alliedmods.net/mmsdrop/2.0/mmsource-2.0.0-git1313-linux.tar.gz",
						Dependencies: nil,
					},
				},
			},
		},
	}
}

func getCs2PracticeMode() Plugin {
	return Plugin{
		Name:        "Cs2PracticeMode",
		Description: "Practice mode for cs2 server based on CounterStrikeSharp",
		URL:         "https://github.com/Phi-S/cs2-practice-mode",
		InstallDir:  "/addons/counterstrikesharp/plugins/",
		Versions: []Version{
			{
				Name:        "0.0.16",
				DownloadURL: "https://github.com/Phi-S/cs2-practice-mode/releases/download/0.0.16/cs2-practice-mode-0.0.16.tar.gz",
				Dependencies: []PluginDependency{
					{
						Name:        "CounterStrikeSharp",
						InstallDir:  "/",
						Version:     "v264",
						DownloadURL: "https://github.com/roflmuffin/CounterStrikeSharp/releases/download/v264/counterstrikesharp-with-runtime-build-264-linux-8f59fd5.zip",
						Dependencies: []PluginDependency{
							{
								Name:         "metamod_source",
								InstallDir:   "/",
								Version:      "2.0.0-git1313",
								DownloadURL:  "https://mms.alliedmods.net/mmsdrop/2.0/mmsource-2.0.0-git1313-linux.tar.gz",
								Dependencies: nil,
							},
						},
					},
				},
			},
		},
	}
}
