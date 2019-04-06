package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func CheckErr(err error){
	if err != nil{
		panic(err)
	}
}


func main() {
	db,err := sql.Open("mysql","root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	CheckErr(err)
	stmt,err := db.Prepare("SELECT `id`,`name` FROM `foo`")
	conn,err := db.Conn(context.Background())
	conn,err = db.Conn(context.Background())

	fmt.Println(conn)
	rows,err := stmt.Query()
	rows,err = stmt.Query()
	stmt.Exec()
	CheckErr(err)
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id,&name)
		fmt.Println(id,name)
		CheckErr(err)
	}
	err = stmt.Close()
	CheckErr(err)
	err = db.Close()
	CheckErr(err)
}
