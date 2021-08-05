package pkg

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("select * from users")
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		/*rows.Next() что-то вроде foreach*/
		var person Person
		err := rows.Scan(&person.ID, &person.Name)
		/*rows.Scan() записывает из строки БД в переменные*/
		checkError(err)
		people = append(people, person)
		//log.Println(person.ID, person.Name)
	}
	err = rows.Err()
	checkError(err)

	json.NewEncoder(w).Encode(people)
}

func GetName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var person Person

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("select id, name from users where id = ?")
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(params["id"])
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		/*rows.Next() что-то вроде foreach*/
		err := rows.Scan(&person.ID, &person.Name)
		/*rows.Scan() записывает из строки БД в переменные*/
		checkError(err)
		//log.Println(person.ID, person.Name)
	}
	err = rows.Err()
	checkError(err)

	json.NewEncoder(w).Encode(person)
}

func CreateName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users(name) VALUES (?)")
	checkError(err)
	defer stmt.Close()

	res, err := stmt.Exec(person.Name)
	checkError(err)

	lastId, err := res.LastInsertId()
	checkError(err)

	rowCnt, err := res.RowsAffected()
	checkError(err)

	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	json.NewEncoder(w).Encode(person)
}

func connectToDB() *sql.DB {
	db, err := sql.Open("mysql",
		"root:root@tcp(localhost:6033)/app_db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
