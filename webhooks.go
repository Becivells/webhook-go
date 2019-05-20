package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var db *sql.DB

func init() {
	db, err := sql.Open("mysql", "ty_hl_seevul_com:F5RXLmNmDXSTCJMH@tcp(10.10.100.60:3306)/ty_hl_seevul_com?charset=utf8")
	if err != nil {
		log.Print(err)
	}

	rows, err := db.Query("SELECT id,URL from webinfos")
	if err != nil {
		log.Print(err)
	}
	for rows.Next() {
		var id, url string
		err = rows.Scan(&id, &url)
		fmt.Println(id, url)

	}
}

func Webhook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func cmd() string {
	cmd := exec.Command("ls", "-la")
	stdout, err := cmd.StdoutPipe()
	cmd.Start()
	content, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println(err)
	}
	return string(content) //输出ls命令查看到的内容
}

func main() {
	router := httprouter.New()
	router.GET("/hooksync/:name", Webhook)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
