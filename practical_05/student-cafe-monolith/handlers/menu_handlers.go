package handlers

import (
	"encoding/json"
	"net/http"
	"student-cafe-monolith/database"
	"student-cafe-monolith/models"

	"github.com/go-chi/chi/v5"
)

func CreateMenuItem(w http.ResponseWriter, r *http.Request) {
    var menuItem models.MenuItem
    if err := json.NewDecoder(r.Body).Decode(&menuItem); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := database.DB.Create(&menuItem).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(menuItem)
}

func GetMenuItems(w http.ResponseWriter, r *http.Request) {
    var menuItems []models.MenuItem
    if err := database.DB.Find(&menuItems).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(menuItems)
}

func GetMenuItem(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    var menuItem models.MenuItem
    if err := database.DB.First(&menuItem, id).Error; err != nil {
        http.Error(w, "Menu item not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(menuItem)
}