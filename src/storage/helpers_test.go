package storage

import (
	"context"
	"fmt"
	"log"
	"strings"

	storagedb "posterr/src/storage/db"
	typesrand "posterr/src/types/rand"
)

func dropDatabase(dbName string) {
	db := storagedb.NewDatabase("")

	conn, err := db.Connect()
	if err != nil {
		log.Printf("Failed to connect to database: %s", err)
		return
	}
	defer conn.Close()

	_, err = conn.Exec(context.Background(), fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, dbName))
	if err != nil {
		log.Printf("Failed to drop database %s: %s", dbName, err)
		return
	}

	log.Printf("Database %s dropped!", dbName)
}

func generateDBName() string {
	rs := typesrand.NewPseudoRandomString()
	// database names are restricted to contain only lowercase chars
	return fmt.Sprintf("db%s", strings.ToLower(rs.GenerateAny(8)))
}
