package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func Calculate(times int) int64 {
	if times == 1 || times == 0 {
		return 1
	}

	return Calculate(times-1) + Calculate(times-2)
}
func main() {
	router := gin.Default()
	router.GET("calculate", func(c *gin.Context) {

		str := c.Query("x")
		value, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("error : %s \n", err)
			return
		}
		result := Calculate(value)
		fmt.Print(result)
		c.JSON(http.StatusOK, result)

	})

	err := router.Run(":8080")
	if err != nil {
		return
	}

}
