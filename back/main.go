package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// 絵の構造体
type Picture struct {
	ID		string  `json:"id"`
	Title	string  `json:"title"`
	Path	string	`json:"path"`
}

// 絵の格納用配列
var pictures []Picture

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pictures)
}

// 特定の絵の取得
func getPicture(w http.ResponseWriter, r *http.Request) {
	// ヘッダをセット
	t, err := template.ParseFiles("template/index.html")

	// エラーの場合
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	params := mux.Vars(r)

	// 一致するものを取得（仮の処理）
	for _, item := range pictures {
		if item.ID == params["id"] {
			err = t.Execute(w, item)

			if err != nil {
				log.Printf("failed to execute template: %v", err)
			}
			return
		}
	}
	json.NewEncoder(w).Encode(&Picture{})

}

// Create a Book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var picture Picture
	_ = json.NewDecoder(r.Body).Decode(&picture)
	picture.ID = strconv.Itoa(rand.Intn(10000)) // Mock ID - not safe in production
	pictures = append(pictures, picture)
	json.NewEncoder(w).Encode(picture)
}

// Update a Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, item := range pictures {
		if item.ID == params["id"] {
			pictures = append(pictures[:index], pictures[index+1:]...)
			var picture Picture
			_ = json.NewDecoder(r.Body).Decode(&picture)
			picture.ID = params["id"]
			pictures = append(pictures, picture)
			json.NewEncoder(w).Encode(picture)
			return
		}
	}
	json.NewEncoder(w).Encode(pictures)
}

// Delete a Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, item := range pictures {
		if item.ID == params["id"] {
			pictures = append(pictures[:index], pictures[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(pictures)
}

func main() {
	// Initiate Router
	r := mux.NewRouter()

	// Mock Data
	pictures = append(pictures, Picture{ID: "1", Title: "Book one"})
	pictures = append(pictures, Picture{ID: "2", Title: "Book Two"})

	// Route Hnadlers / Endpoints
	r.HandleFunc("/", getBooks).Methods("GET")
	r.HandleFunc("/challenge/{id}", getPicture).Methods("GET")
	r.HandleFunc("/challenge", createBook).Methods("POST")
	r.HandleFunc("/gallery/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// ポート指定
	log.Fatal(http.ListenAndServe(":8000", r))
}