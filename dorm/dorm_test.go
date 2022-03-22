package dorm

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func connectSQL() *sql.DB {
	conn, err := sql.Open("sqlite3", "file:test.db?mode=memory")
	if err != nil {
		panic(err)
	}
	return conn
}

func createUserTable(conn *sql.DB) {
	_, err := conn.Exec(`create table user (
		full_name text
	)`)

	if err != nil {
		panic(err)
	}
}

func insertUsers(conn *sql.DB, users []User) {
	for _, uc := range users {
		_, err := conn.Exec(`insert into user
		values
		(?)`, uc.FullName)

		if err != nil {
			panic(err)
		}
	}
}

type User struct {
	FullName string
}

var MockUsers = []User{
	User{FullName: "Test User1"},
}

func TestString(t *testing.T) {
	words := []string{"CamelCase", "EMail", "COSFiles", "camelCase", "OldCOSFiles", "COSFiles"}
	for _, val := range words {
		arr := camelToArray(val)
		fmt.Println(val, arr)
		fmt.Println(arrayToUnderscore(arr))
	}

}

func TestFind(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	results := []User{}
	// fmt.Println("here")
	// db.Find(&results)
	ColumnNames(&results)

	// if len(results) != 1 {
	// 	t.Errorf("Expected 1 users but found %d", len(results))
	// }
}
