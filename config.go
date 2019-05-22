package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Dsn struct {
	Enable   bool   `yaml:"enable"`
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Schema   string `yaml:"schema"`
	Charset  string `yaml:"charset"`
	TimeOut  string `yaml:"timeout"`
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

type Table struct {
	TableName string `yaml:"tableName"`
	Token     string `yaml:"token"`
	Name      string `yaml:"name"`
	Path      string `yaml:"path"`
	Ip        string `yaml:"ip"`
	Cmd       string `yaml:"cmd"`
	User      string `yaml:"user"`
	Interval  string `yaml:"interval"`
}

//SELECT TOKEN
func (table *Table) FormatSQL() string {
	buf := bytes.Buffer{}
	buf.WriteString("SELECT `" + table.Token + "` AS `token`,`" + table.Name + "` as `name`,`" + table.Path + "` as `path`,`")
	buf.WriteString(table.Ip + "` as `ip`,`" + table.Cmd + "` as `cmd`,`" + table.User + "` as `user`,`" + table.Interval + "` as `interval`")
	buf.WriteString(" from " + table.TableName + " where " + table.Token + "=?")
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
	LimitIP     map[string]time.Time
	db          *sql.DB
)

//得到 dsn 字符串
func (dsn *Dsn) FormatDSN() string {
	var buf bytes.Buffer
	buf.WriteString(dsn.User + ":" + dsn.Password + "@")
	buf.WriteString("tcp(" + dsn.Host + ":" + dsn.Port + ")/")
	buf.WriteString(dsn.Schema + "?" + "charset=" + dsn.Charset)
	buf.WriteString("&timeout=" + dsn.TimeOut)
	return buf.String()
}

//TableName string `yaml:"tableName"`
//Token     string `yaml:"token"`
//Name      string `yaml:"name"`
//Path      string `yaml:"path"`
//Ip        string `yaml:"ip"`
//Cmd       string `yaml:"cmd"`
//User      string `yaml:"user"`
//Interval  string `yaml:"interval"`
func enabelMysql() {
	if config.MySQL.Enable {
		log.Println("use MySQL and yaml file...")
		if db != nil {
			db.Close()
		}
		//enabelMysql()
		db, err = sql.Open("mysql", config.MySQL.FormatDSN())
		if err != nil {
			panic(err)
		}
		var version string
		rows, err := db.Query("select version()")
		if err != nil {
			log.Fatal(err)
		} else {
			for rows.Next() {
				rows.Scan(&version)
				log.Println("MySQL Version: " + version)
			}

		}

	} else {
		log.Println("using yaml ...")
	}
}

func Querytoken(tokeninfo string) *Hook {

	var (
		token    string
		name     string
		spath    string
		ip       string
		cmd      string
		user     string
		interval int
	)
	row := db.QueryRow(config.Table.FormatSQL(), tokeninfo)
	err = row.Scan(&token, &name, &spath, &ip, &cmd, &user, &interval)

	if err != nil {

		log.Printf(fmt.Sprintf("%s", err))
	}

	return &Hook{
		Token:    token,
		Name:     name,
		Path:     spath,
		Ip:       strings.Split(ip, ","),
		Cmd:      cmd,
		User:     user,
		Interval: interval,
	}
}
func initConfig() {

	config = &Wconfig{}
	configData, err := ioutil.ReadFile(configPath)

	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(configData, config)
	if err != nil {
		panic(err)
	}
	//配置初始化
	enabelMysql() //数据库
	ListenIP = config.Host
	ListenPort = config.Port
	SyncPath = config.SyncPath
	LimitIP = make(map[string]time.Time, 20) // ip 限制速率
	RepoIP = make(map[string]int, 20)        //初始化
	// git仓库 ip
	for _, rip := range config.RepoIp {
		RepoIP[rip] = 1
	}
	dsn = config.MySQL.FormatDSN() //数据库账号密码配置
	token = make(map[string]*Hook, 20)
	for _, hook := range config.Hooks {
		token[hook.Token] = hook
		//仓库ip 和repo
		for _, rip := range config.RepoIp {
			LimitIP[rip+hook.Token] = getTime()
		}
		for _, rip := range hook.Ip {
			LimitIP[rip+hook.Token] = getTime() //初始化 一开始就可以用
		}
	}

}

func getTime() time.Time {
	now := time.Now()
	h, _ := time.ParseDuration("-1h")
	return now.Add(24 * h)

}

func prinfConfig() {
	d, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t yamlfile:\n%s\n\n", string(d))
	os.Exit(0)
}
