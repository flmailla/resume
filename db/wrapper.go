package db

func (db *DBWrapper) Query(query string, args ...interface{}) (RowsInterface, error) {
	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *DBWrapper) QueryRow(query string, args ...interface{}) RowInterface {
	row := db.db.QueryRow(query, args...)
	return row
}

func (db *DBWrapper) Close() error {
	return db.db.Close()
}
