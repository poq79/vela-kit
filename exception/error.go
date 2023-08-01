package exception

import (
	"bytes"
)

type errKV struct {
	key string
	err error
}

type Cause struct {
	data []errKV
}

func (c *Cause) Len() int {
	return len(c.data)
}

func (c *Cause) Error() string {
	n := c.Len()
	if n == 0 {
		return ""
	}

	var buff bytes.Buffer
	for i := 0; i < n; i++ {
		if i != 0 {
			buff.WriteByte('\n')
		}

		item := c.data[i]
		if item.key != "" {
			buff.WriteString(item.key)
			buff.WriteByte(':')
		}
		buff.WriteString(item.err.Error())

	}

	return buff.String()
}

func (c *Cause) Try(key string, err error) {
	if err == nil {
		return
	}
	c.data = append(c.data, errKV{key, err})
}

func (c *Cause) Wrap() error {
	if len(c.data) == 0 {
		return nil
	}

	return c
}

func New() *Cause {
	return &Cause{}
}
