package inode

type Inodes struct {
	list        []int32
	Total       int64
	SocketTotal int64
	socket      map[uint32]int32 //socket inode InodeMap
	collect     map[int32]int    //pid num inode number
}

func New(v []int32) *Inodes {
	inode := &Inodes{list: v}
	inode.R()
	return inode
}

func All() *Inodes {
	inode := &Inodes{list: nil}
	inode.R()
	return inode
}
