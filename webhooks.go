package main

import (
	"flag"
	"github.com/julienschmidt/httprouter"
	"strconv"

	"log"
	"net/http"
	"os"
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")
	flag.StringVar(&configPath, "c", "webhooks.yaml", "webhook的配置文件必须是")
	flag.BoolVar(&printconfig, "p", false, "打印配置文件")
}

func main() {
	flag.Parse()

	initConfig() //初始化配置 或者重载配置

	if h || len(os.Args) == 1 {
		flag.Usage()
	} else if printconfig {
		prinfConfig()
	}
	router := httprouter.New()
	router.GET("/"+SyncPath+"/:gittoken", Webhook)
	log.Printf("Listen %s:%s", ListenIP, strconv.Itoa(ListenPort))
	log.Fatal(http.ListenAndServe(ListenIP+":"+strconv.Itoa(ListenPort), loggingHandler(router)))
}
