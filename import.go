package mysqlImport

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

func createTable(db *sql.DB, tableName string, columns []string) error {
	// Escape table name
	escapedTableName := fmt.Sprintf("`%s`", tableName)

	// Create column definitions
	columnDefs := []string{"id INT AUTO_INCREMENT PRIMARY KEY"}
	for _, col := range columns {
		columnDefs = append(columnDefs, fmt.Sprintf("`%s` TEXT", col))
	}

	// Create table with sequential ID and index on the first column
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", escapedTableName, strings.Join(columnDefs, ", "))
	_, err := db.Exec(query)
	if err != nil {
		return errors.New("error: " + err.Error() + " query: " + query)
	}
	return nil
}

func ImportCSV(db *sql.DB, tableName, filePath string) error {
	// Open CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Retrieve column names from the first row and replace spaces with underscores
	columns := make([]string, len(records[0]))
	for i, col := range records[0] {
		columns[i] = strings.ReplaceAll(strings.ToLower(strings.ReplaceAll(col, " ", "_")), "id", "id_2")
	}

	// Create table
	err = createTable(db, tableName, columns)
	if err != nil {
		return err
	}

	// Prepare SQL statement
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		"`"+strings.Join(columns, "`, `")+"`",
		strings.Repeat("?, ", len(columns)-1)+"?")
	stmt, err := db.Prepare(query)
	if err != nil {
		return errors.New("error: " + err.Error() + " query: " + query)
	}
	defer stmt.Close()

	// Insert records into the database (skip the first row, as it contains column names)
	for _, record := range records[1:] {
		args := make([]interface{}, len(record))
		for i, v := range record {
			args[i] = sql.NullString{String: v, Valid: v != ""}
		}

		_, err := stmt.Exec(args...)
		if err != nil {
			return err
		}
	}

	return nil
}
