package pkg

import (
	"encoding/json"
	"log"
	"myrest-api/pkg/model"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func Register(w http.ResponseWriter, r *http.Request) {

	var person model.Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users(login, password, name, isAdmin) VALUES (?,?,?,?)")
	checkError(err)
	defer stmt.Close()

	res, err := stmt.Exec(person.Login, person.Password, person.Name, person.IsAdmin)
	checkError(err)

	lastId, err := res.LastInsertId()
	checkError(err)

	log.Printf("User with ID = %d was added", lastId)

	var answer model.Answer
	answer.Message = "You registered. Your user ID = %d"
	json.NewEncoder(w).Encode(answer)
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var personIn model.Person
	json.NewDecoder(r.Body).Decode(&personIn)

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("select id, login, password, name, isAdmin from users where login = ?")
	checkError(err)
	defer stmt.Close()

	var personOut model.Person

	rows, err := stmt.Query(personIn.Login)
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		/*rows.Next() что-то вроде foreach*/
		err := rows.Scan(&personOut.ID, &personOut.Login, &personOut.Password, &personOut.Name, &personOut.IsAdmin)
		/*rows.Scan() записывает из строки БД в переменные*/
		checkError(err)
		//log.Println(person.ID, person.Name)
	}

	if (personIn.Password == personOut.Password) && (personIn.Login == personOut.Login) && (personIn.Password != "") && (personIn.Login != "") {
		cookie := &http.Cookie{
			Name:   "token",
			Value:  "some_token_value",
			MaxAge: 300,
		}
		http.SetCookie(w, cookie)
		var answer model.Answer
		answer.Message = "You log in"
		json.NewEncoder(w).Encode(answer)
	} else {
		var answer model.Answer
		answer.Message = "Wrong auth data"
		json.NewEncoder(w).Encode(answer)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c := http.Cookie{
		Name:   "token",
		MaxAge: -1,
	}
	http.SetCookie(w, &c)

	var answer model.Answer
	answer.Message = "You are logging out now"
	json.NewEncoder(w).Encode(answer)
}
