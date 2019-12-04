package main

import (
	"database/sql"
	"encoding/json"
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

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

func getItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	params := mux.Vars(r)

	rows := db.QueryRow("select * from items where id=$1", params["id"])

	err := rows.Scan(&item.ID, &item.Title, &item.Owner, &item.Year)
	logFatal(err)

	json.NewEncoder(w).Encode(item)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	var itemID int

	json.NewDecoder(r.Body).Decode(&item)

	err := db.QueryRow("insert into items (title, owner, year) values ($1, $2, $3) RETURNING id;",
		item.Title, item.Owner, item.Year).Scan(&itemID)

	logFatal(err)

	json.NewEncoder(w).Encode(itemID)

}

func updateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	json.NewDecoder(r.Body).Decode(&item)
	result, err := db.Exec("update items set title=$1, owner=$2, year=$3 where id=$4 RETURNING id",
		&item.Title, &item.Owner, &item.Year, &item.ID)
	rowsUpdated, err := result.RowsAffected()
	logFatal(err)
	json.NewEncoder(w).Encode(rowsUpdated)

}

func removeItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result, err := db.Exec("delete from items where id = $1", params["id"])
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
