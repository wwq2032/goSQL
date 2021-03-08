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
	//insertRow()
	//queryRow()
	//deleteRows()
	//update()
	queryRows()
	fmt.Printf("-----------------------\n")

	sqlInjectDemo("xxx' or 1=1#")
	sqlInjectDemo("xxx' union select * from user #")
	sqlInjectDemo("xxx' and (select count(*) from user) <10 #")

	//prepareInsert()
	//prepareQuery()

	db.Close()

}

func init_DB() {
	dsn := "test:Ww*****6!@tcp(81.68.251.42:3306)/GOSQLTest?charset=utf8"
	//dsn := "root:Wwq123456_@tcp(81.68.251.42:3306)/GOSQLTest?charset=utf8"
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

func deleteRows() {
	sqlStr := "delete from user where name=?"

	ret, err := db.Exec(sqlStr, "小王总")

	if err != nil {
		fmt.Printf("delete failed,err:%v\n", err)
		return
	}

	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected faild,err:%v\n", err)
		return
	}
	fmt.Printf("delete successfully,the affected rows is %d\n", n)
}

func update() {
	sqlStr := "update user set age=? where id=?"

	ret, err := db.Exec(sqlStr, 39, 10)

	if err != nil {
		fmt.Printf("update failed,err:%v\n", err)
		return
	}

	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected affected faild,err:%v\n", err)
		return
	}
	fmt.Printf("update successfully,the rows is %d\n", n)
}

func prepareQuery() {
	sqlStr := "select id,name,age from user where id>?"

	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare faild,err:%v\n", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(0)
	if err != nil {
		fmt.Printf("querry faild,err:%v\n", err)
		return
	}
	defer rows.Close()

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

func prepareInsert() {
	sqlStr := "insert into user(name,age) values (?,?)"

	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		fmt.Printf("prepare faild,err:%v\n", err)
		return
	}

	defer stmt.Close()
	_, err = stmt.Query("小明", 6)
	if err != nil {
		fmt.Printf("insert faild,err:%v\n", err)
		return
	}

	_, err = stmt.Query("小黄", 6)
	if err != nil {
		fmt.Printf("insert faild,err:%v\n", err)
		return
	}

	fmt.Printf("Insert successfully\n")
}

// sql注入示例
func sqlInjectDemo(name string) {
	sqlStr := fmt.Sprintf("select id, name, age from user where name='%s'", name)
	fmt.Printf("SQL:%s\n", sqlStr)
	var u user
	err := db.QueryRow(sqlStr).Scan(&u.id, &u.name, &u.age)
	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	fmt.Printf("user:%#v\n", u)
}

func transaction() {
	tx, err := db.Begin() //开启事务
	if err != nil {
		fmt.Printf("begin transaction failed,err:%v\n", err)
		return
	}

	sqlStr1 := "update user set age=30 where id=?"
	ret1, err := tx.Exec(sqlStr1, 2)
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("exec sql1 failed,err:%v\n", err)
		return
	}
	affRow1, err := ret1.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("exec ret1.RowsAffected() failed,err:%v\n", err)
		return
	}

	sqlStr2 := "update user set age=30 where id=?"
	ret2, err := tx.Exec(sqlStr2, 23)
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("exec sql1 failed,err:%v\n", err)
		return
	}
	affRow2, err := ret2.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("exec ret2.RowsAffected() failed,err:%v\n", err)
		return
	}

	fmt.Println(affRow1, affRow2)

	if affRow1 == 1 && affRow2 == 1 {
		_ = tx.Commit()
	} else {
		_ = tx.Rollback()

	}

	fmt.Printf("exec trans successfully!")

}
