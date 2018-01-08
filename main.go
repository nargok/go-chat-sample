package main

import (
  "log"
  "net/http"
  "path/filepath"
  "text/template"
  "flag"
  "sync"
  "os"
  "trace"
  "github.com/stretchr/gomniauth"
  "github.com/stretchr/gomniauth/providers/google"
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

  // Gomniauthのセットアップ
  gomniauth.SetSecurityKey("YVjUz7iTXLfyPrZo2PCMx3q")
  gomniauth.WithProviders(
    google.New(GoogleClientId, GoogleSecretId, "http://localhost:8080/auth/callback/google"),
  )

  r := newRoom()
  r.tracer = trace.New(os.Stdout) // os.Stdout ログの出力先を標準出力にする

  http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
  http.Handle("/login", &templateHandler{filename: "login.html"})
  http.HandleFunc("/auth/", loginHandler)
  http.Handle("/room", r)

  // チャットルームを開始します
  go r.run()

  // webサーバを開始します
  log.Println("Webサーバを開始します。ポート: ", *address) // ログを画面に出力する
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
