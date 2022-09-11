package gc

import (
	"database/sql"
	"log"
	"time"
)

func GC(db *sql.DB, tableName string, columnName string, timer time.Duration) {
	for {
		time.Sleep(timer)
		_, err := db.Query("DELETE FROM $1 WHERE expire_time < $1", tableName, columnName, time.Now().Unix())
		if err != nil {
			log.Printf("Error in GC: %v", err)
		}
	}
}
