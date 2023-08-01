package auxlib

import (
	"context"
	"errors"
	"github.com/vela-ssoc/vela-kit/buffer"
	"io"
)

func Push(w io.Writer, v interface{}) (int, error) {
	chunk, err := ToStringE(v)
	if err != nil {
		return 0, err
	}
	return w.Write(S2B(chunk))
}

var (
	errInvalidWrite = errors.New("invalid write result")
	ErrShortBuffer  = errors.New("short buffer")
	ErrShortWrite   = errors.New("short write")
	EOF             = errors.New("EOF")
)

func copyBuffer(ctx context.Context, dst io.Writer, src io.Reader, buf *buffer.Byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	//if wt, ok := src.(io.WriterTo); ok {
	//	return wt.WriteTo(dst)
	//}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	//if rt, ok := dst.(io.ReaderFrom); ok {
	//	return rt.ReadFrom(src)
	//}

	data := buf.Bytes()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			nr, er := src.Read(data)
			if nr > 0 {
				nw, ew := dst.Write(data[0:nr])
				if nw < 0 || nr < nw {
					nw = 0
					if ew == nil {
						ew = errInvalidWrite
					}
				}
				written += int64(nw)
				if ew != nil {
					err = ew
					goto done
				}

				if nr != nw {
					err = ErrShortWrite
				}
			}
			if er != nil {
				if er != EOF {
					err = er
				}
				goto done
			}
		}
	}

done:
	return written, err
}

func Copy(ctx context.Context, dst io.Writer, src io.Reader) (written int64, err error) {
	buf := &buffer.Byte{B: make([]byte, 4096)}
	defer func() {
		buffer.Put(buf)
	}()

	return copyBuffer(ctx, dst, src, buf)
}
