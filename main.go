package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type User struct {
	gorm.Model
	Name  string
	Email string
}

var Database *gorm.DB

func init() {
	var err error

	Database, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic(err.Error())
	}

	// set this to 'true' to see sql logs
	Database.LogMode(true)

	fmt.Println("Database connection successful.")
	// Database.Debug().CreateTable(&User{})
	// var user []User = []User{
	// 	User{Name: "Clifton", Email: "cliftonavil@gmail.com"},
	// 	User{Name: "Austin", Email: "austindsouza41@gmail.com"},
	// }

	// for _, user := range user {
	// 	Database.Debug().Create(&user)
	// }
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/users", allUsers).Methods("GET")
	router.HandleFunc("/user/{id}", oneUser).Methods("GET")
	router.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/update_user/{id}", updateUser)
	router.HandleFunc("/user/", newUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	Database.Find(&users)
	fmt.Println("{}", users)
	json.NewEncoder(w).Encode(users)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("id", id)
	var users User
	Database.Where("id = ?", id).Delete(&users)
}

func oneUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var user []User
	Database.Where("id = ?", id).First(&user).RecordNotFound()
	if len(user) > 0 {
		json.NewEncoder(w).Encode(user)
	} else {
		data := "NOt FOund!!!"
		json.NewEncoder(w).Encode(data)
	}

}

func newUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	email := r.FormValue("email")
	fmt.Println(name, email)
	// Create
	Database.Create(&User{Name: name, Email: email})
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	name := r.FormValue("name")
	email := r.FormValue("email")
	var user User
	Database.Debug().First(&user, id).Update(map[string]interface{}{"Name": name, "Email": email})
	fmt.Printf("completee!!")

}
