package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Cthulhu239/task_manager/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

//to interact using SQL intertwined with go logic (repository)
func CreateTask(pool *pgxpool.Pool,title string,completed bool) (*models.Task,error){
    var ctx context.Context //db connection context
	var cancel context.CancelFunc
	ctx,cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query string = `
	   INSERT INTO task (title,completed)
	   VALUES ($1,$2)
	   RETURNING id,title,completed,created_at,updated_at
	`

	var task models.Task

	var err error = pool.QueryRow(ctx,query,title,completed).Scan(
		&task.ID,
		&task.Title,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil{
		return nil,err
	}

	return &task, nil
}

func GetAllTask(pool *pgxpool.Pool) ([]models.Task,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query string = `
	   SELECT id,title,completed,created_at,updated_at
	   FROM task
	   ORDER BY created_at DESC
	`

	var rows,err =pool.Query(ctx,query)

	if err != nil{
		return nil,err
	}

	defer rows.Close()

	var tasks []models.Task = []models.Task{}
    
	for rows.Next(){
		var task models.Task

		err = rows.Scan(
			&task.ID,
			&task.Title,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		)

		if err != nil{
			return nil,err
		}

		tasks = append(tasks,task)
	}

	if err = rows.Err(); err != nil{
		return nil,err
	}

	return tasks,nil
}

func GetTaskById(pool *pgxpool.Pool,id int) (*models.Task,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	var query string = `
	  SELECT * FROM task
	  WHERE id = $1
	`

	var task models.Task
    var err  error = pool.QueryRow(ctx,query,id).Scan( //actually executes the query
       &task.ID,
	   &task.Title,
	   &task.Completed,
	   &task.CreatedAt,
	   &task.UpdatedAt,
	)

	if err != nil{
		return nil,err
	}

	return &task,nil
    
}

func UpdateTask(pool *pgxpool.Pool,id int, title string, completed bool) (*models.Task,error){
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query string = `
	  UPDATE task
	  SET title = $1,completed = $2,updated_at = CURRENT_TIMESTAMP
      WHERE id = $3
	  RETURNING id,title,completed,created_at,updated_at
	`
	var task models.Task

	var err error = pool.QueryRow(ctx,query,title,completed,id).Scan(
	   &task.ID,
	   &task.Title,
	   &task.Completed,
	   &task.CreatedAt,
	   &task.UpdatedAt,
	)

	if err != nil{
		return nil,err
	}
    return &task,nil
}

func DeleteTask(pool *pgxpool.Pool, id int) error {
    var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()

	var query string = `
      DELETE FROM task
	  WHERE id = $1
	`

	commandTag,err := pool.Exec(ctx,query,id)
	
	if err != nil{
     return err
	}

	if commandTag.RowsAffected() == 0{
		return fmt.Errorf("task with id %d was not found",id)
	}
	return nil
}