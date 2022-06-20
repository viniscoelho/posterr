package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

/*
	DATABASE_URL:               postgres://{user}:{password}@{hostname}:{port}/{database-name}
*/

const (
	DatabaseName = "posterr"

	connectionURL  = "postgres://localhost:5432"
	envDatabaseURL = "DATABASE_URL"

	databaseCreationErrorCode = "SQLSTATE 42P04"
	tableCreationErrorCode    = "SQLSTATE 42P07"
)

type postgresDB struct {
	databaseName string
}

type ConnectDB interface {
	Connect() (*pgxpool.Pool, error)
	InitializeDB()
}

func NewDatabase(dbName string) *postgresDB {
	return &postgresDB{
		databaseName: dbName,
	}
}

// Connect connects to the default database
func (pg *postgresDB) Connect() (*pgxpool.Pool, error) {
	conn, err := connect(pg.databaseName)
	if err != nil {
		return nil, err
	}

	return conn, err
}

// InitializeDB initializes the database if an init flag
// is given in the main function
func (pg *postgresDB) InitializeDB() {
	if err := pg.createDatabase(); err != nil {
		if !databaseExists(err) {
			log.Fatalf("Database creation failed: %s", err)
		}
		log.Print("Database already exists. Skipping...")
	}

	conn, err := connect(pg.databaseName)
	if err != nil {
		log.Fatalf("Database connection failed: %s", err)
	}
	createTables(conn)
}

func connect(dbName string) (*pgxpool.Pool, error) {
	err := os.Setenv(envDatabaseURL, fmt.Sprintf("%s/%s", connectionURL, dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to set %s: %w", envDatabaseURL, err)
	}

	conn, err := pgxpool.Connect(context.Background(), os.Getenv(envDatabaseURL))
	if err != nil {
		return nil, fmt.Errorf("pool connection failed: %w", err)
	}

	return conn, err
}

func (pg *postgresDB) createDatabase() error {
	conn, err := connect("")
	if err != nil {
		log.Fatalf("Database connection failed: %s", err)
	}
	defer conn.Close()

	_, err = conn.Exec(context.Background(), fmt.Sprintf(`CREATE DATABASE %s`, pg.databaseName))
	if err != nil {
		return err
	}

	log.Printf("Database %s created!", pg.databaseName)
	return nil
}

func createTables(conn *pgxpool.Pool) {
	defer conn.Close()

	if err := createUsersTable(conn); err != nil {
		if !tableExists(err) {
			log.Fatalf("Table users creation failed: %s", err)
		}
		log.Print("Table users already exists. Skipping...")
	}

	if err := createPostsTable(conn); err != nil {
		if !tableExists(err) {
			log.Fatalf("Table posts creation failed: %s", err)
		}
		log.Print("Table posts already exists. Skipping...")
	}

	if err := createFollowersTable(conn); err != nil {
		if !tableExists(err) {
			log.Fatalf("Table followers creation failed: %s", err)
		}
		log.Print("Table followers already exists. Skipping...")
	}
}

func createUsersTable(conn *pgxpool.Pool) error {
	table := `CREATE TABLE users(
        username VARCHAR (14) NOT NULL PRIMARY KEY,
        joined_at TIMESTAMPTZ DEFAULT NOW())`

	_, err := conn.Exec(context.Background(), table)
	if err != nil {
		return err
	}

	log.Printf("Table users created!")
	return nil
}

func createPostsTable(conn *pgxpool.Pool) error {
	table := `CREATE TABLE posts(
        post_id SERIAL PRIMARY KEY,
        username VARCHAR (14) NOT NULL REFERENCES users (username),
        content VARCHAR (777),
        reposted_id INTEGER,
        created_at TIMESTAMPTZ DEFAULT NOW())`

	_, err := conn.Exec(context.Background(), table)
	if err != nil {
		return err
	}

	log.Printf("Table posts created!")
	return nil
}

func createFollowersTable(conn *pgxpool.Pool) error {
	table := `CREATE TABLE followers(
        username VARCHAR (14) NOT NULL,
        followed_by VARCHAR (14) NOT NULL,
        FOREIGN KEY (username) REFERENCES users (username),
        FOREIGN KEY (followed_by) REFERENCES users (username),
        PRIMARY KEY(username, followed_by))`

	_, err := conn.Exec(context.Background(), table)
	if err != nil {
		return err
	}

	log.Printf("Table followers created!")
	return nil
}

func databaseExists(err error) bool {
	if strings.Contains(err.Error(), databaseCreationErrorCode) {
		return true
	}
	return false
}

func tableExists(err error) bool {
	if strings.Contains(err.Error(), tableCreationErrorCode) {
		return true
	}
	return false
}
