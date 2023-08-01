package webdav

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/disk"
	"github.com/valyala/fasthttp"
	"github.com/vela-ssoc/vela-kit/vela"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var xEnv vela.Environment

type Matcher func(string) bool

var permission = []Matcher{
	Prefix("C:\\ssoc"),
	Prefix("/usr/local/ssoc"),
	Prefix("/tmp"),
	Prefix("/var/log"),
	Ext(".exe", ".lua", ".log", ".txt", ".sh", ".bat", ".jar", ".ps1", ".json", ".php", ".jsp", ".go"),
}

type FileInfo struct {
	FileName string      `json:"filename"`
	Size     int64       `json:"size"`
	MTime    string      `json:"mtime"`
	CTime    string      `json:"ctime"`
	ATime    string      `json:"atime"`
	Dir      bool        `json:"dir"`
	User     string      `json:"user"`
	Group    string      `json:"group"`
	Perm     fs.FileMode `json:"perm"`
	Root     bool        `json:"root"`
}

func mount(ctx *fasthttp.RequestCtx) error {
	partitions, err := disk.Partitions(false)
	if err != nil {
		xEnv.Errorf("disk partitions found fail %v", err)
	}

	chunk, err := json.Marshal(partitions)
	if err != nil {
		return err
	}
	ctx.Write(chunk)
	return nil
}

func FileList(path string) ([]byte, error) {
	fl, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	var list []FileInfo

	for _, filename := range fl {
		entry, e := os.Stat(filename)
		if e != nil {
			xEnv.Errorf("%s stat read fail %v", filename, e)
			continue
		}

		user, group := owner(entry)
		list = append(list, FileInfo{
			FileName: entry.Name(),
			Size:     entry.Size(),
			MTime:    entry.ModTime().Format("2006-01-02 15:04:05"),
			CTime:    ctime(entry).Format("2006-01-02 15:04:05"),
			ATime:    atime(entry).Format("2006-01-02 15:04:05"),
			Dir:      entry.IsDir(),
			Perm:     entry.Mode().Perm(),
			User:     user,
			Group:    group,
		})
	}

	return json.Marshal(list)
}

func ls(ctx *fasthttp.RequestCtx) error {
	args := ctx.PostArgs()
	path := string(args.Peek("path"))

	size := len(path)
	if size == 0 {
		return fmt.Errorf("invalid path got empty")
	}

	suffix := path[size-1]
	var body []byte
	var err error
	var stat os.FileInfo

	if suffix == '/' {
		body, err = FileList(path + "*")
		goto response
	}

	stat, err = os.Stat(path)
	if err == nil && stat.IsDir() {
		body, err = FileList(path + "/*")
		goto response
	}

	body, err = FileList(path)
response:
	if err != nil {
		return err
	}

	_, err = ctx.Write(body)
	return err
}

func cat(ctx *fasthttp.RequestCtx) error {
	path := string(ctx.QueryArgs().Peek("path"))
	if len(path) == 0 {
		return fmt.Errorf("invalid path got emtpy")
	}

	path = filepath.Clean(path)

	for _, fn := range permission {
		if fn(path) {
			goto read
		}
	}

	return fmt.Errorf("not allow")

read:
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	max := int64(1024 * 1024)

	if n := info.Size(); n > max {
		fd, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fd.Close()

		buff := make([]byte, max)
		position := n - max
		if ctx.QueryArgs().Has("not_reverse") {
			position = 0
		}

		if _, e := fd.ReadAt(buff, position); e != nil {
			return e
		}

		ctx.Write(buff)
		return nil
	}

	chunk, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	ctx.Write(chunk)
	return nil
}

func remove(ctx *fasthttp.RequestCtx) error {
	path := string(ctx.QueryArgs().Peek("path"))
	if len(path) == 0 {
		return fmt.Errorf("invalid path got emtpy")
	}

	path = filepath.Clean(path)

	if !strings.HasPrefix(path, "C:\\ssoc") && !strings.HasPrefix(path, "/usr/local/ssoc") {
		return fmt.Errorf("not allow")
	}

	err := os.Remove(path)
	if err != nil {
		return err
	}

	ctx.WriteString("success")
	return nil
}

func Constructor(env vela.Environment) {
	xEnv = env

	r := xEnv.R()
	r.GET("/api/v1/arr/webdav/mount", xEnv.Then(mount))
	r.POST("/api/v1/arr/webdav/ls", xEnv.Then(ls))
	r.GET("/api/v1/arr/webdav/cat", xEnv.Then(cat))
	r.DELETE("/api/v1/arr/webdav", xEnv.Then(remove))
}
