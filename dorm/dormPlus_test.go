package dorm

import (
	// "database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

//Check that Filter() locates all relevant database records
func TestFilter(t *testing.T) {
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
	db.Filter(&User{FullName: "Frelicia"}, "FullName", "Frelicia")
	//ColumnNames(&results)
}

// Check that Find() panics when no table exists for a query
func TestFilterPanic(t *testing.T) {}

func TestTopN(t *testing.T) {

}

func TestDelete(t *testing.T) {}

func TestGrant(t *testing.T) {}

func TestRevoke(t *testing.T) {}
