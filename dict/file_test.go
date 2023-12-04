package dict

import (
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	f := File{
		filename: "D:\\github.com\\vela-dev\\logs\\pass.dict",
	}

	err := f.ForEach(func(s string) (over bool) {
		if s == "yes" {
			return true
		}

		t.Logf("%s", s)
		return
	})

	t.Logf("over %v", err)

}

func TestFile_Scanner(t *testing.T) {
	f := File{
		filename: "D:\\vela-ssoc\\dict\\top1000.txt",
	}

	it := f.Iterator()
	for it.Next() {
		t.Logf("id:1 text:%s", it.Text())
	}

	it.Reset()
	for it.Next() {
		t.Logf("id:2 text:%s", it.Text())
	}

	time.Sleep(2 * time.Second)
}
