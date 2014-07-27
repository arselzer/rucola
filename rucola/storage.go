package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"encoding/json"
	"fmt"
)

var db *leveldb.DB

func getUser(id string) (User, error) {
	fmt.Println(id)
	jsonUser, err := db.Get([]byte(id), nil)

	if err != nil {
		return User{}, err
	}

	var user User

	json.Unmarshal(jsonUser, &user)

	return user, nil
}

func saveUser(user User) error {
	jsonUser, _ := json.Marshal(user)

	err := db.Put([]byte(user.Name), jsonUser, nil)
	if err != nil {
		return err
	}
	return nil
}

func initDb(dir string) error {
	var err error
	db, err = leveldb.OpenFile(dir, nil)

	if err != nil {
		return err
	}

	return nil
}


func closeDb() error {
	return db.Close()
}
