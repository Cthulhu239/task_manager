// to centralize the db logic
package database

import (
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(databaseURL string) (*pgxpool.Pool,error){
	var ctx context.Context = context.Background()
	var config *pgxpool.Config
	var err error
	config,err = pgxpool.ParseConfig(databaseURL)

	if err != nil{
		log.Printf("unable to parse database_url: %v",err)
		return nil, err
	}
	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx,config)

	if err != nil{
		log.Printf("unable to create a connection pool: %v",err)
		return nil, err
	}

	err = pool.Ping(ctx)

	if err != nil{
		log.Printf("unable to ping the database: %v",err)
		pool.Close()
		return nil,err
	}
    
	log.Println("successfully connected to postgre sql database")
	return pool,nil
}

