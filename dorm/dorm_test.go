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

type User2 struct {
	id       int
	FullName string `mytag:"special field"`
	eMail    string
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

func TestColumnNames(t *testing.T) {
	cols := ColumnNames(&User2{})
	if len(cols) != 1 {
		t.Errorf("Expected 1 col but found %d", len(cols))
		fmt.Println(cols)
	}
}

func TestFind(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	results := []User{}
	db.Find(&results)

	if len(results) != 1 {
		t.Errorf("Expected 1 users but found %d", len(results))
		fmt.Println(results)
	}
}

func TestFirst(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	result := &User{}
	db.First(result)

	mockUser := User{FullName: "Test User1"}

	if *result != mockUser {
		t.Errorf("incorrect")
		fmt.Println(*result)
	}
}

func TestCreate(t *testing.T) {
	//fmt.Println("in test create")
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	res, table_check := db.inner.Query("SELECT * FROM user")
	fmt.Println("SELECT * FROM user", res, table_check)
	/*for res.Next() {
		var full_name string
	r	res.Scan(&full_name)
		fmt.Println("names:    ", full_name)
	}*/
	cols, err := res.ColumnTypes()
	fmt.Println("tbales plz", &cols, cols[len(cols)-1].Name(), cols[0].Name(), err, table_check)
	//rows, table_check := db.inner.Query("select * from " + TableName(&User{FullName: "Frelicia"}))
	fmt.Println()
	fmt.Println("rows", res, table_check, TableName(&User{FullName: "Frelicia"}))
	res.Close()
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

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
