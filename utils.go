package main

import (
	"log"
	"flag"
)

func utilsCreateAdminUser(name string) error {
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

func utilsHandleOptions() {
	flagCreateAdmin := flag.String("createadmin", "",
		"-createadmin \"<name>\"")
	flagListUsers := flag.Bool("users", false, "-users")
	flag.Parse()

	if *flagCreateAdmin != "" {
		utilsCreateAdminUser(*flagCreateAdmin)
	}
	if *flagListUsers {
		utilsListUsers()
	}
}