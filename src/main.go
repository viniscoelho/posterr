package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"posterr/src/router"
	"posterr/src/storage"
	storagedb "posterr/src/storage/db"

	"github.com/rs/cors"
)

var (
	initDB = flag.Bool("init-db", false, "creates a database and its tables")
	port   = flag.Int("port", 3000, "application port ")
)

func main() {
	flag.Parse()

	db := storagedb.NewDatabase(storagedb.DatabaseName)
	if *initDB {
		if err := db.InitializeDB(); err != nil {
			log.Fatalf("An error occurred: %s", err)
		}
	}

	posts := storage.NewPosterrBacked(db)
	users := storage.NewUserBacked(db, posts)

	r := router.CreateRoutes(posts, users)
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
		Addr:         fmt.Sprintf(":%d", *port),
		IdleTimeout:  time.Second * 60,
	}
	log.Fatal(s.ListenAndServe())
}
