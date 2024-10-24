package servicelistener

// HTTP server to listen for service requests

import (
	"fmt"
	"goshard/internal/database"
	"goshard/internal/dbmapper"
	"goshard/lib/service"
	"net/http"
	"strconv"
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
	return service.Request{
		Query:    query,
		Shardid:  uint64(shardid),
		Sharduid: sharduid,
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	serviceRequest := parseUrlToRequest(r)
	// for debugging show the url and the params
	fmt.Println("URL:", r.URL)
	fmt.Println("Params:", r.URL.Query())
	fmt.Println("Request:", serviceRequest)
	if serviceRequest.Query == "" || serviceRequest.Shardid == 0 && serviceRequest.Sharduid == "" {
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
	if err != nil {
		fmt.Println(err)
		return
	}
	if dbConnectionString == "" {
		// need to create db
		// not implemented yet
		fmt.Println("Database creation not implemented")
		return
	}
	fmt.Println("Database connection string:", dbConnectionString)
	// forward the request to the database
	fmt.Println("Forwarding request to database")
	db, err := database.Connect(database.DbTestConnectionString)
	if err != nil {
		fmt.Println(err)
		return
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
