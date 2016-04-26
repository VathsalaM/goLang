package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"path/filepath"
	_ "github/lib/pq"
	"database/sql"
	"strings"
	"strconv"
)

const (
	DB_USER = "postgres"
	DB_PASSWORD = "19Sep"
	DB_NAME = "postgres"
)

func staticHandler(res http.ResponseWriter, req *http.Request) {
	path:=req.URL.Path
	if path == "/" {
		path = "/registrationPage.html"
	}
	var fileName, err = filepath.Abs(path)
	if err != nil {
		fmt.Printf("File not found")
	}
	reader, err := ioutil.ReadFile("public"+fileName)
	checkErr(err)
	extension:=strings.Split(path,".")
	res.Header().Set("Content-Type","text/"+extension[len(extension)-1])
	res.Write(reader)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func valueExtractor(res http.ResponseWriter, req *http.Request) (values []string) {
	req.ParseForm()
	registrationForm := req.Form
	firstName := registrationForm.Get("firstName")
	lastName := registrationForm.Get("lastName")
	Gender := registrationForm.Get("gender")
	E_mail := registrationForm.Get("Email")
	Address1 := registrationForm.Get("address1")
	Address2 := registrationForm.Get("address2")
	UserId := strconv.Itoa(len(firstName + lastName))
	DateOfCreation := res.Header().Get("Date")
	Referer := req.URL.String()
	RequestMethod := req.Method
	RemoteAddress := req.RemoteAddr
	RequestCookie:= req.Header.Get("Cookie")
	Host := req.Host
	UserAgent := req.UserAgent()
	values = []string{firstName, lastName, Gender, E_mail, Address1, Address2, UserId, DateOfCreation, Referer, RequestMethod, RemoteAddress,string(RequestCookie), Host, UserAgent}
	return
}

func registrationHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	firstName:=req.Form.Get("firstName")
	lastName:=req.Form.Get("lastName")
	cookie:=http.Cookie{Name:"userName",Value:firstName+" "+lastName}
	req.Header.Del("Cookie")
	http.SetCookie(res,&cookie)
	//req.Header.Set("Cookie","name="+firstName+" "+lastName+";value="+strconv.Itoa(len(firstName+lastName)))
	values := valueExtractor(res, req)
	fmt.Println("==>",values[len(values)-3],"<==",req.Header.Get("Cookie"))
	cookies:=res.Header().Get("Cookies")
	fmt.Println("cookie length :",cookies)
	//fmt.Println("first cookie",cookies[0])
	//fmt.Println("second cookie",cookies[1])
	//dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
	//	DB_USER, DB_PASSWORD, DB_NAME)
	//db, err := sql.Open("postgres", dbinfo)
	//checkErr(err)
	//row, err := db.Query("INSERT INTO employe(firstName,lastName,Gender,E_mail,Address1,Address2,UserId,DateOfCreation,Referer,RequestMethod,RemoteAddress,RequestCookies,Host,UserAgent)values('" + values[0] + "','" + values[1] + "','" + values[2] + "','" + values[3] + "','" + values[4] + "','" + values[5] + "','" + values[6] + "','" + values[7] + "','" + values[8] + "','" + values[9] + "','" + values[10] + "','" + values[11] + "','" + values[12] + "','" + values[13] + "');")
	//checkErr(err)
	//fmt.Println(row)
}

//func successHandler(res http.ResponseWriter, req *http.Request) {
//	fmt.Println("hey..........success handler is called.....")
//	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
//		DB_USER, DB_PASSWORD, DB_NAME)
//	db, err := sql.Open("postgres", dbinfo)
//	checkErr(err)
//	rows, err := db.Query("SELECT userId,firstName,lastName,Gender,Address1,Address2,RequestCookies FROM employe where firstName='" + req.FormValue("firstName") + "';")
//	checkErr(err)
//	var (
//		firstName string
//		lastName string
//		Gender string
//		Address1 string
//		Address2 string
//		UserId string
//		//ContentEncoding string
//		//ContentType string
//		//ResponseBody string
//		ResponseCookie string
//	)
//	rows.Next()
//	err = rows.Scan(&UserId, &firstName, &lastName, &Gender, &Address1, &Address2, &ResponseCookie)
//	checkErr(err)
//	htmlData, err := ioutil.ReadAll(req.Body)
//	data := string(htmlData)
//	outputData := "\n\nUser Details \n\nUserId : " + UserId + "\nName : " + firstName + " " + lastName + "\nGender : " + Gender + "\nAddress1 : " + Address1 + "\nAddress2 : " + Address2 + "\nContent-Type : " + http.DetectContentType([]byte(data[:len(data)])) + "\nResponseCookie : " + ResponseCookie
//	output := []byte(string(outputData[:len(outputData)]))
//	res.Write(output)
//	fmt.Println("sucess page closed")
//}

func viewHandler(res http.ResponseWriter, req *http.Request) {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	checkErr(err)
	rows, err := db.Query("SELECT userId,firstName,lastName,Gender,Address1,Address2,RequestCookies FROM employe;")
	checkErr(err)
	defer rows.Close()
	var (
		firstName string
		lastName string
		Gender string
		Address1 string
		Address2 string
		UserId string
		//ContentEncoding string
		//ContentType string
		//ResponseBody string
		ResponseCookie string
	)
	outputData := "User Details \n\nUserId \tName\tGender\tAddress1\tAddress2\tContent-Type\tResponseCookie"
	for rows.Next() {
		err = rows.Scan(&UserId, &firstName, &lastName, &Gender, &Address1, &Address2, &ResponseCookie)
		checkErr(err)
		htmlData, newErr := ioutil.ReadAll(req.Body)
		checkErr(newErr)
		data := string(htmlData)
		outputData = outputData+"\n"+UserId+"\t" +firstName + " " + lastName + "\t" + Gender + "\t" + Address1 + "\t" + Address2 + "\t" + http.DetectContentType([]byte(data[:len(data)])) + "\t" + ResponseCookie
	}
	output := []byte(string(outputData[:len(outputData)]))
	res.Write(output)
}

func main() {
	http.HandleFunc("/", staticHandler)
	http.HandleFunc("/register", registrationHandler)
	//http.HandleFunc("/success", successHandler)
	http.HandleFunc("/view", viewHandler)
	http.ListenAndServe(":8000", nil)
}
