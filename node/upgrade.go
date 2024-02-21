package node

import (
	"fmt"
	"net/url"
)

type upgrade struct {
	//Id         int    `json:"id"`
	Semver     string `json:"semver"`
	Customized string `json:"customized"`
}

func (u *upgrade) String() string {
	return fmt.Sprintf("版本号:%s 分支版本:%s", u.Semver, u.Customized)
}

func (u *upgrade) Query() url.Values {
	return url.Values{"version": []string{u.Semver}, "tags": xEnv.Tags(), "customized": []string{u.Customized}}
}
