package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	firebase "firebase.google.com/go"
)

// 絵の構造体
type Picture struct {
	id		string
	title	string
	path	string
}

// 絵の格納用配列
var pictures	[]Picture
var client		*firestore.Client
var ctx			context.Context

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

	// DBからデータを取得
	params := mux.Vars(r)
	dsnap, err := client.Collection("pictures").Doc(params["id"]).Get(ctx)

	// エラー時の処理
	if err != nil {
		fmt.Printf("err\n")
	}

	pictureData := dsnap.Data()

	fmt.Printf("value : %v\n", pictureData["title"].(string))
	fmt.Printf("type : %T\n", dsnap.Data())

	// テンプレート
	err = t.Execute(w, struct{
		Title	string
		Path	string
	}{
		Title: pictureData["title"].(string),
		Path: pictureData["path"].(string),
	})

	json.NewEncoder(w).Encode(&Picture{})

}

// Create a Book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var picture Picture
	_ = json.NewDecoder(r.Body).Decode(&picture)
	picture.id = strconv.Itoa(rand.Intn(10000)) // Mock ID - not safe in production
	pictures = append(pictures, picture)
	json.NewEncoder(w).Encode(picture)
}

// Update a Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	for index, item := range pictures {
		if item.id == params["id"] {
			pictures = append(pictures[:index], pictures[index+1:]...)
			var picture Picture
			_ = json.NewDecoder(r.Body).Decode(&picture)
			picture.id = params["id"]
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
		if item.id == params["id"] {
			pictures = append(pictures[:index], pictures[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(pictures)
}

func getDataTest (w http.ResponseWriter, r *http.Request) {
	dsnap, err := client.Collection("pictures").Doc("nUPs4UpO8K7ErrWICCJc").Get(ctx)
	if err != nil {
		fmt.Printf("err\n")
	}
	fmt.Printf("Document data: %#v\n", dsnap.Data())
}

func main() {
	// Initiate Router
	r := mux.NewRouter()

	// Firebaseの設定
	ctx = context.Background()
	sa := option.WithCredentialsFile("key/secchallenge-aac82-firebase-adminsdk-du7lm-5dd831a3cb.json")
	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Fatalln(err)
	}

	// Firestoreのインスタンスを取得
	client, err = app.Firestore(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	// Mock Data
	pictures = append(pictures, Picture{id: "1", title: "ドラえもん"})
	pictures = append(pictures, Picture{id: "2", title: "アンパンマン"})

	// Route Hnadlers / Endpoints
	r.HandleFunc("/", getBooks).Methods("GET")
	//
	r.HandleFunc("/challenge/{id}", getPicture).Methods("GET")
	// 結果の表示
	r.HandleFunc("/gallery/{id}", getPicture).Methods("GET")
	r.HandleFunc("/challenge", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// テスト
	r.HandleFunc("/test", getDataTest).Methods("GET")


	// ポート指定
	log.Fatal(http.ListenAndServe(":8000", r))

	// 切断
	defer client.Close()
}