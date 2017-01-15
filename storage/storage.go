package storage

import (
	"database/sql"
	"os"
	"strings"
)

var (
	modelName = os.Getenv("VERSIA_MODEL_NAME")
)

type Datastore interface {
	ListModels() ([]Model, error)
	FindVersions(id int) ([]Version, error)
}

type DB struct {
	*sql.DB
}

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

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) ListModels() ([]Model, error) {
	rows, err := db.Query("SELECT id FROM " + modelName + "s ORDER BY id DESC")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	models := []Model{}

	for rows.Next() {
		var i Model

		err := rows.Scan(&i.Id)
		if err != nil {
			return nil, err
		}
		models = append(models, i)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (db *DB) FindVersions(id int) ([]Version, error) {
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
		return nil, err
	}
	defer rows.Close()

	versions := []Version{}

	for rows.Next() {
		var v Version

		err := rows.Scan(&v.Id, &v.Event, &v.Whodunnit, &v.Object, &v.Object_changes)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return versions, nil
}
