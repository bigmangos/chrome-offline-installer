package do

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"github/bigmangos/chrome-offline-installer/internal/model"
	"github/bigmangos/chrome-offline-installer/internal/util"
	"log/slog"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type info struct {
	os  string
	app string
}

// https://source.chromium.org/chromium/chromium/src/+/main:chrome/installer/util/additional_parameters.cc;drc=406947a0f1e0e6b596d387b6b14156f369e8c55d;l=206
var versionInfo = map[string]info{
	"win_stable_x86":   {os: `arch="x86"`, app: `appid="{8A69D345-D564-463C-AFF1-A69D9E530F96}" ap="x86-stable"`},
	"win_stable_x64":   {os: `arch="x64"`, app: `appid="{8A69D345-D564-463C-AFF1-A69D9E530F96}" ap="x64-stable"`},
	"win_stable_arm64": {os: `arch="arm64"`, app: `appid="{8A69D345-D564-463C-AFF1-A69D9E530F96}" ap="arm64-stable"`},
	"win_beta_x86":     {os: `arch="x86"`, app: `appid="{8A69D345-D564-463C-AFF1-A69D9E530F96}" ap="1.1-beta-arch_x86"`},
	"win_beta_x64":     {os: `arch="x64"`, app: `appid="{8A69D345-D564-463C-AFF1-A69D9E530F96}" ap="1.1-beta-arch_x64"`},
	"win_beta_arm64":   {os: `arch="arm64"`, app: `appid="{8A69D345-D564-463C-AFF1-A69D9E530F96}" ap="1.1-beta-arch_arm64"`},
	"win_dev_x86":      {os: `arch="x86"`, app: `appid="{8A69D345-D564-463C-AFF1-A69D9E530F96}" ap="2.0-dev-arch_x86"`},
	"win_dev_x64":      {os: `arch="x64"`, app: `appid="{8A69D345-D564-463C-AFF1-A69D9E530F96}" ap="2.0-dev-arch_x64"`},
	"win_dev_arm64":    {os: `arch="arm64"`, app: `appid="{8A69D345-D564-463C-AFF1-A69D9E530F96}" ap="2.0-dev-arch_arm64"`},
	"win_canary_x86":   {os: `arch="x86"`, app: `appid="{4EA16AC7-FD5A-47C3-875B-DBF4A2008C20}" ap="x86-canary"`},
	"win_canary_x64":   {os: `arch="x64"`, app: `appid="{4EA16AC7-FD5A-47C3-875B-DBF4A2008C20}" ap="x64-canary"`},
	"win_canary_arm64": {os: `arch="arm64"`, app: `appid="{4EA16AC7-FD5A-47C3-875B-DBF4A2008C20}" ap="arm64-canary"`},
}

const (
	updateUrl = "https://tools.google.com/service/update2"
)

func post(os, app string) (string, error) {
	xml := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
    <request protocol="3.0" updater="Omaha" updaterversion="1.3.36.372" shell_version="1.3.36.352" ismachine="0" sessionid="{11111111-1111-1111-1111-111111111111}" installsource="taggedmi" requestid="{11111111-1111-1111-1111-111111111111}" dedup="cr" domainjoined="0">
    <hw physmemory="16" sse="1" sse2="1" sse3="1" ssse3="1" sse41="1" sse42="1" avx="1"/>
    <os platform="win" version="10.0.26100.1742" %s/>
    <app version="" %s>
    <updatecheck/>
    <data name="install" index="empty"/>
    </app>
    </request>`, os, app)

	client := resty.New()
	resp, err := client.R().
		SetBody(xml).
		Post(updateUrl)
	if err != nil {
		return "", errors.New("post error: " + err.Error())
	}

	return resp.String(), nil
}

func decode(input string) (*model.ChromeInstallerInfo, error) {
	r := &model.Response{}
	if err := xml.NewDecoder(strings.NewReader(input)).Decode(r); err != nil {
		return nil, errors.New("decode error: " + err.Error())
	}

	//fmt.Printf("xml decode: %#v\n", r)

	if len(r.App.Updatecheck.Manifest.Packages.Package) == 0 || len(r.App.Updatecheck.URLs.URL) == 0 {
		return nil, errors.New(fmt.Sprintf("no package, response: %+v", r))
	}

	packageFirst := r.App.Updatecheck.Manifest.Packages.Package[0]

	hash := packageFirst.Hash
	decodedHash, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("decode hash error: %+v", err))
	}

	size, err := strconv.Atoi(packageFirst.Size)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("parse size error: %+v", err))
	}

	urls := make([]string, 0, len(r.App.Updatecheck.URLs.URL))
	for _, url := range r.App.Updatecheck.URLs.URL {
		urls = append(urls, url.Codebase+packageFirst.Name)
	}

	return &model.ChromeInstallerInfo{
		Version: r.App.Updatecheck.Manifest.Version,
		Size:    size,
		Sha1:    hex.EncodeToString(decodedHash),
		Sha256:  packageFirst.HashSha256,
		Urls:    urls,
	}, nil
}

func FetchAndUpdateData(data map[string]*model.ChromeInstallerInfo) {
	for k, v := range versionInfo {
		res, err := post(v.os, v.app)
		if err != nil {
			slog.Error(fmt.Sprintf("post error: %v", err))
			continue
		}
		decodedRes, err := decode(res)
		if err != nil {
			slog.Error(fmt.Sprintf("decode error: %v", err))
			continue
		}
		if data[k] == nil || util.IsNewVersion(decodedRes.Version, data[k].Version) {
			data[k] = decodedRes
		}
	}
}
