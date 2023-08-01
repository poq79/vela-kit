package node

import (
	"github.com/vela-ssoc/vela-kit/auxlib"
	"os"
)

var shell = `#!/usr/bin/env bash

prefix="/usr/local/ssoc"

# 安装客户端程序
bin=$prefix/ssc

# 关闭现在的进程
function kill_w() {
  for i in $(ps -lef | grep -v "grep" | grep "/ssoc/ssc" | awk '{print $4}')
  do
      info=$(readlink /proc/$i/exe)
      kill $i
      echo "kill $info succeed"
  done
}

function clear_rc_local() {
    sed -i '/^initctl start ssc/d'       /etc/rc.local
    sed -i '/^cd*ssoc$/d'             /etc/rc.local
    sed -i '/^nohup.*\/ssoc\/ssc.*/d' /etc/rc.local
    sed -i '/^cd -/d'                     /etc/rc.local
}

function install_rc_local() {
    echo "cd $prefix" >> /etc/rc.d/rc.local
    echo "nohup $bin worker -d &" >> /etc/rc.d/rc.local
    echo "cd -"  >> /etc/rc.d/rc.local
}

#判断是否更行
if [ -f $bin ]; then
    $bin uninstall
    mv $bin $bin.old
    echo "$bin backup $bin.old succeed"
fi

# 安装客户端服务
chmod +x $bin

# 安装客户端
$bin install

# 启动程序
if [ $(command -v systemctl) ]; then
    systemctl enable ssc
    systemctl daemon-reload

    kill_w
    systemctl restart ssc
    echo "systemctl ssc running"

elif [ -f /etc/init.d/ssc ]; then
    service  ssc stop
    sleep 1s

    kill_w
    service ssc start
    chkconfig ssc on
    echo "service ssc running"

elif [ -f /etc/init/ssc.conf ]; then
    count=$(initctl status ssc | grep "start/running" | grep -v "grep" | wc -l)
    if [ $count -eq 1 ]; then
       initctl stop ssc
    fi

    sleep 1s

    kill_w
    initctl start ssc
    echo "initctl ssc running"
    chmod +x /etc/rc.d/rc.local
    count=$(cat /etc/rc.d/rc.local | grep "initctl start ssc" | grep -v "grep" | wc -l)
    if [ $count -eq 0 ]; then
        echo "initctl start ssc" >> /etc/rc.d/rc.local
    fi

else
    #自启动
    chmod +x /etc/rc.d/rc.local
    sleep 1s

    count=$(cat /etc/rc.d/rc.local | grep "/ssoc/ssc" | grep -v "grep" | wc -l);
    if [ $count -eq 0 ]; then
        install_rc_local
    else
        clear_rc_local
        sleep 1s
        install_rc_local
    fi

    kill_w
    nohup $bin worker -d &
    echo "nohup ssc running"

fi`

func NotUpgrade(exe string) bool {
	return false
}

func (nd *node) hot(save, abs string, out func(string, ...interface{})) error {
	_, std := auxlib.Stdout()
	defer std.Close()

	//添加执行权限
	if err := os.Chmod(save, 0755); err != nil {
		out("chmod 文件%s错误: %v", abs, err)
		return err
	}

	// 刚刚下载的文件覆盖掉运行的文件名
	if err := os.Remove(abs); err != nil {
		out("删除文件%s错误: %v", abs, err)
		return err
	}

	if err := os.Rename(save, abs); err != nil {
		out("升级包 %s -> %s 覆盖失败: %v", save, abs, err)
		return err
	}

	os.Exit(0)
	return nil
}
