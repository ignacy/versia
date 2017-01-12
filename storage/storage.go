package storage

import (
	"database/sql"
	"log"
)

type Version struct {
	id             string
	event          string
	whodunnit      string
	object         string
	object_changes string
}

func FindVersions(id int) []Version {
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
      `, id)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	versions := []Version{}

	for rows.Next() {
		var v Version

		err := rows.Scan(&v.id, &v.event, &v.whodunnit, &v.object, &v.object_changes)
		if err != nil {
			log.Fatal(err)
		}
		versions = append(versions, v)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return versions
}
