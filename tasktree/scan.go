package tasktree

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"sync"
	"time"
)

type ScanPool struct {
	mu    sync.Mutex
	dirty map[string]*scanTask
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

func (sp *ScanPool) Get(name string) *scanTask {
	s, ok := sp.dirty[name]
	if !ok {
		return nil
	}

	return s
}

func (tt *TaskTree) visit(callback func(key string, s *scanTask) bool) {
	tt.scans.Range(func(key, value interface{}) bool {
		var ok bool

		s, ok := value.(*scanTask)
		if !ok {
			tt.scans.Delete(key)
			return true
		}
		return callback(key.(string), s)
	})
}

func (tt *TaskTree) StopScanAll() {
	tt.visit(func(key string, s *scanTask) bool {
		s.StopScanTask()
		return true
	})
}

func (tt *TaskTree) StopScanById(id int64) {
	tt.visit(func(key string, s *scanTask) bool {
		if s.code.header.id == id {
			s.StopScanTask()
			return false
		}
		return true
	})
}

func (tt *TaskTree) StopScanByName(name string) {
	tt.visit(func(key string, s *scanTask) bool {
		if key == name {
			s.StopScanTask()
			return false
		}
		return true
	})

}

func (tt *TaskTree) ScanList() []*vela.ScanInfo {
	var list []*vela.ScanInfo
	tt.visit(func(key string, s *scanTask) bool {
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

func (tt *TaskTree) Scan(env vela.Environment, id int64, cname string, chunk []byte,
	metadata map[string]interface{}, timeout int) error {

	var s *scanTask

	v, ok := tt.scans.Load(cname)
	if !ok {
		goto run
	}

	s, ok = v.(*scanTask)
	if !ok {
		tt.scans.Delete(cname)
		goto run
	}

	s.StopScanTask()
	time.Sleep(5 * time.Second)

run:
	task := newScanTask(env, id, cname, chunk, metadata, timeout)
	tt.scans.Store(cname, task)
	return task.call(env)
}
