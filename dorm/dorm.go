package dorm

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
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

// source: https://stackoverflow.com/questions/56616196/how-to-convert-camel-case-string-to-snake-case
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

func ColumnVal(v interface{}) []interface{} {
	t := reflect.TypeOf(v).Elem()
	cols := []interface{}{}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() {
			cols = append(cols, reflect.ValueOf(v).Elem().Field(i).Interface())
		}
	}
	return cols
}

func ColumnTypes(v interface{}) []string {
	t := reflect.TypeOf(v).Elem()
	cols := []string{}
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() {
			cols = append(cols, fmt.Sprintf("%v", t.Field(i).Type))
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
	return ToSnakeCase(val)
}

// Write rows resulting from SQL database query to result interface, which has
// the type that interface r has.
func writeRows(r interface{}, rows *sql.Rows, result interface{}) {
	cols := ColumnNames(r)

	v := reflect.Indirect(reflect.New(reflect.ValueOf(r).Type().Elem()))
	fields := make([]interface{}, len(cols))

	for i := 0; i < len(cols); i++ {
		t := v.Field(i).Type()
		field := reflect.New(t).Interface()
		fields[i] = field
	}

	for rows.Next() {
		rows.Scan(fields...)
		resultRow := reflect.New(reflect.TypeOf(r).Elem())
		val := reflect.Indirect(resultRow)
		for i := 0; i < len(fields); i++ {
			if val.Field(i).CanSet() {
				val.Field(i).Set(reflect.ValueOf(fields[i]).Elem())
			}
		}
		res := reflect.ValueOf(result).Elem()
		res.Set(reflect.Append(reflect.Indirect(reflect.ValueOf(result)), val))
	}
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
	r := reflect.New(reflect.ValueOf(result).Type().Elem().Elem()).Interface()
	tableName := TableName(r)

	query := "SELECT * FROM " + tableName
	rows, err := db.inner.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	writeRows(r, rows, result)

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
	cols := ColumnNames(result)

	query := "SELECT * FROM " + tableName + " LIMIT 1"
	rows, err := db.inner.Query(query)

	v := reflect.ValueOf(result).Elem()
	fields := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		t := v.Field(i).Type()
		field := reflect.New(t).Interface()
		fields[i] = field
	}

	if err != nil {
		log.Panic(err)
		return false
	} else if rows.Next() {
		rows.Scan(fields...)
		for i := 0; i < len(fields); i++ {
			v.Field(i).Set(reflect.ValueOf(fields[i]).Elem())
		}
		return true
	}
	return false
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
		numExportedField := 0
		for i := 0; i < len(cols); i++ {

			for !reflect.TypeOf(model).Elem().Field(numExportedField).IsExported() {
				numExportedField++

			}
			fmt.Println("&&&&&&&&&&&&&&&&&&", reflect.TypeOf(model).Elem().Field(numExportedField))
			if val, ok := reflect.TypeOf(model).Elem().Field(numExportedField).Tag.Lookup("dorm"); ok {
				if val == "primary_key" {
					continue
				}
			}
			fmt.Println("field", i, cols[i].Name())
			colNames = append(colNames, cols[i].Name())
			placeholder = append(placeholder, "?")

			fmt.Println("*******************************")
			numExportedField++
			//tag := t.Tag
			//fmt.Println("tag", tag)
		}
		fmt.Println("Placeholders", placeholder)
		fmt.Printf("%v  value again \n", reflect.ValueOf(model).Elem().Kind())
		colVals := []interface{}{}

		var colValsStr []string
		fmt.Println("len of colnames", len(colNames), len(colValsStr))
		//  colTypes,err := rows.ColumnTypes()
		fmt.Println("model", reflect.ValueOf(model).Elem(), reflect.TypeOf(model))
		ele := reflect.TypeOf(model).Elem()
		fmt.Println("posted", ele, ele.NumField(), ele.Field(0))

		numExportedField = 0
		for i := 0; i < ele.NumField() && numExportedField < ele.NumField(); i++ {
			fmt.Println("population col vals", i)
			for !reflect.TypeOf(model).Elem().Field(numExportedField).IsExported() {
				numExportedField++
			}
			if val, ok := ele.Field(numExportedField).Tag.Lookup("dorm"); ok {
				if val == "primary_key" {
					skipped = true
					continue
				}
			}
			fmt.Println("after checks", i)
			colVals = append(colVals, reflect.ValueOf(model).Elem().Field(numExportedField).Interface())
			fmt.Println("after checks", i, fmt.Sprintf("%v", colVals[len(colVals)-1]))
			colValsStr = append(colValsStr, fmt.Sprintf("%v", colVals[len(colVals)-1]))
			//fmt.Println("colnams/val", colNames[i], colVals[len(colVals)-1])
			numExportedField++
		}

		fmt.Println(colVals)
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

func (db *DB) Filter(result interface{}, filter interface{}) {
	r := reflect.New(reflect.ValueOf(result).Type().Elem().Elem()).Interface()
	tableName := TableName(r)
	cols := ColumnNames(filter)
	colTypes := ColumnTypes(filter)
	colVals := ColumnVal(filter)

	totalString := ""

	for i := 0; i < len(cols); i++ {
		if colTypes[i] == "string" {
			totalString += "(" + cols[i] + "=" + fmt.Sprintf("'%v'", colVals[i]) + ")"
		} else {
			totalString += "(" + cols[i] + "=" + fmt.Sprintf("%v", colVals[i]) + ")"
		}
		if i != len(cols)-1 {
			totalString += "\n OR \n"
		}
	}

	query := fmt.Sprintf("SELECT * FROM %v WHERE %v", tableName, totalString)
	rows, err := db.inner.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	writeRows(r, rows, result)

}

// Query the database for the first n rows in a given table
func (db *DB) TopN(result interface{}, n int) {
	r := reflect.New(reflect.ValueOf(result).Type().Elem().Elem()).Interface()
	tableName := TableName(r)

	query := fmt.Sprintf("SELECT * FROM %s LIMIT %d", tableName, n)
	rows, err := db.inner.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	writeRows(r, rows, result)
}

// Query and return database results for a user specified SQL query
func (db *DB) Query(result interface{}, query string) {
	r := reflect.New(reflect.ValueOf(result).Type().Elem().Elem()).Interface()

	rows, err := db.inner.Query(query)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	writeRows(r, rows, result)
}

// Remove row from the database
func (db *DB) Delete(model interface{}) {
	fmt.Println("in delete")

	rows, table_check := db.inner.Query("select * from " + TableName(model) + ";")
	//fmt.Println("rows",rows, table_check)
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
		for rows.Next() {
			var full_name string
			rows.Scan(&full_name)
			fmt.Println("names:    ", full_name)
		}

		//fmt.Println(r.Interface())
		//.Elem()).Interface()

		cols := ColumnNames(model)
		fmt.Println("col names of result:", cols, name)
		colTypes := ColumnTypes(model)
		colVals := []interface{}{}

		//var colValsStr []string
		totalString := ""
		withVal := ""
		word := ""
		//ele := reflect.TypeOf(model).Elem()
		numExportedField := 0
		for i := 0; i < len(cols)-1; i++ {
			fmt.Println("population col vals", i)
			fmt.Println("after checks", i)

			fmt.Println("after checks", i)
			//colValsStr = append(colValsStr, fmt.Sprintf("%v", colVals[len(colVals)-1]))
			//fmt.Println("colnams/val", colNames[i], colVals[len(colVals)-1])
			for !reflect.TypeOf(model).Elem().Field(numExportedField).IsExported() {
				numExportedField++
			}
			fmt.Println(colTypes[i])
			totalString += cols[i] + "=?" + " OR "
			//if colTypes[i] == "string" {
			word = fmt.Sprintf("%v", reflect.ValueOf(model).Elem().Field(numExportedField).Interface())

			/*} else {
				word = fmt.Sprintf("%v", reflect.ValueOf(model).Elem().Field(i).Interface())

			}*/
			colVals = append(colVals, word)
			withVal += cols[i] + "=" + word + " OR "
		}
		fmt.Println("********", cols)
		fmt.Println(colTypes[len(cols)-1])
		totalString += cols[len(cols)-1] + "=?"
		//if colTypes[len(cols)-1] == "string" {
		word = fmt.Sprintf("%v", reflect.ValueOf(model).Elem().Field(numExportedField+1).Interface())
		/*} else {
			word = fmt.Sprintf("%v", reflect.ValueOf(model).Elem().Field(len(cols)-1).Interface())
		}*/
		colVals = append(colVals, word)
		withVal += cols[len(cols)-1] + "=" + word
		//build or string

		query := fmt.Sprintf("DELETE FROM %v WHERE %v", name, totalString)
		fullString := fmt.Sprintf("DELETE FROM %v WHERE %v", name, withVal)
		fmt.Println(query)

		//fmt.Println(query)
		rows.Close()
		stmt, err := db.inner.Prepare(query)
		fmt.Println("stmt", stmt, err, colVals)
		fmt.Println(colVals)
		fmt.Println(colVals...)
		res, errExec := stmt.Exec(colVals...) // also returns uniquw id
		/*for res.Next() {
			var full_name string
			res.Scan(&full_name)
			fmt.Println("names:    ", full_name)
		}*/
		num, err := res.RowsAffected()
		if errExec != nil {
			fmt.Println("0", err, num)
		} else {
			fmt.Println(num, err)
		}
		//query := "SELECT * FROM " + tableName
		fmt.Println("placeholder string I have: ", query)
		fmt.Println("full string it will be: ", fullString)
		//rows, err := db.inner.Query(query)

	} else {
		fmt.Println("table not there")
		log.Panic("table not there", table_check)
	}

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

func sqlType(t string) string {
	if t == "string" {
		return "text"
	} else if strings.Contains(t, "float") || strings.Contains(t, "int") || strings.Contains(t, "bool") || strings.Contains(t, "char") {
		return t
	} else {
		return "text"
	}

}

// REVOKE:Takes back privileges granted from user.
func (db *DB) CreateTable(model interface{}) {
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

	if table_check != nil {

		fmt.Println("table should be created")
		//t := reflect.TypeOf(model).Elem()
		//u := reflect.New(t).Elem().Interface()
		//fmt.Printf("u is %T, %#v\n", u, u)

		/*resultValuePtr := reflect.New(t)
		resultValue := resultValuePtr.Elem()
		fieldCount := t.NumField()
		fields := make([]reflect.StructField, fieldCount)
		for i := 0; i < fieldCount; i++ {
			fields[i] = t.Field(i)
			fmt.Println(fields[i])
		}
		fmt.Println(resultValue)*/

		columns := ColumnNames(model) //make([]string, fieldCount)
		fieldCount := len(columns)
		//fieldAddrs := make([]interface{}, fieldCount)
		types := ColumnTypes(model)

		together := make([]string, fieldCount)
		query := "create table " + name + " (\n"
		for i := 0; i < fieldCount; i++ {
			//columns[i] = ToSnakeCase(fields[i].Name)
			//fieldAddrs[i] = resultValue.Field(i).Addr().Interface()
			//types[i] = sqlType(columns[i].Type.Name())
			if i != fieldCount-1 {
				together[i] = fmt.Sprintf("%v %v,", columns[i], types[i])
			} else {
				together[i] = fmt.Sprintf("%v %v", columns[i], types[i])
			}
			query += together[i] + "\n"
		}
		query += ")"
		fmt.Println("new:", columns, types)
		fmt.Println(together)
		fmt.Println(query)

		_, err := db.inner.Exec(query)

		if err != nil {
			panic(err)
		}
		/*_, err := db.inner.Exec("PRAGMA table_info(table_name);")
		if err != nil {
			panic(err)
		}*/
		//fmt.Println("ans", reflect.TypeOf(ans).Elem())

		db.Create(model)

	} else {
		fmt.Println("table already there")
	}

}
