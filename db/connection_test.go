package db

import (
	"errors"
)

// Mock implementations for testing
type MockDB struct {
	queryFunc    func(query string, args ...interface{}) (RowsInterface, error)
	queryRowFunc func(query string, args ...interface{}) RowInterface
	close        func() error
}

func (m *MockDB) Query(query string, args ...interface{}) (RowsInterface, error) {
	if m.queryFunc != nil {
		return m.queryFunc(query, args...)
	}
	return nil, errors.New("not implemented")
}

func (m *MockDB) QueryRow(query string, args ...interface{}) RowInterface {
	if m.queryRowFunc != nil {
		return m.queryRowFunc(query, args...)
	}
	return &MockRow{scanFunc: func(dest ...interface{}) error {
		return errors.New("not implemented")
	}}
}

func (m *MockDB) Close() error {
	return m.close()
}

type MockRows struct {
	nextFunc  func() bool
	scanFunc  func(dest ...interface{}) error
	closeFunc func() error
	errFunc   func() error
	nextCalls int
}

func (m *MockRows) Next() bool {
	m.nextCalls++
	if m.nextFunc != nil {
		return m.nextFunc()
	}
	return false
}

func (m *MockRows) Scan(dest ...interface{}) error {
	if m.scanFunc != nil {
		return m.scanFunc(dest...)
	}
	return nil
}

func (m *MockRows) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func (m *MockRows) Err() error {
	if m.errFunc != nil {
		return m.errFunc()
	}
	return nil
}

type MockRow struct {
	scanFunc func(dest ...interface{}) error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	if m.scanFunc != nil {
		return m.scanFunc(dest...)
	}
	return nil
}
