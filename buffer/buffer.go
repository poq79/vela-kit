package buffer

import "io"

// Byte provides byte buffer, which can be used for minimizing
// memory allocations.
//
// Byte may be used with functions appending data to the given []byte
// slice. See example code for details.
//
// Use Get for obtaining an empty byte buffer.
type Byte struct {

	// B is a byte buffer to use in append-like workloads.
	// See example code for details.
	B []byte
}

// Len returns the size of the byte buffer.
func (b *Byte) Len() int {
	return len(b.B)
}

// ReadFrom implements io.ReaderFrom.
//
// The function appends all the data read from r to b.
func (b *Byte) ReadFrom(r io.Reader) (int64, error) {
	p := b.B
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}
	for {
		if n == nMax {
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:])
		n += int64(nn)
		if err != nil {
			b.B = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}

// WriteTo implements io.WriterTo.
func (b *Byte) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.B)
	return int64(n), err
}

// Bytes returns b.B, i.e. all the bytes accumulated in the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
func (b *Byte) Bytes() []byte {
	return b.B
}

// Write implements io.Writer - it appends p to Byte.B
func (b *Byte) Write(p []byte) (int, error) {
	b.B = append(b.B, p...)
	return len(p), nil
}

// WriteByte appends the byte c to the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
//
// The function always returns nil.
func (b *Byte) WriteByte(c byte) error {
	b.B = append(b.B, c)
	return nil
}

// WriteString appends s to Byte.B.
func (b *Byte) WriteString(s string) (int, error) {
	b.B = append(b.B, s...)
	return len(s), nil
}

// Set sets Byte.B to p.
func (b *Byte) Set(p []byte) {
	b.B = append(b.B[:0], p...)
}

// SetString sets Byte.B to s.
func (b *Byte) SetString(s string) {
	b.B = append(b.B[:0], s...)
}

// String returns string representation of Byte.B.
func (b *Byte) String() string {
	return string(b.B)
}

// Reset makes Byte.B empty.
func (b *Byte) Reset() {
	b.B = b.B[:0]
}
