package kind

import (
	"context"
	"fmt"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/exception"
	"github.com/vela-ssoc/vela-kit/vela"
	"net"
	"sync"
	"sync/atomic"
)

type Accept func(context.Context, net.Conn) error

type Listener struct {
	done uint32
	refs uint32
	xEnv vela.Environment
	bind auxlib.URL
	fd   []net.Listener
	pool *sync.Map
	ctx  context.Context
	stop context.CancelFunc
}

type Session struct {
	ctx    context.Context
	cancel context.CancelFunc
	conn   net.Conn
}

func (ln *Listener) CloseActiveConn() {
	ln.stop()

	ln.ctx, ln.stop = context.WithCancel(ln.xEnv.Context())

	//ln.pool.Range(func(key, value any) bool {
	//	sess := value.(*Session)
	//	conn := sess.conn
	//	ref := key.(uint32)
	//	sess.cancel()

	//	if err := conn.Close(); err != nil {
	//		ln.xEnv.Errorf("ref=%d %s -> %s close fail %v", ref, conn.RemoteAddr().String(), conn.LocalAddr().String(), err)
	//	} else {
	//		ln.xEnv.Errorf("ref=%d %s -> %s close succeed", ref, conn.RemoteAddr().String(), conn.LocalAddr().String())

	//	}
	//	ln.pool.Delete(ref)
	//	return false
	//})
}

func (ln *Listener) Done() bool {
	return atomic.LoadUint32(&ln.done) == 1
}

func (ln *Listener) shutdown() {
	atomic.StoreUint32(&ln.done, 1)
}

func (ln *Listener) store(sess *Session) uint32 {
	ref := atomic.AddUint32(&ln.refs, 1)
	ln.pool.Store(ref, sess)
	return ref
}

func (ln *Listener) destroy(ref uint32) {
	ln.pool.Delete(ref)
}

func (ln *Listener) session(accept Accept, conn net.Conn) {

	ctx, cancel := context.WithCancel(ln.ctx)
	defer cancel()

	//ref := ln.store(&Session{
	//	ctx:    ctx,
	//	conn:   conn,
	//	cancel: cancel,
	//})

	if e := accept(ctx, conn); e != nil {
		ln.xEnv.Errorf("%s listen handler failure , error %v", conn.LocalAddr().String(), e)
	} else {
		//ln.destroy(ref)
		ln.xEnv.Errorf("%s -> %s accept handle exit", conn.RemoteAddr().String(), conn.LocalAddr().String())
	}
}

func (ln *Listener) loop(accept Accept, sock net.Listener) {
	defer sock.Close()

	for {

		if ln.Done() {
			ln.xEnv.Errorf("%s listen done.", sock.Addr().String())
			return
		}

		conn, err := sock.Accept()
		if err == nil {
			go ln.session(accept, conn)
			continue
		}

		ln.xEnv.Errorf("%s listen accept fail %v", sock.Addr().String(), err)
		return
	}

}

func (ln *Listener) multipleH(accept Accept) error {
	for _, sock := range ln.fd {
		go ln.loop(accept, sock)
	}
	<-ln.ctx.Done()
	ln.xEnv.Errorf("%s multiple handle exit", ln.bind.String())
	return nil
}

func (ln *Listener) OnAccept(accept Accept) error {

	n := len(ln.fd)
	if n < 1 {
		return fmt.Errorf("not found ative listen fd")
	}

	if n == 1 {
		go ln.loop(accept, ln.fd[0])
		return nil
	} else {
		return ln.multipleH(accept)
	}
}

func (ln *Listener) Close() error {
	if ln == nil {
		return nil
	}

	ln.shutdown()

	ln.stop()
	errs := exception.New()
	for _, fd := range ln.fd {
		errs.Try(fd.Addr().String(), fd.Close())
	}
	ln.fd = nil
	return errs.Wrap()
}

func (ln *Listener) single() error {
	fd, err := net.Listen(ln.bind.Scheme(), ln.bind.Host())
	if err != nil {
		return err
	}
	ln.fd = []net.Listener{fd}
	return nil
}

// multiple tcp://192.168.0.1/?port=1024,65535&exclude=1,2,3,4
func (ln *Listener) multiple() error {
	ps := ln.bind.Ports()
	n := len(ps)
	if n == 0 {
		return fmt.Errorf("%s not found listen", ln.bind.String())
	}

	for i := 0; i < n; i++ {
		port := ps[i]
		fd, e := net.Listen(ln.bind.Scheme(), fmt.Sprintf("%s:%d", ln.bind.Hostname(), port))
		if e != nil {
			return fmt.Errorf("listen %s://%s:%d error %v", ln.bind.Scheme(), ln.bind.Hostname(), port, e)
			continue
		}
		ln.fd = append(ln.fd, fd)
	}

	if len(ln.fd) == 0 {
		return fmt.Errorf("%s listen fail", ln.bind.String())
	}

	return nil
}

func (ln *Listener) Start() error {
	if ln.bind.Port() != 0 {
		return ln.single()
	}

	return ln.multiple()
}

func Listen(env vela.Environment, bind auxlib.URL) (*Listener, error) {
	ctx, stop := context.WithCancel(env.Context())
	ln := &Listener{
		bind: bind,
		ctx:  ctx,
		stop: stop,
		done: 0,
		xEnv: env,
		pool: new(sync.Map),
	}

	return ln, ln.Start()
}
