package main

import (
	"fmt"
	"supermarket/internal/application"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// env
	// ...

	// app
	// - config
	// -- default store
	// app := application.NewApplicationDefault("", "./docs/db/json/products.json")
	// -- mysql store
	app := application.NewApplicationMySQL("", mysql.Config{
		User:      "user1",
		Passwd:    "secret_password",
		Net:       "tcp",
		Addr:      "localhost:3306",
		DBName:    "supermarket",
		ParseTime: true,
	})

	// - tear down
	defer app.TearDown()
	// - set up
	if err := app.SetUp(); err != nil {
		fmt.Println(err)
		return
	}
	// - run
	if err := app.Run(); err != nil {
		fmt.Println(err)
		return
	}
}
