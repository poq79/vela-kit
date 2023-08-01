package grep

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/vela-ssoc/vela-kit/grep/compiler"
	"github.com/vela-ssoc/vela-kit/grep/syntax"
	"strings"
)

// Grep represents compiled glob pattern.
type Grep interface {
	Match(string) bool
}

// Compile creates Grep for given pattern and strings (if any present after pattern) as separators.
// The pattern syntax is:
//
//    pattern:
//        { term }
//
//    term:
//        `*`         matches any sequence of non-separator characters
//        `**`        matches any sequence of characters
//        `?`         matches any single non-separator character
//        `[` [ `!` ] { character-range } `]`
//                    character class (must be non-empty)
//        `{` pattern-list `}`
//                    pattern alternatives
//        c           matches character c (c != `*`, `**`, `?`, `\`, `[`, `{`, `}`)
//        `\` c       matches character c
//
//    character-range:
//        c           matches character c (c != `\\`, `-`, `]`)
//        `\` c       matches character c
//        lo `-` hi   matches character c for lo <= c <= hi
//
//    pattern-list:
//        pattern { `,` pattern }
//                    comma-separated (without spaces) patterns
//

func key(pattern string, sep ...byte) Key {
	h := md5.New()
	h.Write([]byte(pattern))
	h.Write([]byte{'_'})
	h.Write(sep)
	return Key(hex.EncodeToString(h.Sum(nil)))
}

func Compile(pattern string, separators []byte) (Grep, error) {
	k := key(pattern, separators...)
	g, ok := pool.Get(k)
	if ok {
		return g, nil
	}

	ast, err := syntax.Parse(pattern)
	if err != nil {
		return nil, err
	}

	matcher, err := compiler.Compile(ast, bytes.Runes(separators))
	if err != nil {
		return nil, err
	}

	pool.Add(k, matcher)

	return matcher, nil
}

func New(pattern string) func(string) bool {
	switch pattern {
	case "", "*":
		return func(_ string) bool { return true }

	default:
		switch pattern[0] {
		case '=':
			raw := pattern[1:]
			return func(val string) bool { return val == raw }

		case '~':
			pattern = strings.ToLower(pattern[1:])
		}

		gx, err := Compile(pattern, nil)
		if err != nil {
			return func(val string) bool { return pattern == val }
		}
		return gx.Match
	}
}

// QuoteMeta returns a string that quotes all glob pattern meta characters
// inside the argument text; For example, QuoteMeta(`{foo*}`) returns `\[foo\*\]`.
func QuoteMeta(s string) string {
	b := make([]byte, 2*len(s))

	// a byte loop is correct because all meta characters are ASCII
	j := 0
	for i := 0; i < len(s); i++ {
		if syntax.Special(s[i]) {
			b[j] = '\\'
			j++
		}
		b[j] = s[i]
		j++
	}

	return string(b[0:j])
}
