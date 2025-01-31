package main

import (
	"database/sql"
	"fmt"
	avatica "github.com/apache/calcite-avatica-go/v5"
	"log"
)

func checkErr(remark string, err error) {
	if err != nil {
		log.Println(remark + ":" + err.Error())
	}
}

func main() {
	databaseUrl := "http://localhost:30060" // 这里的链接地址与lindorm-cli的链接地址比，需要去掉http之前的字符串
	conn := avatica.NewConnector(databaseUrl).(*avatica.Connector)
	conn.Info = map[string]string{
		"user":     "sql",     // 数据库用户名
		"password": "test",    // 数据库密码
		"database": "default", // 初始化连接指定的默认database
	}

	db := sql.OpenDB(conn)

	// 创建表
	_, err := db.Exec("create table if not exists user_test(id int, name varchar,age int, primary key(id))")
	checkErr("create table", err)

	// 添加数据
	_, err = db.Exec("upsert into user_test(id,name,age) values(?,?,?)", 1, "zhangsan", 17)
	_, err = db.Exec("upsert into user_test(id,name,age) values(?,?,?)", 2, "lisi", 18)
	_, err = db.Exec("upsert into user_test(id,name,age) values(?,?,?)", 3, "wanger", 19)
	checkErr("insert", err)

	// 删除数据
	_, err = db.Exec("delete from user_test where id=?", 2)
	checkErr("delete", err)

	// 修改数据
	_, err = db.Exec("upsert into user_test(id,name,age) values(?,?,?)", 1, "zhangsan", 7)
	checkErr("update", err)

	// 获取一行数据
	row := db.QueryRow("select * from user_test where id=?", 1)
	var id int
	var name string
	var age int
	err = row.Scan(&id, &name, &age)
	checkErr("scan", err)
	fmt.Println("id:", id, "name:", name, "age:", age)

	// 查询数据
	querySql := "select * from user_test"
	rows, err := db.Query(querySql)
	fmt.Println(querySql)
	checkErr("select", err)
	for rows.Next() {
		err = rows.Scan(&id, &name, &age)
		checkErr("scan", err)
		fmt.Println("id:", id, "name:", name, "age:", age)
	}
	checkErr("close rows", rows.Close())
}
