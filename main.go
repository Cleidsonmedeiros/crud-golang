package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type People struct {
	Id    int
	Name  string
	Phone int
	Cpf int
}

type Address struct {
	Id int
	Street string
	Cep string
	Number int
	PeopleId int
}

func dbConn() (db *sql.DB) {

	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "Matohehe123a.2"
	dbName := "mydb" 

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("tmpl/*"))

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	selDB, err := db.Query("SELECT * FROM people ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	n := People{}

	res := []People{}

	for selDB.Next() {
		var id int
		var name string
		var phone int
		var cpf int

		err = selDB.Scan(&id, &name, &phone, &cpf)
		if err != nil {
			panic(err.Error())
		}

		n.Id = id
		n.Name = name
		n.Phone = phone
		n.Cpf = cpf

		res = append(res, n)

	}

	tmpl.ExecuteTemplate(w, "Index", res)

	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM address WHERE people_id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	n := Address{}

	for selDB.Next() {
		var id int
		var street string
		var cep string
		var number int
		var people_id int

		err = selDB.Scan(&id, &street, &cep, &number, &people_id)
		if err != nil {
			panic(err.Error())
		}

		n.Id = id
		n.Street = street
		n.Cep = cep
		n.Number = number
	}

	tmpl.ExecuteTemplate(w, "Show", n)

	defer db.Close()

}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func InsertAddress(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {

		street := r.FormValue("street")
		cep := r.FormValue("cep")
		number := r.FormValue("number")
		people_id := r.FormValue("people_id")

		insForm, err := db.Prepare("INSERT INTO address(street, cep, number, people_id) VALUES(?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(street, cep, number, people_id)

		log.Println("INSERT: Street: " + street + " | Cep: " + cep + " | Number: " + number + " | PeopleId" + people_id)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM people WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	n := People{}

	for selDB.Next() {
		var id int
		var name string
		var phone int
		var cpf int

		err = selDB.Scan(&id, &name, &phone, &cpf)
		if err != nil {
			panic(err.Error())
		}

		n.Id = id
		n.Name = name
		n.Phone = phone
		n.Cpf = cpf

	}

	tmpl.ExecuteTemplate(w, "Edit", n)

	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		phone := r.FormValue("phone")
		cpf := r.FormValue("cpf")

		insForm, err := db.Prepare("INSERT INTO people(name, phone, cpf) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, phone, cpf)

		log.Println("INSERT: Name: " + name + " | Phone: " + phone + " | CPF: " + cpf)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		phone := r.FormValue("phone")
		cpf := r.FormValue("cpf")
		id := r.FormValue("uid")

		insForm, err := db.Prepare("UPDATE people SET name=?, phone=?, cpf=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, phone, cpf, id)

		log.Println("UPDATE: Name: ", name, " | Phone: " ,phone, " | CPF: ", cpf)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	nId := r.URL.Query().Get("id")

	delForm, err := db.Prepare("DELETE FROM people WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	delForm.Exec(nId)

	log.Println("DELETE")

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func main() {

	log.Println("Server started on: http://localhost:9000")

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)

	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
	http.HandleFunc("/insert-address", InsertAddress)

	http.ListenAndServe(":9000", nil)

}

//S
