# COS316Final:, an extention of Assignment 4

## Due: Deans Date

# Adding selection and projection operations to (DORM)

This assignment will extend the functionality of Dorm


## API

The `dormPlus` package exposes the same API as 'dorm' plus the following:

```go
// Specify value for a field and return rows that match the value 
//be super general
func (db *DB) Filter(result interface{}, field string, value string)

// Query the database for the first n rows in a given table
func (db *DB) TopN(result interface{}, n int)

// Query and return database results for a user specified SQL query
func (db *DB) Query(result interface{}, query string)

// Remove row from the database
func (db *DB) Delete(result interface{},field string, value string)

// limit what data in the table a user has access to 
//Might be useful to look at assignment 6
//DCL - Data Control Language
//GRANT:Gives a privilege to user.
func (db *DB) Grant(userid string, permission string)

// REVOKE:Takes back privileges granted from user.
func (db *DB) Revoke(userid string, permission string)

```



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
