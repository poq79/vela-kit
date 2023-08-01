package webdav

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
	"time"
)

func ctime(v os.FileInfo) time.Time {
	stat := v.Sys().(*syscall.Stat_t)
	return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
}

func atime(v os.FileInfo) time.Time {
	stat := v.Sys().(*syscall.Stat_t)
	return time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
}

func owner(v os.FileInfo) (userName string, groupName string) {
	// 获取文件所有者信息
	s, ok := v.Sys().(*syscall.Stat_t)
	if !ok {
		return
	}

	u, err := user.LookupId(fmt.Sprint(s.Uid))
	if err != nil {
		xEnv.Errorf("not found file owner %s uid:%d", v.Name(), u.Uid)
		userName = ""
	} else {
		userName = u.Username
	}

	group, err := user.LookupGroupId(fmt.Sprint(u.Gid))
	if err != nil {
		xEnv.Errorf("not found file group %s uid:%d", v.Name(), u.Uid)
		groupName = ""
	} else {
		groupName = group.Name
	}
	return
}
