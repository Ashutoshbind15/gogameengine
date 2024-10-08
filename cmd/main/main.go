package main

import (
	"fmt"
	"net/http"

	"github.com/Ashutoshbind15/gogameengine/internal/data"
	"github.com/Ashutoshbind15/gogameengine/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func initHandler (w http.ResponseWriter, rq *http.Request) {
	fmt.Fprintf(w, "init")
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

	r := mux.NewRouter()
	r.HandleFunc("/", initHandler)
	r.HandleFunc("/login", handlers.LoginEpHandler)
	r.HandleFunc("/protected", handlers.ProtectedEPHandler)
	http.ListenAndServe(":3000", r)
}
