package main

import (
	"fmt"
	"goshard/internal/config"
	"goshard/internal/dbmapper"
	"goshard/internal/servicelistener"
)

func main() {
	fmt.Println("Running")
	err := config.LoadDbMappings()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Mappings:\n%v", dbmapper.DbMapsByUserId)
	servicelistener.ListenAndServe()
}
