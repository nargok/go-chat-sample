package main

import (
  "log"
  "net/http"

  "github.com/gorilla/websocket"
)

type room struct {
  // forwardは他のクライアントに転送するためのメッセージを保持するチャンネルです。
  forward chan []byte
  // joinはチャットルームに参加しようとしているクライアントのためのチャンネルです。
  join chan *client
  // leaveはチャットルームから退室しようとしているクライアントのためのチャンネルです
  leave chan *client
  // clientsには在室しているすべてのクライアントが保持されます。
  clients map[*client]bool

}

// newRoomはすぐに利用できるチャットルームを生成して返します。
func newRoom() *room {
  return &room {
    forward: make(chan []byte),
    join: make(chan *client),
    leave: make(chan *client),
    clients: make(map[*client]bool),
  }
}

func (r *room) run() {
  for {
    select {
      case client := <- r.join:
        // 参加
        r.clients[client] = true
      case client := <- r.leave:
        // 退室
        delete(r.clients, client)
        close(client.send)

    case msg := <- r.forward:
      // すべてのクライアントにメッセージを送信
      for client := range r.clients {
        select {
        case client.send <- msg:
          // メッセージ送信
        default:
          // 送信に失敗
          delete(r.clients, client)
          close(client.send)
        }
      }
    }
  }
}

const (
  socketBufferSize = 1024
  messageBufferSize = 256
)

var upgrader = &websocket.Upgrader { ReadBufferSize:
    socketBufferSize, WriteBufferSize: socketBufferSize }
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  socket, err := upgrader.Upgrade(w, req, nil)  // Websocketコネクションを取得する
  if err != nil {
    log.Fatal("ServeHTTP:", err)
    return
  }
  client := &client {  // コネクション取得成功の場合、clientを生成する
    socket: socket,
    send: make(chan []byte, messageBufferSize),
    room: r,
  }
  r.join <- client     //  joinチャンネルにclientを渡す
  defer func() { r.leave <- client } () // clientの終了時に退出の処理を行うように指定する
  go client.write() // go keywordに続けて書くと、goroutineとなる。
  client.read()     // メインスレッドで呼ぶ。接続が保持され終了するまで他の処理はブロックする。
}