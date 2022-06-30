package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

var db *sql.DB

type BearMemory struct {
	CreationDate       int    `json:"creationDate"`
	Base64StringOfFile string `json:"base64StringOfFile"`
}

var allBearMemories = []BearMemory{
	{Base64StringOfFile: "hello world", CreationDate: 19391},
}

func getBearMemory(c *gin.Context) {

}

func getBearMemories(c *gin.Context) {
	fmt.Println(allBearMemories)
	c.IndentedJSON(http.StatusOK, dbFunctionForGetMemory(db))
}

func postBearMemory(c *gin.Context) {
	var newBearMemory = BearMemory{}

	err := c.BindJSON(&newBearMemory)
	if err != nil {
		fmt.Println("IS IT HERE")
		log.Fatalln(err)
	}
	res, er := db.Exec("INSERT INTO BearMemories (creationDate, base64StringOfFile) VALUES (" + fmt.Sprint(newBearMemory.CreationDate) + ", '" + newBearMemory.Base64StringOfFile + "');")
	if er != nil {
		fmt.Println("BEAR BEHAVIOR")
		log.Fatalln(er)
	}
	fmt.Println(res.LastInsertId())
	c.IndentedJSON(http.StatusCreated, newBearMemory)
}

func dbFunctionForGetMemory(db *sql.DB) []BearMemory {

	rows, err := db.Query("SELECT * FROM BearMemories;")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	var allBearMemories []BearMemory
	for rows.Next() {
		var creationDate int
		var base64StringOfFile string
		err := rows.Scan(&creationDate, &base64StringOfFile)
		if err != nil {
			log.Fatalln(err)
		}
		var bearMemoryFromDB BearMemory = BearMemory{Base64StringOfFile: base64StringOfFile, CreationDate: creationDate}
		allBearMemories = append(allBearMemories, bearMemoryFromDB)
	}
	allBearMemories = append(allBearMemories, BearMemory{Base64StringOfFile: "", CreationDate: 30})
	return allBearMemories

}
func main() {
	//fmt.Println("Hello, World!")
	fmt.Println("This is my Go Server!")
	var resp *http.Response
	var err error
	resp, err = http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	//  fmt.Println(string(bytes))
	var arrOfJson []map[string]interface{}
	err = json.Unmarshal(body, &arrOfJson)
	if err != nil {
		log.Fatalln(err)
	}
	if len(arrOfJson) > 0 {
		fmt.Println(arrOfJson[0]["userId"])
	}
	db, err = sql.Open("postgres", "postgres://rmfbqwqhgpyhsf:1a74f7e2002ba419e93f09831576d267f7212f8e349e3225ad9e31f411049ed9@ec2-23-23-151-191.compute-1.amazonaws.com:5432/daduj4br7pupke")
	if err != nil {
		log.Fatalln(err)
	}
	if testPing := db.Ping(); testPing != nil {
		log.Fatalln(testPing)
	}
	_, dberr := db.Exec("CREATE TABLE IF NOT EXISTS BearMemories (creationDate bigint, base64StringOfFile text);")
	if dberr != nil {
		log.Fatalln(dberr)
	}
	router := gin.Default()
	router.GET("/bear-memory", getBearMemories)
	router.POST("/bear-memory", postBearMemory)
	router.GET("/bear-memory/:id", getBearMemory)
	router.Run()
}
