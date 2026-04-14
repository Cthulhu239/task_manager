package main

import (
	//"fmt"
	"log"

	"github.com/Cthulhu239/task_manager/internal/config"
	"github.com/Cthulhu239/task_manager/internal/database"
	"github.com/Cthulhu239/task_manager/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main(){

	var cfg *config.Config
	var err error

	cfg,err = config.Load()

	if err != nil{
		log.Fatal("failed to load configuration",err)
	}
	var pool *pgxpool.Pool
	pool,err = database.Connect(cfg.DatabaseURL)

    if err != nil{
		log.Fatal("failed to connect to database:",err)
	}

    defer pool.Close()

    var router *gin.Engine = gin.Default() //pointer to avoid creating copies
	 
	 router.GET("/",func(c *gin.Context){    //the context stores both the request and the response
		   //map[string]any
           c.JSON(200,gin.H{
			"message":"task manager api is running",
			"status":"success",
			"database":"connected",
		   })
	 }) 
	 router.POST("/task",handler.CreateTaskHandler(pool))
	 router.GET("/tasks",handler.GetAllTasksHandler(pool))
	 router.GET("/task/:id",handler.GetTaskByIdHandler(pool))
	 router.PUT("/task/:id",handler.UpdateTaskHandler(pool))
	 router.DELETE("/task/:id",handler.DeleteTaskHandler(pool))
	 router.POST("/auth/register",handler.CreateUserHandler(pool))
	 router.POST("/auth/login",handler.LoginHandler(pool,cfg))
	 router.Run(":" + cfg.Port)

}