package store

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
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

func (app *AppDatabase) CheckIfUserIsBlocked(message *events.Message) bool {
	jid := message.Info.Sender.String()
  number := strings.Split(jid, ":")[0]

  var blocked int
  query := "SELECT blocked FROM users WHERE number = ?"

  app.db.QueryRow(query, number).Scan(&blocked)

  return blocked == 1
}

func (app *AppDatabase) BlockUserByNumber(number string) {
  query := "UPDATE users SET blocked = 1 where number like ?"

	_, err := app.db.Exec(query, "%" + number + "%")
	if err != nil {
		fmt.Println("Error Blocking user", err)
	}
}

func (app *AppDatabase) UnblockUserByNumber(number string) {
  query := "UPDATE users SET blocked = 0 where number like ?"

	_, err := app.db.Exec(query, "%" + number + "%")
	if err != nil {
		fmt.Println("Error Blocking user", err)
	}
}

func (app *AppDatabase) SaveUser(message *events.Message) {
	jid := message.Info.Sender.String()
	name := message.Info.PushName
  number := strings.Split(jid, ":")[0]

  var foundNumber string
  app.db.QueryRow("SELECT number FROM users where number = ?", number).Scan(&foundNumber)

  if foundNumber != ""  {
    return
  }

  fmt.Println("Saving new user", jid, name, number);

	insertQuery := `INSERT INTO users (number, jid, name, blocked) 
					VALUES (?, ?, ?, 0)`

  _, err := app.db.Exec(insertQuery, number, jid, name)
	if err != nil {
		fmt.Println("Error Inserting on database", err)
	}
}

func (app *AppDatabase) ExecSql(query string) (string, error) {
	rows, err := app.db.Query(query)

	if err != nil {
		fmt.Println("Error running query")
		return "", nil
	}

	defer rows.Close()

	var resultStrings []string
	var rowStrings []string

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		// Create a slice of empty interfaces to hold the column values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		// Assign pointers to each value in the slice
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the pointers
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}

		// Convert the row into a string
		for i, col := range columns {
			val := values[i]

			// Handle NULL values
			if val == nil {
				rowStrings = append(rowStrings, fmt.Sprintf("%s: NULL", col))
			} else {
				rowStrings = append(rowStrings, fmt.Sprintf("%s: %v", col, val))
			}
		}
	}

	resultStrings = append(resultStrings, strings.Join(rowStrings, ", "))
	return strings.Join(resultStrings, "\n"), nil
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
