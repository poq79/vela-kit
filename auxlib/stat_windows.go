package auxlib

import (
	"os"
	"syscall"
	"time"
)

func FileStat(fi os.FileInfo) (atime, mtime, ctime time.Time) {
	mtime = fi.ModTime()
	stat := fi.Sys().(*syscall.Win32FileAttributeData)
	atime = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	ctime = time.Unix(0, stat.CreationTime.Nanoseconds())
	return
}

func FileStatByFile(name string) (atime, mtime, ctime time.Time, err error) {
	var fi os.FileInfo
	fi, err = os.Stat(name)
	if err != nil {
		return
	}
	mtime = fi.ModTime()
	stat := fi.Sys().(*syscall.Win32FileAttributeData)
	atime = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	ctime = time.Unix(0, stat.CreationTime.Nanoseconds())
	return
}
