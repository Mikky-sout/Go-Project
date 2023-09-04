package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Ledger struct {
	ID       int     `json:"id"`
	Status   string  `json:"status"`
	Detail   string  `json:"detail"`
	Amount   float64 `json:"amount"`
	DateTime string  `json:"dateTime"`
}

func getValueFromParams(params string, contentType string) any {
	if params != "" {
		switch contentType {
		case "int":
			output, _ := strconv.Atoi(params)
			return output
		case "float":
			output, _ := strconv.ParseFloat(params, 64)
			return output
		case "string":
			return params
		}
	} else {
		return ""
	}
	return nil
}

func findIndexWithId(id int) int {
	recordList := &[]Ledger{}
	outputs := execute("SELECT * FROM `ledger`")
	for outputs.Next() {
		record := Ledger{}
		err := outputs.Scan(&record.ID, &record.Status, &record.Detail, &record.Amount, &record.DateTime)
		if err != nil {
			panic(err)
		}
		*recordList = append(*recordList, record)
	}
	outputs.Close()

	for i, element := range *recordList {
		if element.ID == id {
			return i
		}
	}
	return -1
}

func getID() int {
	recordList := &[]Ledger{}
	outputs := execute("SELECT * FROM `ledger`")
	for outputs.Next() {
		record := Ledger{}
		err := outputs.Scan(&record.ID, &record.Status, &record.Detail, &record.Amount, &record.DateTime)
		if err != nil {
			panic(err)
		}
		*recordList = append(*recordList, record)
	}
	outputs.Close()

	if len(*recordList) <= 0 {
		return 1
	} else {
		id := 1
		isFound := true
		for isFound {
			isFound = false
			for _, val := range *recordList {
				if val.ID == id {
					isFound = true
					id = id + 1
					break
				}
			}
		}
		return id
	}
}

func getRecord(c *gin.Context) {
	recordList := &[]Ledger{}
	outputs := execute("SELECT * FROM `ledger`")
	if err := c.ShouldBind(&Ledger{}); err != nil {
		log.Fatal(err)
	}
	for outputs.Next() {
		record := Ledger{}
		err := outputs.Scan(&record.ID, &record.Status, &record.Detail, &record.Amount, &record.DateTime)
		if err != nil {
			panic(err)
		}
		*recordList = append(*recordList, record)
	}
	outputs.Close()
	c.JSON(200, recordList)
}

func concatString(stringSet []string) string {
	outputString := stringSet[0]
	if len(stringSet) > 1 {
		for i, val := range stringSet {
			if i > 0 {
				outputString = outputString + "," + val
			}
		}
	}
	return outputString
}

func addRecord(c *gin.Context) {
	var led Ledger

	curruentTime := fmt.Sprint(time.Now().Format("02/Jan/2006,15:04:05"))

	status := c.DefaultQuery("status", "Unknown")
	detail := c.DefaultQuery("detail", "0")
	amount, _ := strconv.ParseFloat(c.DefaultQuery("amount", "0.00"), 64)

	if err := c.ShouldBind(&led); err != nil {
		log.Fatal(err)
	}
	led = Ledger{ID: getID(), Status: status, Detail: detail, Amount: amount, DateTime: curruentTime}
	query, err := db.Prepare("INSERT INTO `ledger` (id,status,detail,amount,dateTime) VALUES (?,?,?,?,?)")
	if err != nil {
		panic(err)
	}

	query.Exec(led.ID, led.Status, led.Detail, led.Amount, led.DateTime)
	c.JSON(200, "Successfully adding")
}

func updateRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	status := getValueFromParams(c.DefaultQuery("status", ""), "string")
	detail := getValueFromParams(c.DefaultQuery("detail", ""), "string")
	amount := getValueFromParams(c.DefaultQuery("amount", ""), "float")
	params := []string{}

	if status != "" && status != nil {
		params = append(params, fmt.Sprintf("status = '%v'", status.(string)))
		fmt.Println(params)
		// LedgerList[i].Status = status.(string)
	}
	if detail != "" && detail != nil {
		params = append(params, fmt.Sprintf("detail = '%v'", detail.(string)))
		fmt.Println(params)
		// LedgerList[i].Detail = detail.(string)
	}
	if amount != "" && amount != nil {
		params = append(params, fmt.Sprintf("amount = %v", amount.(float64)))
		fmt.Println(params)
		// LedgerList[i].Amount = amount.(float64)
	}

	fmt.Print(params)
	qString := fmt.Sprintf("UPDATE ledger SET %s WHERE id=%v", concatString(params), id)
	fmt.Print(qString)
	outputs := execute(qString)
	outputs.Close()
	c.JSON(200, "Successfully updating")
}

func deleteRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if index := findIndexWithId(id); index != -1 {
		outputs, err := db.Query(fmt.Sprintf("DELETE FROM ledger WHERE id=%v", id))
		if err != nil {
			panic(err)
		}
		outputs.Close()
		c.JSON(200, "Successfully Deleting")
	} else {
		c.JSON(200, "ID Not found.")
	}
}

var db *sql.DB
var dbErr error

func databaseConnect() {
	db, dbErr = sql.Open("mysql", "root@tcp(localhost:3306)/mikdb")

	if dbErr != nil {
		fmt.Println("error while valid")
		panic(dbErr)
	}
}

func execute(query string) *sql.Rows {
	output, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	return output
}

func main() {

	r := gin.Default()

	databaseConnect()

	r.GET("/get", getRecord)

	r.POST("/add", addRecord)

	r.PATCH("/update/:id/", updateRecord)

	r.DELETE("/remove/:id", deleteRecord)

	r.Run(":8000")
}
