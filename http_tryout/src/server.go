package main

import (
	"net/http"
	"path/filepath"
	"fmt"
	"io/ioutil"
	_ "github/lib/pq"
	"database/sql"
	"strings"
	"log"
	"strconv"
)

const (
	DB_USER = "postgres"
	DB_PASSWORD = "19Sep"
	DB_NAME = "postgres"
)

func staticHandler(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	fmt.Println("Path  : ",path)
	if path == "/" {
		path = "/inputForm.html"
	}
	renderFile(res, "src" + path)
}

func renderFile(res http.ResponseWriter, file string) {
	var fileName, fileNotFound = filepath.Abs(file)
	if fileNotFound != nil {
		fmt.Printf("File not found")
	}
	var reader, err = ioutil.ReadFile(fileName)
	checkErr(err)
	res.Write(reader)
}

func addNameHandler(res http.ResponseWriter, req *http.Request) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	rows,err := db.Query("CREATE SCHEMA IF NOT EXISTS ListOfEmployees;")
	checkErr(err)
	fmt.Println("created schema if not exist",rows)
	values :=strings.Split(req.FormValue("Name")," ")
	row := db.QueryRow("INSERT INTO employelist(name,id)values('"+ values[0]+"',"+ values[1]+");")
	fmt.Println("selected list ",row)
	http.Redirect(res,req,"/",308)
}

func deleteNameHandler(res http.ResponseWriter, req *http.Request) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	rows,err := db.Query("CREATE SCHEMA IF NOT EXISTS ListOfEmployees;")
	checkErr(err)
	fmt.Println("created schema if not exist",rows)
	values :=strings.Split(req.FormValue("Name")," ")
	row := db.QueryRow("delete from employelist where name ='"+values[0]+"';")
	fmt.Println("selected list ",row)
	http.Redirect(res,req,"/",308)
}

func viewHandler(res http.ResponseWriter, req *http.Request) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	rows,err := db.Query("CREATE SCHEMA IF NOT EXISTS ListOfEmployees;")
	checkErr(err)
	rows,err = db.Query("SELECT * FROM employelist")
	checkErr(err)
	defer rows.Close()
	var (
		id int
		name string
	)
	var output []byte
	for rows.Next() {
		err := rows.Scan(&name, &id)
		if err != nil {
			log.Fatal(err)
		}
		output=[]byte(string(output[:len(output)])+"\n"+name+" : "+strconv.Itoa(id))
	}
	res.Write(output)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", staticHandler)
	http.HandleFunc("/addName", addNameHandler)
	http.HandleFunc("/deleteName", deleteNameHandler)
	http.HandleFunc("/view", viewHandler)
	http.ListenAndServe(":8000", nil)
}


