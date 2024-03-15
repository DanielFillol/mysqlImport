go get github.com/DanielFillol/mysqlImport 

# MySQL Import
It is a Go implementation for importing csv on mysql databases

## Install
You can install mysqlImport using Go's built-in package manager, go get:
``` 
go get github.com/DanielFillol/mysqlImport 
```
## Usage
Here's a simple example of how to use mysqlImport:
```go
package main

import (
	"database/sql"
	"github.com/DanielFillol/mysqlImport"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	DatabaseName = "YOUR_DATABASE"
	TableName    = "TABLE_NAME"
	CsvFilePath  = "FILE_PATH"
)

func main() {
	// Connect to MySQL database
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/"+DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Import CSV into the table
	err = mysqlImport.ImportCSV(db, TableName, CsvFilePath)
	if err != nil {
		log.Fatal(err)
	}
}
```
