package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Ashutoshbind15/gogameengine/internal/data"
	"github.com/Ashutoshbind15/gogameengine/internal/extrascripts"
	"github.com/Ashutoshbind15/gogameengine/internal/scriptingmappers"
	"github.com/Ashutoshbind15/gogameengine/internal/types"
	"github.com/Ashutoshbind15/gogameengine/internal/utils"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	lua "github.com/yuin/gopher-lua"
)


func test() {
	L := lua.NewState()
	defer L.Close()

	fratk := extrascripts.FrontalAttackScript
	
	// Sample GameState
	gameState := &types.GameState{
		Enemies: []types.Enemy{
			{Health: 100, Px: 2, Py: 3},
			{Health: 80, Px: 2, Py: 4},
		},
		Time:     120,
		TimeLeft: 60,
		Arena:    types.Arena{X: 100, Y: 200},
		Players: []types.Player{
			{Name: "Player1", Px: 2, Py: 3, Dynamics: map[string]string{"attack": fratk}},
		},
	}

	modifiedGameState := scriptingmappers.GameStateScriptRunner(L, fratk, "attack", gameState)

	fmt.Printf("Modified GameState: %+v\n", modifiedGameState)

}

func initHandler (w http.ResponseWriter, rq *http.Request) {
	fmt.Fprintf(w, "init")
}

func protectedEPHandler (w http.ResponseWriter, rq *http.Request) {
	if rq.Method == "GET" {
		cookies := rq.Cookies();

		for _, cook:= range cookies {
			if(cook.Name == "sid") {
				// parse, unsign and then check against the db session
				sid := cook.Value
				
				var sess types.Session
				row := data.DbConn.QueryRow(`SELECT * FROM DBSESSIONS WHERE id = $1`, sid)
				err := row.Scan(&sess.Id, &sess.UserId, &sess.ValidTo)

				// todo: handle the errors better, not just print and return

				if err != nil {
					fmt.Println("err: ", err)
					return
				}

				var user types.User

				row = data.DbConn.QueryRow(`SELECT * FROM USERS WHERE id = $1`, sess.UserId)
				err = row.Scan(&user.Id, &user.Username, &user.Password)

				if err != nil {
					fmt.Println("err: ", err)
					return
				}

				fmt.Println("sess: ", sess)
				fmt.Println("user: ", user)
				

			}
		}
	
		fmt.Fprintf(w, "protected route")
	}
}

func loginEpHandler (w http.ResponseWriter, rq *http.Request) {

	if rq.Method == "POST" {
		var loginIp types.LoginIp;

		body := rq.Body
	
		data, err := io.ReadAll(body)
	
		if err != nil {
			fmt.Println("err: ", err)
		}
	
		json.Unmarshal(data, &loginIp)
	
		uname, pwd := loginIp.Username, loginIp.Password
	
		if uname == "test" && pwd == "test" {
			// set the set-cookie header here after signing and possible enc the cookie, set the db session as well
	
			sid := utils.GenerateSessionId()
	
			// manually set the cookie
			cookieHeaderVal := fmt.Sprintf("%s=%s", "sid", sid)
			fmt.Println(cookieHeaderVal)
	
			// use the util
	
			// todo: modify the cookie opts acc to the env, dev or prod
	
			cookie := http.Cookie{
				Name:     "sid",
				Value:    sid,
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
				Domain: "",
			}

			utils.CreateDBSession(1, sid)
			http.SetCookie(w, &cookie)
	
			fmt.Fprintf(w, "ok")
		}	
		// err out the resp as 401
	}
}

func main() {
	
	fmt.Println("Game init")

	// test()

	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	data.InitDB()

	defer func(){
		data.DbConn.Close()
	}()

	// data.InitTables()

	r := mux.NewRouter()
	r.HandleFunc("/", initHandler)
	r.HandleFunc("/login", loginEpHandler)
	r.HandleFunc("/protected", protectedEPHandler)
	http.ListenAndServe(":3000", r)
}
