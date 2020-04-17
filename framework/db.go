package framework

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	// db driver
	"crypto/md5"
	"crypto/rand"

	guuid "github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

// InitDB db
func InitDB() {
	if d, err := sql.Open("sqlite3", "user.db"); err == nil {
		db = d
		rows, _ := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='user'")
		if !rows.Next() {
			db.Exec(`
			CREATE TABLE "user" (
				"i" INTEGER PRIMARY KEY AUTOINCREMENT,
				"unique_id" TEXT NOT NULL,
				"id" TEXT NOT NULL,
				"password" TEXT NOT NULL,
				"name" TEXT NOT NULL,
				"number" TEXT NOT NULL,
				"room" INTEGER NOT NULL,
				"type" TEXT NOT NULL
			)
			`)
		}
		rows.Close()

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='token'")
		if !rows.Next() {
			db.Exec(`
			CREATE TABLE "token" (
				"i" INTEGER PRIMARY KEY AUTOINCREMENT,
				"token" TEXT NOT NULL,
				"unique_id" TEXT NOT NULL,
				"id" TEXT NOT NULL,
				"name" TEXT NOT NULL,
				"number" TEXT NOT NULL,
				"room" INTEGER NOT NULL,
				"type" TEXT NOT NULL
			)
			`)
		}
		rows.Close()

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='point'")
		if !rows.Next() {
			db.Exec(
				`CREATE TABLE "point" (
					"i" INTEGER PRIMARY KEY AUTOINCREMENT,
					"unique_id" TEXT NOT NULL,
					"number" TEXT NOT NULL,
					"name" TEXT NOT NULL,
					"point" INTEGER NOT NULL,
					"reason" TEXT NOT NULL
				)
				`)
		}
		rows.Close()

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='meal'")
		if !rows.Next() {
			db.Exec(
				`CREATE TABLE "meal" (
					"i" INTEGER PRIMARY KEY AUTOINCREMENT,
					"date" DATE NOT NULL,
					"type" INTEGER NOT NULL,
					"menu" TEXT NOT NULL
				)
				`)
		}
		rows.Close()

		rows, err = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='board'")
		if !rows.Next() {
			_, err = db.Exec(
				`CREATE TABLE "board" (
					"i" INTEGER PRIMARY KEY AUTOINCREMENT,
					"author" TEXT NOT NULL,
					"type" TEXT NOT NULL,
					"title" TEXT NOT NULL,
					"content" TEXT NOT NULL,
					"file" TEXT NOT NULL
				)`)
			if err != nil {
				panic(err)
			}
		}
		if err != nil {
			fmt.Println(err)
		}
		rows.Close()

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='alarm'")
		if !rows.Next() {
			db.Exec(`CREATE TABLE "alarm" (
					"i" INTEGER PRIMARY KEY AUTOINCREMENT,
					"name" TEXT NOT NULL
				)`)
		}
		rows.Close()

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='declaration'")
		if !rows.Next() {
			fmt.Println("MAKE SUCCESSFUL!")
			db.Exec(`CREATE TABLE "declaration" (
					"i" INTEGER PRIMARY KEY AUTOINCREMENT,
					"author" TEXT NOT NULL,
					"title" TEXT NOT NULL,
					"content" TEXT NOT NULL
					)`)
		}
		rows.Close()
		fmt.Println("MAKE SUCCESSFUL!")
	} else {
		panic(err)
	}
}

// Insert is insert
func Insert(table string, args ...string) (sql.Result, error) {
	fmt.Println(strings.Join(args, ", "))
	res, err := db.Exec(fmt.Sprintf("INSERT INTO %s VALUE (, %s)", table, strings.Join(args, ", ")))
	return res, err
}

// InsertPOST
func InsertPOST(table string, author string, posttype string, title string, content string, jsondata string) {
	statement, err := db.Prepare("INSERT INTO board(author, type, title, content, file) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println("First Error")
		fmt.Println(err)
	}
	_, err = statement.Exec(author, posttype, title, content, jsondata)
	defer statement.Close()

	// return res, err
}

type FileUrl struct {
	FILENAME string `json:"filename"`
	URL      string `json:"fileurl"`
}

type PostInfo struct {
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	File    []FileUrl `json:"file"`
}

func GetPOST(id int) string {
	rows, err := db.Query(`SELECT author, title, content, file FROM board WHERE i=?`, id)
	if err != nil {
		panic(err)
	}
	if rows.Next() {
		var jsondata string
		var author string
		var title string
		var content string
		var files []FileUrl
		rows.Scan(&author, &title, &content, &jsondata)
		json.Unmarshal([]byte(jsondata), &files)
		var postinfo = PostInfo{Author: author, Title: title, Content: content, File: files}
		fmt.Println(postinfo)
		data, err := json.Marshal(postinfo)
		fmt.Println(string(data))
		if err != nil {
			panic(err)
		}
		return string(data)
	} else {
		return "nodata"
	}
}

type PostList struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func GetPostList(typ string, page int) string {

	rows, err := db.Query(`SELECT i, title FROM board WHERE type=? OFFSET ? LIMIT 10`, typ, 10*(page-1))
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	var postlist []PostList
	for rows.Next() {
		var title string
		var id string
		rows.Scan(&title, &id)
		idi, err := strconv.Atoi(id)
		if err != nil {
			panic(err)
		}
		postlist = append(postlist, PostList{Id: idi, Title: title})
	}
	if len(postlist) == 0 {
		return ""
	} else {
		data, err := json.Marshal(postlist)
		if err != nil {
			return "Unexpected Error"
		} else {
			return string(data)
		}
	}
}

func isVaildToken(token string) bool {
	rows, err := db.Query(`SELECT unique_id FROM token WHERE token=?`, token)
	defer rows.Close()
	if err != nil {
		return false
	} else {
		if !rows.Next() {
			return false
		} else {
			return true
		}
	}
}

type UserInfo struct {
	Uniqueid string `json:"unique_id"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Number   string `json:"number"`
	Room     int    `json:"room"`
	Type     string `json:"user_type"`
}

func GetTokenFromId(id string, passwd string) (UserInfo, string, bool) {
	hasher := md5.New()
	hasher.Write([]byte(passwd))
	hash := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println(string(hash))
	rows, err := db.Query(`SELECT unique_id, id, name, number, room, type FROM user WHERE id=? AND password=?`, id, string(hash))
	if err != nil {
		panic(err)
	} else {
		if !rows.Next() {
			fmt.Println("HALO")
			return UserInfo{}, "", false
		} else {
			var (
				uniqueid string
				id       string
				name     string
				number   string
				room     int
				types    string
			)
			rows.Scan(&uniqueid, &id, &name, &number, &room, &types)
			rows.Close()
			statement, err := db.Prepare("INSERT INTO token( token, unique_id, id, name, number, room, type ) VALUES (?, ?, ?, ?, ?, ?, ?)")

			if err != nil {
				panic(err)
			}
			token := GenerateToken()
			_, err = statement.Exec(token, uniqueid, id, name, number, room, types)
			statement.Close()
			fmt.Println(err)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			return UserInfo{Uniqueid: uniqueid, Id: id, Name: name, Number: number, Room: room, Type: types}, token, true
		}
	}
}

// GetInfoWithToken
func GetInfoWithToken(token string) (UserInfo, bool) {
	var (
		uniqueid string
		id       string
		name     string
		number   string
		room     int
		types    string
	)
	rows, err := db.Query(`SELECT unique_id, id, name, number, room, type FROM token WHERE token=?`, token)
	defer rows.Close()
	if err != nil {
		panic(err)
	} else {
		if !rows.Next() {
			rows.Scan(&uniqueid, &id, &name, &number, &room, &types)
			return UserInfo{Uniqueid: uniqueid, Id: id, Name: name, Number: number, Room: room, Type: types}, true
		} else {
			return UserInfo{}, false
		}
	}
}

// GenerateToken is fuck you
func GenerateToken() string {
	b := make([]byte, 20)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// CreateUser is CreateUser
func CreateUser(id string, passwd string, name string, number string, room int, types string) (sql.Result, error) {
	hasher := md5.New()
	hasher.Write([]byte(passwd))
	hash := hex.EncodeToString(hasher.Sum(nil))
	uniqueid := guuid.New().String()

	statement, err := db.Prepare("INSERT INTO user(unique_id, id, password, name, number, room, type) VALUES (?, ?, ?, ?, ?, ?, ?)")
	defer statement.Close()
	if err != nil {
		fmt.Println("First Error")
		fmt.Println(err)
		return nil, err
	}
	res, err := statement.Exec(uniqueid, id, string(hash), name, number, room, types)

	return res, err
}
