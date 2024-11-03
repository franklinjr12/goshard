package main

// this will connect to the service running on localhost:8080 and send a query on the /query endpoint

import (
	"fmt"
	"goshard/internal/database"
	"goshard/lib/service"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	sendQuery()
	sendSchema()
}

func sendQuery() {
	request := service.Request{
		Query:     "SELECT id, name FROM users",
		Shardid:   3,
		Sharduid:  "",
		UserToken: "12345678",
	}
	queryUrlEncoded := url.PathEscape(request.Query)
	const url = "http://localhost:8080/query?"
	urlParams := fmt.Sprintf("query=%s&shardid=%d&sharduid=%s&usertoken=%s", queryUrlEncoded, request.Shardid, request.Sharduid, request.UserToken)
	fmt.Println("URL:", url+urlParams)
	// build the get request
	req, err := http.NewRequest("GET", url+urlParams, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	// read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func sendSchema() {
	// build the post request
	const url = "http://localhost:8080/schema?"
	userToken := "12345678"
	urlParams := fmt.Sprintf("usertoken=%s", userToken)
	fmt.Println("URL:", url+urlParams)
	schema, err := database.ReadSchemaFromFile("sql/schema.sql")
	if err != nil {
		fmt.Println(err)
		return
	}
	// send the request
	req, err := http.Post(url+urlParams, "text/plain", strings.NewReader(schema))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer req.Body.Close()
	// read the response
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
