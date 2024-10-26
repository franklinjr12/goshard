package main

import (
	"fmt"
	"goshard/internal/database"
	"goshard/internal/dbmapper"
	"goshard/internal/servicelistener"
)

func main() {
	fmt.Println("Running")

	dbDefault := database.DbTestConnectionParams

	dbParams1 := dbDefault
	dbParams1.Dbname = dbParams1.Dbname + "1"
	db1Dsn := database.BuildConnectionString(dbParams1)

	dbParams2 := dbDefault
	dbParams2.Dbname = dbParams2.Dbname + "2"
	db2Dsn := database.BuildConnectionString(dbParams2)

	// init the mapper with some data for testing
	dbmapper.AddDbMapId(0, database.DbTestConnectionString)
	dbmapper.AddDbMapId(1, db1Dsn)
	dbmapper.AddDbMapId(2, db2Dsn)

	servicelistener.ListenAndServe()
}
