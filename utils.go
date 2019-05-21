package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// 执行外部命令 第一个是 path ,第二个是命令 其余的是参数
func OsShell(command []string) string {

	cmd := exec.Command(command[1], command[2:]...)
	cmd.Dir = command[0]
	stdout, err := cmd.StdoutPipe()
	cmd.Start()
	content, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println(err)
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

func pullCode(hook *Hook) string {
	var cmdstr string

	return cmdstr
}

//var db *sql.DB
//
//func init() {
//	db, err := sql.Open("mysql", "ty_hl_seevul_com:F5RXLmNmDXSTCJMH@tcp(10.10.100.60:3306)/ty_hl_seevul_com?charset=utf8")
//	if err != nil {
//		log.Print(err)
//	}

//	rows, err := db.Query("SELECT id,URL from webinfos")
//	if err != nil {
//		log.Print(err)
//	}
//	for rows.Next() {
//		var id, url string
//		err = rows.Scan(&id, &url)
//		fmt.Println(id, url)
//
//	}
//}
