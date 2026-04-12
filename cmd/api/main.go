package main


import (
	//"fmt"
	"github.com/gin-gonic/gin"
)

func main(){
     var router *gin.Engine = gin.Default() //pointer to avoid creating copies
	 router.GET("/",func(c *gin.Context){    //the context stores both the request and the response
		   //map[string]any
           c.JSON(200,gin.H{
			"message":"task manager api is running",
			"status":"success",
		   })
	 }) 
	 
	 router.Run(":3000")

}