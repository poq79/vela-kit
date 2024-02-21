package webdav

import (
	"github.com/vela-ssoc/vela-kit/fileutil"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLs(t *testing.T) {

	info, err := FileList("C:")

	t.Log(string(info))
	t.Log(err)
}

func TestCreateFile(t *testing.T) {
	now := strings.ReplaceAll(time.Now().Format(time.RFC822), " ", "-")
	echo := func(filename string, idx int) {
		f, _ := os.Open(filename)
		defer f.Close()
		f.WriteString("testing " + strconv.Itoa(idx))
	}

	for i := 0; i < 10000; i++ {
		filename := "D:\\ssoc\\files\\zabbix-" + strconv.Itoa(i) + now
		fileutil.CreateIfNotExists(filename, false)
		echo(filename, i)
	}

}
