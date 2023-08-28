package vkit

import (
	account "github.com/vela-ssoc/vela-account"
	arp "github.com/vela-ssoc/vela-arp"
	awk "github.com/vela-ssoc/vela-awk"
	capture "github.com/vela-ssoc/vela-capture"
	component "github.com/vela-ssoc/vela-component"
	cond "github.com/vela-ssoc/vela-cond"
	cpu "github.com/vela-ssoc/vela-cpu"
	crack "github.com/vela-ssoc/vela-crack"
	crontab "github.com/vela-ssoc/vela-crontab"
	crypto "github.com/vela-ssoc/vela-crypto"
	cvs "github.com/vela-ssoc/vela-cvs"
	disk "github.com/vela-ssoc/vela-disk"
	vdns "github.com/vela-ssoc/vela-dns"
	elastic "github.com/vela-ssoc/vela-elastic"
	elkeid "github.com/vela-ssoc/vela-elkeid"
	engine "github.com/vela-ssoc/vela-engine"
	extract "github.com/vela-ssoc/vela-extract"
	fasthttp "github.com/vela-ssoc/vela-fasthttp"
	file "github.com/vela-ssoc/vela-file"
	fsnotify "github.com/vela-ssoc/vela-fsnotify"
	group "github.com/vela-ssoc/vela-group"
	host "github.com/vela-ssoc/vela-host"
	ifconfig "github.com/vela-ssoc/vela-ifconfig"
	ip2region "github.com/vela-ssoc/vela-ip2region"
	kfk "github.com/vela-ssoc/vela-kfk"
	"github.com/vela-ssoc/vela-kit/vela"
	logon "github.com/vela-ssoc/vela-logon"
	memory "github.com/vela-ssoc/vela-memory"
	vnet "github.com/vela-ssoc/vela-net"
	osquery "github.com/vela-ssoc/vela-osquery"
	process "github.com/vela-ssoc/vela-process"
	psnotify "github.com/vela-ssoc/vela-psnotify"
	request "github.com/vela-ssoc/vela-request"
	risk "github.com/vela-ssoc/vela-risk"
	sbom "github.com/vela-ssoc/vela-sbom"
	service "github.com/vela-ssoc/vela-service"
	ss "github.com/vela-ssoc/vela-ss"
	vswitch "github.com/vela-ssoc/vela-switch"
	syslog "github.com/vela-ssoc/vela-syslog"
	vtag "github.com/vela-ssoc/vela-tag"
	tail "github.com/vela-ssoc/vela-tail"
	vtime "github.com/vela-ssoc/vela-time"
	track "github.com/vela-ssoc/vela-track"
)

func (dly *Deploy) withAll(xEnv vela.Environment) {
	if !dly.all {
		return
	}
	vela.WithEnv(xEnv)
	awk.WithEnv(xEnv)
	crypto.WithEnv(xEnv)
	file.WithEnv(xEnv)
	awk.WithEnv(xEnv)
	vswitch.WithEnv(xEnv)
	vtag.WithEnv(xEnv)
	risk.WithEnv(xEnv)
	service.WithEnv(xEnv)
	ifconfig.WithEnv(xEnv)
	cpu.WithEnv(xEnv)
	memory.WithEnv(xEnv)
	disk.WithEnv(xEnv)
	host.WithEnv(xEnv)
	ss.WithEnv(xEnv)
	process.WithEnv(xEnv)
	track.WithEnv(xEnv)
	account.WithEnv(xEnv)
	group.WithEnv(xEnv)
	ip2region.WithEnv(xEnv)
	vtime.WithEnv(xEnv)
	vnet.WithEnv(xEnv)
	cond.WithEnv(xEnv)
	tail.WithEnv(xEnv)
	fsnotify.WithEnv(xEnv)
	psnotify.WithEnv(xEnv)
	fasthttp.WithEnv(xEnv)
	request.WithEnv(xEnv)
	osquery.WithEnv(xEnv)
	component.WithEnv(xEnv)
	vdns.WithEnv(xEnv)
	crontab.WithEnv(xEnv)
	kfk.WithEnv(xEnv)
	crack.WithEnv(xEnv)
	syslog.WithEnv(xEnv)
	elastic.WithEnv(xEnv)
	capture.WithEnv(xEnv)
	logon.WithEnv(xEnv)
	engine.WithEnv(xEnv)
	extract.WithEnv(xEnv)
	sbom.WithEnv(xEnv)
	arp.WithEnv(xEnv)
	cvs.WithEnv(xEnv)
	elkeid.WithEnv(xEnv)
}
