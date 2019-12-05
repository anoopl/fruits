package main

import (
	"database/sql"
	"fmt"
	"items/controllers"
	"items/driver"
	"items/models"
	"log"
	"net/http"

	"github.com/subosito/gotenv"

	"github.com/gorilla/mux"
)

var items []models.Item
var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {
	db = driver.ConnectDB()
	controller := controllers.Controller{}

	router := mux.NewRouter()

	router.HandleFunc("/items", controller.GetItems(db)).Methods("GET")
	router.HandleFunc("/items/{id}", controller.GetItem(db)).Methods("GET")
	router.HandleFunc("/items", controller.AddItem(db)).Methods("POST")
	router.HandleFunc("/items", controller.UpdateItem(db)).Methods("PUT")
	router.HandleFunc("/items/{id}", controller.RemoveItem(db)).Methods("DELETE")

	fmt.Println("Server is running at port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
