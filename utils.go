package main

import (
	"log"
)

func utilsCreateAdminUser(name sting) error {
	err := dbCreateUser(name, true, true, true)
	if err != nil {
		panic(err)
	}
	return nil
}

func utilsListUsers() {
	var users []User
	err := db.All(&users)
	if err != nil {
		panic(err)
	}
	log.Printf("%s", users)
}
