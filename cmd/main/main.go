package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Ashutoshbind15/gogameengine/internal/handlers"
	"github.com/Ashutoshbind15/gogameengine/internal/realtime"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"golang.org/x/net/websocket"
)

type ctxkey struct{}


func initHandler (w http.ResponseWriter, rq *http.Request) {
	fmt.Fprintf(w, "init")
}

func mware (nxt http.HandlerFunc) http.HandlerFunc {

	hf := func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("running the mware")
		var uid = uuid.New()
		ctx := context.WithValue(r.Context(), ctxkey{}, uid)

        nxt.ServeHTTP(w, r.WithContext(ctx))
    }

	return hf
}

func wsmware(wsHandler websocket.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wsHandler.ServeHTTP(w, r)
	}
}

func uidfromwsrq (ws *websocket.Conn) uuid.UUID {
	return ws.Request().Context().Value(ctxkey{}).(uuid.UUID)
}


func WsHandler (conn *websocket.Conn) {

	// todo: get the uname from the db
	uid := uidfromwsrq(conn)
	uname := uid

	closeChannel := make(chan bool)

	clientInfo := realtime.NewClientInfo{
		UserId: uid,
		Username: uname.String(),
		Conn: conn,
		CloseChan: closeChannel,
	}

	realtime.Manager.AddClient <- clientInfo
	// todo: why do we need this in this ws implementation?
	<- closeChannel
}

func main() {
	
	fmt.Println("Game init")

	// extrascripts.TestScript()
	// data.InitTables()

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	// data.InitDB()

	// defer func(){
	// 	data.DbConn.Close()
	// }()

	gameManager := realtime.GameManager{
		Games: []*realtime.Game{},
		Closegame: make(chan []byte),
		Creategame: make(chan []byte),
		AddClient: make(chan realtime.NewClientInfo),
		Startgame: make(chan string),
		LoadedGamesMetadata: make(map[string]*realtime.GameMeta),
		ConnectionMeta : map[*websocket.Conn]*realtime.GameClient{},
		JoinGame: make(chan realtime.JoinRelay),
		RemoveClient: make(chan *websocket.Conn),
	}

	realtime.Manager = &gameManager
	go realtime.Manager.Manage()

	r := mux.NewRouter()
	r.HandleFunc("/", initHandler)
	r.HandleFunc("/login", handlers.LoginEpHandler)
	r.HandleFunc("/protected", handlers.ProtectedEPHandler)

	wsh := websocket.Handler(WsHandler)
	r.Handle("/ws", mware(wsmware(wsh)))

	handler := cors.Default().Handler(r)
	http.ListenAndServe(":3000", handler)
}
