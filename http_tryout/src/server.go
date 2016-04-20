package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"path/filepath"
	"os"
)

func staticHandler(res http.ResponseWriter, req *http.Request) {
	path:=req.URL.Path
	if path=="/"{
		path="/inputForm.html"
	}
	println(path)
	renderFile(res,"src"+path)
}

func renderFile(res http.ResponseWriter, file string) {
	var fileName, fileNotFound = filepath.Abs(file)
	if fileNotFound != nil {
		fmt.Printf("File not found")
	}
	var reader, err = ioutil.ReadFile(fileName)
	if (err != nil) {
		res.WriteHeader(404)
		errorMessage:=[]byte{'4','0','4'}
		res.Write(errorMessage)
	}
	res.Write(reader)
}

func nameHandler(res http.ResponseWriter, req *http.Request) {
	var fileName, err = filepath.Abs("src/data.txt")
	file, err := os.OpenFile(fileName, os.O_APPEND | os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err = file.WriteString(req.FormValue("Name") + "\n"); err != nil {
		panic(err)
	}
	http.Redirect(res, req, "/", 301)
}

func viewHandler(res http.ResponseWriter, req *http.Request) {
	renderFile(res, "src/data.txt")
}

func main() {
	http.HandleFunc("/", staticHandler)
	http.HandleFunc("/name", nameHandler)
	http.HandleFunc("/view", viewHandler)
	http.ListenAndServe(":8000", nil)
}
