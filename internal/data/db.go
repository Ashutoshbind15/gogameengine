package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var (
	DbConn *sql.DB
)

func InitDB() {
	connstr := os.Getenv("DB_URI")
	db, err := sql.Open("postgres", connstr)

	if err != nil {
		fmt.Println("errr in db conn: ", err)
		return
	}

	
	if err := db.Ping(); err != nil {
		fmt.Println("CANNOT PING")
		return
	}
	
	fmt.Println("CONNECTED TO THE DB")
	DbConn = db
}

func InitTables() {

	dir, pwderr := os.Getwd()

	if pwderr != nil {
		fmt.Println("err showing the pwd")
	}

	fmt.Println(dir)

	queries, err := os.ReadFile("../../internal/data/tables.sql")
	if err != nil {
		panic(err)
	}

	querystr := string(queries)
	fmt.Println(querystr)

	res, dberr := DbConn.Exec(querystr)

	if dberr != nil {
		fmt.Println("DB ERROR: ", dberr)
	}

	fmt.Println(res)
}