package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	phoneTask "gophercises/phone"
	phoneDB "gophercises/phone/db"
)

func marshalDB(db *sql.DB) {
	out, err := json.Marshal(db)
	if err != nil {
		fmt.Printf("error occured: %v\n", err)
		return
	}

	fmt.Printf("%v\n", string(out))
}

func insertData(db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}
	_, errInserting := phoneDB.Insert(db, 1, "123-456-7892", "contacts")
	if errInserting != nil {
		// fmt.Printf("%v\n", errInserting)
		return errInserting
	}
	if err := tx.Commit(); err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}

	tx, err = db.Begin()
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}
	_, errInserting = phoneDB.Insert(db, 2, "1234567890", "contacts")
	if errInserting != nil {
		// fmt.Printf("%v\n", errInserting)
		return errInserting
	}
	if err = tx.Commit(); err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}

	tx, err = db.Begin()
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}
	_, errInserting = phoneDB.Insert(db, 3, "(123) 456 7892", "contacts")
	if errInserting != nil {
		// fmt.Printf("%v\n", errInserting)
		return errInserting
	}
	if err := tx.Commit(); err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}

	tx, err = db.Begin()
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}
	_, errInserting = phoneDB.Insert(db, 4, "(123) 456-7893", "contacts")
	if errInserting != nil {
		// fmt.Printf("%v\n", errInserting)
		return errInserting
	}
	if err = tx.Commit(); err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}

	tx, err = db.Begin()
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}
	_, errInserting = phoneDB.Insert(db, 5, "123-456-7894", "contacts")
	if errInserting != nil {
		// fmt.Printf("%v\n", errInserting)
		return errInserting
	}
	if err := tx.Commit(); err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}

	tx, err = db.Begin()
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}
	_, errInserting = phoneDB.Insert(db, 6, "123-456-7890", "contacts")
	if errInserting != nil {
		// fmt.Printf("%v\n", errInserting)
		return errInserting
	}
	if err := tx.Commit(); err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}

	tx, err = db.Begin()
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}
	_, errInserting = phoneDB.Insert(db, 7, "1234567892", "contacts")
	if errInserting != nil {
		// fmt.Printf("%v\n", errInserting)
		return errInserting
	}
	if err := tx.Commit(); err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}

	tx, err = db.Begin()
	if err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}
	_, errInserting = phoneDB.Insert(db, 8, "(123)456-7892", "contacts")
	if errInserting != nil {
		// fmt.Printf("%v\n", errInserting)
		return errInserting
	}
	if err := tx.Commit(); err != nil {
		// fmt.Printf("%v\n", err)
		return err
	}

	return nil
}

func resetDB(db *sql.DB) error {
	db, errReset := phoneDB.Reset(db, "contacts")
	if errReset != nil {
		return errReset
	}
	return phoneDB.CreateTable(db, "contacts")
}

func reportData(db *sql.DB) ([]string, error) {
	rows, errReporting := phoneDB.Report(db, "contacts")
	if errReporting != nil {
		// fmt.Printf("%v\n", errReporting)
		return nil, errReporting
	}

	defer rows.Close()

	var (
		id    int
		phone string
	)

	var result = make([]string, 0)

	for rows.Next() {
		err := rows.Scan(&id, &phone)
		if err != nil {
			// fmt.Printf("%v\n", err)
			return nil, err
		}
		// fmt.Printf("%d, %s\n", id, phone)
		result = append(result, phoneTask.Normalize(phone))

	}

	if err := rows.Err(); err != nil {
		// fmt.Printf("%v\n", err)
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		// fmt.Printf("%v\n", err)
		return nil, err
	}
	err = resetDB(db)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		// fmt.Printf("%v\n", err)
		return nil, err
	}
	return result, nil
}

func main() {
	//instantiating the db
	db, errorCreatingDB := phoneDB.DB()
	if errorCreatingDB != nil {
		fmt.Printf("%v\n", errorCreatingDB)
		return
	}

	defer db.Close()

	//check if connection is alive
	if err := db.Ping(); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	//creating table contacts
	errCreatingTable := phoneDB.CreateTable(db, "contacts")
	if errCreatingTable != nil {
		fmt.Printf("%v\n", errCreatingTable)
		return
	}

	// insert seed data
	if err := insertData(db); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	//fetching and normalizing the data
	res, err := reportData(db)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	//de-duplicating the normalized data
	unique := make(map[string]bool)
	normalizedData := make([]string, 0)
	for _, v := range res {
		if _, done := unique[v]; done {
			continue
		}
		unique[v] = true
		normalizedData = append(normalizedData, v)
	}

	//inserting the de-duplicated normalized data
	for i, v := range normalizedData {
		tx, err := db.Begin()
		if err != nil {
			fmt.Printf("%v\n", err)
			break
		}
		_, errInserting := phoneDB.Insert(db, i, v, "contacts")
		if errInserting != nil {
			fmt.Printf("%v\n", errInserting)
			break
		}
		if err := tx.Commit(); err != nil {
			fmt.Printf("%v\n", err)
			break
		}
	}

	//Final fetch from db
	rows, errReporting := phoneDB.Report(db, "contacts")
	if errReporting != nil {
		fmt.Printf("%v\n", errReporting)
		// return nil, errReporting
	}

	defer rows.Close()

	var (
		id    int
		phone string
	)

	for rows.Next() {
		err := rows.Scan(&id, &phone)
		if err != nil {
			fmt.Printf("%v\n", err)
			// return nil, err
		}
		fmt.Printf("%d, %s\n", id, phone)
		// result = append(result, phoneTask.Normalize(phone))

	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("%v\n", err)
		// return nil, err
	}
	err = resetDB(db)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	return
}
