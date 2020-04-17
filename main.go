package main

import (
	"dslife-backend/framework"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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

type FileUrl struct {
	FILENAME string `json:"filename"`
	URL      string `json:"fileurl"`
}

type TokenStr struct {
	TOKEN    string             `json:"token"`
	USERINFO framework.UserInfo `json:"userinfo"`
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
	userinfo, token, isCorrect := framework.GetTokenFromId(UserInfo.ID, UserInfo.PASSWORD)
	if isCorrect {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(TokenStr{TOKEN: token, USERINFO: userinfo}); err != nil {
			panic(err)
		}
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(404)
		if err := json.NewEncoder(w).Encode("Incorrect ID/PW"); err != nil {
			panic(err)
		}
		return
	}

}

type RegisterForm struct {
	Id       string `json:"id"`
	PASSWORD string `json:"password"`
	Name     string `json:"name"`
	Number   string `json:"number"`
	Room     int    `json:"room"`
	Type     string `json:"user_type"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var userinfo = RegisterForm{}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(reqBody, &userinfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	fmt.Println(userinfo)
	_, err = framework.CreateUser(userinfo.Id, userinfo.PASSWORD, userinfo.Name, userinfo.Number, userinfo.Room, userinfo.Type)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(503)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode("{ 'status': 'successful' }"); err != nil {
		panic(err)
	}
}

// PostNotice is Bad Server Fuck You
func PostNotice(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000)
	fileList := []FileUrl{}
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	author := r.Header.Values("Auth")
	formdata := r.MultipartForm
	files := formdata.File["multiplefiles"]
	content := formdata.Value["content"][0]
	title := formdata.Value["title"][0]
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		path := "./static/" + b64.StdEncoding.EncodeToString([]byte(content+title))
		os.Mkdir(path, 0666)
		out, err := os.OpenFile(path+"/"+files[i].Filename, os.O_WRONLY|os.O_CREATE, 0666)
		defer out.Close()
		fileList = append(fileList, FileUrl{FILENAME: files[i].Filename, URL: "http://localhost:8080/static/" + b64.StdEncoding.EncodeToString([]byte(content+title)) + "/" + files[i].Filename})
		if err != nil {
			fmt.Fprintf(w, "Fuck you man")
		}
		_, _ = io.Copy(out, file)
		_, _ = io.WriteString(w, "File Upload Successful!")

	}
	jsondata, err := json.Marshal(fileList)
	if err != nil {
		fmt.Fprintf(w, "Fuck you man")
	}
	framework.InsertPOST("board", author[0], "notice", title, content, string(jsondata))
	if err != nil {
		fmt.Fprintf(w, "DB ERROR")
		panic(err)
	}

}

// GetPost is
func GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.FormValue("id"))
	fmt.Printf(r.FormValue("id"))
	if err != nil {
		panic(err)
	}
	res := framework.GetPOST(id)
	if res == "" {

	} else {
		fmt.Fprintf(w, res)
	}

}
func main() {
	framework.InitDB()
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/login", login).Methods("POST")
	api.HandleFunc("/upload", PostNotice).Methods("POST")
	api.Path("/post").Queries("id", "{[0-9]*?}").HandlerFunc(GetPost).Methods("GET")
	api.HandleFunc("/register", Register).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	log.Fatal(http.ListenAndServe(":8080", router))
}
