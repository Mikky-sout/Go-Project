package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Age        int     `json:"age"`
	Salary     float64 `json:"salary"`
	Department string  `json:"department"`
	Tel        string  `json:"tel"`
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

func findIndexWithId(id int, arr []Employee) int {
	for i, element := range arr {
		if element.ID == id {
			return i
		}
	}
	return -1
}

func removeWithIndex(arr []Employee, index int) []Employee {
	newArr := make([]Employee, 0)
	newArr = append(newArr, arr[:index]...)
	return append(newArr, arr[index+1:]...)
}

func main() {
	r := gin.Default()
	empList := []Employee{}

	r.GET("/get", func(c *gin.Context) {
		var emp Employee

		if err := c.ShouldBind(&emp); err != nil {
			log.Fatal(err)
		}
		c.JSON(200, empList)
	})

	r.POST("/add", func(c *gin.Context) {
		var emp Employee

		name := getValueFromParams(c.DefaultQuery("name", "Unknown"), "string")
		age, _ := strconv.Atoi(c.DefaultQuery("age", "0"))
		salary, _ := strconv.ParseFloat(c.DefaultQuery("salary", "0.00"), 64)
		department := c.DefaultQuery("department", "Unknown")
		tel := c.DefaultQuery("tel", "099-000-0000")

		if err := c.ShouldBind(&emp); err != nil {
			log.Fatal(err)
		}

		emp = Employee{ID: len(empList) + 1, Name: name.(string), Age: age, Salary: salary, Department: department, Tel: tel}
		empList = append(empList, emp)
		c.JSON(200, empList)
	})

	r.PATCH("/update/:id/", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		name := getValueFromParams(c.DefaultQuery("name", ""), "string")
		age := getValueFromParams(c.DefaultQuery("age", ""), "int")
		salary := getValueFromParams(c.DefaultQuery("salary", ""), "float")
		department := getValueFromParams(c.DefaultQuery("department", ""), "string")
		tel := getValueFromParams(c.DefaultQuery("tel", ""), "string")

		for i, element := range empList {
			if element.ID == id {
				if name != "" {
					empList[i].Name = name.(string)
				}
				if age != "" {
					empList[i].Age = age.(int)
				}
				if salary != "" {
					empList[i].Salary = salary.(float64)
				}
				if department != "" {
					empList[i].Department = department.(string)
				}
				if tel != "" {
					empList[i].Tel = tel.(string)
				}
				break
			}
		}
		c.JSON(200, empList)
	})

	r.DELETE("/remove/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		if index := findIndexWithId(id, empList); index != -1 {
			empList = removeWithIndex(empList, index)
			c.JSON(200, empList)
		} else {
			c.JSON(200, "ID Not found.")
		}
	})

	r.Run(":8000")
}
