package dict

import (
	"bufio"
	"os"
	"sync/atomic"
)

type FileDict struct {
	done    uint32
	fd      *os.File
	value   *File
	scanner *bufio.Scanner
}

func (f *FileDict) Close() error {
	return nil
}

func (f *FileDict) Reset() error {
	_, err := f.fd.Seek(0, 0)

	sc := bufio.NewScanner(f.fd)
	sc.Split(bufio.ScanLines)
	f.scanner = sc
	return err
}

func (f *FileDict) Next() bool {
	if atomic.LoadUint32(&f.done) != 0 || f.fd == nil {
		return false
	}
	return f.scanner.Scan()
}

func (f *FileDict) Done() {
	atomic.StoreUint32(&f.done, 1)
	if f.fd != nil {
		f.fd.Close()
		f.fd = nil
	}
}

func (f *FileDict) Text() string {
	if atomic.LoadUint32(&f.done) != 0 {
		return ""
	}

	if f.fd == nil {
		return ""
	}

	return f.scanner.Text()
}
