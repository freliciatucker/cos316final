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

func connectSQLAuth() *sql.DB {
	conn, err := sql.Open("sqlite3", "file:test.db?mode=memory&_auth&_auth_user=admin&_auth_pass=admin")
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

func createUser2Table(conn *sql.DB) {
	_, err := conn.Exec(`create table user2 (
		full_name text,
		e_mail text
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

func insertUsers2(conn *sql.DB, users []User2) {
	for _, uc := range users {
		_, err := conn.Exec(`insert into user2
		values
		(?,?)`, uc.FullName, uc.EMail)

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
	EMail    string
}

var MockUsers = []User{
	User{FullName: "Alice Apple"},
	User{FullName: "Bob Smith"},
	User{FullName: "Kyra Acquah"},
	User{FullName: "Frelicia Tucker"},
	User{FullName: "Carol Crisp"},
	User{FullName: "Devon Donald"},
}

var MockUsers2 = []User{
	User{FullName: "Test User1"},
	User{FullName: "Frelicia"},
}

var MockUsers3 = []User2{
	User2{FullName: "Test User1"},
	User2{FullName: "Frelicia", EMail: "f@t"},
}

func TestColumnNames(t *testing.T) {
	cols := ColumnNames(&User2{})
	if len(cols) != 2 {
		t.Errorf("Expected 2 col but found %d", len(cols))
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

	if len(results) != len(MockUsers) {
		t.Errorf("Expected %d users but found %d", len(MockUsers), len(results))
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

	mockUser := User{FullName: "Alice Apple"}

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
