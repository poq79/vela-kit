package logger

import (
	"os"
	"path/filepath"
)

func (z *zapState) reload(new *zapState) {
	z.sugar.Sync()
	z.stop()

	z.cfg = new.cfg
	z.stop = new.stop
	z.sugar = new.sugar
	z.rotate = new.rotate
}

func (z *zapState) clean() {
	dir := filepath.Dir(z.cfg.Filename)

	items, err := filepath.Glob(dir + "/*.log")
	if err != nil {
		xEnv.Errorf("clean %s logger fail %v", z.cfg.Filename, err)
		return
	}

	if len(items) == 0 {
		return
	}

	for _, item := range items {
		if filepath.Base(item) == filepath.Base(z.cfg.Filename) {
			continue
		}

		if e := os.Remove(item); e != nil {
			xEnv.Errorf("clean %s logger fail %v", item, e)
			continue
		}
		xEnv.Errorf("clean %s logger succeed", item)
	}
}
