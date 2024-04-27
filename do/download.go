package do

import (
	"errors"
	"fmt"
	"github/bigmangos/chrome-offline-installer/internal/util"
	"log/slog"
	"os"
	"strings"

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

func checkUpdate(currentVersion string) bool {
	lastVer, err := getLatestVersion()
	if err != nil {
		slog.Error(fmt.Sprintf("getLatestVersion err: %v", err))
		return false
	}
	return util.IsNewVersion(currentVersion, lastVer)
}

func Download(currentVersion string, url string) error {
	if !checkUpdate(currentVersion) {
		slog.Info("no need to download")
		return nil
	}

	// 更新last_version.txt
	f, err := os.Create(lastVersionPath)
	defer f.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("create file error: %v", err))
	}

	if _, err = f.WriteString(currentVersion); err != nil {
		return errors.New(fmt.Sprintf("write file error: %v", err))
	}

	all := strings.Split(url, "/")
	fileName := all[len(all)-1]
	_, err = os.Stat(fileName)
	if os.IsExist(err) {
		slog.Info(fmt.Sprintf("%s already exists", fileName))
		return nil
	}
	resp, err := resty.New().R().
		SetOutput(fileName).Get(url)
	if err != nil {
		return errors.New(fmt.Sprintf("download %s err: %v", fileName, err))
	}

	// 检查响应状态码
	if resp.IsError() {
		return errors.New(fmt.Sprintf("download %s err: %v", fileName, resp.Status()))
	}

	return nil
}
