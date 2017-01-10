package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var (
	amount string
	id     int
)

func main() {
	db, err := sql.Open("postgres", "dbname=advanon_development sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, amount FROM invoices WHERE status = $1", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &amount)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(id, amount)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
