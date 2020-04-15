package framework

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	// db driver
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
