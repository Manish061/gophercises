package db

import (
	"database/sql"
	"errors"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3" //blank import
	"github.com/mitchellh/go-homedir"
)

var db *sql.DB

// DB initalises a sqlite3 DB and returns it
func DB() (*sql.DB, error) {
	if db == nil {
		// sql.Register("sqlite3", &sqlite3.SQLiteDriver{})
		home, errFindingHomeDir := homedir.Dir()
		if errFindingHomeDir != nil {
			return db, errFindingHomeDir
		}
		dbPath := filepath.Join(home, "phone.db")
		dbTemp, errOpeningDB := sql.Open("sqlite3", dbPath)
		if errOpeningDB != nil {
			return db, errOpeningDB
		}
		db = dbTemp
	}
	return db, nil
}

// Report returns the present state of the table 'tableName'
func Report(db *sql.DB, tableName string) (*sql.Rows, error) {
	if db == nil {
		return nil, errors.New("no database exists")
	}
	statement := "SELECT id, phone from " + tableName
	rows, errExecuting := db.Query(statement)
	if errExecuting != nil {
		return nil, errExecuting
	}
	return rows, nil
}

//Reset resets the given database 'db'
func Reset(db *sql.DB, tableName string) (*sql.DB, error) {
	if db == nil {
		return db, nil
	}
	statement := "DROP TABLE IF EXISTS " + tableName
	_, errDroping := db.Exec(statement)
	if errDroping != nil {
		return db, errDroping
	}
	return db, nil
}

// CreateTable creates the table 'tableName' in the database 'db'
func CreateTable(db *sql.DB, tableName string) error {
	if db == nil {
		return errors.New("no database exists")
	}
	statement := "CREATE TABLE IF NOT EXISTS `" + tableName + "` (`id` INTEGER PRIMARY KEY, `phone` TEXT NOT NULL) WITHOUT ROWID;"
	_, errExecuting := db.Exec(statement)
	if errExecuting != nil {
		return errExecuting
	}
	return nil
}

//Insert inserts phoneNumber into db
func Insert(db *sql.DB, id int, phoneNumber, tableName string) (int, error) {
	if db == nil {
		return -1, errors.New("no database exists")
	}
	statement := "INSERT INTO " + tableName + "(id, phone)" + "VALUES(\"" + strconv.Itoa(id) + "\",\"" + phoneNumber + "\")"
	_, errExecuting := db.Exec(statement)
	if errExecuting != nil {
		return -1, errExecuting
	}
	return 0, nil
}
