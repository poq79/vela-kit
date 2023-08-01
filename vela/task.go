package vela

import "time"

type Task struct {
	ID      int64     `json:"id"`
	Dialect bool      `json:"dialect"`
	Name    string    `json:"name"`
	Link    string    `json:"link"`
	Status  string    `json:"status"`
	Hash    string    `json:"hash"`
	From    string    `json:"from"`
	Uptime  time.Time `json:"uptime"`
	Failed  bool      `json:"failed"`
	Cause   string    `json:"cause"`
	Runners []*Runner `json:"runners"`
}

func (t *Task) Range(fn func(r *Runner)) {
	n := len(t.Runners)
	if n == 0 {
		return
	}

	for i := 0; i < n; i++ {
		fn(t.Runners[i])
	}
}

func (t *Task) Exist(name string) bool {
	n := len(t.Runners)
	if n == 0 {
		return false
	}

	for i := 0; i < n; i++ {
		if t.Runners[i].Name == name {
			return true
		}
	}
	return false
}
