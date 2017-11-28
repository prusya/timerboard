package main

import (
	"github.com/asdine/storm"
	"os"
	"github.com/gorilla/mux"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var db *storm.DB

type Config struct {
	Port            string `json:"port"`
	EveClientId     string `json:"eve_client_id"`
	EveClientSecret string `json:"eve_client_secret"`
	EveCallback     string `json:"eve_callback"`
}

var config Config

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

	r := mux.NewRouter()
	http.ListenAndServe(":"+config.Port, r)
}
