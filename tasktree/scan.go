package tasktree

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"sync"
)

type ScanPool struct {
	mu    sync.Mutex
	dirty map[string]*ScanTask
}

func (sp *ScanPool) StopAll() {
	for _, task := range sp.dirty {
		task.StopScanTask()
	}
}

func (sp *ScanPool) List() []*vela.ScanInfo {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	n := len(sp.dirty)
	tuple := make([]*vela.ScanInfo, 0, n)
	for _, task := range sp.dirty {
		info := &vela.ScanInfo{
			Name:    task.code.Key(),
			Link:    task.code.Link(),
			Hash:    task.code.Hash(),
			From:    task.code.From(),
			Status:  task.code.Status(),
			Uptime:  task.code.header.uptime,
			ID:      task.code.header.id,
			Dialect: task.code.header.dialect,
		}
		info.Runners = task.code.List()

		if task.code.header.err == nil {
			info.Failed = false
		} else {
			info.Cause = task.code.header.err.Error()
		}

		tuple = append(tuple, info)
	}
	return tuple
}

func (sp *ScanPool) Get(name string) *ScanTask {
	s, ok := sp.dirty[name]
	if !ok {
		return nil
	}

	return s
}

func (tt *TaskTree) visit(callback func(key string, s *ScanTask) bool) {
	tt.scans.Range(func(key, value interface{}) bool {
		var ok bool

		s, ok := value.(*ScanTask)
		if !ok {
			tt.scans.Delete(key)
			return true
		}
		return callback(key.(string), s)
	})
}

func (tt *TaskTree) StopScanAll() {
	tt.visit(func(key string, s *ScanTask) bool {
		s.StopScanTask()
		return true
	})
}

func (tt *TaskTree) StopScanById(id int64) {
	tt.visit(func(key string, s *ScanTask) bool {
		if s.code.header.id == id {
			s.StopScanTask()
			return false
		}
		return true
	})
}

func (tt *TaskTree) StopScanByName(name string) {
	tt.visit(func(key string, s *ScanTask) bool {
		if key == name {
			s.StopScanTask()
			return false
		}
		return true
	})

}

func (tt *TaskTree) ScanList() []*vela.ScanInfo {
	var list []*vela.ScanInfo
	tt.visit(func(key string, s *ScanTask) bool {
		info := &vela.ScanInfo{
			Name:    s.code.Key(),
			Link:    s.code.Link(),
			Hash:    s.code.Hash(),
			From:    s.code.From(),
			Status:  s.code.Status(),
			Uptime:  s.code.header.uptime,
			ID:      s.code.header.id,
			Dialect: s.code.header.dialect,
		}

		info.Runners = s.code.List()

		if s.code.header.err == nil {
			info.Failed = false
		} else {
			info.Cause = s.code.header.err.Error()
		}
		list = append(list, info)
		return true
	})

	return list

}

func (tt *TaskTree) Scan(env vela.Environment, id int64, name string, chunk []byte,
	result map[string]interface{}, timeout int) error {

	var s *ScanTask

	v, ok := tt.scans.Load(name)
	if !ok {
		goto run
	}

	s, ok = v.(*ScanTask)
	if !ok {
		tt.scans.Delete(name)
		goto run
	}

	s.StopScanTask()

run:
	task := newScanTask(env, id, name, chunk, result, timeout)
	tt.scans.Store(name, task)
	return task.call(env)
}
