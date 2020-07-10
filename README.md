# qb


## *qb is still in active development. The API is unstable. Use with caution.*

`qb` is a simple SQL query builder. Every method returns a struct of the same type allowing an entire query to be built in one function call chain. This library's only job is to generate properly formatted SQL queries.

Every query builder struct implements the Builder interface defined as

```go
type Builder interface {
    Build() (string, []interface{}, error)
}
```

After setting all of the query's values, call `Build()` to retrieve the query and a slice of correctly ordered parameters.

## Install

```
go get -u github.com/mattmeyers/qb
```

## Generating Queries

Queries are built using unexported structs.  Therefore, methods should be chained following the initialization function for the given type of query.

### Select

A select query can be initialized with the `Select(cols ...string)` function.  The struct returned from this function call can then call the following functions:

- `Select(cols ...string) *selectQuery`
- `Distinct(cols ...string) *selectQuery`
- `SetCols(cols ...string) *selectQuery`
- `From(table interface{}) *selectQuery`
- `InnerJoin(table string, condition interface{}) *selectQuery`
- `LeftJoin(table string, condition interface{}) *selectQuery`
- `RightJoin(table string, condition interface{}) *selectQuery`
- `FullJoin(table string, condition interface{}) *selectQuery`
- `CrossJoin(table string, condition interface{}) *selectQuery`
- `Where(pred Builder) *selectQuery`
- `Limit(l int) *selectQuery`
- `ClearLimit() *selectQuery`
- `Offset(o int) *selectQuery`
- `ClearOffset() *selectQuery`
- `GroupBy(cols ...string) *selectQuery`
- `Having(pred Builder) *selectQuery`
- `OrderBy(col string, dir OrderDir) *selectQuery`
- `RebindWith(r Rebinder) *selectQuery`
- `String() string`
- `Build() (string, []interface{}, error)`


For example, in order to generate the query

```sql
SELECT id, display_name FROM products WHERE (item_number=? OR item_number=?) AND backordered=?
```

use the following code:

```go
qb.Select("id", "display_name").
   From("products").
   Where("item_number", "=", "a123").
   Where(qb.Or{
      qb.Eq("item_number", "a123"),
      qb.Eq("item_number", "b456"),
   }).
   Where(qb.Eq("backordered", false)),
   Build()
```

### Insert

An insert query can be initialized with the `InsertInto(table string)` function.  The struct returned from this function call can then call the following functions:

- `Col(col string, val interface{}) *insertQuery`
- `Cols(cols []string, vals ...interface{}) *insertQuery`
- `OnConflict(target, action interface{}) *insertQuery`
- `Returning(cols ...string) *insertQuery`
- `RebindWith(r Rebinder) *insertQuery`
- `String() string`
- `Build() (string, []interface{}, error)`

Calling `Columns` or `Values` mulitple times will append the passed values to the columns and values arrays.  This can be handy when inserting optional columns. For example, in order to generate the query

```sql
INSERT INTO products (name, qty) VALUES (?, ?)
```

use the following code:

```go
qb.InsertInto("products").
   Columns("name", "qty").
   Values("Hammer", 5).
   String()
```

If using PostgreSQL, the `OnConflict` function can be used to generate an `ON CONFLICT target action` clause.  The provided target should be of type `TargetColumn`, `TargetConstraint`, or `whereClause`.  The provided action should be of type `ActionDoNothing` or `*updateQuery`.  For example, to generate the query

```sql
INSERT INTO products (name, item_number) VALUES (?, ?) ON CONFLICT (item_number) DO UPDATE SET item_number=123
```

use the following code:

```go
qb.InsertInto("products").
   Col("name", "Hammer").
   Col("item_number", 456).
   OnConflict(
     qb.TargetColumn("item_number"),
     qb.Update("").Set("item_number", 123),
   ).
   String()
```

### Update

An update query can be initialized with the `Update(table string)` function.  The struct returned from this function call can then call the following functions:

- `Set(col string, val interface{})`
- `Where(col, cmp string, val interface{})`
- `OrWhere(col, cmp string, val interface{})`

Calling `Set` with the same col value will update the previous value.  For example, in order to generate the query

```sql
UPDATE products SET name=?, qty=? WHERE item_id=?
```

use the following code:

```go
qb.Update("products").
   Set("name", "Screwdriver").
   Set("qty", 10).
   Where("item_id", "=", "a123").
   String()
```

### Delete

A delete query can be initialized with the `DeleteFrom(table string)` function.  The struct returned from this function call can then call the following functions:

- `Where(col, cmp string, val interface{})`
- `OrWhere(col, cmp string, val interface{})`

For example, in order to generate the query

```sql
DELETE FROM products WHERE item_number=? AND qty<? OR backordered=?
```

use the following code:

```go
qb.DeleteFrom("products").
   Where("item_number", "=", "a123").
   Where("qty", "<", 5).
   OrWhere("backordered", "=", true).
   String()
```

## Error Handling

Calling the `String()` function returns an `error` as its third return value.  This error will describe any missing values.  The following error constants are defined in the package and can be used with `errors.Is()` if using Go 1.13+.

```go
ErrMissingTable    = Error("no table specified")
ErrMissingSetPairs = Error("no set pairs provided")
ErrColValMismatch  = Error("the number of columns and values do not match")
ErrInvalidConflictTarget = Error("invalid conflict target")
ErrInvalidConflictAction = Error("invalid conflict action")
```

## Acknowledgments

This library was heavily inspired by the following libraries:
* [squirrel](https://github.com/masterminds/squirrel)
* [dbr](https://github.com/gocraft/dbr)
