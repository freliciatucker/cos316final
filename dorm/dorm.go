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
	if reflect.ValueOf(v).Kind() != reflect.Struct {
		log.Panic("requires struct")
	}
	cols := []string{}
	// placeholder:= []string{}
	val := reflect.ValueOf(v)
	for i := 0; i < val.NumField(); i++ {
		colname := val.Type().Field(i).Name
		colname_fixed := strings.ToLower(string(colname[0])) + colname[1:]
		cols = append(cols, colname_fixed)
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

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func TableName(result interface{}) string {
	name := reflect.TypeOf(result).Name()
	if reflect.ValueOf(result).Kind() != reflect.Struct {
		log.Panic("requires struct")
	}
	str := ToSnakeCase(name)

	return str
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
	// rows,err := db.inner.Query(query)

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
