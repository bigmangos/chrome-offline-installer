package main

import (
	"fmt"
	"github/bigmangos/chrome-offline-installer/do"
	"log/slog"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	data, err := do.LoadSavedData()
	if err != nil {
		slog.Error(fmt.Sprintf("load saved data error: %v", err))
		return
	}
	do.FetchAndUpdateData(data)
	if err = do.SaveData(data); err != nil {
		slog.Error(fmt.Sprintf("save data error: %v", err))
	}
	if err = do.SaveMarkdown(data); err != nil {
		slog.Error(fmt.Sprintf("save markdown error: %v", err))
	}

	// 下载最新x64稳定版本
	downloadVer := data["win_stable_x64"]
	if err = do.Download(downloadVer.Version, downloadVer.Urls[0]); err != nil {
		slog.Error(fmt.Sprintf("download error: %v", err))
	}
}
