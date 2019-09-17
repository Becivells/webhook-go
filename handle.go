package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"time"
)

func printVersion() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Compile: %s\n", Compile)
	fmt.Printf("Branch: %s\n", Branch)
	fmt.Printf("GitDirty: %s\n", GitDirty)
	fmt.Printf("DevPath: %s\n", DevPath)
	fmt.Println("\n")
}

func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("Comleted %s in %v", r.URL.Path, time.Since(start))
	})
}

func ShellWebhook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gittoken := ps.ByName("gittoken")
	//判断token长度
	if len(gittoken) > config.TokenMaxLength || len(gittoken) < config.TokenMinLength {
		w.Write([]byte("token 长度不符合"))
		w.WriteHeader(400)

	} else {
		if mtoken, ok := token[gittoken]; ok {
			remoteIP := parseIP(r.RemoteAddr) //获取 IP 地址
			flag := false

			//允许 git 仓库的IP访问范围
			if _, ok := RepoIP[remoteIP]; ok {
				flag = true

			} else {
				//判断是否在该token的访问范围
				for _, val := range mtoken.Ip {
					if val == remoteIP {
						flag = true
						break
					}
				}
			}

			//在访问IP范围内
			if flag == true {
				//判断访问频率是否符合
				if slimit, ok := LimitIP[remoteIP+gittoken]; ok {
					elimit := time.Now()
					itime64 := int64(elimit.Sub(slimit)) / int64(time.Second)

					if itime64 > int64(mtoken.Interval) {

						w.Write([]byte("-------------git同步------------\r\n"))
						w.Write([]byte(pullCode(mtoken)))
						LimitIP[remoteIP+gittoken] = elimit
					} else {
						w.Write([]byte("访问过于频繁"))
						w.WriteHeader(403)
					}

				} else { //判断访问频率是否符合
					w.Write([]byte("禁止访问"))
					w.WriteHeader(403)
				}

			} else { //在访问IP范围内
				w.Write([]byte("禁止访问"))
				w.WriteHeader(403)
			}
		} else {
			//同步判断一下
			if config.MySQL.Enable {
				Qtoken := Querytoken(gittoken)
				if Qtoken.Token == gittoken {
					remoteIP := parseIP(r.RemoteAddr) //获取 IP 地址
					flag := false

					//允许 git 仓库的IP访问范围
					if _, ok := RepoIP[remoteIP]; ok {
						flag = true

					} else {
						//判断是否在该token的访问范围
						for _, val := range Qtoken.Ip {
							if val == remoteIP {
								flag = true
								break
							}
						}
					}
					if flag == true {
						w.Write([]byte("-------------git同步------------\r\n"))
						w.Write([]byte(pullCode(Qtoken)))
					} else {
						w.Write([]byte("禁止访问"))
						w.WriteHeader(403)
					}

				} else {
					w.Write([]byte("token 无效"))
					w.WriteHeader(404)
				}
			} else {
				w.Write([]byte("token 无效"))
				w.WriteHeader(404)
			}

		}
	}
}
