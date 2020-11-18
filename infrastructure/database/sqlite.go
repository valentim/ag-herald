package database

import (
	"database/sql"
	"fmt"
	"log"

	// This is the pattern used for the database drivers
	_ "github.com/mattn/go-sqlite3"
)

// Database represents the database source
type Database struct {
	Name string
}

// Query represents the query statement
type Query struct {
	Statement string
}

// LocalAccount is the data structure to save accounts
type LocalAccount struct {
	TeamID      string
	AccountGUID string
	AccessToken string
}

// Setup is responsable for create all the database structure
func (d Database) Setup() {
	db := open(d)

	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS account (id integer not null primary key, teamID text, accountGUID text, accessToken text);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

// Insert is the generic method to insert new data
func (a LocalAccount) Insert(d Database) (int64, error) {
	db := open(d)

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO account(teamID, accountGUID) values(?,?)")
	if err != nil {
		fmt.Println("stmtErr", err)
		return 0, err
	}

	result, insertErr := stmt.Exec(a.TeamID, a.AccountGUID)

	if insertErr != nil {
		fmt.Println("insertErr", insertErr)
		return 0, insertErr
	}

	return result.RowsAffected()
}

// UpdateAccessToken is the specific method to update access token
func (a LocalAccount) UpdateAccessToken(d Database) (int64, error) {
	result, err := Update(d, "UPDATE account SET accessToken = ? WHERE teamID = ?", []interface{}{a.AccessToken, a.TeamID})

	if err != nil {
		fmt.Println("updateErr", err)
		return 0, err
	}

	return result, nil
}

// UpdateAccountGUID is the specific method to update account guid
func (a LocalAccount) UpdateAccountGUID(d Database) (int64, error) {

	result, err := Update(d, "UPDATE account SET accountGUID = ? WHERE teamID = ?", []interface{}{a.AccountGUID, a.TeamID})

	if err != nil {
		fmt.Println("updateErr", err)
		return 0, err
	}

	return result, nil
}

// Update is the generic method to update
func Update(d Database, prepared string, values []interface{}) (int64, error) {
	db := open(d)

	defer db.Close()

	stmt, err := db.Prepare(prepared)
	if err != nil {
		fmt.Println("stmtErr", err)
		return 0, err
	}
	fmt.Println(values...)
	result, updateErr := stmt.Exec(values...)
	if updateErr != nil {
		fmt.Println("updateErr", updateErr)
		return 0, updateErr
	}

	return result.RowsAffected()
}

// Get is the generic method to get database data
func (q Query) Get(d Database) string {
	db := open(d)

	defer db.Close()

	log.Println(q.Statement)
	rows, stmtErr := db.Query(q.Statement)
	if stmtErr != nil {
		log.Fatal("stmtErr", stmtErr)
	}

	defer rows.Close()

	var accountGUID string

	if rows.Next() {
		rows.Scan(&accountGUID)
	}

	return accountGUID
}

func open(d Database) *sql.DB {
	db, err := sql.Open("sqlite3", fmt.Sprintf("infrastructure/database/%s", d.Name))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
