package pkg

import (
	"database/sql"
	"encoding/json"
	"log"
	"myrest-api/pkg/model"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func GetNames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var people []model.Person

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
		var person model.Person
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
	var person model.Person

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
	var person model.Person
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

func UpdateName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var person model.Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	params := mux.Vars(r)
	person.ID = params["id"]

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE users SET name =? where id=?")
	checkError(err)
	defer stmt.Close()

	_, err = stmt.Exec(person.Name, person.ID)
	checkError(err)

	json.NewEncoder(w).Encode(person)
}

func DeleteName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM users where id=?")
	checkError(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkError(err)

	var answer model.Answer
	answer.Message = "Name doesn't exist or you've deleted it"

	json.NewEncoder(w).Encode(answer)
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
