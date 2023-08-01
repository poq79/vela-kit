package binary

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	validator "github.com/go-playground/validator/v10"
	"hash/adler32"
	"time"
)

type hide struct {
	Servername string    `json:"servername"`  // wss/https TLS 证书校验时的 servername
	LAN        []string  `json:"lan"`         // broker/manager 的内网地址
	VIP        []string  `json:"vip"`         // broker/manager 的外网地址
	Edition    string    `json:"edition"`     // semver 版本号
	Hash       string    `json:"hash"`        // 文件原始 hash
	Size       int       `json:"size"`        // 文件原始 size
	DownloadAt time.Time `json:"download_at"` // 下载时间
}

var (
	ErrDecrypt = errors.New("decrypt failed")
)

func (h *hide) valid() error {
	v := validator.New()
	err := v.Struct(h)
	if err != nil {
		return err
	}

	if len(h.LAN)+len(h.VIP) == 0 {
		return fmt.Errorf("not found valid remote peer")
	}

	return nil
}

// decrypt 解密
func decrypt(enc []byte) ([]byte, error) {
	// base64解码
	dsz := base64.StdEncoding.DecodedLen(len(enc))
	raw := make([]byte, dsz)
	n, err := base64.StdEncoding.Decode(raw, enc)
	if err != nil || n < 7 {
		return nil, ErrDecrypt
	}
	mask := raw[n-1]
	for i, b := range raw {
		raw[i] = mask ^ b
	}

	rsz := n - 7
	load := raw[rsz : n-1]
	// vc[2] mask2 vc[3] vc[0] vc[1] mask1
	vc := make([]byte, 4)
	var mask1, mask2 byte
	vc[0], vc[1], vc[2], vc[3], mask1, mask2 = load[3], load[4], load[0], load[2], load[5], load[1]
	for i := 0; i < rsz; i += 2 {
		if i+1 >= rsz {
			raw[i] ^= mask1
			break
		}
		// 奇数位mask1 偶数位mask2
		raw[i], raw[i+1] = raw[i+1]^mask1, raw[i]^mask2
	}
	raw = raw[:rsz]
	sum := binary.BigEndian.Uint32(vc)
	if adler32.Checksum(raw) != sum {
		return nil, ErrDecrypt
	}

	return raw, nil
}

// decryptJSON 解密后反序列化为struct
func decryptJSON(enc []byte, v interface{}) error {
	raw, err := decrypt(enc)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
}
