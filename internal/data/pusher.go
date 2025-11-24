package data

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/claytonphenderson/scope/internal/models"
	"go.uber.org/zap"
)

func Write(receivedEvent models.DbRecord, tables map[string]bool, db *sql.DB, logger *zap.Logger) {
	// create new table if not exists
	if _, ok := tables[receivedEvent.Event]; !ok {
		_, err := db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s"(
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					timestamp TEXT,
					payload JSON
				)`, receivedEvent.Event))

		if err != nil {
			log.Fatal("Could not create new table "+receivedEvent.Event, err)
		}

		tables[receivedEvent.Event] = true

		logger.Info("Created new table " + receivedEvent.Event)
	}

	// insert the new record
	command := fmt.Sprintf(`INSERT INTO "%s"(timestamp,payload) VALUES(?,?)`, receivedEvent.Event)
	_, insertErr := db.Exec(command, receivedEvent.Timestamp, receivedEvent.Payload)
	if insertErr != nil {
		log.Fatal("Could not insert into "+receivedEvent.Event, insertErr)
	}
	logger.Info("Inserted into " + receivedEvent.Event)
}
