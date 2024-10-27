package main

import (
	"fmt"
	"goshard/internal/database"
)

func insertSchema() {
	connectionParams := database.DbTestConnectionParams
	connectionParams.Dbname = "goshardconfig"
	dsn := database.BuildConnectionString(connectionParams)
	fmt.Println("conencting to db")
	db, err := database.Connect(dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	schemaStr, err := database.ReadSchemaFromFile("sql/schema.sql")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("inserting schema")
	queryStr := "INSERT INTO user_schemas (user_id, schema) VALUES ($1, $2)"
	_, err = db.Exec(queryStr, 1, schemaStr)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	insertSchema()
	return

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
