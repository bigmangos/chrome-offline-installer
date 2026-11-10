package model

type ChromeInstallerInfo struct {
	Version string
	Size    int
	Sha1    string
	Sha256  string
	Urls    []string
}
