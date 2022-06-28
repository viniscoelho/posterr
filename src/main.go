package main

import (
	"flag"
	"log"

	storagedb "posterr/src/storage/db"
)

var initDB = flag.Bool("init-db", false, "creates a database and its tables")

func main() {
	flag.Parse()

	db := storagedb.NewDatabase(storagedb.DatabaseName)
	if *initDB {
		if err := db.InitializeDB(); err != nil {
			log.Fatalf("An error occurred: %s", err)
		}
	}

	/*
		ls, err := storage.NewLottoBacked()
		if err != nil {
			log.Fatalf("could not initialize storage: %s", err)
		}

		r := router.CreateRoutes(ls)
		c := cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodGet,
				http.MethodPost,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"Accept", "Content-Type"},
			AllowCredentials: false,
		})
		handler := c.Handler(r)

		s := &http.Server{
			Handler:      handler,
			ReadTimeout:  0,
			WriteTimeout: 0,
			Addr:         ":3000",
			IdleTimeout:  time.Second * 60,
		}
		log.Fatal(s.ListenAndServe())
	*/
}
