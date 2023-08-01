package third

import (
	"encoding/json"
	"github.com/vela-ssoc/vela-kit/vela"
	"os"
	"path/filepath"
	"sync"
)

type third struct {
	dir    string
	mutex  sync.RWMutex
	cache  map[string]*vela.ThirdInfo
	bucket vela.Bucket
}

func (th *third) sync() {
	th.cache, _ = th.table()
}

func (th *third) online() {
	xEnv.OnConnect("third.sync", func() error {
		th.thirdSyncDir()
		th.sync()
		return nil
	})
}

func (th *third) info(name string) (*vela.ThirdInfo, bool) {
	th.mutex.Lock()
	defer th.mutex.Unlock()

	info, ok := th.cache[name]
	return info, ok
}

func (th *third) recovery(name string) {
	th.mutex.Lock()
	defer th.mutex.Unlock()

	delete(th.cache, name)

	tab := make(map[string]*vela.ThirdInfo, len(th.cache))
	for key, info := range th.cache {
		tab[key] = info
	}
	th.cache = tab
}

func (th *third) publish(info *vela.ThirdInfo) {
	th.mutex.Lock()
	defer th.mutex.Unlock()
	th.cache[info.Name] = info
}

func (th *third) table() (map[string]*vela.ThirdInfo, []string) {
	var names []string
	current := make(map[string]*vela.ThirdInfo, 32)
	th.bucket.ForEach(func(name string, data []byte) {
		names = append(names, name)
		info := &vela.ThirdInfo{}
		err := json.Unmarshal(data, info)
		if err != nil {
			xEnv.Errorf("%s json unmarshal fail %v", name, err)
		}
		current[name] = info
	})
	return current, names
}

func (th *third) clear() {
	current, _ := th.table()
	for _, info := range current {
		if e := os.Remove(info.File()); e != nil {
			xEnv.Errorf("%s third %+v remove fail %v", info.Name, info, e)
		} else {
			xEnv.Errorf("%s third %+v remove success", info.Name, info)
		}
		th.bucket.Delete(info.Name)
	}
	th.cache = make(map[string]*vela.ThirdInfo, 16)
}

func (th *third) once(current map[string]*vela.ThirdInfo, fileInfo *vela.ThirdInfo) {
	info, ok := current[fileInfo.Name]
	if !ok {
		th.drop(fileInfo)
		return
	}

	if info.IsNull() {
		th.drop(fileInfo)
		return
	}

	if info.Size != fileInfo.Size {
		th.update(fileInfo.Name, fileInfo.Hash)
		return
	}

	if info.MTime != fileInfo.MTime {
		th.update(fileInfo.Name, fileInfo.Hash)
		return
	}

	if !fileInfo.Compression() && fileInfo.Hash != info.Hash {
		th.update(fileInfo.Name, fileInfo.Hash)
		return
	}

	th.publish(info)
	th.update(info.Name, info.Hash)
}

func (th *third) thirdSyncDir() {
	current, _ := th.table()
	if len(current) == 0 {
		err := os.Remove(th.dir)
		if err != nil {
			xEnv.Errorf("third remove %s fail %v", th.dir, err)
		}

		th.check()
		return
	}

	dir, err := os.ReadDir(th.dir)
	if err != nil {
		xEnv.Errorf("third sync %s file info fail %v", th.dir, err)
		th.remove(current)
		return
	}

	for _, entry := range dir {
		info := &vela.ThirdInfo{
			Name: entry.Name(),
		}

		if entry.IsDir() {
			info.Name = info.Name + ".zip"
		}

		ff, e := entry.Info()
		if e != nil {
			xEnv.Errorf("third %s read file info fail %v", info.File(), e)
			th.once(current, info)
			continue
		}

		info.MTime = ff.ModTime().Unix()
		info.Size = ff.Size()

		hash, e := info.CheckSum()
		if e != nil {
			xEnv.Errorf("third %s read file hash read fail %v", info.File(), e)
			th.once(current, info)
			continue
		}

		info.Hash = hash
		info.Extension = filepath.Ext(info.Name)
		th.once(current, info)
		xEnv.Errorf("on connect sync third %s succeed", info.File())
	}

}

func (th *third) remove(expires map[string]*vela.ThirdInfo) {
	for _, info := range expires {
		th.drop(info)
	}
}

func (th *third) Info(name string) *vela.ThirdInfo {
	entry, ok := th.info(name)
	if ok {
		return entry
	}
	return nil
}
