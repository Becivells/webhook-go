package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

type Dsn struct {
	Enable   bool   `yaml:enable`
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Schema   string `yaml:"schema"`
	Charset  string `yaml:"charset"`
}

type Hook struct {
	Token    string   `yaml:"token"`
	Name     string   `yaml:"name"`
	Path     string   `yaml:"path"`
	Ip       []string `yaml:"ip,flow"`
	Cmd      string   `yaml:"cmd"`
	User     string   `yaml:"user"`
	Interval int      `yaml:"interval"`
}
type Wconfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	SyncPath      string `yaml:"hookPath"`
	MySQL         *Dsn   `yaml:"mysql"`
	Table         *Table
	RepoIp        []string `yaml:"repoIp,flow"`
	PathWhiteList []string `yaml:"pathWhiteList,flow"`
	ExecWhiteList []string `yaml:"execWhiteList,flow"`
	Hooks         []*Hook
}

//table:
//token: "token"
//name: "name"
//path: "path"
//ip: "ip"
//cmd: "cmd"
//user: "user"
//interval: "interval"

type Table struct {
	TableName    string `yaml:"tableName"`
	Token        string `yaml:"gittoken"`
	Name         string `yaml:"gitname"`
	Path         string `yaml:"pullpath"`
	pullip       string `yaml:"pullip"`
	pullcmd      string `yaml:"pullcmd"`
	pulluser     string `yaml:"pulluser"`
	pullinterval string `yaml:"pullinterval"`
}

//SELECT TOKEN
func (table *Table) FormatSQL() string {
	buf := bytes.Buffer{}
	buf.WriteString("SELECT " + table.Token + " AS token," + table.Name + " as name," + table.Path + " as path,")
	buf.WriteString(table.ip + " as ip," + table.cmd + " as cmd," + table.user + " as user," + table.interval + " as interval")
	buf.WriteString(" from " + table.TableName)
	return buf.String()

}

var (
	BaseDir     string           //路径
	h           bool             //帮助信息
	configPath  string           //配置文件路径
	ListenIP    string           //侦听IP
	ListenPort  int              // 侦听 端口
	config      *Wconfig         // 配置信息
	err         error            // 错误信息
	printconfig bool             // 打印配置
	dsn         string           // 数据库链接
	RepoIP      map[string]int   // git 仓库地址 使用字典为啥要用字典 自己不想写
	token       map[string]*Hook //token
	SyncPath    string
)

//得到 dsn 字符串
func (dsn *Dsn) FormatDSN() string {
	var buf bytes.Buffer
	buf.WriteString(dsn.User + ":" + dsn.Password + "@")
	buf.WriteString("tcp(" + dsn.Host + ":" + dsn.Port + ")/")
	buf.WriteString(dsn.Schema + "?" + "charset=" + dsn.Charset)
	return buf.String()
}

func initConfig() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	BaseDir = filepath.Dir(ex)
	configPath = path.Join(BaseDir, configPath)
	config = &Wconfig{}
	configData, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		panic(err)
	}

	ListenIP = config.Host
	ListenPort = config.Port
	SyncPath = config.SyncPath
	fmt.Println(SyncPath)
	RepoIP = make(map[string]int, 20) //初始化

	//git仓库 ip
	for _, rip := range config.RepoIp {
		fmt.Println(rip)
		RepoIP[rip] = 1
	}
	dsn = config.MySQL.FormatDSN() //数据库账号密码配置
	token = make(map[string]*Hook, 20)
	for _, hook := range config.Hooks {
		token[hook.Token] = hook
	}

	fmt.Println(config.Table.FormatSQL())

}

func prinfConfig() {
	d, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
	os.Exit(0)
}
