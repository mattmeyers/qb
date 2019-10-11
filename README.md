# qb

`qb` is a query builder library heavily inspired by [squirrel](https://github.com/masterminds/squirrel).  Every method returns a struct of the same type allowing an entire query to be built in one function call chain.  This library's only job is to generate properly formatted SQL queries.  It will not use a database driver to use the proper placeholders, but rather every query will use `?`s.  As such, it is the job of the user to rebind the query for their uses.

Every query builder struct implements the QueryBuilder interface defined as

```go
type QueryBuilder interface {
    String() (string, []interface, error)
}
```

After setting all of the query's values, call `String()` to retrieve the query and a slice of correctly ordered parameters.

## Select

A select query can be initialized with the `Select(cols ...string)` function.  The struct returned from this function call can then call the following functions:

- `From(table string)`
- `Where(col, cmp string, val interface{})`
- `OrWhere(col, cmp string, val interface{})`

For example, in order to generate the query 

```sql
SELECT id FROM products WHERE item_number=? AND in_stock=? OR backordered=?
```

use the following code:

```go
qb.Select("id")
  .From("products")
  .Where("item_number", "=", "a123")
  .Where("in_stock", "=", true)
  .OrWhere("backordered", "=", false)
  .String()
```

## Insert

An insert query can be initialized with the `InsertInto(table string)` function.  The struct returned from this function call can then call the following functions:

- `Columns(cols ...string)`
- `Values(vals ...interface{})`

Calling `Columns` or `Values` mulitple times will append the passed values to the columns and values arrays.  This can be handy when inserting optional columns. For example, in order to generate the query 

```sql
INSERT INTO products (name, qty) VALUES (?, ?)
```

use the following code:

```go
qb.InsertInto("products")
  .Columns("name", "qty")
  .Values("Hammer", 5)
  .String()
```