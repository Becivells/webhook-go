package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Comleted %s in %v", r.URL.Path, time.Since(start))
	})
}

func Webhook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gittoken := ps.ByName("gittoken")
	if mtoken, ok := token[gittoken]; ok {
		remoteIP := parseIP("[::1]:50380") //获取 IP 地址
		flag := false
		//判断是否在该token的访问范围
		for _, val := range mtoken.Ip {
			if val == remoteIP {
				flag = true
			}
		}
		//允许git仓库的IP访问范围
		if _, ok := RepoIP[remoteIP]; ok {
			flag = true
		}
		if flag == true {
			w.Write([]byte("-------------git同步------------\r\n"))
			w.Write([]byte(pullCode(mtoken)))

		} else {
			w.Write([]byte("禁止访问"))
		}
	} else {
		//同步判断一下
		w.Write([]byte("token 无效"))
	}

	//fmt.Fprint(w, "Welcome!\n"+r.RemoteAddr +gittoken)
}
