package proxy

import (
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type Conn struct {
	mu       sync.Mutex
	ctx      context.Context
	cancel   context.CancelFunc
	finished uint32
	reader   io.Reader
	socket   *websocket.Conn
}

func (c *Conn) finish() bool {
	return atomic.LoadUint32(&c.finished) == 1
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if c.finish() {
		return 0, io.EOF
	}

	if c.reader == nil {
		return 0, fmt.Errorf("invalid reader")
	}

	n, err = c.reader.Read(b)
	return
}

func (c *Conn) ReadFrom(r io.Reader) (int64, error) {
	var n int64
	buf := make([]byte, 1024)
	for {
		rn, err := r.Read(buf)
		if err != nil {
			return n, err
		}
		wn, err := c.Write(buf[:rn])
		if err != nil {
			return n, err
		}
		n += int64(wn)
	}
}

func (c *Conn) WriteTo(w io.Writer) (int64, error) {
	var n int64
	buf := make([]byte, 1024)
	for {
		rn, err := c.Read(buf)
		if err != nil {
			return n, err
		}
		wn, err := w.Write(buf[:rn])
		if err != nil {
			return n, err
		}
		n += int64(wn)
	}

	return n, nil
}

func (c *Conn) Write(b []byte) (n int, err error) {
	if c.finish() {
		return 0, io.EOF
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	err = c.socket.WriteMessage(websocket.BinaryMessage, b)
	if err == nil {
		return len(b), nil
	}

	return 0, err
}

func (c *Conn) Close() error {
	if c.finish() {
		return nil
	}
	atomic.StoreUint32(&c.finished, 1)
	return c.socket.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	if c.socket == nil {
		return nil
	}

	return c.socket.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
	if c.socket == nil {
		return nil
	}
	return c.socket.RemoteAddr()
}

func (c *Conn) SetDeadline(t time.Time) error {
	if c.socket == nil {
		return fmt.Errorf("not found websocket tunnel")
	}
	var err error
	err = c.socket.SetReadDeadline(t)
	err = c.socket.SetWriteDeadline(t)
	return err
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	if c.socket == nil {
		return fmt.Errorf("not found websocket tunnel")
	}
	return c.socket.SetReadDeadline(t)
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	if c.socket == nil {
		return fmt.Errorf("not found websocket tunnel")
	}
	return c.socket.SetWriteDeadline(t)
}
