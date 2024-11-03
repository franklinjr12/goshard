package servicelistener

// HTTP server to listen for service requests

import (
	"fmt"
	"goshard/internal/config"
	"goshard/internal/database"
	"goshard/internal/dbmapper"
	"goshard/lib/service"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ListenAndServe starts the HTTP server
func ListenAndServe() {
	http.HandleFunc("/query", queryHandler)
	http.HandleFunc("/schema", schemaHandler)
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

// will post schema to goshardconfig database
func schemaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("URL:", r.URL)
	fmt.Println("Params:", r.URL.Query())
	// accept only POST requests
	if r.Method != http.MethodPost {
		fmt.Fprintln(w, "Invalid request method. Use POST")
		return
	}
	userToken := r.URL.Query().Get("usertoken")
	if len(userToken) == 0 {
		fmt.Fprintln(w, "No user token provided")
		return
	}
	userId, err := config.QueryUserIdFromDbConfig(userToken)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "failed to query user id from db")
		return
	}
	fmt.Println("User id:", userId)
	// get schema from post body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "failed to read request body")
		return
	}
	schema := string(body)
	fmt.Println("Schema:", schema)
	exists, err := config.SchemaExists(userId)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "failed to check if schema exists")
		return
	}
	if exists {
		// update schema
		fmt.Println("Updating schema")
		err = config.UpdateSchemaInDbConfig(userId, schema)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, "failed to update schema")
			return
		}
	} else {
		// create new schema
		fmt.Println("Creating new schema")
		err = config.WriteSchemaToDbConfig(userId, schema)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, "failed to write schema")
			return
		}
	}
	fmt.Fprintln(w, "Schema written successfully")
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
	userToken := serviceRequest.UserToken
	fmt.Println("token:", userToken)
	if len(userToken) == 0 {
		fmt.Fprintln(w, "No user token provided")
		return
	}
	// get the query param from the get request
	query := serviceRequest.Query
	fmt.Println("Query:", query)
	if len(query) == 0 {
		fmt.Fprintln(w, "No query provided")
		return
	}
	userId, err := config.QueryUserIdFromDbConfig(userToken)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "failed to query user id from db")
		return
	}
	fmt.Println("User id:", userId)
	dbConnectionString, err := fetchConnectionString(userId, &serviceRequest)
	fmt.Println("Database connection string:", dbConnectionString)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, "failed to fetch connection string")
		return
	}
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
	userName, err := config.QueryUserNameFromDbConfig(requestParams.UserToken)
	if err != nil {
		db.Close()
		return "", err
	}
	if requestParams.Shardid != 0 {
		dbParams.Dbname = fmt.Sprintf("token%s%d", requestParams.UserToken+userName, requestParams.Shardid)
	} else {
		dbParams.Dbname = "token" + requestParams.UserToken + userName + requestParams.Sharduid
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

func fetchConnectionString(userId uint64, serviceRequest *service.Request) (dbmapper.DbConnectionString, error) {
	dbConnectionString, err := dbmapper.GetDbConnectionStringByUserId(userId, serviceRequest.Shardid, serviceRequest.Sharduid)
	if err != nil && !strings.Contains(err.Error(), dbmapper.DbMapNotFoundStr) {
		return "", err
	}
	if err != nil && strings.Contains(err.Error(), dbmapper.DbMapNotFoundStr) {
		fmt.Println("Database dsn not found in mapper. Creating new")
		dsn, err := createDatabase(serviceRequest)
		if err != nil {
			return "", err
		}
		fmt.Printf("New database created: %s\n", dsn)
		err = config.WriteNewMapping(userId, serviceRequest.Shardid, serviceRequest.Sharduid, dsn)
		if err != nil {
			return "", err
		}
		return dbmapper.DbConnectionString(dsn), nil
	}
	return dbConnectionString, nil
}
