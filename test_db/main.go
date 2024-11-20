package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	_ "github.com/mattn/go-sqlite3"
)
func main() {
  fmt.Println("Hello World")
  db, err := sql.Open("sqlite3", "./database.db")

  if err != nil {
    fmt.Println("Error", err)
  }
  defer db.Close()
  initDB()
}

func initDB() (*sql.DB, error) {
  db, err := sql.Open("sqlite3", "./database.db")

  if err != nil {
    fmt.Println("Error opening database")
    return nil, err
  }

  createTablesSQL, err := os.Open("../sql/create_tables.sql")

  if err != nil {
    fmt.Println("Error reading sql file")
    return nil, err
  }

  defer createTablesSQL.Close()

  contentSQLBytes, err := io.ReadAll(createTablesSQL)

  if err != nil {
    fmt.Println("Error reading bytes")
    return nil, err
  }
  contentSql := string(contentSQLBytes)

  _, err = db.Exec(contentSql)

  if err != nil {
    fmt.Println("Error creating table", err)
    return nil, err
  }
  
  return db, nil
}
