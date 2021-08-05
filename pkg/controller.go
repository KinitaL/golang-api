package pkg

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetNames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var people []Person

	db, err := sql.Open("mysql",
		"root:root@tcp(localhost:6033)/app_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select * from users")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		/*rows.Next() что-то вроде foreach*/
		var person Person
		err := rows.Scan(&person.ID, &person.Name)
		/*rows.Scan() записывает из строки БД в переменные*/
		if err != nil {
			log.Fatal(err)
		}
		people = append(people, person)
		//log.Println(person.ID, person.Name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(people)
}

func GetName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var person Person
	db, err := sql.Open("mysql",
		"root:root@tcp(localhost:6033)/app_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select id, name from users where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(params["id"])
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		/*rows.Next() что-то вроде foreach*/
		err := rows.Scan(&person.ID, &person.Name)
		/*rows.Scan() записывает из строки БД в переменные*/
		if err != nil {
			log.Fatal(err)
		}
		//log.Println(person.ID, person.Name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(person)
}

func CreateName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = strconv.Itoa(rand.Intn(1000))

	db, err := sql.Open("mysql",
		"root:root@tcp(localhost:6033)/app_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO users(id,name) VALUES (?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(person.ID, person.Name)
	if err != nil {
		log.Fatal(err)
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	json.NewEncoder(w).Encode(person)
}
