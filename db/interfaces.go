// db/interfaces.go
package db

// DBInterface abstracts sql.DB operations
type DBInterface interface {
	Query(query string, args ...interface{}) (RowsInterface, error)
	QueryRow(query string, args ...interface{}) RowInterface
}

// RowsInterface abstracts sql.Rows operations
type RowsInterface interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// RowInterface abstracts sql.Row operations
type RowInterface interface {
	Scan(dest ...interface{}) error
}
