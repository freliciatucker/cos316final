package dorm

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

// DB handle
type DB struct {
	inner *sql.DB
}

// NewDB returns a new DB using the provided `conn`,
// an sql database connection.
// This function is provided for you. You DO NOT need to modify it.
func NewDB(conn *sql.DB) DB {
	return DB{inner: conn}
}

// Close closes db's database connection.
// This function is provided for you. You DO NOT need to modify it.
func (db *DB) Close() error {
	return db.inner.Close()
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ColumnNames analyzes a struct, v, and returns a list of strings,
// one for each of the public fields of v.
// The i'th string returned should be equal to the name of the i'th
// public field of v, converted to underscore_case.
// Refer to the specification of underscore_case, below.

// Example usage:
// type MyStruct struct {
//    ID int64
//    UserName string
// }
// ColumnNames(&MyStruct{})    ==>   []string{"id", "user_name"}
func ColumnNames(v interface{}) []string {
	t := reflect.TypeOf(v).Elem()
	cols := []string{}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() {
			cols = append(cols, ToSnakeCase(t.Field(i).Name))
		}
	}
	return cols
}

// TableName analyzes a struct, v, and returns a single string, equal
// to the name of that struct's type, converted to underscore_case.
// Refer to the specification of underscore_case, below.

// Example usage:
// type MyStruct struct {
//    ...
// }
// TableName(&MyStruct{})    ==>  "my_struct"
func TableName(result interface{}) string {
	val := reflect.TypeOf(result).Elem().Name()
	fmt.Printf("%v", val)
	//fmt.Println(val.Name())
	fmt.Printf("%v  value", reflect.ValueOf(result))
	fmt.Printf("%v  type", reflect.TypeOf(result))
	/*if reflect.TypeOf(result).Elem().Name() != reflect.Struct {
		log.Panic("requires struct")
	}*/

	fmt.Println("passed!")
	//str := val.String()

	return ToSnakeCase(val)
}

// Find queries a database for all rows in a given table,
// and stores all matching rows in the slice provided as an argument.

// The argument `result` will be a pointer to an empty slice of models. // To be explicit, it will have type: *[]MyStruct,
// where MyStruct is any arbitrary struct subject to the restrictions
// discussed later in this document.
// You may assume the slice referenced by `result` is empty.

// Example usage to find all UserComment entries in the database:
//    type UserComment struct = { ... }
//    result := []UserComment{}
//    db.Find(&result)
func (db *DB) Find(result interface{}) {
	val := reflect.ValueOf(result).Kind()
	fmt.Printf("%v  kind", val)
	fmt.Printf("%v  value", reflect.ValueOf(result))
	fmt.Printf("%v  type", reflect.TypeOf(result))
	fmt.Println(val.String())

	// str := val.String()

	// defer rows.Close()

	// fields := make([]interface{}, len(cols))
	// for i := 0; i < v.NumField; i++ {
	// 	field := reflect.New(v.Field(i).Type()).Interface()
	// 	fields[i] = field
	// }
	// for rows.Next() {
	// 	err := rows.Scan(fields...)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	for i := 0; i < len(fields); i++ {
	// 		fmt.Println(reflect.ValueOf(fields[i]).Elem())
	// 	}
	// }
}

// First queries a database for the first row in a table,
// and stores the matching row in the struct provided as an argument.
// If no such entry exists, First returns false; else it returns true.

// The argument `result` will be a pointer to a model.
// To be explicit, it will have type: *MyStruct,
// where MyStruct is any arbitrary struct subject to the restrictions
// discussed later in this document.

// Example usage to find the first UserComment entry in the database:
//    type UserComment struct = { ... }
//    result := &UserComment{}
//    ok := db.First(result)
// with the argument), otherwise return true.
func (db *DB) First(result interface{}) bool {
	tableName := TableName(result)
	query := "SELECT * FROM " + tableName + "LIMIT 1"
	row := db.inner.QueryRow(query)
	switch err := row.Scan(result); err {
	case nil:
		return true
	default:
		return false
	}
}

// Create adds the specified model to the appropriate database table.
// The table for the model *must* already exist, and Create() should
// panic if it does not.

// Optionally, at most one of the fields of the provided `model`
// might be annotated with the tag `dorm:"primary_key"`. If such a
// field exists, Create() should ignore the provided value of that
// field, overwriting it with the auto-incrementing row ID.
// This ID is given by the value of last_inserted_rowid(),
// returned from the underlying sql database.
func (db *DB) Create(model interface{}) {
	rows, table_check := db.inner.Query("select * from " + TableName(model) + ";")
	fmt.Println("rows", rows, table_check)
	fmt.Println("-------------------------------")
	val := reflect.ValueOf(model).Elem().Kind()
	fmt.Printf("%v  kind \n", val)

	fmt.Printf("%v  value \n", reflect.ValueOf(model).Elem())
	fmt.Printf("%v  type\n", reflect.TypeOf(model).Elem())

	fmt.Println(val.String())
	fmt.Println("-------------------------------")
	// fmt.Println(db.inner)
	fmt.Println("calling tablenane")
	name := TableName(model)
	fmt.Println("_________________________")
	fmt.Println("select * from " + name + ";")

	if table_check == nil {
		skipped := false
		fmt.Println("table is there")
		fmt.Println(rows.ColumnTypes())
		cols, _ := rows.ColumnTypes()
		fmt.Println("cols", cols)
		rows.Close()
		colNames := []string{}
		placeholder := []string{}
		for i := 0; i < len(cols); i++ {
			fmt.Println("&&&&&&&&&&&&&&&&&&", reflect.TypeOf(model).Elem().Field(i))
			if val, ok := reflect.TypeOf(model).Elem().Field(i).Tag.Lookup("dorm"); ok {
				if val == "primary_key" {
					continue
				}
			}
			fmt.Println("field", i, cols[i].Name())
			colNames = append(colNames, cols[i].Name())
			placeholder = append(placeholder, "?")

			fmt.Println("*******************************")
			//tag := t.Tag
			//fmt.Println("tag", tag)
		}

		fmt.Printf("%v  value again \n", reflect.ValueOf(model).Elem().Kind())
		colVals := []interface{}{}

		var colValsStr []string
		fmt.Println("len of colnames", len(colNames), len(colValsStr))
		//  colTypes,err := rows.ColumnTypes()
		fmt.Println("model", reflect.ValueOf(model).Elem(), reflect.TypeOf(model))
		ele := reflect.TypeOf(model).Elem()
		fmt.Println("posted", ele, ele.NumField(), ele.Field(1))

		for i := 0; i < ele.NumField(); i++ {
			fmt.Println("population col vals", i)
			if val, ok := ele.Field(i).Tag.Lookup("dorm"); ok {
				if val == "primary_key" {
					skipped = true
					continue
				}
			}
			fmt.Println("after checks", i)
			colVals = append(colVals, reflect.ValueOf(model).Elem().Field(i).Interface())
			fmt.Println("after checks", i)
			colValsStr = append(colValsStr, fmt.Sprintf("%v", colVals[len(colVals)-1]))
			//fmt.Println("colnams/val", colNames[i], colVals[len(colVals)-1])
		}

		query := fmt.Sprintf("INSERT OR REPLACE INTO %v(%v) VALUES(%v)", name, strings.Join(colNames, ","), strings.Join(placeholder, ","))
		fmt.Println(query)
		rows.Close()
		stmt, err := db.inner.Prepare(query)
		fmt.Println("stmt", stmt, err)
		fmt.Println(colVals...)
		res, errExec := stmt.Exec(colVals...) // also returns uniquw id

		if errExec != nil {
			log.Panic(errExec)
		}
		if errExec == nil {
			fmt.Println("somehow nil??", "should go on then")
		}

		lastinsert, _ := res.LastInsertId()
		rowsAffected, _ := res.RowsAffected()
		fmt.Println("res", lastinsert, rowsAffected)
		if skipped {
			//reflect.TypeOf(model).Elem()
			reflect.ValueOf(model).Elem().Field(0).Set(reflect.ValueOf(lastinsert))
			fmt.Println("model", reflect.ValueOf(model).Elem(), reflect.TypeOf(model))

		}
	} else {
		fmt.Println("table not there")
		log.Panic("table not there", table_check)
	}

}

func camelToArray(word string) []string {
	allWords := []string{}

	currentWord := ""
	// endOfWord := false
	beginnningOfWord := true
	isUpper := false

	for i := 0; i < len(word); i++ {
		//currentWord += string(word[i])
		if beginnningOfWord {
			if i < len(word) {
				isUpper = unicode.IsUpper(rune(word[i+1]))
			}
			beginnningOfWord = false
			currentWord += string(word[i])
		} else {

			if isUpper && unicode.IsUpper(rune(word[i])) {
				if i < len(word)-2 && !unicode.IsUpper(rune(word[i+2])) {
					currentWord += string(word[i])
					allWords = append(allWords, currentWord)
					currentWord = ""
					beginnningOfWord = true
					continue

				} else if i == len(word)-1 || unicode.IsUpper(rune(word[i+1])) {
					currentWord += string(word[i])
					continue

				} else {
					allWords = append(allWords, currentWord)
					fmt.Println(currentWord)
					currentWord = ""
					beginnningOfWord = true
					//allWords = append(allWords, currentWord)
					continue
				}
			}

			if !unicode.IsUpper(rune(word[i])) {
				isUpper = false
				currentWord += string(word[i])
				continue
			}
			allWords = append(allWords, currentWord)
			//fmt.Println(currentWord)
			currentWord = ""
			beginnningOfWord = true
			i--
			//allWords = append(allWords, currentWord)
		}

	}
	allWords = append(allWords, currentWord)
	return allWords
}

func arrayToUnderscore(arr []string) string {
	fmt.Println("in to lower", len(arr), arr)
	if len(arr) < 1 {
		return ""
	}
	fmt.Println("in to lower2")
	str := ""
	for i := 0; i < len(arr)-1; i++ {
		fmt.Print("in to lower3")
		val := strings.ToLower(string(arr[i])) //+ string(arr[i][1:])
		str += val + "_"
		fmt.Println(str)
	}
	fmt.Println("in to lower4", str)
	val := strings.ToLower(string(arr[len(arr)-1])) // + string(arr[len(arr)-1][1:])
	str += val
	return str
}

// Specify value for a field and return rows that match the value
//be super general
//stores the matching row in the struct provided as an argument.
func (db *DB) Filter(result interface{}, field string, value string) {
	rows, table_check := db.inner.Query("select * from " + TableName(result) + " WHERE " + field + " = '" + value + "';")
	fmt.Println(rows, table_check)
}

// Query the database for the first n rows in a given table
func (db *DB) TopN(result interface{}, n int) {

}

// Query and return database results for a user specified SQL query
func (db *DB) Query(result interface{}, query string) {

}

// Remove row from the database
func (db *DB) Delete(result interface{}, field string, value string) {

}

// limit what data in the table a user has access to
//Might be useful to look at assignment 6
//DCL - Data Control Language
//GRANT:Gives a privilege to user.
func (db *DB) Grant(userid string, permission string) {

}

// REVOKE:Takes back privileges granted from user.
func (db *DB) Revoke(userid string, permission string) {

}
