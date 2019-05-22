package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// 执行外部命令 第一个是 path ,第二个是命令 其余的是参数
func OsShell(command []string) string {

	cmd := exec.Command(command[1], command[2:]...)
	cmd.Dir = command[0]
	content, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return string(content) //输出ls命令查看到的内容
}

func parseIP(ipstr string) string {
	ipdata := strings.Split(ipstr, ":")
	var ip string
	if len(ipdata) > 2 {
		ip = strings.Join(ipdata[0:len(ipdata)-1], ":")
	} else {
		ip = ipdata[0]
	}
	return ip
}

func Validate(reg []string, cmd string) bool {

	for _, val := range reg {
		if regexp.MustCompile(val).MatchString(cmd) {
			return true
			break
		}
	}
	return false
}

func pullCode(hook *Hook) string {
	//检查路径是否合法
	if !Validate(config.PathWhiteList, hook.Path) {
		return fmt.Sprintf("非法操作不在路径白名单内")
	}
	//检查命令是否合法
	if !Validate(config.ExecWhiteList, hook.Cmd) {
		return "非法操作不在命令白名单内"
	}

	var shell []string
	shell = append(shell, hook.Path)
	shell = append(shell, strings.Split(hook.Cmd, " ")...)

	return OsShell(shell)
}
