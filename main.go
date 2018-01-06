package main

import (
  "log"
  "net/http"
  "path/filepath"
  "text/template"
)

// temp1は1つのテンプレートを表します
type templ struct {
  source string
  templ  *template.Template
}

// ServerHTTPはHTTPリクエストを処理します
func (t *templ) Handle(w http.ResponseWriter, r *http.Request) {
  if t.templ == nil {
    t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.source)))
  }
  t.templ.Execute(w, nil)
}

func main() {
  http.HandleFunc("/", (&templ{source: "chat.html"}).Handle)
    // webサーバを開始します
    if err := http.ListenAndServe(":8080", nil); err != nil {
      log.Fatal("ListenAndServe:", err)
    }
}
