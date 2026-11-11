package do

import (
	"errors"
	"fmt"
	"github/bigmangos/chrome-offline-installer/internal/util"
	"log/slog"
	"os"
	"path"

	"github.com/go-resty/resty/v2"
)

const (
	lastVersionPath = "last_download.txt"
)

func getLatestVersion() (string, error) {
	f, err := os.ReadFile(lastVersionPath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error reading %s: %s", lastVersionPath, err))
	}
	return string(f), nil
}

func CheckUpdate(currentVersion string) bool {
	lastVer, err := getLatestVersion()
	if err != nil {
		slog.Error(fmt.Sprintf("getLatestVersion err: %v", err))
		return false
	}

	if !util.IsNewVersion(currentVersion, lastVer) {
		return false
	}

	// 更新last_version.txt
	f, err := os.Create(lastVersionPath)
	defer f.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("create file error: %v", err))
		return true
	}

	if _, err = f.WriteString(currentVersion); err != nil {
		slog.Error(fmt.Sprintf("write file error: %v", err))
		return true
	}

	return true
}

func Download(arch, url string) error {
	fileName := path.Base(url)
	if fileName == "." || fileName == "/" {
		return errors.New(fmt.Sprintf("can not find file name: %v", url))
	}

	ext := path.Ext(fileName)
	fileName = fileName[:len(fileName)-len(ext)] + "_" + arch + ext

	slog.Info("download", "name", fileName)

	resp, err := resty.New().R().SetOutput(fileName).Get(url)
	if err != nil {
		return errors.New(fmt.Sprintf("download %s err: %v", fileName, err))
	}

	if resp.IsError() {
		return errors.New(fmt.Sprintf("download %s err: %v", fileName, resp.Status()))
	}

	return nil
}
