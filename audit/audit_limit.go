package audit

import (
	auxlib "github.com/vela-ssoc/vela-kit/auxlib"
)

type inhibitMatch func([]string, *Event) bool // ev.a .. ev.bb ..ev.cc ..ev.dd / 500

type inhibit struct {
	offset int
	tag    string
	fnc    []func(*Event) string
}

func newInhibit(tag string) *inhibit {
	return &inhibit{tag: tag, offset: 0}
}

func (inh *inhibit) compile() {
	n := len(inh.tag)
	if n == 0 {
		return
	}

	for idx := 0; idx < n; idx++ {
		ch := inh.tag[idx]
		if ch != '$' {
			continue
		}

		if inh.offset != idx {
			item := inh.tag[inh.offset:idx]
			inh.append(func(ev *Event) string {
				return item
			})
			inh.Offset(idx)
		}

		switch inh.tag[idx : idx+4] {
		case "$msg":
			inh.append(func(ev *Event) string { return ev.msg })
			idx += 4
			inh.Offset(idx)
			goto NEXT

		case "$err":
			inh.append(func(ev *Event) string {
				if ev.err == nil {
					return ""
				}
				return ev.err.Error()
			})
			idx += 4
			inh.Offset(idx)
			goto NEXT
		}

		switch inh.tag[idx : idx+5] {
		case "$time":
			inh.append(func(ev *Event) string {
				return ev.time.Format("2006-01-02.15:04:05")
			})
			idx += 5
			inh.Offset(idx)
			goto NEXT
		case "$inet":
			inh.append(func(ev *Event) string { return ev.inet })
			idx += 5
			inh.Offset(idx)
			goto NEXT
		case "$from":
			inh.append(func(ev *Event) string { return ev.from })
			idx += 5
			inh.Offset(idx)
			goto NEXT
		case "$user":
			inh.append(func(ev *Event) string { return ev.user })
			idx += 5
			inh.Offset(idx)
			goto NEXT
		case "$auth":
			inh.append(func(ev *Event) string { return ev.auth })
			idx += 5
			inh.Offset(idx)
			goto NEXT
		}

		switch inh.tag[idx : idx+6] {
		case "$level":
			inh.append(func(ev *Event) string {
				return ev.level
			})
			idx += 6
			inh.Offset(idx)
			goto NEXT
		case "$alert":
			inh.append(func(ev *Event) string {
				if ev.alert {
					return "true"
				}
				return "false"
			})
			idx += 6
			inh.Offset(idx)
			goto NEXT
		}

		switch inh.tag[idx : idx+7] {
		case "$typeof":
			inh.append(func(ev *Event) string {
				return ev.typeof
			})
			idx += 7
			inh.Offset(idx)
			goto NEXT

		case "region":
			inh.append(func(ev *Event) string {
				return ev.region
			})
			idx += 7
			inh.Offset(idx)
			goto NEXT

		case "$upload":
			inh.append(func(ev *Event) string {
				if ev.upload {
					return "true"
				}
				return "false"
			})
			idx += 7
			inh.Offset(idx)
			goto NEXT
		}

		switch {
		case inh.tag[idx:idx+3] == "$id":
			inh.append(func(ev *Event) string {
				return ev.id
			})
			idx += 3
			inh.Offset(idx)
			goto NEXT

		case inh.tag[idx:idx+8] == "$subject":
			inh.append(func(ev *Event) string {
				return ev.subject
			})
			idx += 8
			inh.Offset(idx)
			goto NEXT

		case inh.tag[idx:idx+12] == "$remote_addr":
			inh.append(func(ev *Event) string {
				return ev.rAddr
			})
			idx += 12
			inh.Offset(idx)
			goto NEXT

		case inh.tag[idx:idx+12] == "$remote_port":
			inh.append(func(ev *Event) string {
				return auxlib.ToString(ev.rPort)
			})
			idx += 12
			inh.Offset(idx)
		}
	NEXT:
	}

	if inh.offset < n {
		item := inh.tag[inh.offset:]
		inh.append(func(_ *Event) string {
			return item
		})
	}
}

func (inh *inhibit) append(fn func(*Event) string) {
	inh.fnc = append(inh.fnc, fn)
}

func (inh *inhibit) Offset(idx int) {
	inh.offset = idx
}

func (inh *inhibit) Key(ev *Event) string {
	var buf []byte
	for _, fn := range inh.fnc {
		buf = append(buf, fn(ev)...)
	}
	return auxlib.B2S(buf)
}

// newInhibitMatch helo.$id.$inet.$bb.xx => helo.id1.192.179.1.1.vela.cc
func newInhibitMatch(tag string, ttl int) inhibitMatch {
	inh := newInhibit(tag)
	inh.compile()

	return func(bkt []string, ev *Event) bool {
		if len(bkt) == 0 {
			return false
		}

		db := xEnv.Bucket(bkt...)
		key := inh.Key(ev)
		count, err := db.Incr(key, 1, ttl)
		if err != nil {
			xEnv.Errorf("%v incr %s fail %v", bkt, count, err)
			return false
		}

		if count >= 1 {
			return true
		}
		return false
	}
}
