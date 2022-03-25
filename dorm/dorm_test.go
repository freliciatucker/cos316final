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
	// res, err := conn.Query("SELECT name FROM sqlite_schema WHERE type ='table' AND name NOT LIKE 'sqlite_%';")
	// fmt.Println("in connect so shouldd work", res, err)
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

func TestCreate(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)

	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	res, _ := db.inner.Query("SELECT * FROM user;")
	cols, err := res.ColumnTypes()
	fmt.Println("tbales plz", &cols, cols[len(cols)-1].Name(), cols[0].Name(), err)
	defer db.Close()

	//results := MockUsers
	// fmt.Println("here")
	db.Create(&User{FullName: "Frelicia"})
	//ColumnNames(&results)
}

func TestKey(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)

	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	res, _ := db.inner.Query("SHOW TABLES")
	fmt.Println("tbales plz", res)
	defer db.Close()

	//results := MockUsers
	// fmt.Println("here")
	db.Create(&User{FullName: "Frelicia "})
	//ColumnNames(&results)
}
