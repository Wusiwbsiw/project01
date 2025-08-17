package main

import (
	"fmt"
	"https_operation/handlers"
	"log"
	"net/http"
	"sql_operation/db"
)

func main() {
	database, err := db.InitDatabase("login_project01", "Aa!123456", "localhost:3306", "PROJECT01")
	if err != nil {
		log.Fatalf("无法初始化数据库连接:%v", err)
	}
	defer database.Close()
	fmt.Println("数据库连接初始化成功")
	mux := http.NewServeMux()
	Userheadler := &handlers.UserHandler{DB: database}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "欢迎来到我们的 HTTPS API 服务器！")
	})
	mux.HandleFunc("/api/user/register", Userheadler.RegisterHTTP)
	mux.HandleFunc("/api/user/login", Userheadler.LoignHTTP)
	mux.HandleFunc("/api/user/reset-name", Userheadler.ResetNameHTTP)
	mux.HandleFunc("/api/user/reset-password", Userheadler.ResetPasswordHTTP)
	addr := ":8443" // 监听 8443 端口
	fmt.Printf("服务器正在启动，监听地址: https://localhost%s\n", addr)
	err = http.ListenAndServeTLS(addr, "server.crt", "server.key", mux)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
