package do

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/bigmangos/chrome-offline-installer/internal/model"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	dataPath   = "data.json"
	readmePath = "README.md"
)

var verList = []string{
	"win_stable_x64",
	"win_stable_x86",
	"win_beta_x64",
	"win_beta_x86",
	"win_dev_x64",
	"win_dev_x86",
	"win_canary_x64",
	"win_canary_x86",
}

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

	r, err := json.Marshal(&newData)
	if err != nil {
		return errors.New(fmt.Sprintf("unmarshal error: %v", err))
	}

	if _, err = f.Write(r); err != nil {
		return errors.New(fmt.Sprintf("write file error: %v", err))
	}

	return nil
}

func SaveMarkdown(data map[string]*model.ChromeInstallerInfo) error {
	indexUrl := "https://github.com/bigmangos/chrome-offline-installer?tab=readme-ov-file#"
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
	buf.WriteString("最后检测更新时间\n")
	buf.WriteString(time.Now().In(loc).Format("2006-01-02 15:04:05"))

	buf.WriteString("\n\n## 目录\n")
	for _, v := range verList {
		link := indexUrl + v
		buf.WriteString(fmt.Sprintf("* [%s](%s)\n", v, link))
	}
	buf.WriteString("\n")
	for _, v := range verList {
		d, ok := data[v]
		if !ok {
			slog.Error(fmt.Sprintf("not found: %s", v))
			continue
		}
		buf.WriteString(fmt.Sprintf("## %s\n", v))
		buf.WriteString(fmt.Sprintf("**最新版本**： %s  \n", d.Version))
		buf.WriteString(fmt.Sprintf("**文件大小**： %0.2f MB  \n", formatSize(d.Size)))
		buf.WriteString(fmt.Sprintf("**校验值（Sha256）**： %s  \n", d.Sha256))
		buf.WriteString("**下载链接**：\n")
		for _, url := range d.Urls {
			buf.WriteString(fmt.Sprintf("%s\n", url))
		}
	}

	if _, err := f.WriteString(buf.String()); err != nil {
		return errors.New(fmt.Sprintf("write file error: %v", err))
	}

	return nil
}

func formatSize(in int) float64 {
	return float64(in) / 1024 / 1024
}
