package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
)

type Store interface {
	Connect(args ...interface{}) error
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Select(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func NewStore(config map[string]string) (Store, error) {
	switch config["type"] {
	case "filesystem":
		return &FileSystem{basePath: config["basePath"]}, nil
	case "mysql":
		//return &MySQL{dsn: config["dsn"]}, nil
		return &MySQL{}, nil
	case "genji":
		//return &Genji{dbPath: config["dbPath"]}, nil
		return &Genji{}, nil
	case "badger":
		//return &Badger{dir: config["dir"]}, nil
		return &Badger{}, nil
	case "s3":
		//return &S3{bucket: config["bucket"], region: config["region"]}, nil
		//return &S3{}, nil
		return nil, errors.New(fmt.Sprintf("S3 is not implemented yet"))
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", config["type"])
	}
}

type FileSystem struct {
	basePath string
}

func (fs *FileSystem) Connect(args ...interface{}) error {
	// code to write the value to a file on the filesystem
	return nil
}

func (fs *FileSystem) Put(key string, value []byte) error {
	// code to write the value to a file on the filesystem
	return nil
}

func (fs *FileSystem) Get(key string) ([]byte, error) {
	// code to read the value from a file on the filesystem
	return []byte{}, nil
}

func (fs *FileSystem) Delete(key string) error {
	// code to delete the file on the filesystem
	return nil
}

func (fs *FileSystem) Select(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New(fmt.Sprintf("Function Select(query string, args ...interface{}) (*sql.Rows, error) is not implemented for this Storage class"))
}

func (fs *FileSystem) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, errors.New(fmt.Sprintf("Function Exec(query string, args ...interface{}) (sql.Result, error) is not implemented for this Storage class"))
}

// MySQL struct
type MySQL struct {
	db        *sql.DB
	dbConnErr error
}

func (m *MySQL) Connect(args ...interface{}) error {
	// code to write the value to a file on the filesystem
	return nil
}

func (m *MySQL) Put(key string, value []byte) error {
	// code to insert the key-value pair into the MySQL database
	_, err := m.db.Exec("INSERT INTO storage (key, value) VALUES (?, ?)", key, value)
	return err
}

func (m *MySQL) Get(key string) ([]byte, error) {
	// code to select the value from the MySQL database by key
	var value []byte
	err := m.db.QueryRow("SELECT value FROM storage WHERE key = ?", key).Scan(&value)
	return value, err
}

func (m *MySQL) Delete(key string) error {
	// code to delete the key-value pair from the MySQL database
	_, err := m.db.Exec("DELETE FROM storage WHERE key = ?", key)
	return err
}

func (m *MySQL) Select(query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.Query(query, args...)
}

func (m *MySQL) Exec(query string, args ...interface{}) (sql.Result, error) {
	return m.db.Exec(query, args)
}

// Genji struct
type Genji struct {
	db        *genji.DB
	dbConnErr error
}

// Connect - takes two arguments dbPath(string) and ctx(context.Context) in that order
func (g *Genji) Connect(args ...interface{}) error {
	// code to open a genjji DB file on filesystem
	dbPath := args[0].(string)
	ctx := args[0].(context.Context)
	if dbPath != "" {
		g.db, g.dbConnErr = genji.Open(dbPath)
		if g.dbConnErr != nil {
			return g.dbConnErr
		}
		g.db.WithContext(ctx)
	}

	return errors.New(fmt.Sprintln("Invalid DB Path ", args, "Valid DB path should be a string"))
}

func (g *Genji) Put(key string, value []byte) error {
	// code to insert the key-value pair into the Genji database
	return g.db.Exec("INSERT INTO storage (key, value) VALUES (?, ?)", key, value)
}

func (g *Genji) Get(key string) ([]byte, error) {
	// code to select the value from the Genji database by key
	var value []byte
	result, resultErr := g.db.Query("SELECT value FROM storage WHERE key = ?", key)
	if resultErr != nil {
		return []byte{}, nil
	}

	err := result.Iterate(func(d types.Document) error {
		err := document.Scan(d, &value)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return []byte{}, err
	}
	return value, nil
}

func (g *Genji) Delete(key string) error {
	// code to delete the key-value pair from the Genji database
	return g.db.Exec("DELETE FROM storage WHERE key = ?", key)
}

func (g *Genji) Select(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New(fmt.Sprintf("Function Select(query string, args ...interface{}) (*sql.Rows, error) is not implemented for this Storage class"))
}

func (g *Genji) Exec(query string, args ...interface{}) (sql.Result, error) {
	if g.dbConnErr != nil {
		return nil, g.dbConnErr
	}
	return nil, g.db.Exec(query, args)
}

// Badger struct
type Badger struct {
	db        *badger.DB
	dbConnErr error
}

func (b *Badger) Connect(args ...interface{}) error {
	// Open the Badger database located in the dbPath directory.
	// It will be created if it doesn't exist.
	dbPath := args[0].(string)
	b.db, b.dbConnErr = badger.Open(badger.DefaultOptions(dbPath))
	if b.dbConnErr != nil {
		return b.dbConnErr
	}
	return nil
}

func (b *Badger) Put(key string, value []byte) error {
	// code to insert the key-value pair into the Badger database
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

func (b *Badger) Get(key string) ([]byte, error) {
	// code to select the value from the Badger database by key
	var value []byte
	return value, nil
}

func (b *Badger) Delete(key string) error {
	// code to delete the file on the filesystem
	return nil
}

func (b *Badger) Select(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, errors.New(fmt.Sprintf("Function Select(query string, args ...interface{}) (*sql.Rows, error) is not implemented for this Storage class"))
}

func (b *Badger) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, errors.New(fmt.Sprintf("Function Exec(query string, args ...interface{}) (sql.Result, error) is not implemented for this Storage class"))
}

//type S3 struct {
//}
//
//func (s *S3) Connect(args ...interface{}) error {
//	// code to write the value to a file on the filesystem
//	return nil
//}
//
//func (s *S3) Put(key string, value []byte) error {
//	// code to insert the key-value pair into the Badger database
//	return s.Put(func(txn *badger.Txn) error {
//		return txn.Set([]byte(key), value)
//	})
//}
//
//func (s *S3) Get(key string) ([]byte, error) {
//	// code to select the value from the Badger database by key
//	var value []byte
//	return value, nil
//}
//
//func (s *S3) Delete(key string) error {
//	// code to delete the file on the filesystem
//	return nil
//}
//
//func (s *S3) Select(query string, args ...interface{}) (*sql.Rows, error) {
//	return nil, errors.New(fmt.Sprintf("Function Select(query string, args ...interface{}) (*sql.Rows, error) is not implemented for this Storage class"))
//}
