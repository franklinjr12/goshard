package main

import (
	"fmt"
	"goshard/internal/database"
	"goshard/internal/dbmapper"
	"goshard/internal/servicelistener"
)

func main() {
	fmt.Println("Running")

	// init the mapper with some data for testing
	dbmapper.AddDbMapId(1, database.DbTestConnectionString)

	connectionParams := database.DbTestConnectionParams
	connectionParams.Dbname = ""
	dsn := database.BuildConnectionString(connectionParams)
	fmt.Println("conencting to db")
	// db, err := database.Connect("host=localhost port=5432 user=postgres password=postgres sslmode=disable")
	db, err := database.Connect(dsn)
	if err != nil {
		fmt.Println(err)
		database.Close(db)
		return
	}
	fmt.Println("creating database")
	dbname := "testapplication1"
	res, err := db.Exec("CREATE DATABASE " + dbname)
	if err != nil {
		fmt.Println(err)
		database.Close(db)
		return
	}
	fmt.Println(res)
	database.Close(db)
	connectionParams.Dbname = dbname
	dsn = database.BuildConnectionString(connectionParams)
	fmt.Println("conencting to db")
	db, err = database.Connect(dsn)
	if err != nil {
		fmt.Println(err)
		database.Close(db)
		return
	}
	fmt.Println("creating from schema")
	schemaStr, err := database.ReadSchemaFromFile("sql/schema.sql")
	if err != nil {
		fmt.Println(err)
		database.Close(db)
		return
	}
	err = database.CreateDatabaseFromSchema(db, schemaStr)
	if err != nil {
		fmt.Println(err)
	}
	database.Close(db)
	return // just for testing

	servicelistener.ListenAndServe()
}
