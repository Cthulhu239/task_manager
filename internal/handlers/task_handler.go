package handler

import (
	"net/http"
	"strconv"

	"github.com/Cthulhu239/task_manager/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTaskInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTaskInput struct{
	Title     *string  `json:"title"`
	Completed *bool   `json:"completed"`
	// &true ---------------> set completed as -> true
	// &false ---------------> set completed as -> false
	// nil ---------------> set completed as -> not provided
}

func CreateTaskHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func (c *gin.Context){
		var input CreateTaskInput

		if err := c.ShouldBindJSON(&input); err != nil{ // json to struct fail
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}

		task,err := repository.CreateTask(pool,input.Title,input.Completed) //creating the task and save to db

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()}) //problem with database
			return
		}

		c.JSON(http.StatusCreated,task)
	}
}

func GetAllTasksHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func (c *gin.Context){
		tasks,err := repository.GetAllTask(pool)

		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
			return
		}
   

		c.JSON(http.StatusOK,tasks)
	}
}

func GetTaskByIdHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func (c *gin.Context){
        idStr := c.Param("id")
		id,err := strconv.Atoi(idStr)
	    if err != nil{
			c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
			return
		}	

		task,err := repository.GetTaskById(pool,id)

		if err != nil{
			if err == pgx.ErrNoRows{
				c.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		}

		c.JSON(http.StatusOK,task)
	}
}

func UpdateTaskHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task_id"})
			return
		}

		var input UpdateTaskInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
        
		if input.Title == nil && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field should be provided"})
			return
		}

		existing, err := repository.GetTaskById(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return 
		}

		title := existing.Title
		if input.Title != nil {
			title = *input.Title
		}

		completed := existing.Completed
		if input.Completed != nil {
			completed = *input.Completed
		}

		task, err := repository.UpdateTask(pool, id, title, completed)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, task)
	}
}

func DeleteTaskHandler(pool *pgxpool.Pool) gin.HandlerFunc{
	return func(c *gin.Context) {
		idStr := c.Param("id")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task_id"})
			return
		}
       
		err = repository.DeleteTask(pool,id)
		if err != nil{
			if err.Error() == "task with id "+idStr+" was not found"{
				c.JSON(http.StatusNotFound,gin.H{"error":"task not found"}) //to send it to the frontend
				return
			}
			c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		}
		c.JSON(http.StatusOK,gin.H{"message":"task deleted successfully"})
	}
}