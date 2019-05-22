package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"

	"log"
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
		flag.Usage()
		os.Exit(0)
	} else if printconfig {
		prinfConfig()
	}
	router := httprouter.New()
	router.GET("/"+SyncPath+"/:gittoken", Webhook)
	log.Printf("Listen %s:%s", ListenIP, strconv.Itoa(ListenPort))
	log.Fatal(http.ListenAndServe(ListenIP+":"+strconv.Itoa(ListenPort), loggingHandler(router)))
}
