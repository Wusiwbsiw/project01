package main

import (
	"log"
	"sql_operation/db" // 导入我们自己写的包 (模块名/包名)
)

// User 结构体用于映射 users 表
type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	// 1. 初始化数据库连接
	// !!! 重要：请将这里的密码替换成你自己的MySQL root密码 !!!
	database, err := db.NewDatabase("root", "Aa001111", "localhost:3306", "") // 初始连接时不指定数据库
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

}
