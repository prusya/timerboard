package main

import (
	"github.com/asdine/storm"
	"os"
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/markbates/goth/gothic"
	"github.com/gorilla/sessions"
)

var db *storm.DB

type Config struct {
	Port            string `json:"port"`
	EveClientId     string `json:"eve_client_id"`
	EveClientSecret string `json:"eve_client_secret"`
	EveCallback     string `json:"eve_callback"`
	SessionsKey     string `json:"sessions_key"`
	GothicKey       string `json:"gothic_key"`
}

var config Config
var cookieStore *sessions.CookieStore
var gothicStore *sessions.FilesystemStore

func main() {
	db, err := storm.Open("my.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//In case of cmd args handle them and exit
	if len(os.Args) > 1 {
		utilsHandleOptions()
		return
	}

	//Read config
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(raw, &config)
	if err != nil {
		panic(err)
	}
	cookieStore = sessions.NewCookieStore([]byte(config.SessionsKey))
	gothicStore = sessions.NewFilesystemStore(".", []byte(config.GothicKey))

	r := mux.NewRouter()
	r.HandleFunc("/", GetIndexHandler).Methods("GET")

	r.HandleFunc("/auth/{provider}", gothic.BeginAuthHandler)
	r.HandleFunc("/eve_callback", GetEveCallbackHandler).Methods("GET")
	r.HandleFunc("/logout", GetLogoutHandler).Methods("GET")

	r.HandleFunc("/users", GetUsersHandler).Methods("GET")
	r.HandleFunc("/users", PostUsersHandler).Methods("POST")

	r.HandleFunc("/timers", GetTimersHandler).Methods("GET")
	r.HandleFunc("/timers", PostTimersHandler).Methods("POST")
	r.HandleFunc("/timers/{id}", DeleteTimersHandler).Methods("DELETE")

	http.ListenAndServe(":"+config.Port, r)
}
