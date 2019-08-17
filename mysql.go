package hemmingway

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

// I am using a struct here to represent everything to do with connecting to a db
// In this case MySQL - and then using methods off of this to connect, query, insert etc
type MySQL struct {
	Host string
	Port string
	Username string
	Password string
	Database string
	connection *sql.DB
}


func (m *MySQL) connect() string {
	url := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", m.Username, m.Password, m.Host, m.Port, m.Database)

	db, err := sql.Open(
		"mysql",
		url,
	)
	if err != nil {
		log.Fatal(err)
	}

	m.connection = db

	return "ok"
}

// This struct is used as part of the below GetAllTransactions() Query Function
// This is the
/*
CREATE TABLE transactions (
	id TEXT,
	description TEXT,
	amount TEXT
);
*/

type TransactionsTable struct {
	ID string
	Description string
	Amount string
}

// GetAllTransactions is an example of a very specific function that can be written
// to achieve some desired result. Above we essentially create a Data Structure
// that looks like the data we're querying for, and here we're writing that query
// then using `Scan()` to transpose the results into a struct for each row and appending
// to an array of those results.
// TODO: Figure out a way to be able to pass a query into this function and automatically Scan into a map
func (m *MySQL) GetAllTransactions() []TransactionsTable {
	stmt, err := m.connection.Prepare("SELECT * FROM transactions")
	FailOnError(err, "Prepare Statement Failed:")

	rows, err := stmt.Query()
	FailOnError(err, "Query Failed:")

	output := []TransactionsTable{}
	var row TransactionsTable

	for rows.Next() {
		err := rows.Scan(&row.ID, &row.Description, &row.Amount)
		FailOnError(err, "Assignment to struct failed:")
		output = append(output, row)

	}
	return output
}

// Stub for querying without knowing columns/types/rows
//func (m *MySQL) Query() {
//	cols, err := rows.Columns() // Remember to check err afterwards
//	vals := make([]interface{}, len(cols))
//	for i, _ := range cols {
//		vals[i] = new(sql.RawBytes)
//	}
//	for rows.Next() {
//		err = rows.Scan(vals...)
//		// Now you can check each element of vals for nil-ness,
//		// and you can use type introspection and type assertions
//		// to fetch the column into a typed variable.
//	}
//}

func (m *MySQL) SaveTransaction(transaction TransactionsTable) {
	stmt, err := m.connection.Prepare(
		"INSERT INTO transactions(id, description, amount) VALUES(?, ?, ?)",
	)
	FailOnError(err, "Preparing Insert Failed:")

	_, err = stmt.Exec(transaction.ID, transaction.Description, transaction.Amount)
	FailOnError(err, "Insert Operation Failed:")
}



func (m *MySQL) CleanUp() {
	err := m.connection.Close()
	FailOnError(err, "Connection Close err:")
}
