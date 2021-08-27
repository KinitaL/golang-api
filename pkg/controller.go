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

	var stocks []model.Stock

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("select * from stocks")
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		/*rows.Next() что-то вроде foreach*/
		var stock model.Stock
		err := rows.Scan(&stock.ID, &stock.Fullname, &stock.Shortname, &stock.Price)
		/*rows.Scan() записывает из строки БД в переменные*/
		checkError(err)
		stocks = append(stocks, stock)
		//log.Println(person.ID, person.Name)
	}
	err = rows.Err()
	checkError(err)

	json.NewEncoder(w).Encode(stocks)
}

func GetName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var stock model.Stock

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("select id, fullname, shortname, price from stocks where id = ?")
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query(params["id"])
	checkError(err)
	defer rows.Close()

	for rows.Next() {
		/*rows.Next() что-то вроде foreach*/
		err := rows.Scan(&stock.ID, &stock.Fullname, &stock.Shortname, &stock.Price)
		/*rows.Scan() записывает из строки БД в переменные*/
		checkError(err)
		//log.Println(person.ID, person.Name)
	}
	err = rows.Err()
	checkError(err)

	json.NewEncoder(w).Encode(stock)
}

func CreateName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var stock model.Stock
	_ = json.NewDecoder(r.Body).Decode(&stock)

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO stocks(fullname, shortname, price) VALUES (?, ?, ?)")
	checkError(err)
	defer stmt.Close()

	res, err := stmt.Exec(stock.Fullname, stock.Shortname, stock.Price)
	checkError(err)

	lastId, err := res.LastInsertId()
	checkError(err)

	rowCnt, err := res.RowsAffected()
	checkError(err)

	log.Printf("ID = %d, affected = %d\n", lastId, rowCnt)

	json.NewEncoder(w).Encode(stock)
}

func UpdateName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var stock model.Stock
	_ = json.NewDecoder(r.Body).Decode(&stock)

	params := mux.Vars(r)
	stock.ID = params["id"]

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("UPDATE stocks SET fullname=?, shortname=?, price=? where id=?")
	checkError(err)
	defer stmt.Close()

	_, err = stmt.Exec(stock.Fullname, stock.Shortname, stock.Price, stock.ID)
	checkError(err)

	json.NewEncoder(w).Encode(stock)
}

func DeleteName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	db := connectToDB()
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM stocks where id=?")
	checkError(err)
	defer stmt.Close()

	_, err = stmt.Exec(id)
	checkError(err)

	var answer model.Answer
	answer.Message = "Stock doesn't exist or you've deleted it"

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
