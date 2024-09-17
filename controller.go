package crud

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	util "github.com/polyglot-chameleon/goutil"
)

type CRUD struct {
	db *sql.DB
}

func (mc *CRUD) Connect() error {
	var err error
	mc.db, err = sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_URL"))
	util.Check(err)
	return err
}

func (mc *CRUD) CloneModel() error {
	rows, err := mc.db.Query("pragma table_info(posts);")
	util.Check(err)

	rb := resourceBuilder{}
	rb.init()

	var curCol tableInfo

	for rows.Next() {
		rows.Scan(&curCol.cid, &curCol.name, &curCol.ctype, &curCol.notnull, &curCol.dflt_val, &curCol.pk)
		rb.add(curCol)
	}
	rb.close()
	rb.build()

	return err
}

func (mc *CRUD) Create(new Resource) (sql.Result, error) {
	result, err := mc.db.Exec(fmt.Sprintf("INSERT INTO posts(title, body) VALUES ('%s', '%s')", new.title, new.body))

	util.Check(err)

	lastInsertId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()

	log.Printf("Inserted post(id=%v)\tRows affected: #%v", lastInsertId, rowsAffected)

	return result, err
}

func (mc *CRUD) Read(resourceID int64) (Resource, error) {
	rows, err := mc.db.Query(fmt.Sprintf("SELECT title, body FROM posts WHERE id = %v", resourceID))
	util.Check(err)

	defer rows.Close()

	resource := Resource{}

	for rows.Next() {
		rows.Scan(&resource.title, &resource.body)
	}

	return resource, err
}

func (mc *CRUD) All() ([]Resource, error) {
	rows, err := mc.db.Query("SELECT title, body FROM posts;")
	util.Check(err)

	defer rows.Close()

	var stored []Resource
	resource := Resource{}

	for rows.Next() {
		rows.Scan(&resource.title, &resource.body)
		stored = append(stored, resource)
	}

	log.Printf("Read %v posts", len(stored))

	return stored, err
}

func (mc *CRUD) Delete(resourceId int64) (sql.Result, error) {
	result, err := mc.db.Exec(fmt.Sprintf("DELETE FROM posts WHERE id = %v;", resourceId))
	util.Check(err)
	return result, err
}
