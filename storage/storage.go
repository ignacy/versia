package storage

import (
	"database/sql"
	"log"
	"os"
)

var (
	pgString = os.Getenv("VERSIA_PG_STRING")
)

type Version struct {
	Id             string
	Event          string
	Whodunnit      string
	Object         string
	Object_changes string
}

type Invoice struct {
	Id int
}

func ListInvoices() []Invoice {
	db, err := sql.Open("postgres", pgString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(`
       SELECT
           id
       FROM invoices
       ORDER BY id DESC
      `)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	invoices := []Invoice{}

	for rows.Next() {
		var i Invoice

		err := rows.Scan(&i.Id)
		if err != nil {
			log.Fatal(err)
		}
		invoices = append(invoices, i)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return invoices
}

func FindVersions(id int) []Version {
	db, err := sql.Open("postgres", pgString)
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
       ORDER BY id DESC
      `, id)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	versions := []Version{}

	for rows.Next() {
		var v Version

		err := rows.Scan(&v.Id, &v.Event, &v.Whodunnit, &v.Object, &v.Object_changes)
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
