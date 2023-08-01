package vela

import "time"

type ScanInfo struct {
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
