package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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

func deleteBearMemory(c *gin.Context) {
	db.Query("DELETE from BearMemories WHERE creationDate=10000;")
	c.IndentedJSON(http.StatusOK, BearMemory{Base64StringOfFile: "", CreationDate: 0})
}

func postBearMemory(c *gin.Context) {
	//start := time.Now()
	var newBearMemory = BearMemory{}

	err := c.BindJSON(&newBearMemory)
	if err != nil {
		//fmt.Println("IS IT HERE")
		log.Fatalln(err)
	}
	var imageStringBeforeChangingFileType string = newBearMemory.Base64StringOfFile
	now := time.Now()
	dec, err := base64.StdEncoding.DecodeString(imageStringBeforeChangingFileType)
	if err != nil {
		panic(err)
	}
	var newFileName = "ImageForBearMemory" + fmt.Sprint(now.Nanosecond()) + fmt.Sprint(rand.Intn(100000)) + ".png"
	f, err := os.Create(newFileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	f.Write(dec)
	openedNewConvertedFile, _ := os.Open(newFileName)

	// Read entire png into new slice
	reader := bufio.NewReader(openedNewConvertedFile)
	content, _ := ioutil.ReadAll(reader)
	encoded := base64.StdEncoding.EncodeToString(content)
	newBearMemory.Base64StringOfFile = encoded
	openedNewConvertedFile.Close()
	defer os.Remove(newFileName)
	var sqlQueryString string = "INSERT INTO BearMemories (creationDate, base64StringOfFile) VALUES (" + fmt.Sprint(newBearMemory.CreationDate) + ", '" + newBearMemory.Base64StringOfFile + "');"
	res, er := db.Exec(sqlQueryString)
	if er != nil {
		//fmt.Println("BEAR BEHAVIOR")
		log.Fatalln(er)
	}
	fmt.Println(res.LastInsertId())
	//elapsed := time.Since(start)
	//log.Printf("Binomial took %s", elapsed)
	c.IndentedJSON(http.StatusCreated, newBearMemory)
}

func dbFunctionForGetMemory(db *sql.DB) []BearMemory {

	rows, err := db.Query("SELECT * FROM BearMemories ORDER BY creationDate ASC;")
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
	router.DELETE("/bear-memory", deleteBearMemory)
	router.Run()
}
