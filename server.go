package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/nvdarekar/go-hello-world/db"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("input")
	if input != "" {
		insert, err := db.DBConn.Prepare("INSERT INTO inputs(input) VALUES(?)")
		if err != nil {
			log.Fatal(err)
		}
		insert.Exec(input)
		fmt.Fprintf(w, "ok")
	}
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DBConn.Query("SELECT input FROM inputs")
	if err != nil {
		log.Fatal(err)
	}
	var inputs []string
	for rows.Next() {
		var input string
		err = rows.Scan(&input)
		if err != nil {
			panic(err.Error())
		}
		inputs = append(inputs, input)
	}
	js, err := json.Marshal(inputs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func dbConnect() (db *sql.DB) {
	dbConnectionStr := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER_NAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))
	db, err := sql.Open("mysql", dbConnectionStr)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	// load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.DBConn = dbConnect()
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/retrieve", retrieveHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
