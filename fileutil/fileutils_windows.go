package fileutil

import (
	"os"
	"syscall"
	"time"
)

// GetTotalUsedFds Returns the number of used File Descriptors. Not supported
// on Windows.
func GetTotalUsedFds() int {
	return -1
}

func State(name string) (atime, mtime, ctime time.Time, size int64, err error) {
	var fi os.FileInfo
	fi, err = os.Stat(name)
	if err != nil {
		return
	}
	size = fi.Size()
	mtime = fi.ModTime()
	stat := fi.Sys().(*syscall.Win32FileAttributeData)
	atime = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	ctime = time.Unix(0, stat.CreationTime.Nanoseconds())
	return
}

func StateByInfo(fi os.FileInfo) (atime, mtime, ctime time.Time, size int64, err error) {
	size = fi.Size()
	mtime = fi.ModTime()
	stat := fi.Sys().(*syscall.Win32FileAttributeData)
	atime = time.Unix(0, stat.LastAccessTime.Nanoseconds())
	ctime = time.Unix(0, stat.CreationTime.Nanoseconds())
	return
}
