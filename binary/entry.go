package binary

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Decode 从文件的payload读取加密数据并反序列化到struct
func Decode(path string) (*hide, error) {
	self, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = self.Close() }()

	// 读取文件的总大小
	var size int64
	if stat, err := self.Stat(); err == nil {
		size = stat.Size()
	}
	// 4位校验码 3个mask 4位payload长度
	if size < 11 {
		return nil, ErrDecrypt
	}

	// 只读取长度
	if _, err = self.Seek(-4, io.SeekEnd); err != nil {
		return nil, err
	}
	data := make([]byte, 4)
	if _, err = self.Read(data); err != nil {
		return nil, err
	}
	// 扰乱顺序 3012 回归正常
	data[0], data[1], data[2], data[3] = data[1], data[2], data[3], data[0]

	psz := binary.BigEndian.Uint32(data)

	if size-4 < int64(psz) {
		return nil, ErrDecrypt
	}

	if _, err = self.Seek(-int64(psz+4), io.SeekEnd); err != nil {
		return nil, err
	}

	enc := make([]byte, psz)
	if _, err = self.Read(enc); err != nil {
		return nil, err
	}

	ret := &hide{}
	err = decryptJSON(enc, ret)
	if err != nil {
		return nil, err
	}

	if err = ret.valid(); err != nil {
		return nil, err
	}

	if int64(psz+4)+int64(ret.Size) != size {
		return nil, fmt.Errorf("binary config %+v , but got size %d", ret, size)
	}

	chunk, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	hash := fmt.Sprintf("%x", md5.Sum(chunk[:ret.Size]))

	if strings.ToLower(hash) != ret.Hash {
		return nil, fmt.Errorf("invalid hash")
	}

	return ret, err
}

// Encode 将struct序列化为JSON后加密
func Encode(v interface{}) ([]byte, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	enc := encrypt(raw)
	return enc, nil
}
