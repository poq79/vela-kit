package env

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/vela"
)

func (env *Environment) Scan(id int64, cname string, chunk []byte, metadata map[string]interface{}, timeout int) error {
	if env.sub.task == nil {
		return fmt.Errorf("not found task tree")
	}

	return env.sub.task.Scan(env, id, cname, chunk, metadata, timeout)
}

func (env *Environment) ScanList() []*vela.ScanInfo {
	if env.sub.task == nil {
		return nil
	}

	return env.sub.task.ScanList()

}

func (env *Environment) StopScanAll() {
	if env.sub.task == nil {
		return
	}

	env.sub.task.StopScanAll()
}

func (env *Environment) StopScanById(id int64) {
	if env.sub.task == nil {
		return
	}

	env.sub.task.StopScanById(id)
}

func (env *Environment) StopScanByName(name string) {
	if env.sub.task == nil {
		return
	}

	env.sub.task.StopScanByName(name)
}
