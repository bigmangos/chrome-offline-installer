package main

import (
	"fmt"
	"github/bigmangos/chrome-offline-installer/do"
	"log/slog"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	data := do.LoadSavedData()
	do.FetchAndUpdateData(data)
	if err := do.SaveData(data); err != nil {
		slog.Error(fmt.Sprintf("save data error: %v", err))
	}
	if err := do.SaveMarkdown(data); err != nil {
		slog.Error(fmt.Sprintf("save markdown error: %v", err))
	}

	downloadVer := data["win_stable_x64"]

	if !do.CheckUpdate(downloadVer.Version) {
		slog.Info("no need to download")
		return
	}

	if err := do.Download("x64", downloadVer.Urls[0]); err != nil {
		slog.Error(fmt.Sprintf("download x86 error: %v", err))
	}

	downloadVer = data["win_stable_arm64"]
	if err := do.Download("arm64", downloadVer.Urls[0]); err != nil {
		slog.Error(fmt.Sprintf("download arm64 error: %v", err))
	}
}
