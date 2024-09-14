package controller

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	util "github.com/polyglot-chameleon/goutil"
)

type Resource struct {
	Title string
	Body  string
}

type Controller struct {
	db *sql.DB
}

func (mc *Controller) Connect() error {
	var err error
	mc.db, err = sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_URL"))
	util.Check(err)
	return err
}

func (mc *Controller) Create(new Resource) (sql.Result, error) {
	result, err := mc.db.Exec(fmt.Sprintf("INSERT INTO posts(title, body) VALUES ('%s', '%s')", new.Title, new.Body))

	util.Check(err)

	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()

	log.Printf("Inserted post(id=%v)\tRows affected: #%v", lastInsertId, rowsAffected)

	return result, err
}

func (mc *Controller) Read(resourceID int64) (Resource, error) {
	rows, err := mc.db.Query(fmt.Sprintf("SELECT title, body FROM posts WHERE id = %v", resourceID))
	util.Check(err)

	defer rows.Close()

	post := Resource{Title: "", Body: ""}

	for rows.Next() {
		rows.Scan(&post.Title, &post.Body)
	}

	return post, err
}

func (mc *Controller) All() ([]Resource, error) {
	rows, err := mc.db.Query("SELECT title, body FROM posts;")
	util.Check(err)

	defer rows.Close()

	var stored []Resource
	post := Resource{Title: "", Body: ""}

	for rows.Next() {
		rows.Scan(&post.Title, &post.Body)
		stored = append(stored, post)
	}

	log.Printf("Read %v posts", len(stored))

	return stored, err
}

func (mc *Controller) Delete(resourceId int64) (sql.Result, error) {
	result, err := mc.db.Exec(fmt.Sprintf("DELETE FROM posts WHERE id = %v;", resourceId))
	util.Check(err)
	return result, err
}
