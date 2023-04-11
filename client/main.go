package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

type User struct {
	Gender   string `json:"gender"`
	Name     Name   `json:"name"`
	Location struct {
		Street      Street      `json:"street"`
		City        string      `json:"city"`
		State       string      `json:"state"`
		Country     string      `json:"country"`
		Postcode    interface{} `json:"postcode"`
		Coordinates Coordinates `json:"coordinates"`
		Timezone    Timezone    `json:"timezone"`
	} `json:"location"`
	Email      string `json:"email"`
	Login      Login  `json:"login"`
	Dob        Dob    `json:"dob"`
	Registered struct {
		Date string `json:"date"`
		Age  int    `json:"age"`
	} `json:"registered"`
	Phone   string `json:"phone"`
	Cell    string `json:"cell"`
	Id      Id     `json:"id"`
	Picture struct {
		Large     string `json:"large"`
		Medium    string `json:"medium"`
		Thumbnail string `json:"thumbnail"`
	} `json:"picture"`
	Nat string `json:"nat"`
}

type Name struct {
	Title string `json:"title"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type Street struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

type Coordinates struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Timezone struct {
	Offset      string `json:"offset"`
	Description string `json:"description"`
}

type Login struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Md5      string `json:"md5"`
	Sha1     string `json:"sha1"`
	Sha256   string `json:"sha256"`
}

type Dob struct {
	Date string `json:"date"`
	Age  int    `json:"age"`
}

type Id struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Results struct {
	Users []User `json:"results"`
}

type Info struct {
	Seed    string `json:"seed"`
	Results int    `json:"results"`
	Page    int    `json:"page"`
	Version string `json:"version"`
}

type Response struct {
	Results []User `json:"results"`
	Info    Info   `json:"info"`
}

const host = "localhost"
const dbName = "postgres"
const user = "postgres"
const password = "13134777"

func GetJson() Response {
	var response Response

	resp, err := http.Get("https://randomuser.me/api/")
	if err != nil {
		fmt.Println(err.Error())
		return response
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err.Error())
		return response
	}

	return response
}

func main() {
	connectString := fmt.Sprintf(
		"host=%s dbname=%s sslmode=disable user=%s password=%s", host, dbName, user, password)

	conn, err := sql.Open("postgres", connectString)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sqlCommand := "CREATE TABLE IF NOT EXISTS userdata (uuid text PRIMARY KEY , login text, data json);"
	_, err = conn.Query(sqlCommand)
	if err != nil {
		fmt.Println(err.Error())
	}

	userInfo := GetJson()

	for _, v := range userInfo.Results {
		dataJson, err := json.Marshal(v)
		if err != nil {
			fmt.Println(err.Error())
		}

		// uuid Будет использовать для поиска уже прошедших вход/регистрацию пользователей
		// а поле login будет использовать для поиска пользователя во время входа в свою учетную запись
		sqlInsertCommand := fmt.Sprintf("insert into userdata (uuid,login,data) VALUES('%s','%s','%s')", v.Login.Uuid, v.Login.Username, string(dataJson))
		_, err = conn.Query(sqlInsertCommand)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
