package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/guruorgoru/newsguru/pkg/logs"
	"github.com/guruorgoru/newsguru/pkg/models"
	"github.com/jackc/pgx/v5"
)

func RootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		_, err := fmt.Fprintln(w, "Welcome to news Guru.")
		if err != nil {
			http.Error(w, "Internal Server Error Occured", http.StatusInternalServerError)
		}
	}
}
func GetNewsHandler(a *models.NewsModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		news, err := a.GetAll()
		if err != nil {
			http.Error(w, "Internal Server Error Occured", http.StatusInternalServerError)
			return
		}
		for _, news := range news {
			fmt.Fprintf(w, "%+v\n", news)
		}
	}
}
func GetNewsByIdHandler(a *models.NewsModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			logs.Error.Println("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			logs.Error.Println("ID is Invalid")
			http.Error(w, "Proper Integer ID is required", http.StatusBadRequest)
			return
		}
		news, err := a.GetByID(int64(id))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, logs.ErrorNotFound, http.StatusNotFound)
				return
			}
			http.Error(w, "Internal Server Error Occured", http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "%+v", news)
	}
}
func PostNewsHandler(a *models.NewsModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logs.Error.Println("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var news models.News
		err := json.NewDecoder(r.Body).Decode(&news)
		if err != nil {
			logs.Error.Println("JSON decode error:", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		newsID, err := a.Insert(news)
		if err != nil {
			logs.Error.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"id":      newsID,
			"message": "News created",
		})
	}
}
func DeleteNewsHandler(a *models.NewsModel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			logs.Error.Println("Proper ID is Required")
			http.Error(w, "Proper Integer ID is required", http.StatusBadRequest)
			return
		}
		if r.Method != http.MethodDelete {
			logs.Error.Println("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		err = a.Delete(id)
		if err != nil {
			if err == logs.SErrorNotFound {
				http.Error(w, logs.ErrorNotFound, http.StatusNotFound)
				return
			}
			http.Error(w, "Internal Server Error Occured", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
