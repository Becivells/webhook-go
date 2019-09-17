package main

import (
	"flag"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&configPath, "c", "webhooks.yaml", "webhook 的配置文件")
	flag.BoolVar(&printconfig, "p", false, "打印配置文件需要配合 -c 使用")
}

func main() {
	flag.Parse()

	initConfig() //初始化配置 或者重载配置

	if h || len(os.Args) == 1 {
		printVersion()
		flag.Usage()
		os.Exit(0)
	} else if printconfig {
		prinfConfig()
	}
	router := httprouter.New()
	router.GET(fmt.Sprintf("/%s/:gittoken", SyncPath), ShellWebhook)  //使用 shell 方式
	router.POST(fmt.Sprintf("/%s/:gittoken", SyncPath), ShellWebhook) //gogs 等临时使用
	log.Printf(fmt.Sprintf("Listen: %s:%d", ListenIP, ListenPort))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", ListenIP, ListenPort), loggingHandler(router)))
}
