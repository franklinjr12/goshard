package main

// this will connect to the service running on localhost:8080 and send a query on the /query endpoint

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func fomatStrToUrl(str string) string {
	// replace the spaces with %20
	str = strings.Replace(str, " ", "%20", -1)
	// replace the commas with %2C
	str = strings.Replace(str, ",", "%2C", -1)
	// replace the quotes with %22
	str = strings.Replace(str, "\"", "%22", -1)
	// replace the brackets with %5B and %5D
	str = strings.Replace(str, "[", "%5B", -1)
	str = strings.Replace(str, "]", "%5D", -1)
	return str
}

func main() {
	query := fomatStrToUrl("SELECT id, name FROM users")
	const url = "http://localhost:8080/query?"
	urlParams := fmt.Sprintf("query=%s", query)
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
