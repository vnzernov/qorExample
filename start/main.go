package main

import (
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qor/admin"
)

// Define a GORM-backend model
type User struct {
	gorm.Model
	Name string
}

// Define another GORM-backend model
type Product struct {
	gorm.Model
	Name        string
	Description string
}

func main() {
	// Set up the database
	DB, _ := gorm.Open("sqlite3", "demo.db")
	DB.AutoMigrate(&User{}, &Product{})

	// Initalize
	Admin := admin.New(&admin.AdminConfig{DB: DB})

	// Create resources from GORM-backend model
	Admin.AddResource(&User{})
	Admin.AddResource(&Product{})

	// Initalize an HTTP request multiplexer
	mux := http.NewServeMux()

	// Mount admin to the mux
	Admin.MountTo("/admin", mux)

	fmt.Println("Listening on: 9000")
	http.ListenAndServe(":9000", mux)
}
