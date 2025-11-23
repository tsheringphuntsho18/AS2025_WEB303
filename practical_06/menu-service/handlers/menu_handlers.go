package handlers

import (
	"encoding/json"
	"net/http"
	"menu-service/database"
	"menu-service/models"

	"github.com/go-chi/chi/v5"
)

func GetMenu(w http.ResponseWriter, r *http.Request) {
	var menu models.Menu
	if err := database.DB.First(&menu, "id = ?", chi.URLParam(r, "id")).Error; err != nil {
		http.Error(w, "Menu not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func CreateMenu(w http.ResponseWriter, r *http.Request) {
	var menu models.Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.DB.Create(&menu).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(menu)
}