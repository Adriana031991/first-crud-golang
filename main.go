package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"text/template"
)

func connectionDB() (connection *sql.DB) {
	Driver := "mysql"
	User := "root"
	Password := "Adrianita-02"
	Name := "crud-estudiante"

	connection, err := sql.Open(Driver, User+":"+Password+"@tcp(127.0.0.1)/"+Name)
	if err != nil {
		panic(err.Error())
	}
	return connection
}

var template1 = template.Must(template.ParseGlob("templates/*"))

func main() {

	http.HandleFunc("/", Init)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/edit", Editar)
	http.HandleFunc("/update", Update)

	log.Println("Server running")
	http.ListenAndServe(":8080", nil)
}

type Employer struct {
	Id    int
	Name  string
	Email string
}

func Init(w http.ResponseWriter, r *http.Request) {
	connectionEstablished := connectionDB()
	registers, err := connectionEstablished.Query("SELECT * FROM emploees ORDER BY id DESC")

	if err != nil {
		panic(err.Error())

	}
	// log.Println("conecction ok")

	employer := Employer{}
	arrayEmployee := []Employer{}

	for registers.Next() {
		var id int
		var name, email string
		err = registers.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())

		}
		employer.Id = id
		employer.Name = name
		employer.Email = email

		arrayEmployee = append(arrayEmployee, employer)
	}
	// fmt.Println(arrayEmployee)
	template1.ExecuteTemplate(w, "init", arrayEmployee)
	defer connectionEstablished.Close()

}

func Create(w http.ResponseWriter, r *http.Request) {
	template1.ExecuteTemplate(w, "create", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {

	connectionEstablished := connectionDB()
	fmt.Println("daata", r.Method)

	if r.Method == "POST" {
		name := r.FormValue("name")
		email := r.FormValue("email")
		fmt.Println("daata2", name, email)

		insertRegister, err := connectionEstablished.Prepare("INSERT INTO emploees (userName, email) VALUES (?,?)")

		if err != nil {
			panic(err.Error())

		}
		insertRegister.Exec(name, email)
		fmt.Println("daata3")
	}

	defer connectionEstablished.Close()
	http.Redirect(w, r, "/", 301)
}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idEmployer := r.URL.Query().Get("id")

	connectionEstablished := connectionDB()
	deleteRegister, err := connectionEstablished.Prepare("DELETE FROM emploees WHERE id = ?")

	if err != nil {
		panic(err.Error())

	}
	deleteRegister.Exec(idEmployer)
	defer connectionEstablished.Close()

	http.Redirect(w, r, "/", 301)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idEmployer := r.URL.Query().Get("id")
	// fmt.Println("idEmployer", idEmployer)

	connectionEstablished := connectionDB()
	register, err := connectionEstablished.Query("SELECT * FROM emploees WHERE id = ?", idEmployer)

	employer := Employer{}

	for register.Next() {
		var id int
		var userName, email string
		err = register.Scan(&id, &userName, &email)
		if err != nil {
			panic(err.Error())

		}
		employer.Id = id
		employer.Name = userName
		employer.Email = email

	}
	// fmt.Println("employer", employer)

	template1.ExecuteTemplate(w, "edit", employer)
	defer connectionEstablished.Close()

}

// https://www.youtube.com/watch?v=G58gN0lIbyI&ab_channel=Develoteca
func Update(w http.ResponseWriter, r *http.Request) {
	connectionEstablished := connectionDB()
	if r.Method == "POST" {

		id := r.FormValue("id")
		name := r.FormValue("name")
		email := r.FormValue("email")

		updateRegister, err := connectionEstablished.Prepare("UPDATE emploees SET userName=?,email=? WHERE id =?")

		if err != nil {
			panic(err.Error())
		}

		updateRegister.Exec(name, email, id)
	}

	defer connectionEstablished.Close()
	http.Redirect(w, r, "/", 301)

}
