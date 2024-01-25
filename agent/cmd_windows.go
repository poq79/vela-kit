package agent

import (
	"fmt"
	"github.com/shirou/gopsutil/process"
	"github.com/vela-ssoc/vela-kit/fileutil"
	"github.com/vela-ssoc/vela-kit/stdutil"
	tunnel "github.com/vela-ssoc/vela-tunnel"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	name = "ssc"
)

type program struct {
	shutdown bool
	cmd      exec.Cmd
	fn       construct
	output   *stdutil.Output
}

func (p *program) stop() {
	p.shutdown = true
	if p.cmd.Process == nil {
		p.output.ERR("not found ssc worker")
		return
	}

	err := p.cmd.Process.Kill()
	if err != nil {
		p.output.ERR("stop ssc worker fail %v", err)
		return
	}

	p.output.Info("stop ssc worker succeed")
}

func (p *program) Execute(args []string, r <-chan svc.ChangeRequest, s chan<- svc.Status) (sec bool, errno uint32) {
	const accepts = svc.AcceptStop | svc.AcceptShutdown

	go run(p.fn, &p.cmd, &p.shutdown)

	s <- svc.Status{State: svc.Running, Accepts: accepts}

	for {
		c := <-r
		p.output.ERR("ssc service svc signal %v", c.Cmd)
		switch c.Cmd {
		case svc.Interrogate:
			s <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			s <- svc.Status{State: svc.StopPending}
			p.stop()
			s <- svc.Status{State: svc.Stopped}
			return
		}
	}

	return
}

func Install(_ construct) {
	output := stdutil.New(stdutil.Daemon())
	defer output.Close()

	conn, err := mgr.Connect()
	if err != nil {
		output.ERR(`connet windows service error %v`, err)
		return
	}

	defer func() { _ = conn.Disconnect() }()

	if sc, erx := conn.OpenService(name); erx == nil {
		_ = sc.Close()
		return
	}

	exe, erx := os.Executable()
	if erx != nil {
		output.ERR(`ssc filepath got fail %v`, erx)
		return
	}

	cfg := mgr.Config{
		DisplayName:      "SSOC Sensor",
		Description:      "EastMoney Security Management Platform",
		StartType:        mgr.StartAutomatic,
		DelayedAutoStart: true,
	}

	ss, ers := conn.CreateService(name, exe, cfg, "service")
	if ers != nil {
		output.ERR(`ssc create service error %v`, ers)
		return
	}
	defer func() { _ = ss.Close() }()

	ras := []mgr.RecoveryAction{{Type: mgr.ServiceRestart, Delay: 5 * time.Second}}

	if err = ss.SetRecoveryActions(ras, 0); err != nil {
		output.ERR(`ssc create recovery action %v`, err)
		return
	}

	eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	output.ERR(`ssc install %s succeed`, exe)
}

func Uninstall(_ construct) {
	cnn, _ := mgr.Connect()
	if cnn == nil {
		return
	}
	defer func() { _ = cnn.Disconnect() }()

	ss, _ := cnn.OpenService(name)
	if ss == nil {
		return
	}
	defer func() { _ = ss.Close() }()
	ss.Delete()
}

func Service(fn construct) {
	output := stdutil.New(stdutil.Daemon())
	defer output.Close()

	p := &program{fn: fn, shutdown: false, output: output}

	ok, err := svc.IsWindowsService()
	if err != nil {
		p.output.ERR("ssc service not windows %v", err)
		return
	}

	if !ok {
		return
	}

	err = svc.Run(name, p)
	if err == nil {
		p.output.ERR("ssc service exit error %v", err)
		return
	}

	p.output.ERR("ssc service exit")
	return

}

func NewSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}

func exeWalk(current string) []string {
	var ret []string
	var mask []fs.FileInfo

	filepath.Walk(current, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if filepath.Ext(path) != ".exe" ||
			!strings.HasPrefix(info.Name(), "ssc-") ||
			filepath.Dir(path) != current {
			return nil
		}

		cAttr := info.Sys().(*syscall.Win32FileAttributeData)
		for i, stat := range mask {
			fAttr := stat.Sys().(*syscall.Win32FileAttributeData)
			if fAttr.CreationTime.Nanoseconds() < cAttr.CreationTime.Nanoseconds() {
				e := append([]string{}, ret[i:]...)
				s := append(ret[0:i], path)
				ret = append(s, e...)

				em := append([]os.FileInfo{}, mask[i:]...)
				sm := append(mask[0:i], info)
				mask = append(sm, em...)
				return nil
			}
		}

		ret = append(ret, path)
		mask = append(mask, info)
		return nil
	})

	return ret
}

func executable(output *stdutil.Output) string {
	exe, err := os.Executable()
	if err != nil {
		output.ERR(`ssc executable got fail %v`, err)
		return ""
	}
	current := filepath.Dir(exe)
	files := exeWalk(current)

	if len(files) == 0 {
		output.ERR(`not found ssc file`)
		return ""
	}

	if len(files) == 1 && files[0] == exe {
		path := filepath.Join(current, "ssc-worker.exe")
		fileutil.CopyFile(exe, path)
		files[0] = path
	}

	for _, path := range files {
		hi, e := tunnel.ReadHide(path)
		if e == nil {
			output.ERR(`ssc %s binary code succeed %+v`, path, hi)
			return path
		}
		output.ERR(`ssc %s binary decode error %v`, path, e)
	}
	output.ERR(`%+v not found valid ssc exe`, files)
	return ""
}

func Exe() string {
	exe, _ := os.Executable()
	return exe
}

func killall(output *stdutil.Output) {
	ps, err := process.Pids()
	if err != nil {
		output.ERR("not found process %v", err)
		return
	}

	for _, pid := range ps {
		pr, er := process.NewProcess(pid)
		if er != nil {
			output.ERR("not found %d process %v", pid, er)
			continue
		}

		pName, er := pr.Name()
		if er != nil {
			output.ERR("not found %d process name %v", pid, er)
			continue
		}

		if !strings.HasPrefix(pName, "ssc-") {
			continue
		}

		if int(pid) == os.Getpid() {
			continue
		}

		if e := pr.Kill(); e != nil {
			output.ERR("process %d %s kill fail %v", pid, pName, e)
		}
	}
}

func StartService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("service")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

func ControlService(name string, c svc.Cmd, to svc.State) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}

func Upgrade() {
	output := stdutil.New(stdutil.Upgrade())
	defer output.Close()

	output.Info("ssc upgrade start ...")
	exe, err := os.Executable()
	if err != nil {
		output.ERR("executable %v", err)
		return
	}

	err = ControlService(name, svc.Stop, svc.Stopped)
	if err != nil {
		output.ERR("control service fail %v", err)
		return
	}

	output.Info("sc stop scc succeed")
	killall(output)

	size, err := fileutil.CopyFile(exe, "ssc-mgt.exe")
	if err != nil {
		output.ERR("ssc-mgt.exe copy file fail %v", err)
		return
	}
	output.Info("ssc-mgt.exe upgrade succeed size:%d", size)

	StartService(name)
	output.Info("sc start scc succeed")
}
