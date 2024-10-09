package store

import (
	"database/sql"
	"fmt"
	"io"
	"os"
)

type AppDatabase struct {
	db *sql.DB
}

func (app *AppDatabase) Close() {
	app.db.Close()
}

type SaveMessageInput struct {
	JID         string
	Name        string
	ChannelJID  string
	Message     string
	MessageType string
	Command     string
	Timestamp   string
	IsGroup     bool
}

func (app *AppDatabase) SaveMessage(msg SaveMessageInput) {
	insertQuery := `INSERT INTO messages (jid, name, channel_jid, message, type, command, timestamp, is_group) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	isGroupInt := 0
	if msg.IsGroup {
		isGroupInt = 1
	}

	_, err := app.db.Exec(insertQuery, msg.JID, msg.Name, msg.ChannelJID, msg.Message, msg.MessageType, msg.Command, msg.Timestamp, isGroupInt)
	if err != nil {
		fmt.Println("Error Inserting on database", err)
	}
}

func NewAppDatabase() (*AppDatabase, error) {
	db, err := initDB()

	if err != nil {
		fmt.Println("Error starting database", err)
		return &AppDatabase{}, err
	}

	return &AppDatabase{
		db: db,
	}, nil
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./databases/app.db")

	if err != nil {
		fmt.Println("Error opening database")
		return nil, err
	}

	createTablesSQL, err := os.Open("./pkg/store/sql/create_tables.sql")

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
