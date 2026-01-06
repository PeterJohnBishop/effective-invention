package pgdb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func PostGres() *sql.DB {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening db: ", err)
		return nil
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to the db: ", err)
		return nil
	}

	log.Println("Connected to Postgres volume in Docker!")
	return db

}

// one...
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ...to many
type Post struct {
	ID      int
	UserID  int // This is the Foreign Key linking to User.ID
	Title   string
	Content string
}

func Example(db *sql.DB) {

	// CREATE TABLE
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// CREATE
	var newID int
	err = db.QueryRow("INSERT INTO users(name, email) VALUES($1, $2) RETURNING id", "Alice", "alice@example.com").Scan(&newID) // .Scan() scans the return value or the row and assigns them to pointers
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created User ID: %d\n", newID)

	// READ
	var u User
	err = db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", newID).Scan(&u.ID, &u.Name, &u.Email)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Read User: %+v\n", u)

	// UPDATE
	_, err = db.Exec("UPDATE users SET name = $1 WHERE id = $2", "Alice Smith", newID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User updated successfully.")

	// DELETE
	_, err = db.Exec("DELETE FROM users WHERE id = $1", newID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User deleted successfully.")
}

func ExampleMultipleROws(db *sql.DB) {
	rows, err := db.Query("SELECT name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		rows.Scan(&name) // Scans the current row's name
		fmt.Println(name)
	}
}

func ExampleRelationshop(db *sql.DB, postID int) {
	query := `
		SELECT p.id, p.title, u.name 
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = $1`

	var postIDResult int
	var title string
	var authorName string

	err := db.QueryRow(query, postID).Scan(&postIDResult, &title, &authorName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Post: %s | Written by: %s\n", title, authorName)
}
