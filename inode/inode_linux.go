package inode

import (
	"github.com/shirou/gopsutil/v3/process"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"os"
	"strconv"
	"strings"
)

func (inode *Inodes) Value(pid int32, link string) {

	inode.Total++
	if !strings.HasPrefix(link, "socket:[") {
		return
	}
	inode.SocketTotal++

	node, err := strconv.ParseInt(link[8:len(link)-1], 10, 64)
	if err != nil {
		return
	}
	inode.socket[uint32(node)] = pid
}

func (inode *Inodes) read(pid int32) {
	path := "/proc" + "/" + auxlib.ToString(pid) + "/fd/"
	d, err := os.Open(path)
	if err != nil {
		return
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return
	}

	for _, name := range names {
		pathLink := path + name
		target, er := os.Readlink(pathLink)
		if er != nil {
			continue
		}
		inode.Value(pid, target)
	}

	inode.collect[pid] = len(names)
}

func (inode *Inodes) List() []int32 {
	if len(inode.list) > 0 {
		return inode.list
	}

	pids, err := process.Pids()
	if err != nil {
		return nil
	}

	return pids
}

func (inode *Inodes) R() error {
	inode.socket = make(map[uint32]int32)
	inode.collect = make(map[int32]int)
	list := inode.List()
	n := len(list)
	for i := 0; i < n; i++ {
		inode.read(list[i])
	}

	return nil
}

func (inode *Inodes) FindPid(id uint32) int32 {
	return inode.socket[id]
}
