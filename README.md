# COS316, Assignment 4: Dopey Object Relational Mapper (DORM)

## Due: March 23rd at 11pm

# Dopey Object Relational Mapper (DORM)

This assignment asks you to build a generic Object Relational Mapper (ORM). An 
ORM translates language-level objects---Go structs in our case---to and from a
"relational mapping"---a SQLite database in our case. This allows application
developers to use normal language constructs and idioms to manipulate data
stored in a database.

For example, an application that keeps track of posts by different users might
have a Go struct modeling each post:

```go
type Post struct {
    ID     int64
    Author string
    Posted time.Time
    Likes  int
    Body   string
}
```

If we want to store this in a SQL database, the schema might look like:

```sql
create table post (
    id integer primary key autoincrement,
    author text,
    posted timestamp,
    likes integer,
    body text
)
```

where `id` is a primary key that is auto-incremented by 1 for each new record 
(tuple) inserted into the table named `post`.

An object relational mapper (ORM) allows an application developer to "talk" to the 
database through the Go struct:

```go
// Create a new post
post1 := &Post{
    Author: "alevy",
    Posted: time.Now(),
    Likes: 0,
    Body: "Hello fellow kids! This post will surely be viral"
}

// Insert the record into the database
dorm.Create(post1)
...

// Get and display all posts
allPosts := []Post{}
dorm.Find(&allPosts)

for _, post := range allPosts {
    fmt.Printf("%s said: %s\n", post.Author, post.Body)
}
```

In this assignment you'll build a simple ORM for mapping Go structs to a
SQLite database that can add new rows from a struct, fetch all rows of a given 
type, and fetch the first row of a given type, if it exists.

You will also implement some helper functions to analyze provided structs and 
return a string or strings representing the named fields of that struct.

## API

The `dorm` package exposes the following API:

```go
type DB struct {
	inner *sql.DB
}

// NewDB returns a new DB using the provided `conn`,
// an sql database connection.
// This function is provided for you. You DO NOT need to modify it.
func NewDB(conn *sql.DB) DB

// Close closes db's database connection.
// This function is provided for you. You DO NOT need to modify it.
func (db *DB) Close() error

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
func ColumnNames(v interface{}) []string

// TableName analyzes a struct, v, and returns a single string, equal
// to the name of that struct's type, converted to underscore_case.
// Refer to the specification of underscore_case, below.

// Example usage:
// type MyStruct struct {
//    ...
// }
// TableName(&MyStruct{})    ==>  "my_struct"
func TableName(result interface{}) string

// The function Find queries a database for all rows in a given table,
// and stores all matching rows in the slice provided as an argument.

// The argument `result` will be a pointer to an empty slice of models.
// To be explicit, it will have type *[]MyStruct,
// where MyStruct is any arbitrary struct subject to the restrictions
// discussed later in this document.
// You may assume the slice referenced by `result` is empty.

// Example usage to find all UserComment entries in the database:
//    type UserComment struct = { ... }
//    result := []UserComment{}
//    db.Find(&result)
func (db *DB) Find(result interface{})

// The function First queries a database for the first row in a table
// and stores the matching row in the struct provided as an argument.
// If no such entry exists, First returns false; otherwise it returns true.

// The argument `result` will be a pointer to a model.
// To be explicit, it will have type *MyStruct,
// where MyStruct is any arbitrary struct subject to the restrictions
// discussed later in this document.

// Example usage to find the first UserComment entry in the database:
//    type UserComment struct = { ... }
//    result := &UserComment{}
//    ok := db.First(result)
func (db *DB) First(result interface{}) bool

// Create adds the specified model to the appropriate database table.
// The table for the model *must* already exist, and Create() should
// panic if it does not.

// Optionally, at most one of the fields of the provided `model`
// might be annotated with the tag `dorm:"primary_key"`. If such a
// field exists, Create() should ignore the provided value of that
// field, overwriting it with the auto-incrementing row ID.
// This ID is given by the value of last_inserted_rowid(),
// returned from the underlying sql database.
func (db *DB) Create(model interface{})
```

You will be required to implement the following functions:
`TableName()`, `ColumnNames()`,  `Find()`, `First()`, `Create()`.

The functions `NewDB` and `Close` are provided for you.
There is no need to modify these functions for your implementation,
although you are welcome to if it will help your implementation.

### CamelCase and Underscore Case Specification

You will need to devise a way to convert between `camelCase` identifiers (structs in Go)
and `underscore_case` identifiers (columns in SQL). This conversion should apply whenever 
your Go program interacts directly with the SQL database. You may find the 
[strings](https://golang.org/pkg/strings/) package useful.

Below is a formal specification of `CamelCase` and `underscore_case`:

#### CamelCase
In CamelCase, we define a "word" to be either:
1.  any sequence of uppercase letters that is *not* followed
    by a non-uppercase character.
2.  a capitalized word (one uppercase letter followed by any number
    of non-uppercase characters)
3.  an initial sequence of non-uppercase characters.

For our purposes, uppercase letters are as defined by
[unicode.IsUpper](https://golang.org/pkg/unicode/#IsUpper).

Consider the following examples, with CamelCase identifiers on the
left, and their component words on the right.

```
CamelCase     ==>    ["Camel", "Case"]
EMail         ==>    ["E", "Mail"]
COSFiles      ==>    ["COS", "Files"]
camelCase     ==>    ["camel", "Case"]
OldCOSFiles   ==>    ["Old", "COS", "Files"]
COSFilesX     ==>    ["COS", "Files", "X"]
```

Non-alphabetical characters like numbers will never cause a word split.
In other words, you may consider them 'non-uppercase characters' for
the purposes of the definitions above.

#### Underscore_Case

Underscore case consists of a sequence of lowercased words joined
together by underscores. All of the names of tables and columns in your
database will be in underscore_case, which means you will need to translate
Golang's CamelCase identifiers to underscore_case identifiers.

We recommend mapping a CamelCase identifier to a list of component words,
and then lowercasing and joining the words together to obtain
your underscore case identifier.

### Restrictions on Structs

For the purposes of this assignment, we will make several simplifying
assumptions about the sorts of structs that make valid DB models.

In particular:
* You may assume that the fields of structs will all be primitive types
  (e.g. `string`, `int`, `int64`, `bool`, ...). This means you need not
  worry about `map`, `slice`, or `struct` types being included as
  fields.
  Primitive types are handled natively by the `sql` library we are
  using, so you should *not* have to do any special work to support
  these types. In contrast, `map`, `slice`, or nested `struct` types
  add complexity to the ORM implementation, so you are not responsible
  for supporting them.
* You may assume there will be no nested structs. For example, your
  implementation will not be tested against a model like the following:
  ```golang
  type StudentData struct {
    Name string
    ID int64
    Enrolled bool
  }

  type Roster struct {
    FirstStudent StudentData
    NumStudents int
    OtherStudents []StudentData
  }
  ```
* You may assume that all field names will be in a valid `camelCase` or
  `CamelCase` format. You may assume the same of all named struct types.
  Consider the following examples:

  ```golang
  type CamelCase struct {...}   // OK - valid camel case
  type camelCase struct {...}   // OK - valid camel case
  type camel_case struct {...}  // NOT OK - no underscores in camel case
  type CAMEL_CASE struct {...}  // NOT OK - no underscores in camel case
  ```

### Additional specifications

* `First()` should return the first row in a table based on the table's
  natural order. That is, your `First()` should respect the order of
  rows as returned by the underlying SQL database.
* In the event that a table contains zero rows, `First()` must return
  false, and the value `First()` assigns to its argument `result`
  is unspecified. Your implementation may leave `result` unchanged, or
  assign it some other reasonable default value.
* `Create()` is only responsible for adding a row to an existing
  table, *but it is not responsible for creating new tables*.
  If an attempt is made to add a row to a table that does not exist,
  `Create()` should panic with an explanatory error message.
* If any of the fields of the `model` provided to `Create()` are
  tagged with `dorm:"primary_key"`, you should assume that the type
  of that field will be `int64`.

## SQL Resources

As part of this assignment, you will need to write code that composes SQL
queries. We recommend you consult the 
[SQL precept slides](https://cos316.princeton.edu/precepts/SQL.pdf) or 
[SQLite documentation](https://www.sqlite.org/index.html) for a refresher on 
how you might accomplish this.

You may also find the golang 
[database/sql](https://golang.org/pkg/database/sql/) library useful.

As always, you are free to use the internet and any other resources you would
like to learn more information about SQL.

## Unit testing

Recall Go uses the [testing package](https://golang.org/pkg/testing/) to create
unit tests for Go packages.

For this assignment, you are provided with dorm_test.go, which contains very 
basic unit tests. You are encouraged to extend this file to create your own 
unit tests.

## Sample Application

A sample application, based on the example above, is provided in example/main.go

## Submission & Grading

Your assignment will be automatically submitted every time you push your changes
to your GitHub repo. Within a couple minutes of your submission, the
autograder will make a comment on your commit listing the output of our testing
suite when run against your code. **Note that you will be graded only on your
changes to the `dorm` package**, and not on your changes to any other files,
though you may modify any files you wish.

You may submit and receive feedback in this way as many times as you like,
whenever you like, but a substantial lateness penalty will be applied to
submissions past the deadline.
