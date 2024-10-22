package servicelistener

// HTTP server to listen for service requests

import (
	"fmt"
	"goshard/internal/database"
	"net/http"
)

// ListenAndServe starts the HTTP server
func ListenAndServe() {
	http.HandleFunc("/query", queryHandler)
	http.ListenAndServe(":8080", nil)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	// for debugging show the url and the params
	fmt.Println("URL:", r.URL)
	fmt.Println("Params:", r.URL.Query())
	// get the query param from the get request
	query := r.URL.Query().Get("query")
	fmt.Println("Query:", query)
	if len(query) == 0 {
		fmt.Fprintln(w, "No query provided")
		return
	}
	// forward the request to the database
	fmt.Println("Forwarding request to database")
	db, err := database.Connect()
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
