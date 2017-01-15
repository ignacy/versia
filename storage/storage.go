package storage

import (
	"database/sql"
	"log"
	"os"
	"strings"
)

var (
	modelName = os.Getenv("VERSIA_MODEL_NAME")
)

type Version struct {
	Id             string
	Event          string
	Whodunnit      string
	Object         string
	Object_changes string
}

type Model struct {
	Id int
}

var db *sql.DB

func InitDB(pgString string) {
	db, err := sql.Open("postgres", pgString)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

func ListModels() []Model {
	rows, err := db.Query("SELECT id FROM " + modelName + "s ORDER BY id DESC")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	models := []Model{}

	for rows.Next() {
		var i Model

		err := rows.Scan(&i.Id)
		if err != nil {
			log.Fatal(err)
		}
		models = append(models, i)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return models
}

func FindVersions(id int) []Version {
	rows, err := db.Query(`
       SELECT
           id,
           COALESCE(event, '') as event,
           COALESCE(whodunnit, '') as whodunnit,
           COALESCE(object, '') as object,
           COALESCE(object_changes, '') as object_changes
       FROM versions WHERE item_type = $1 AND item_id = $2
       ORDER BY id DESC
      `, strings.Title(modelName), id)

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
