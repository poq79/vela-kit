package webdav

import (
	"os"
	"syscall"
	"time"
)

func ctime(v os.FileInfo) time.Time {
	stat := v.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, stat.CreationTime.Nanoseconds())
}

func atime(v os.FileInfo) time.Time {
	stat := v.Sys().(*syscall.Win32FileAttributeData)
	return time.Unix(0, stat.LastAccessTime.Nanoseconds())
}

func owner(v os.FileInfo) (string, string) {
	return "", ""
}
