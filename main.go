package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

type user struct {
	id   int
	age  int
	name string
}

func main() {
	init_DB()
	insertRow()
	//queryRow()
	queryRows()

}

func init_DB() {
	dsn := "root:Wwq123456_@tcp(81.68.251.42:3306)/GOSQLTest?charset=utf8"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		fmt.Printf("ok\n")
	}
}
func queryRow() {
	sqlStr := "select id,name,age from user where id=?"
	var u user
	err = db.QueryRow(sqlStr, 1).Scan(&u.id, &u.name, &u.age)

	if err != nil {
		fmt.Printf("query faild,err:%v\n", err)
		return
	}
	fmt.Printf("id:%d,name:%s,age:%d", u.id, u.name, u.age)
}
func queryRows() {
	sqlStr := "select id,name,age from user where id>?"

	rows, err := db.Query(sqlStr, 0)
	if err != nil {
		fmt.Printf("query faild,err:%v\n", err)
		return
	}

	for rows.Next() {
		var u user
		err = rows.Scan(&u.id, &u.name, &u.age)
		if err != nil {
			fmt.Printf("scan faild,err:%v\n", err)
			return
		}
		fmt.Printf("id:%d,name:%s,age:%d\n", u.id, u.name, u.age)
	}
}

func insertRow() {
	sqlStr := "insert into user(name, age) values(?,?)"

	ret, err := db.Exec(sqlStr, "小王总", 28)


	if err != nil {
		fmt.Printf("insert failed,err:%v\n", err)
		return
	}
	var theID int64
	theID, err = ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastInserID faild,err:%v\n", err)
		return
	}
	fmt.Printf("Insert successfully,the id is %d\n", theID)
}
