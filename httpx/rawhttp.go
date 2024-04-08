package httpx

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/vektra/errors"
	"github.com/vela-ssoc/vela-kit/strutil"
	"io"
	"net"
	"strconv"
	"unicode"
)

const OptimalBufferSize = 1500

type header struct {
	Name  []byte
	Value []byte
}

type HTTPParser struct {
	Raw               []byte
	Method            string
	Path              string
	Version           string
	Body              []byte
	host              []byte
	Headers           []header
	TotalHeaders      int
	hostRead          bool
	BodyLen           int
	contentLength     int64
	contentLengthRead bool
}

const DefaultHeaderSlice = 50

// Create a new parser
func NewRawHTTP() *HTTPParser {
	return NewSizedHTTPParser(DefaultHeaderSlice)
}

// Create a new parser allocating size for size headers
func NewSizedHTTPParser(size int) *HTTPParser {
	return &HTTPParser{
		Headers:       make([]header, size),
		TotalHeaders:  size,
		contentLength: -1,
	}
}

var (
	ErrNotFound    = errors.New("not found")
	ErrBadProto    = errors.New("bad protocol")
	ErrMissingData = errors.New("missing data")
	ErrUnsupported = errors.New("unsupported http feature")
)

const (
	eNextHeader int = iota
	eNextHeaderN
	eHeader
	eHeaderValueSpace
	eHeaderValue
	eHeaderValueN
	eMLHeaderStart
	eMLHeaderValue
)

// Parse the buffer as an HTTP Request. The buffer must contain the entire
// request or Parse will return ErrMissingData for the caller to get more
// data. (this thusly favors getting a completed request in a single Read()
// call).
//
// Returns the number of bytes used by the header (thus where the body begins).
// Also can return ErrUnsupported if an HTTP feature is detected but not supported.

func (hp *HTTPParser) trim(input []byte) []byte {
	return bytes.TrimLeftFunc(input, unicode.IsSpace)
}

func (hp *HTTPParser) Parse(buf []byte) (int, error) {
	input := hp.trim(buf)

	n, err := hp.parse(input)
	if err != nil {
		return n, err
	}

	body := hp.Raw[n:]
	hp.Body = body
	hp.BodyLen = len(body)

	return n, err
}

func (hp *HTTPParser) parse(input []byte) (int, error) {
	var headers int
	var path int
	var ok bool

	hp.Raw = input
	total := len(input)
	if total == 0 {
		return 0, ErrNotFound
	}

method:
	for i := 0; i < total; i++ {
		switch input[i] {
		case ' ', '\t':
			hp.Method = strutil.B2S(input[0:i])
			ok = true
			path = i + 1
			break method
		}
	}

	if !ok {
		return 0, ErrMissingData
	}

	var version int

	ok = false

path:
	for i := path; i < total; i++ {
		switch input[i] {
		case ' ', '\t':
			ok = true
			hp.Path = strutil.B2S(input[path:i])
			version = i + 1
			break path
		}
	}

	if !ok {
		return 0, ErrMissingData
	}

	var readN bool

	ok = false
loop:
	for i := version; i < total; i++ {
		c := input[i]

		switch readN {
		case false:
			switch c {
			case '\r':
				hp.Version = strutil.B2S(input[version:i])
				readN = true
			case '\n':
				hp.Version = strutil.B2S(input[version:i])
				headers = i + 1
				ok = true
				break loop
			}
		case true:
			if c != '\n' {
				return 0, errors.Context(ErrBadProto, "missing newline in version")
			}
			headers = i + 1
			ok = true
			break loop
		}
	}

	if !ok {
		return 0, ErrMissingData
	}

	var h int

	var headerName []byte

	state := eNextHeader

	start := headers

	for i := headers; i < total; i++ {
		switch state {
		case eNextHeader:
			switch input[i] {
			case '\r':
				state = eNextHeaderN
			case '\n':
				return i + 1, nil
			case ' ', '\t':
				state = eMLHeaderStart
			default:
				start = i
				state = eHeader
			}
		case eNextHeaderN:
			if input[i] != '\n' {
				return 0, ErrBadProto
			}

			return i + 1, nil
		case eHeader:
			if input[i] == ':' {
				headerName = input[start:i]
				state = eHeaderValueSpace
			}
		case eHeaderValueSpace:
			switch input[i] {
			case ' ', '\t':
				continue
			}

			start = i
			state = eHeaderValue
		case eHeaderValue:
			switch input[i] {
			case '\r':
				state = eHeaderValueN
			case '\n':
				state = eNextHeader
			default:
				continue
			}

			hp.Headers[h] = header{headerName, input[start:i]}
			h++

			if h == hp.TotalHeaders {
				newHeaders := make([]header, hp.TotalHeaders+10)
				copy(newHeaders, hp.Headers)
				hp.Headers = newHeaders
				hp.TotalHeaders += 10
			}
		case eHeaderValueN:
			if input[i] != '\n' {
				return 0, ErrBadProto
			}
			state = eNextHeader

		case eMLHeaderStart:
			switch input[i] {
			case ' ', '\t':
				continue
			}

			start = i
			state = eMLHeaderValue
		case eMLHeaderValue:
			switch input[i] {
			case '\r':
				state = eHeaderValueN
			case '\n':
				state = eNextHeader
			default:
				continue
			}

			cur := hp.Headers[h-1].Value

			newheader := make([]byte, len(cur)+1+(i-start))
			copy(newheader, cur)
			copy(newheader[len(cur):], []byte(" "))
			copy(newheader[len(cur)+1:], input[start:i])

			hp.Headers[h-1].Value = newheader
		}
	}

	return 0, ErrMissingData
}

// Return a value of a header matching name.
func (hp *HTTPParser) FindHeader(name []byte) []byte {
	for _, header := range hp.Headers {
		if bytes.Equal(header.Name, name) {
			return header.Value
		}
	}

	for _, header := range hp.Headers {
		if bytes.EqualFold(header.Name, name) {
			return header.Value
		}
	}

	return nil
}

// Return all values of a header matching name.
func (hp *HTTPParser) FindAllHeaders(name []byte) [][]byte {
	var headers [][]byte

	for _, header := range hp.Headers {
		if bytes.EqualFold(header.Name, name) {
			headers = append(headers, header.Value)
		}
	}

	return headers
}

func (hp *HTTPParser) URL(scheme string) string {
	return fmt.Sprintf("%s://%s%s", scheme, strutil.B2S(hp.Host()), hp.Path)
}

var cHost = []byte("Host")

// Return the value of the Host header
func (hp *HTTPParser) Host() []byte {
	if hp.hostRead {
		return hp.host
	}

	hp.hostRead = true
	hp.host = hp.FindHeader(cHost)
	return hp.host
}

var cContentLength = []byte("Content-Length")

// Return the value of the Content-Length header.
// A value of -1 indicates the header was not set.
func (hp *HTTPParser) ContentLength() int64 {
	if hp.contentLengthRead {
		return hp.contentLength
	}

	header := hp.FindHeader(cContentLength)
	if header != nil {
		i, err := strconv.ParseInt(string(header), 10, 0)
		if err == nil {
			hp.contentLength = i
		}
	}

	hp.contentLengthRead = true
	return hp.contentLength
}

func (hp *HTTPParser) BodyReader(rest []byte, in io.ReadCloser) io.ReadCloser {
	return RawBodyReader(hp.ContentLength(), rest, in)
}

var cGet = "GET"

func (hp *HTTPParser) Get() bool {
	return hp.Method == cGet
}

var cPost = "POST"

func (hp *HTTPParser) Post() bool {
	return hp.Method == cPost
}

// Do performs the HTTP request for the given Requester and returns
// a *Response and any error that occured
func Call(req Requester) (*RawResponse, error) {
	var err error
	var conn net.Conn

	conn, err = req.Connection()
	if err == nil {
		goto DO
	}

	// This needs timeouts because it's fairly likely
	// that something will go wrong :)
	if req.IsTLS() {
		roots, e := x509.SystemCertPool()
		if e != nil {
			return nil, e
		}

		// This library is meant for doing stupid stuff, so skipping cert
		// verification is actually the right thing to do
		conf := &tls.Config{RootCAs: roots, InsecureSkipVerify: true}
		conn, err = tls.DialWithDialer(&net.Dialer{
			Timeout: req.GetTimeout(),
		}, "tcp", req.Address(), conf)

		if err != nil {
			return nil, err
		}
		defer conn.Close()

	} else {
		d := net.Dialer{Timeout: req.GetTimeout()}
		addr := req.Address()
		conn, err = d.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		defer conn.Close()
	}

DO:
	fmt.Fprint(conn, req.String())
	fmt.Fprint(conn, "\r\n")

	return newResponse(conn)
}
