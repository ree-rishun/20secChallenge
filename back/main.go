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
	"time"

	firebase "firebase.google.com/go"
)

// 絵の構造体
type Picture struct {
	ID		string
	Title	string
	Path	string
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

// 特定の絵の情報取得
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

	// テンプレート
	err = t.Execute(w, Picture{
		ID: pictureData["id"].(string),
		Title: pictureData["title"].(string),
		Path: pictureData["path"].(string),
	})
}

// 作品の保存
func savePicture(w http.ResponseWriter, r *http.Request) {
	// Firestoreへのデータ格納
	_, _, err := client.Collection("pictures").Add(ctx, map[string]interface{}{
		"title":	r.FormValue("title"),
		"path":		r.FormValue("path"),
		"createdAt":time.Now(),
	})

	// エラー
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
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

func test(w http.ResponseWriter, r *http.Request) {
	// ヘッダをセット
	t, err := template.ParseFiles("template/test.html")

	// エラーの場合
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	// テンプレート
	err = t.Execute(w, "")
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

	// Route Hnadlers / Endpoints
	r.HandleFunc("/", getBooks).Methods("GET")
	r.HandleFunc("/challenge", createBook).Methods("POST")

	// 結果の表示
	r.HandleFunc("/gallery/{id}", getPicture).Methods("GET")

	// 結果の保存
	r.HandleFunc("/save", savePicture).Methods("POST")

	// 結果の保存
	r.HandleFunc("/test", test).Methods("GET")


	// ポート指定
	log.Fatal(http.ListenAndServe(":8000", r))

	// 切断
	defer client.Close()
}