package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

func Connect(host, user, password, dbname string) (*sql.DB, error) {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable",
		host, user, password, dbname)

	var db *sql.DB
	var err error

	// Retry logic
	for i := 0; i < 5; i++ {
		db, err = sql.Open("sqlserver", connString)
		if err != nil {
			log.Printf("Error creating connection pool: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err == nil {
			log.Println("Successfully connected to Azure SQL")
			return db, nil
		}

		log.Printf("Failed to ping DB (attempt %d/5): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after retries: %v", err)
}
