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
	"mime/multipart"
	"net/http"
	"os"
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
var app			*firebase.App

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pictures)
}

// 特定の絵の情報取得
func getPicture(w http.ResponseWriter, r *http.Request) {
	// ヘッダをセット
	t, err := template.ParseFiles("template/gallery.html")

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

// ファイルの保存
func saveFile() ( w http.ResponseWriter, r *http.Request) {
	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if  (r.Method != "POST") {
		fmt.Fprintln(w, "許可したメソッドとはことなります。")
		return
	}
	var file multipart.File
	var fileHeader *multipart.FileHeader
	var e error
	var uploadedFileName string
	var img []byte = make([]byte, 1024)
	// POSTされたファイルデータを取得する
	file , fileHeader , e = r.FormFile ("image")
	if (e != nil) {
		fmt.Fprintln(w, "ファイルアップロードを確認できませんでした。")
		return
	}
	uploadedFileName = fileHeader.Filename
	// サーバー側に保存するために空ファイルを作成
	var saveImage *os.File
	saveImage, e = os.Create("./" + uploadedFileName)
	if (e != nil) {
		fmt.Fprintln(w, "サーバ側でファイル確保できませんでした。")
		return
	}
	defer saveImage.Close()
	defer file.Close()
	var tempLength int64 =0
	for {
		// 何byte読み込んだかを取得
		n , e := file.Read(img)
		// 読み混んだバイト数が0を返したらループを抜ける
		if (n == 0) {
			fmt.Println(e)
			break
		}
		if (e != nil) {
			fmt.Println(e)
			fmt.Fprintln(w, "アップロードされたファイルデータのコピーに失敗。")
			return
		}
		saveImage.WriteAt(img, tempLength)
		tempLength = int64(n) + tempLength
	}
	fmt.Fprintf(w, "文字列HTTPとして出力させる")
	return
}

func test (w http.ResponseWriter, r *http.Request) {
	// ヘッダをセット
	t, err := template.ParseFiles("template/test.html")

	// エラーの場合
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	// テンプレート
	err = t.Execute(w, nil)
}

func drawPicture (w http.ResponseWriter, r *http.Request) {

	// ヘッダをセット
	t, err := template.ParseFiles("template/index.html")

	// エラーの場合
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	// テンプレート
	err = t.Execute(w, nil)
}

func main() {
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


	// Initiate Router
	r := mux.NewRouter()

	// CSSの関連付け
	// r.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	// r.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	r.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))

	// Route Hnadlers / Endpoints
	r.HandleFunc("/", drawPicture).Methods("GET")

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