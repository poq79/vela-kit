package node

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func socket(addr *string) error {
	s := strings.Split(*addr, ":")
	n := len(s)
	if n == 1 && net.ParseIP(s[0]) != nil {
		*addr = fmt.Sprintf("%s:53", s[0])
		return nil
	}

	if n > 2 {
		return errors.New("invalid socket ip" + *addr)
	}

	if net.ParseIP(s[0]) == nil {
		return errors.New("invalid socket ip " + *addr)
	}

	port, err := strconv.Atoi(s[1])
	if err != nil {
		return errors.New("invalid socket port type " + *addr)
	}

	if port < 1 || port > 65535 {
		return errors.New("invalid socket " + *addr)
	}

	return nil
}
