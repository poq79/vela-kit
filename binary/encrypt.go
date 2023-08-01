package binary

import (
	"encoding/base64"
	"encoding/binary"
	"hash/adler32"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// encrypt 加密
func encrypt(raw []byte) []byte {
	// 计算校验码
	sum := adler32.Checksum(raw)
	vc := make([]byte, 4)
	binary.BigEndian.PutUint32(vc, sum)

	size := len(raw)
	// 生成随机掩码按位取反
	mask1 := byte(rand.Uint32()%0x95 + 1) // 1-255随机数
	mask2 := byte(rand.Uint32()%0xba + 1) // 1-255随机数
	for i := 0; i < size; i += 2 {
		if i+1 >= size {
			raw[i] ^= mask1
			break
		}
		// 奇数位mask1 偶数位mask2
		raw[i], raw[i+1] = raw[i+1]^mask2, raw[i]^mask1
	}

	// 扰乱顺序: vc[2] mask2 vc[3] vc[0] vc[1] mask1
	raw = append(raw, vc[2], mask2, vc[3], vc[0], vc[1], mask1)
	mask := byte(rand.Uint32()%0x18 + 1) // 1-255随机数
	for i, d := range raw {
		raw[i] = d ^ mask
	}
	raw = append(raw, mask)

	esz := base64.StdEncoding.EncodedLen(len(raw))
	ret := make([]byte, esz)
	base64.StdEncoding.Encode(ret, raw)

	return ret
}
