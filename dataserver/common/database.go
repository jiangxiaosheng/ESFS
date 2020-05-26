package common

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

/**
@author js
获取数据库url
*/
func getURL() string {
	name := GetConfig("db.name")
	username := GetConfig("db.username")
	host := GetConfig("db.host")
	password := GetConfig("db.password")
	port := GetConfig("db.port")
	url := fmt.Sprintf("%s:%s@(%s:%s)/%s", username, password, host, port, name)

	return url
}

/**
@author js
获取数据库连接
*/
func GetDBConnection() (*sql.DB, error) {
	driver := GetConfig("db.driver")
	url := getURL()

	var err error
	db, err := sql.Open(driver, url)
	if err != nil {
		fmt.Println("数据库连接失败", err.Error())
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Second * 30)

	err = db.Ping()
	if err != nil {
		fmt.Println("连接数据库失败", err.Error())
		return nil, err
	}
	return db, nil
}

/**
@author js
执行exec事务
*/
func DoExecTx(sql string, db *sql.DB) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("事务开启失败 %v", err.Error())
		return nil, err
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Printf("事务准备失败 %v", err.Error())
		return nil, err
	}
	res, err := stmt.Exec()
	if err != nil {
		log.Printf("事务执行失败 %v", err.Error())
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("事务提交失败 %v", err.Error())
		return nil, err
	}
	return res, nil
}

/**
@author js
执行query,不需要加事务
*/
func DoQuery(sql string, db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("数据库查询失败 %s\n%v", sql, err.Error())
		return nil, err
	}
	return rows, nil
}
