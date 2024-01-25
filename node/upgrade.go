package node

import "net/url"

type upgrade struct {
	Id         int    `json:"id"`
	Semver     string `json:"semver"`
	Customized string `json:"customized"`
}

func (u *upgrade) Query() url.Values {
	return url.Values{"version": []string{u.Semver}, "tags": xEnv.Tags(), "customized": []string{u.Customized}}
}
