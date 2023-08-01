package webdav

import "testing"

func TestLs(t *testing.T) {

	info, err := FileList("C:")

	t.Log(string(info))
	t.Log(err)
}
