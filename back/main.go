package main

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/appengine" // go get "google.golang.org/appengine"
)

var (
	storageClient *storage.Client

	// Set this in app.yaml when running in production.
	bucket = "secchallenge-aac82.appspot.com"
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
		ID: params["id"],
		Title: pictureData["title"].(string),
		Path: pictureData["path"].(string),
	})
}

// 作品情報の保存
func savePicture(w http.ResponseWriter, r *http.Request) {
	// 画像ファイルの保存
	path, err := uploadHandler(w, r)

	// エラー
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}

	// Firestoreへのデータ格納
	id, _, err := client.Collection("pictures").Add(ctx, map[string]interface{}{
		"title":	r.FormValue("title"),
		"path":		path,
		"createdAt":time.Now(),
	})

	// エラー
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(id.ID)
}

// 画像のアップロード
func uploadHandler (w http.ResponseWriter, r *http.Request) (string, error) {
	// Base64の画像データを復元して変数に格納
	base64img := strings.Replace(r.FormValue("file"), " ", "+", -1)

	// デコードしてバイナリに
	dec, err:= base64.StdEncoding.DecodeString(base64img)

	if err != nil {
		msg := fmt.Sprintf("ERROR : %v", err)
		return "", errors.New(msg)
	}

	image := bytes.NewReader(dec)
	ctx := appengine.NewContext(r)

	// 画像ファイル名の作成
	filePath := "pictures/" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".png"

	sw := storageClient.Bucket(bucket).Object(filePath).NewWriter(ctx)
	if _, err := io.Copy(sw, image);

	err != nil {
		msg := fmt.Sprintf("Could not write file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		println(msg)
		return "", errors.New(msg)
	}

	if err := sw.Close(); err != nil {
		msg := fmt.Sprintf("Could not put file: %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		println(msg)
		return "", errors.New(msg)
	}

	return "https://firebasestorage.googleapis.com/v0/b/" + bucket + "/o/" + url.QueryEscape(filePath) + "?alt=media", nil
}

// 絵の新規作成ページ
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

// メイン処理
func main() {
	// Firebaseの設定
	ctx = context.Background()
	sa := option.WithCredentialsFile("./key/secchallenge-aac82-firebase-adminsdk-du7lm-5dd831a3cb.json")
	app, err := firebase.NewApp(ctx, nil, sa)

	if err != nil {
		log.Fatalln(err)
	}

	// Firestoreのインスタンスを取得
	client, err = app.Firestore(ctx)
	storageClient, err = storage.NewClient(ctx, sa)

	if err != nil {
		log.Fatalln(err)
	}

	// 新規Routerの作成
	r := mux.NewRouter()

	// 静的ファイルのルーティング
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// トップページ（新規作成ページ）
	r.HandleFunc(	"/", drawPicture).Methods("GET")

	// 結果の表示
	r.HandleFunc("/gallery/{id}", getPicture).Methods("GET")

	// 結果の保存
	r.HandleFunc("/save", savePicture).Methods("POST")

	// ポート指定
	log.Fatal(http.ListenAndServe(":8000", r))

	// 切断
	defer client.Close()
}