package main



import (format "fmt")
import "net/http"
import "os"
import "html/template"
import "mime/multipart"
func main () {

	var mux *http.ServeMux;
	mux = http.NewServeMux();
	mux.HandleFunc("/hello", func (writer http.ResponseWriter, r *http.Request) {
		format.Fprintln(writer, "HandleFuncメソッドにServeHTTPを関数と同じシグネチャをそのまま渡す");
	})
	var hf http.HandlerFunc;
	hf = func (writer http.ResponseWriter, request *http.Request) {
		format.Fprintln(writer, "HandlerFunc型を定義->Handleメソッドに渡す");
	}
	mux.Handle("/hf", hf);
	mux . HandleFunc("/upload", upload);
	mux . HandleFunc("/", index);
	var myHandler = new(MyHandle);

	mux.Handle("/handle", myHandler);
	// http.Server構造体のポインタを宣言
	var server *http.Server;
	// http.Serverのオブジェクトを確保
	// &をつけること構造体ではなくポインタを返却
	server = &http.Server{}; // or new (http.Server);
	server.Addr = ":11180";
	server.Handler = mux;
	server.ListenAndServe();
}

type MyHandle struct {
	//とりあえず中身は空っぽ
}
func (this *MyHandle) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	format.Fprintln(writer, "http.Handler型インターフェースを実装したオブジェクトをhttp.Handleに渡す");
};


func index (writer http.ResponseWriter , request *http.Request) {
	var t *template.Template;
	/*
	<!-- テンプレートの中身 -->
	<!DOCTYPE html>
	<!-- template用のhtmlファイル -->
	<html>
	<head>
	    <title>ファイルアップロードテスト</title>
	</head>
	<body>
	<div class="container">
	    <h1>ファイルアップロードテスト</h1>
	    <form method="post" action="http://localhost:11180/upload" enctype="multipart/form-data">
	        <fieldset>
	            <input type="file" name="image" id="upload_files" multiple="multiple">
	            <input type="submit" name="submit" value="アップロード開始">
	        </fieldset>
	    </form>
	</div>
	</body>
	</html>
	*/
	// テンプレートをロード
	t, _ = template.ParseFiles("template/fileTest.html");
	t.Execute(writer, struct{}{});
}
func upload ( w http.ResponseWriter, r *http.Request) {
	// このハンドラ関数へのアクセスはPOSTメソッドのみ認める
	if  (r.Method != "POST") {
		format.Fprintln(w, "許可したメソッドとはことなります。");
		return;
	}
	var file multipart.File;
	var fileHeader *multipart.FileHeader;
	var e error;
	var uploadedFileName string;
	var img []byte = make([]byte, 1024);
	// POSTされたファイルデータを取得する
	file , fileHeader , e = r.FormFile ("image");
	if (e != nil) {
		format.Fprintln(w, "ファイルアップロードを確認できませんでした。");
		return;
	}
	uploadedFileName = fileHeader.Filename;
	// サーバー側に保存するために空ファイルを作成
	var saveImage *os.File;
	saveImage, e = os.Create("./" + uploadedFileName);
	if (e != nil) {
		format.Fprintln(w, "サーバ側でファイル確保できませんでした。");
		return;
	}
	defer saveImage.Close();
	defer file.Close();
	var tempLength int64 =0;
	for {
		// 何byte読み込んだかを取得
		n , e := file.Read(img);
		// 読み混んだバイト数が0を返したらループを抜ける
		if (n == 0) {
			format.Println(e);
			break;
		}
		if (e != nil) {
			format.Println(e);
			format.Fprintln(w, "アップロードされたファイルデータのコピーに失敗。");
			return;
		}
		saveImage.WriteAt(img, tempLength);
		tempLength = int64(n) + tempLength;
	}
	format.Fprintf(w, "文字列HTTPとして出力させる");
}