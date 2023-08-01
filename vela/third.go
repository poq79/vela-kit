package vela

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"encoding/json"
	"github.com/vela-ssoc/vela-kit/lua"
	"io"
	"os"
	"path/filepath"
)

type ThirdEnv interface {
	Third(name string) (*ThirdInfo, error) //同步接口
	ThirdInfo(name string) *ThirdInfo      //查看缓存
}

type ThirdInfo struct {
	Dir       string `json:"dir"`
	Hash      string `json:"hash,omitempty"`
	Name      string `json:"name,omitempty"`
	Size      int64  `json:"size,omitempty"`
	MTime     int64  `json:"mtime"`
	CTime     int64  `json:"ctime"`
	Expire    int    `json:"expire,omitempty"`
	Extension string `json:"extension"`
	Changed   bool   `json:"changed"`
}

func (info *ThirdInfo) String() string                         { return lua.B2S(info.Byte()) }
func (info *ThirdInfo) Type() lua.LValueType                   { return lua.LTObject }
func (info *ThirdInfo) AssertFloat64() (float64, bool)         { return 0, false }
func (info *ThirdInfo) AssertString() (string, bool)           { return "", false }
func (info *ThirdInfo) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (info *ThirdInfo) Peek() lua.LValue                       { return info }

func (info *ThirdInfo) IsNull() bool {
	return info.Name == ""
}

func (info *ThirdInfo) Filepath() string {
	return filepath.Join(info.Dir, info.Name)

}

func (info *ThirdInfo) IsZip() bool {
	return filepath.Ext(info.Name) == ".zip"
}

func (info *ThirdInfo) unzip(from, dst string) error {
	defer func() {
		err := os.Remove(info.Filepath())
		if err != nil {
			_G.Errorf("%s remove fail %v", info.Name, err)
		}
	}()

	rc, err := zip.OpenReader(from)
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	files := rc.File
	for _, file := range files {
		if err = info.extract(dst, file); err != nil {
			break
		}
	}

	return err
}

func (info *ThirdInfo) extract(dir string, file *zip.File) error {
	ff := file.FileInfo()
	full := filepath.Join(dir, file.Name)
	if ff.IsDir() {
		return os.MkdirAll(full, ff.Mode())
	}
	_ = os.MkdirAll(filepath.Dir(full), 0644)

	df, err := os.OpenFile(full, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer func() { _ = df.Close() }()

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	_, err = io.Copy(df, rc)
	return err
}

func (info *ThirdInfo) Modified(log Log) bool {
	stat, err := os.Stat(info.File())
	if err != nil {
		log.Errorf("%s third update case stat read fail %v", info.Name, err)
		return true
	}

	size := stat.Size()
	if size != info.Size {
		log.Errorf("%s third update case size change file-size=%d cache-size=%d", info.Name, size, info.Size)
		return true
	}

	mtime := stat.ModTime().Unix()
	if mtime != info.MTime {
		log.Errorf("%s third update case mtime change size=%d cache-mtime=%d  file-mtime= %d hash=%s",
			info.Name, size, info.MTime, mtime, info.Hash)
		return true
	}
	return false
}

func (info *ThirdInfo) CheckSum() (string, error) {
	filename := info.File()
	fd, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("open %v error %v", filename, err)
	}
	defer fd.Close()

	hub := md5.New()
	io.Copy(hub, fd)
	return hex.EncodeToString(hub.Sum(nil)), nil
}

func (info *ThirdInfo) decompression() error {
	ext := filepath.Ext(info.Name)
	if ext != ".zip" {
		return nil
	}

	_, err := os.Stat(info.File())
	if err != nil {
		return info.unzip(info.Filepath(), info.File())
	}

	if e := os.RemoveAll(info.File()); e != nil {
		_G.Errorf("%s remove fail %v", info.Name, e)
	}

	return info.unzip(info.Filepath(), info.File())

}

func (info *ThirdInfo) Compression() bool {
	ext := filepath.Ext(info.Name)
	if ext != ".zip" {
		return false
	}
	return true
}

func (info *ThirdInfo) File() string {
	ext := filepath.Ext(info.Name)
	if ext == ".zip" {
		return filepath.Join(info.Dir, info.Name[:len(info.Name)-len(ext)])
	}
	return filepath.Join(info.Dir, info.Name)
}

func (info *ThirdInfo) Byte() []byte {
	chunk, _ := json.Marshal(info)
	return chunk
}

func (info *ThirdInfo) FlushStat() error {
	file := info.File()
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}

	info.MTime = stat.ModTime().Unix()
	info.Size = stat.Size()
	return nil
}

func (info *ThirdInfo) Flush() error {
	if e := info.decompression(); e != nil {
		return e
	}

	file := info.File()
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}

	info.MTime = stat.ModTime().Unix()
	info.Size = stat.Size()
	return nil
}

func (info *ThirdInfo) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "hash":
		return lua.S2L(info.Hash)
	case "name":
		return lua.S2L(info.Name)
	case "size":
		return lua.LInt(info.Size)
	case "mtime":
		return lua.LNumber(info.MTime)
	case "ctime":
		return lua.LNumber(info.CTime)
	case "expire":
		return lua.LNumber(info.Expire)
	case "ext":
		return lua.S2L(info.Extension)
	case "file":
		return lua.S2L(info.File())
	}

	return lua.LNil
}
