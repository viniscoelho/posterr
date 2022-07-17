package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"posterr/src/router"
	storagedb "posterr/src/storage/db"
	"posterr/src/storage/posterr"
	storageusers "posterr/src/storage/users"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
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
			logrus.Fatalf("An error occurred: %s", err)
		}
	}

	posts := posterr.NewPosterrBacked(db)
	users := storageusers.NewUserBacked(db)

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
	logrus.Fatal(s.ListenAndServe())
}
