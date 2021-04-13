package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	host     = "172.17.0.2"
	port     = 5432
	user     = "postgres"
	password = "//*p05tgr355//*"
	dbname   = "postgres"
)

type Person struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

func main() {
	http.HandleFunc("/", getHandler)
	http.HandleFunc("/insert", postHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	db := openConn()

	rows, err := db.Query("SELECT * FROM PERSON")
	if err != nil {
		log.Fatal(err)
	}

	var people []Person

	for rows.Next() {
		var person Person
		rows.Scan(&person.Name, &person.Nickname)
		people = append(people, person)
	}

	peopleBytes, _ := json.MarshalIndent(people, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
	defer db.Close()
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	db := openConn()

	var p Person
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO person (name, nickname) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, p.Name, p.Nickname)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func openConn() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		"password=%s dbname=%s sslmode=disable", host, port, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
