package main

import (
	"log"
	"testing"

	util "github.com/polyglot-chameleon/goutil"
)

var newPost Resource

func init() {
	util.LoadDotEnv(".env.test")
	if err := Controller.Connect(); err != nil {
		log.Fatal(err)
	}
	newPost = Resource{Title: "NewTestPostTitle", Body: "NewTestPostBody"}
}

func TestCreate(t *testing.T) {
	res, err := Controller.Create(newPost)

	testError(err, t)

	nRowsAffected, err := res.RowsAffected()

	testError(err, t)

	if nRowsAffected != 1 {
		t.Fatalf("RowsAffected: %v != 1", nRowsAffected)
	}

	lastInsertId, _ := res.LastInsertId()
	storedPost, err := Controller.Read(lastInsertId)

	testError(err, t)

	if storedPost.Title != newPost.Title {
		t.Fatalf("%v != %v", storedPost.Title, newPost.Title)
	}

	if storedPost.Body != newPost.Body {
		t.Fatalf("%v != %v", storedPost.Body, newPost.Body)
	}

	Controller.Delete(lastInsertId)
}

func TestRead(t *testing.T) {
	res, _ := Controller.Create(newPost)
	lastInsertId, _ := res.LastInsertId()

	storedPost, err := Controller.Read(lastInsertId)
	testError(err, t)

	if storedPost.Title != newPost.Title {
		t.Fatalf("%v != %v", storedPost.Title, newPost.Title)
	}

	if storedPost.Body != newPost.Body {
		t.Fatalf("%v != %v", storedPost.Body, newPost.Body)
	}
	Controller.Delete(lastInsertId)
}

func TestAll(t *testing.T) {
	res, _ := Controller.Create(newPost)
	lastInsertId, _ := res.LastInsertId()

	storedPosts, err := Controller.All()
	testError(err, t)

	nPosts := len(storedPosts)

	if nPosts != 1 {
		t.Fatalf("%v != 1", nPosts)
	}
	Controller.Delete(lastInsertId)
}

func TestDelete(t *testing.T) {
	res, _ := Controller.Create(newPost)
	lastInsertId, _ := res.LastInsertId()

	res, err := Controller.Delete(lastInsertId)
	testError(err, t)

	nRowsAffected, err := res.RowsAffected()

	testError(err, t)

	if nRowsAffected != 1 {
		t.Fatalf("RowsAffected: %v != 1", nRowsAffected)
	}

	storedPost, err := Controller.Read(lastInsertId)
	testError(err, t)

	if storedPost.Title != "" {
		t.Fatalf("storedPost.Title = %s", storedPost.Title)
	}
	if storedPost.Body != "" {
		t.Fatalf("storedPost.Body = %s", storedPost.Body)
	}
}

func testError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
