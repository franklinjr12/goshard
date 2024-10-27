package servicelistener

// HTTP server to listen for service requests

import (
	"fmt"
	"goshard/internal/config"
	"goshard/internal/database"
	"goshard/internal/dbmapper"
	"goshard/lib/service"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ListenAndServe starts the HTTP server
func ListenAndServe() {
	http.HandleFunc("/query", queryHandler)
	http.ListenAndServe(":8080", nil)
}

func parseUrlToRequest(r *http.Request) service.Request {
	query := r.URL.Query().Get("query")
	shardid, err := strconv.Atoi(r.URL.Query().Get("shardid"))
	if err != nil {
		shardid = 0
	}
	sharduid := r.URL.Query().Get("sharduid")
	userToken := r.URL.Query().Get("usertoken")
	return service.Request{
		Query:     query,
		Shardid:   uint64(shardid),
		Sharduid:  sharduid,
		UserToken: userToken,
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	serviceRequest := parseUrlToRequest(r)
	// for debugging show the url and the params
	fmt.Println("URL:", r.URL)
	fmt.Println("Params:", r.URL.Query())
	fmt.Println("Request:", serviceRequest)
	if serviceRequest.Query == "" || (serviceRequest.Shardid == 0 && serviceRequest.Sharduid == "") || serviceRequest.UserToken == "" {
		fmt.Fprintln(w, "Invalid request")
		return
	}
	// get the query param from the get request
	query := serviceRequest.Query
	fmt.Println("Query:", query)
	if len(query) == 0 {
		fmt.Fprintln(w, "No query provided")
		return
	}
	// find the database on the map
	dbConnectionString, err := dbmapper.GetDbConnectionString(serviceRequest.Shardid, serviceRequest.Sharduid)
	if err != nil && !strings.Contains(err.Error(), "dbMap not found") {
		fmt.Println(err)
		return
	}
	fmt.Println("Database connection string:", dbConnectionString)
	if dbConnectionString == "" {
		fmt.Println("Database dsn not found in mapper. Creating new")
		dsn, err := createDatabase(&serviceRequest)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, "new database creation failed")
			return
		}
		dbConnectionString = dbmapper.DbConnectionString(dsn)
	}
	fmt.Println("Database connection string:", dbConnectionString)
	// forward the request to the database
	fmt.Println("Forwarding request to database")
	db, err := database.Connect(string(dbConnectionString))
	if err != nil {
		if strings.Contains(err.Error(), "failed to ping") {
			fmt.Println("Database creation not implemented")
			return
		} else {
			fmt.Println(err)
			return
		}
	}
	rows, err := database.Query(db, query)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	// write the response to the client
	fmt.Println("Writing response to client")
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintln(w, id, name)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func createDatabase(requestParams *service.Request) (string, error) {
	dbParams := database.DbTestConnectionParams
	dbParams.Dbname = ""
	dsn := database.BuildConnectionString(dbParams)
	db, err := database.Connect(dsn)
	if err != nil {
		db.Close()
		return "", err
	}
	if requestParams.Shardid != 0 {
		dbParams.Dbname = fmt.Sprintf("%s%d", database.DbTestConnectionParams.Dbname, requestParams.Shardid)
	} else {
		dbParams.Dbname = fmt.Sprintf("%s%s", database.DbTestConnectionParams.Dbname, requestParams.Sharduid)
	}
	if !isValidDatabaseName(dbParams.Dbname) {
		db.Close()
		return "", fmt.Errorf("invalid database name: %s", dbParams.Dbname)
	}
	fmt.Println("Creating database:", dbParams.Dbname)
	_, err = db.Exec("CREATE DATABASE " + dbParams.Dbname)
	if err != nil {
		db.Close()
		return "", err
	}
	db.Close()
	dsn = database.BuildConnectionString(dbParams)
	db, err = database.Connect(dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()
	userId, err := config.QueryUserIdFromDbConfig(requestParams.UserToken)
	if err != nil {
		return "", err
	}
	schemaStr, err := config.ReadSchemaFromDbConfig(userId)
	if err != nil {
		return "", err
	}
	err = database.CreateDatabaseFromSchema(db, schemaStr)
	if err != nil {
		return "", err
	}
	if requestParams.Shardid != 0 {
		dbmapper.AddDbMapId(requestParams.Shardid, dsn)
	} else {
		dbmapper.AddDbMapUid(requestParams.Sharduid, dsn)
	}
	return dsn, nil
}

// isValidDatabaseName checks if the database name contains only valid characters
func isValidDatabaseName(name string) bool {
	// Allow only alphanumeric characters and underscores
	var validName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return validName.MatchString(name)
}
