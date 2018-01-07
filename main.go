package main

import (
  "log"
  "net/http"
  "path/filepath"
  "text/template"
  "flag"
  "sync"
)

// temp1は1つのテンプレートを表します
type templateHandler struct {
  once sync.Once
  filename string
  templ  *template.Template
}

// ServerHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.once.Do(func() {
    t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
  })
  t.templ.Execute(w, r)
}

func main() {
  var address = flag.String("address", ":8080", "アプリケーションのアドレス")
  flag.Parse() // フラグを解釈します

  r := newRoom()

  http.Handle("/", &templateHandler{filename: "chat.html"})
  http.Handle("/room", r)

  // チャットルームを開始します
  go r.run()

  // webサーバを開始します
  log.Println("Webサーバを開始します。ポート: ", *address) // ログを画面に出力する
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
