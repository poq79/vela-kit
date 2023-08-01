package auxlib

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"os"
)

func FileMd5(filename string) (string, error) {
	path := fmt.Sprintf("./%s", filename)
	pFile, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open %v error %v", filename, err)
	}
	defer pFile.Close()

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	hub := md5.New()
	Copy(ctx, hub, pFile)
	return hex.EncodeToString(hub.Sum(nil)), nil
}

func Checksum(path string, hash string) error {
	if len(hash) != 32 {
		return fmt.Errorf("invalid osquery file hash")
	}

	if path == "" {
		return fmt.Errorf("invalid osquery path")
	}

	fmd5, err := FileMd5(path)
	if err != nil {
		return err
	}

	if fmd5 != hash {
		return fmt.Errorf("hash not match file=%s signate=%s", fmd5, hash)
	}

	return nil
}

func CheckWriter(val lua.LValue, L *lua.LState) lua.Writer {
	if val.Type() != lua.LTVelaData {
		L.RaiseError("must be writer , got %s", val.Type().String())
		return nil
	}

	return lua.CheckWriter(val.(*lua.VelaData))

}
