package do

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/bigmangos/chrome-offline-installer/internal/model"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"
)

const (
	dataPath   = "data.json"
	readmePath = "README.md"
)

func LoadSavedData() map[string]*model.ChromeInstallerInfo {
	data := make(map[string]*model.ChromeInstallerInfo)
	f, err := os.ReadFile(dataPath)
	if err != nil {
		slog.Error(fmt.Sprintf("read file error: %v", err))
		return data
	}

	if err = json.Unmarshal(f, &data); err != nil {
		slog.Error(fmt.Sprintf("unmarshal error: %v", err))
		return data
	}

	//fmt.Printf("json load: %#v\n", data)
	return data
}

func SaveData(newData map[string]*model.ChromeInstallerInfo) error {
	f, err := os.Create(dataPath)
	defer f.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("create file error: %v", err))
	}

	r, err := json.MarshalIndent(&newData, "", "    ")
	if err != nil {
		return errors.New(fmt.Sprintf("unmarshal error: %v", err))
	}

	if _, err = f.Write(r); err != nil {
		return errors.New(fmt.Sprintf("write file error: %v", err))
	}

	return nil
}

func SaveMarkdown(data map[string]*model.ChromeInstallerInfo) error {
	f, err := os.Create(readmePath)
	defer f.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("create file error: %v", err))
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return errors.New(fmt.Sprintf("load location error: %v", err))
	}

	var buf strings.Builder
	buf.WriteString("# Google Chrome 离线安装包\n")
	buf.WriteString("本工程是 [Bush2021/chrome_installer](https://github.com/Bush2021/chrome_installer) Go实现，感谢原作者\n\n")
	buf.WriteString("稳定版存档：<https://github.com/bigmangos/chrome-offline-installer/releases>\n\n")
	buf.WriteString("最近一次检测更新时间（UTC+8）：\n")
	buf.WriteString(time.Now().In(loc).Format("2006-01-02 15:04:05"))

	channelOrder := []string{"stable", "beta", "dev", "canary"}
	channelNames := map[string]string{
		"stable": "Stable",
		"beta":   "Beta",
		"dev":    "Dev",
		"canary": "Canary",
	}

	buf.WriteString("\n\n## Contents\n\n")
	for _, channel := range channelOrder {
		name := channelNames[channel]
		if channelNames[channel] != "" {
			buf.WriteString(fmt.Sprintf("- [%v](#%v)\n\n", name, channel))
		}
	}
	buf.WriteString("\n")

	for _, channel := range channelOrder {
		channelName := channelNames[channel]

		buf.WriteString(fmt.Sprintf("## %v\n\n", channelName))
		buf.WriteString("| Architecture | Version | Size | SHA-256 | Download |\n")
		buf.WriteString("|--------------|---------|------|---------|----------|\n")

		archOrder := []string{"x64", "arm64", "x86"}
		archNames := map[string]string{
			"x64":   "X64",
			"arm64": "ARM64",
			"x86":   "X86",
		}

		for _, arch := range archOrder {
			key := strings.ToLower("win_" + channelName + "_" + arch)
			d, ok := data[key]
			if !ok || d == nil {
				slog.Error(fmt.Sprintf("not found: %s", key))
				continue
			}
			sha256Short := d.Sha256
			if len(sha256Short) > 8 {
				sha256Short = sha256Short[0:8]
			}

			urls := ""
			for i, url := range d.Urls {
				urls += fmt.Sprintf("[url-%d](%v) ", i, url)
			}

			buf.WriteString(fmt.Sprintf("| **%v** | %v | %v | %v | %v |\n", archNames[arch], d.Version, formatSize(d.Size), sha256Short, urls))
		}

		buf.WriteString("\n")
		buf.WriteString("<details>\n")
		buf.WriteString("<summary>Full SHA-256 (sha256sum -c)</summary>\n\n")
		buf.WriteString("```\n")

		for _, arch := range archOrder {
			key := strings.ToLower("win_" + channelName + "_" + arch)
			d, ok := data[key]
			if !ok || d == nil {
				slog.Error(fmt.Sprintf("not found: %s", key))
				continue
			}
			name := path.Base(d.Urls[0])

			buf.WriteString(fmt.Sprintf("%v %v\n", name, d.Sha256))
		}

		buf.WriteString("```\n\n")
		buf.WriteString("</details>\n\n")
	}

	if _, err := f.WriteString(buf.String()); err != nil {
		return errors.New(fmt.Sprintf("write file error: %v", err))
	}

	return nil
}

func formatSize(in int) float64 {
	return float64(in) / 1024 / 1024
}
