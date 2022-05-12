package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cpotdevin/my-backend/pkg/websocket"
)

func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/sketch-pad", func(w http.ResponseWriter, r *http.Request) {
		serveWS(pool, w, r)
	})
	http.HandleFunc("/", helloWorldHandler)
}

func main() {
	fmt.Println("My Backend running")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
