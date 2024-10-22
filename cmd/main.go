package main

import (
	"fmt"
	_ "goshard/internal/database"
	"goshard/internal/servicelistener"
)

func main() {
	fmt.Println("Running")
	// db, err := database.Connect()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// rows, err := db.Query("SELECT id, name FROM users")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer rows.Close()
	// fmt.Println("Showing users rows")
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err = rows.Scan(&id, &name)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return
	// 	}
	// 	fmt.Println(id, name)
	// }
	// err = rows.Err()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	servicelistener.ListenAndServe()
}
