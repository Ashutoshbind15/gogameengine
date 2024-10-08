package main

import (
	"fmt"
	"net/http"

	"github.com/Ashutoshbind15/gogameengine/internal/data"
	"github.com/Ashutoshbind15/gogameengine/internal/handlers"
	"github.com/Ashutoshbind15/gogameengine/internal/realtime"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

var (
	Manager *realtime.GameManager
)

func initHandler (w http.ResponseWriter, rq *http.Request) {
	fmt.Fprintf(w, "init")
}

func mware (nxt http.HandlerFunc) http.HandlerFunc {
	hf := func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("running the mware")
        nxt.ServeHTTP(w, r)
    }

	return hf
}

func wsmware(wsHandler websocket.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsHandler.ServeHTTP(w, r)
	}
}

func WsHandler (conn *websocket.Conn) {
	defer func(){
		conn.Close()
	}()
}

func main() {
	
	fmt.Println("Game init")

	// extrascripts.TestScript()
	// data.InitTables()

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	data.InitDB()

	defer func(){
		data.DbConn.Close()
	}()

	gameManager := realtime.GameManager{
		Games: []realtime.Game{},
		Closegame: make(chan bool),
		Startgame: make(chan []byte),
		LoadedGamesMetadata: make(map[string]*realtime.GameMeta),
	}

	Manager = &gameManager

	r := mux.NewRouter()
	r.HandleFunc("/", initHandler)
	r.HandleFunc("/login", handlers.LoginEpHandler)
	r.HandleFunc("/protected", handlers.ProtectedEPHandler)

	wsh := websocket.Handler(WsHandler)

	http.Handle("/ws", mware(wsmware(wsh)))
	http.ListenAndServe(":3000", r)
}
