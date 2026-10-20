package model

type Response struct {
	Protocol string   `xml:"protocol,attr"`
	Server   string   `xml:"server,attr"`
	Daystart Daystart `xml:"daystart"`
	App      App      `xml:"app"`
}

type Daystart struct {
	ElapsedDays    int `xml:"elapsed_days,attr"`
	ElapsedSeconds int `xml:"elapsed_seconds,attr"`
}

type App struct {
	Appid       string      `xml:"appid,attr"`
	Cohort      string      `xml:"cohort,attr"`
	Cohortname  string      `xml:"cohortname,attr"`
	Status      string      `xml:"status,attr"`
	Updatecheck Updatecheck `xml:"updatecheck"`
}

type Updatecheck struct {
	Status   string   `xml:"status,attr"`
	URLs     URLs     `xml:"urls"`
	Manifest Manifest `xml:"manifest"`
}

type URLs struct {
	URL []URL `xml:"url"`
}

type URL struct {
	Codebase string `xml:"codebase,attr"`
}

type Manifest struct {
	Version  string   `xml:"version,attr"`
	Actions  Actions  `xml:"actions"`
	Packages Packages `xml:"packages"`
}

type Actions struct {
	Action []Action `xml:"action"`
}

type Action struct {
	Arguments string `xml:"arguments,attr"`
	Event     string `xml:"event,attr"`
	Run       string `xml:"run,attr"`
	Version   string `xml:"Version,attr"`
	Onsuccess string `xml:"onsuccess,attr"`
}

type Packages struct {
	Package []Package `xml:"package"`
}

type Package struct {
	Fp         string `xml:"fp,attr"`
	Hash       string `xml:"hash,attr"`
	HashSha256 string `xml:"hash_sha256,attr"`
	Name       string `xml:"name,attr"`
	Required   string `xml:"required,attr"`
	Size       string `xml:"size,attr"`
}
