package node

import (
	"sync"
)

var (

	//默认远程地址
	resolve = "114.114.114.114:53"

	//初始化一次
	once sync.Once

	//节点实体
	_G *node
)
