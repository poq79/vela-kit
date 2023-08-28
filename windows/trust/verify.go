package trust

import (
	"github.com/gentlemanautomaton/wintrust"
)

func Verify(file string) error {
	return wintrust.VerifyFile(file)
}
