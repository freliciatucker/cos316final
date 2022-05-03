package dorm

import (
	// "database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

//Check that Filter() locates all relevant database records
func TestFilter(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	results := []User{}
	db.Filter(&results, &User{FullName: "Frelicia Tucker"})
	expectedResults := []User{
		{FullName: "Frelicia Tucker"},
	}
	if !reflect.DeepEqual(results, expectedResults) {
		t.Errorf("incorrect")
		fmt.Println(results)
	}
}

func Test_Filter_Multiple(t *testing.T) {
	conn := connectSQL()
	createUser2Table(conn)
	insertUsers2(conn, MockUsers3)

	db := NewDB(conn)
	defer db.Close()

	results := []User2{}
	filter := &User2{FullName: "Frelicia", EMail: "f@t"}
	db.Filter(&results, filter)
}

// Check that Filter() panics when no table exists for a query
func Test_Filter_Panic(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers2)

	db := NewDB(conn)
	defer db.Close()
	type UserFake struct {
		FullName string
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	results := []UserFake{}
	db.Filter(&results, &UserFake{FullName: "Frelicia"})

}

func TestTopN(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	for n := 0; n <= len(MockUsers); n++ {
		results := []User{}
		db.TopN(&results, n)

		if len(results) != n {
			t.Errorf("Expected %d users but found %d", n, len(results))
			fmt.Println(results)
		}
	}

	results := []User{}
	db.TopN(&results, len(MockUsers)+1)

	if len(results) != len(MockUsers) {
		t.Errorf("Expected %d users but found %d", len(MockUsers), len(results))
		fmt.Println(results)
	}
}

func TestQuery(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers)

	db := NewDB(conn)
	defer db.Close()

	results := []User{}
	db.Query(&results, "select * from user where full_name < 'D'")

	if len(results) != 3 {
		t.Errorf("Expected 3 users but found %d", len(results))
		fmt.Println(results)
	}

}

func TestDelete(t *testing.T) {
	//fmt.Println("in test create")
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers2)

	db := NewDB(conn)
	defer db.Close()

	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	res, table_check := db.inner.Query("SELECT * FROM user")
	fmt.Println("SELECT * FROM user", res, table_check)

	cols, err := res.ColumnTypes()
	fmt.Println("tbales plz", &cols, cols[len(cols)-1].Name(), cols[0].Name(), err, table_check)
	//rows, table_check := db.inner.Query("select * from " + TableName(&User{FullName: "Frelicia"}))
	fmt.Println()
	fmt.Println("rows", res, table_check, TableName(&User{FullName: "Frelicia"}))
	fmt.Println("starting Database: ")
	for res.Next() {
		var full_name string
		res.Scan(&full_name)
		fmt.Println("names:    ", full_name)
	}
	res.Close()
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	fmt.Println("after Deleteing 'Test User1' from database")

	//results := MockUsers
	// fmt.Println("here")
	db.Delete(&User{FullName: "Test User1"})
	//ColumnNames(&results)

	//fmt.Println("after Delete")
	str := ""
	rows, _ := db.inner.Query("SELECT * FROM user")
	for rows.Next() {
		var full_name string
		rows.Scan(&full_name)
		str += fmt.Sprintf("names:%v\n", full_name)
		fmt.Println("names:    ", full_name)
	}
	rows.Close()

	ans2 := "names:Frelicia\n"

	if str != ans2 {
		fmt.Println(str)
		fmt.Println(ans2)
		log.Fatal("didnt delete Test User1", strings.Compare(str, ans2))

	} else {
		fmt.Println("successfully delete Test User1")
	}
	fmt.Println("________________________________________")
	db.Delete(&User{FullName: "Test User1"})
	//ColumnNames(&results)

	fmt.Println("after Delete2 ,everything should be the same")
	rows2, _ := db.inner.Query("SELECT * FROM user")
	str = ""
	for rows2.Next() {
		var full_name string
		rows2.Scan(&full_name)
		str += fmt.Sprintf("names:%v\n", full_name)
		fmt.Println("names:    ", full_name)
	}
	rows2.Close()
	if str != ans2 {
		log.Fatal("error trying to delete same entry 2x")
	} else {
		fmt.Println("successfully did nothing")
	}
	fmt.Println("________________________________________")
	db.Delete(&User{FullName: "Frelicia"})
	fmt.Println("after Delete3 ,should be no rows left")
	rows3, _ := db.inner.Query("SELECT * FROM user")
	str = ""
	for rows3.Next() {
		var full_name string
		rows3.Scan(&full_name)
		str += fmt.Sprintf("names:%v\n", full_name)
		fmt.Println("names:    ", full_name)
	}
	rows3.Close()

	if str != "" {
		log.Fatal("did not delete last entry in table")
	} else {
		fmt.Println("successfully deleted last entry in table")
	}

	fmt.Println("________________________________________")
	fmt.Println("DONE WITH JUST DELETE")
}

func Test_Delete_MultipleCases(t *testing.T) {
	conn := connectSQL()
	createUser2Table(conn)
	insertUsers2(conn, MockUsers3)

	db := NewDB(conn)
	defer db.Close()
	res := &User2{FullName: "Frelicia", EMail: "f@t"}
	db.Delete(res)
	fmt.Println("after Delete")
	rows, _ := db.inner.Query("SELECT * FROM user2")
	for rows.Next() {
		var full_name string
		rows.Scan(&full_name)
		fmt.Println("names:    ", full_name)
	}
	rows.Close()
}

func Test_Delete_Panic(t *testing.T) {
	conn := connectSQL()
	createUserTable(conn)
	insertUsers(conn, MockUsers2)

	db := NewDB(conn)
	defer db.Close()
	type UserFake struct {
		FullName string
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	db.Delete(&UserFake{FullName: "Frelicia"})

}
func TestGrant(t *testing.T) {
	conn := connectSQLAuth()
	createUserTable(conn)
	insertUsers(conn, MockUsers2)

	db := NewDB(conn)
	defer db.Close()

	fmt.Println("show database")
	rows, _ := db.inner.Query("SELECT * FROM user")
	for rows.Next() {
		var full_name string
		rows.Scan(&full_name)
		fmt.Println("names:    ", full_name)
	}
	rows.Close()
}

func TestRevoke(t *testing.T) {
	conn := connectSQLAuth()
	createUserTable(conn)
	insertUsers(conn, MockUsers2)

	db := NewDB(conn)
	defer db.Close()

	fmt.Println("show database")
	rows, _ := db.inner.Query("SELECT * FROM user")
	for rows.Next() {
		var full_name string
		rows.Scan(&full_name)
		fmt.Println("names:    ", full_name)
	}
	rows.Close()

}

func TestCreateTable(t *testing.T) {
	conn := connectSQLAuth()
	db := NewDB(conn)
	defer db.Close()

	type NewTable struct {
		id       int
		Name     string
		Quantity int
	}

	db.CreateTable(&NewTable{id: 5, Name: "Frelicia", Quantity: 2})
	fmt.Println("show database")
	rows, _ := db.inner.Query("SELECT * FROM new_table")
	for rows.Next() {
		var full_name string
		rows.Scan(&full_name)
		fmt.Println("names:    ", full_name)
	}
	rows.Close()

}

func TestCreateTable2(t *testing.T) {
	conn := connectSQLAuth()
	db := NewDB(conn)
	defer db.Close()

	type NewTable struct {
		id       int
		name     string
		Quantity int
	}

	db.CreateTable(&NewTable{id: 5, name: "Frelicia", Quantity: 2})
	fmt.Println("show database")
	rows, _ := db.inner.Query("SELECT * FROM new_table")
	for rows.Next() {
		var full_name string
		rows.Scan(&full_name)
		fmt.Println("names:    ", full_name)
	}
	rows.Close()

}
