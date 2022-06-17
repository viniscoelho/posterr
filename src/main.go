package main

import (
	"flag"

	"posterr/src/db"
)

var initDB = flag.Bool("init-db", false, "creates a database and its tables")

func main() {
	flag.Parse()
	if *initDB {
		db.InitializeDatabase()
	}
}
