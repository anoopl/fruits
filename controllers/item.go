package controllers

import (
	"database/sql"
	"encoding/json"
	"items/models"
	itemRepository "items/repository/item"
	"items/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Controller struct{}

var items []models.Item

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c Controller) GetItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item models.Item
		var error models.Error

		items = []models.Item{}
		itemRepo := itemRepository.ItemRepository{}
		items, err := itemRepo.GetItems(db, item, items)

		if err != nil {
			error.Message = "Server Error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, items)

	}
}

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item models.Item
		var error models.Error

		params := mux.Vars(r)

		items = []models.Item{}
		itemRepo := itemRepository.ItemRepository{}

		id, _ := strconv.Atoi(params["id"])

		item, err := itemRepo.GetItem(db, item, id)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "Not Found"
				utils.SendError(w, http.StatusNotFound, error)
				return
			} else {
				error.Message = "Server Error"
				utils.SendError(w, http.StatusInternalServerError, error)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, item)

	}
}

func (c Controller) AddItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item models.Item
		var itemID int
		var error models.Error

		json.NewDecoder(r.Body).Decode(&item)

		if item.Owner == "" || item.Title == "" || item.Year == "" {
			error.Message = "Enter missing fields."
			utils.SendError(w, http.StatusBadRequest, error) //400
			return
		}

		itemRepo := itemRepository.ItemRepository{}
		itemID, err := itemRepo.AddItem(db, item)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, itemID)
	}
}

func (c Controller) UpdateItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item models.Item
		var error models.Error

		json.NewDecoder(r.Body).Decode(&item)

		if item.ID == 0 || item.Owner == "" || item.Title == "" || item.Year == "" {
			error.Message = "All fields are required."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		itemRepo := itemRepository.ItemRepository{}
		rowsUpdated, err := itemRepo.UpdateItem(db, item)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsUpdated)
	}
}

func (c Controller) RemoveItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		params := mux.Vars(r)
		itemRepo := itemRepository.ItemRepository{}
		id, _ := strconv.Atoi(params["id"])

		rowsDeleted, err := itemRepo.RemoveItem(db, id)

		if err != nil {
			error.Message = "Server error."
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		if rowsDeleted == 0 {
			error.Message = "Not Found"
			utils.SendError(w, http.StatusNotFound, error) //404
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsDeleted)
	}
}
