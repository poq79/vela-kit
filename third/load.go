package third

import (
	"fmt"
	"github.com/vela-ssoc/vela-common-mba/netutil"
	"github.com/vela-ssoc/vela-kit/vela"
	tunnel "github.com/vela-ssoc/vela-tunnel"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type thirdHttpReply struct {
	Data []struct {
		Name string `json:"name"`
		Hash string `json:"hash"`
	} `json:"data"`
}

func (th *third) drop(info *vela.ThirdInfo) {
	//清除内存缓存
	th.recovery(info.Name)

	//删除本地缓存
	th.bucket.Delete(info.Name)

	//删除文件
	s, err := os.Stat(info.File())
	if err != nil {
		return
	}

	if s.IsDir() {
		err = os.RemoveAll(info.File())
	} else {
		err = os.Remove(info.File())
	}

	if err != nil {
		xEnv.Errorf("%s remove %s fail info:%+v %v", info.Name, info.File(), info, err)
		return
	}
	xEnv.Errorf("%s drop success info:%+v", info.Name, info)

}

func (th *third) success(info *vela.ThirdInfo) {
	th.bucket.Push(info.Name, info.Byte(), 0)
	th.publish(info)
	xEnv.Errorf("%s third update success info:%+v", info.Name, info)
}

func (th *third) Delay() {
	delay := time.Second * 10
	time.Sleep(delay)
}

// http 请求下载接口 name=aaa.lua&hash=123
func (th *third) download(name string, hash string, retry int) (*vela.ThirdInfo, error) {
	uri := th.uri(name, hash)
	var att *tunnel.Attachment
	var err error

	retryN := 0
RETRY:
	att, err = xEnv.Attachment(uri)
	switch e := err.(type) {
	case nil:
		defer att.Close()

	case *netutil.HTTPError:
		switch e.Code {
		case http.StatusNotFound:
			return nil, err
		case http.StatusTooManyRequests:
			th.Delay()
			goto RETRY
		default:
			if retryN > retry {
				return nil, err
			}
			retryN++
			th.Delay()
			goto RETRY
		}

	default:
		if retryN > retry {
			return nil, err
		}
		retryN++
		th.Delay()
		goto RETRY
	}

	info := &vela.ThirdInfo{
		Dir:  th.dir,
		Name: name,
		Hash: att.Hash(),
	}

	if att.NotModified() {
		return info, info.FlushStat()
	}

	save := filepath.Join(th.dir, filepath.Clean(name))
	md5, err := att.File(save)
	if err != nil {
		return nil, err
	}

	if md5 != att.Hash() {
		th.drop(&vela.ThirdInfo{Dir: th.dir, Name: name})
		return nil, fmt.Errorf("file.md5=%s http.hash=%s not equal", md5, att.Hash())
	}

	err = info.Flush()
	if err != nil {
		return nil, fmt.Errorf("third=%v flush fail %v", info, err)
	}

	th.success(info)
	return info, nil
}

func (th *third) update(name string, checksum string) {
	info, err := th.download(name, checksum, 3) //hash=
	if err != nil {
		th.drop(info)
		xEnv.Errorf("update %s third fail %+v hash=%s", info.Name, info, checksum)
		return
	}
}

func (th *third) uri(name string, hash string) string {
	return fmt.Sprintf("/api/v1/broker/third?name=%s&hash=%s", name, hash)

}

func (th *third) load(name string) (*vela.ThirdInfo, error) {
	if len(name) < 2 {
		return nil, fmt.Errorf("invalid empty third name")
	}

	return th.download(name, "", 5) //hash empty download
}

func (th *third) Load(name string) (*vela.ThirdInfo, error) {
	info, ok := th.info(name)
	if !ok {
		xEnv.Errorf("%s third update case not found cache table", name)
		return th.load(name)
	}

	if info.Modified(xEnv) {
		return th.load(info.Name)
	}

	return info, nil
}
