package main

import (
	"fmt"
	"goshard/internal/database"
)

func main() {
	connectionParams := database.DbTestConnectionParams
	connectionParams.Dbname = ""
	dsn := database.BuildConnectionString(connectionParams)
	fmt.Println("conencting to db")
	db, err := database.Connect(dsn)
	if err != nil {
		fmt.Println(err)
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
}
