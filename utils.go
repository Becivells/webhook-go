package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
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

type Wconfig struct {
	Host     string `yaml:host`
	Port     int    `yaml:port`
	SyncPath string `yaml:syncPath`
	MySQL    struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Password string `yaml:password`
		Port     int    `yaml:port`
		Schema   string `yaml:schema`
	}
	Repoip        []string `yaml:"best_authors,flow"`
	PathWhiteList []string `yaml:pathWhiteList`
	ExecWhiteList []string `yaml:execWhiteList`
	Hooks         []struct {
		Token    string   `yaml:token`
		Name     string   `yaml:name`
		Path     string   `yaml:path`
		Ip       []string `yaml:ip,flow`
		Cmd      string   `yaml:cmd`
		User     string   `yaml:user`
		Interval string   `yaml:interval`
	}
}

func getConfig() {

}

func main() {

	configStr := `
host: 0.0.0.0
port: 21332
syncPath: hookssync
# 数据库配置 仅支持 MySQL
mysql:
  enable: false       # 是否开启数据库支持
  host: "10.10.11.23"   # 主机名
  user: root          # 用户名
  password: Test      # 密码
  port： 3306         # 端口号
  schema: test        # 数据库
repoip: 
  - "123.0.0.1"
  - "177.177.188.2"  # 仓库的 IP 地址
# 路径白名单 要求路径在下列目录
pathWhiteList:
  - "/app/web/test"
  - "/www/wwwroot/.*"

# 运行命令的白名单 采用正则
execWhiteList:
  - '^git pull \w+ \w+$',
    # git fetch origin master && git reset --hard origin/master
  -  '^git fetch \w+ \w+ && git reset --hard \w+/\w+$',
    # git pull origin master && supervisorctl restart webhooks
  -  '^git pull \w+ \w+ && supervisorctl restart \w+$',
    # git fetch origin master && git reset --hard origin/master && supervisorctl restart webhooks
  -  '^git fetch \w+ \w+ && git reset --hard \w+/\w+ && supervisorctl restart \w+$'

# 项目 token 配置
hooks:
  - token: werwerwerwerwerwerwe  # 项目同步的 token
    name: webhooks               # 项目名称
    path: testsdfds              #项目路径
    ip: ['123.0.00.1','122.0.0.1']  # 只允许某个 IP 访问单个项目
    cmd: ''                         # 执行的命令
    user: 'sfs'                     # 执行命令的用户
    interval: 3                     #  间隔执行时间
`

	config := Wconfig{}
	err := yaml.Unmarshal([]byte(configStr), &config)
	if err != nil {
		log.Fatal(err)
	}
	d, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))

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
