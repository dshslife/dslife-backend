package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type user struct {
	TOKEN string `json:"TOKEN"`
}

// LoginUser is LoginUser
type LoginUser struct {
	ID       string `json:"id"`
	PASSWORD string `json:"password"`
}

func login(w http.ResponseWriter, r *http.Request) {
	var UserInfo LoginUser
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(reqBody, &UserInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	var usertoken user
	usertoken.TOKEN = "aksdljaslkdjalskdjalsk123"
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(usertoken); err != nil {
		panic(err)
	}

}
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/api/login", login).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
