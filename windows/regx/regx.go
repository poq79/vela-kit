package regx

import (
	"golang.org/x/sys/windows/registry"
)

type Btt struct {
	path  string
	root  registry.Key
	reply registry.Key
	Error error
}

func New(root registry.Key, path string) *Btt {
	return &Btt{
		root: root,
		path: path,
	}
}

func (b *Btt) Have() bool {
	key, err := registry.OpenKey(b.root, b.path, registry.QUERY_VALUE)
	if err != nil {

		return false
	}

	b.reply = key
	return true
}
