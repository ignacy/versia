package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var (
	event          string
	whodunnit      string
	object         string
	object_changes string
	id             int
)

// 29

func main() {
	db, err := sql.Open("postgres", "dbname=advanon_development sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`
       SELECT
           id,
           COALESCE(event, '') as event,
           COALESCE(whodunnit, '') as whodunnit,
           COALESCE(object, '') as object,
           COALESCE(object_changes, '') as object_changes
       FROM versions WHERE item_type = 'Invoice' AND item_id = $1
      `, 29)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &event, &whodunnit, &object, &object_changes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Version id: %d Event type %s Authored by %s\n", id, event, whodunnit)
		fmt.Printf("Object before changes: \n %s \n", object)
		fmt.Printf("Changes made: \n %s \n\n", object_changes)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}
