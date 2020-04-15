package framework

import (
	"database/sql"
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
	if d, err := sql.Open("sqlite3", "bot.db"); err == nil {
		db = d
		rows, _ := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='user'")
		if !rows.Next() {
			db.Exec(`
			CREATE TABLE "user" (
				"i" INTEGER NOT NULL AUTOINCREMENT,
				"unique_id" TEXT NOT NULL,
				"id" TEXT NOT NULL,
				"password" TEXT NOT NULL,
				"name" TEXT NOT NULL,
				"number" TEXT NOT NULL,
				"room" INTEGER NOT NULL,
				"type" TEXT NOT NULL,
				PRIMARY KEY("i")
			)
			`)
		}

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='token'")
		if !rows.Next() {
			db.Exec(`
			CREATE TABLE "token" (
				"i" INTEGER NOT NULL AUTOINCREMENT,
				"token" TEXT NOT NULL,
				"unique_id" TEXT NOT NULL,
				"id" TEXT NOT NULL,
				"name" TEXT NOT NULL,
				"number" TEXT NOT NULL,
				"room" INTEGER NOT NULL,
				"type" TEXT NOT NULL,
				PRIMARY KEY("i")
			)
			`)
		}

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='point'")
		if !rows.Next() {
			db.Exec(
				`CREATE TABLE "point" (
					"i" INTEGER NOT NULL AUTOINCREMENT,
					"unique_id" TEXT NOT NULL,
					"number" TEXT NOT NULL,
					"name" TEXT NOT NULL,
					"point" INTEGER NOT NULL,
					"reason" TEXT NOT NULL,
					PRIMARY KEY("i")
				)
				`)
		}

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='meal'")
		if !rows.Next() {
			db.Exec(
				`CREATE TABLE "meal" (
					"i" INTEGER NOT NULL AUTOINCREMENT,
					"date" DATE NOT NULL,
					"type" INTEGER NOT NULL,
					"menu" TEXT NOT NULL,
					PRIMARY KEY("i")
				)
				`)
		}

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='board'")
		if !rows.Next() {
			db.Exec(
				`CREATE TABLE "board" (
					"i" INTEGER NOT NULL AUTOINCREMENT,
					"author" TEXT NOT NULL,
					"type" TEXT NOT NULL,
					"content" TEXT NOT NULL,
					"file" TEXT NOT NULL,
					PRIMARY KEY("i")
				)`)
		}

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='alarm'")
		if !rows.Next() {
			db.Exec(
				`CREATE TABLE "alarm" (
					"i" INTEGER NOT NULL AUTOINCREMENT,
					"name" TEXT NOT NULL,
					PRIMARY KEY("i")
				)`)
		}

		rows, _ = db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='declaration'")
		if !rows.Next() {
			db.Exec(
				`CREATE TABLE "declaration" (
					"i" INTEGER NOT NULL AUTOINCREMENT,
					"author" TEXT NOT NULL,
					"content" TEXT NOT NULL,
					PRIMARY KEY("i")
					)`)
		}

	}
}

func insert(table string, args ...string) (sql.Result, error) {
	res, err := db.Exec(fmt.Sprintf("INSERT INTO %s VALUE (, %s)", table, strings.Join(args, ", ")))
	return res, err
}
