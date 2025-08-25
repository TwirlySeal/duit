package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	Title string `json:"title"`
	Done bool `json:"done"`
}

type Project struct {
	Name string
	Id string
}

const tasksQuery string = "SELECT title, done FROM tasks WHERE project_id=$1"

func homeData(dbpool *pgxpool.Pool, ctx context.Context, id string) ([]Project, []Task, error) {
	batch := pgx.Batch{}
	batch.Queue("SELECT * FROM projects")
	batch.Queue(tasksQuery, id)

	results := dbpool.SendBatch(ctx, &batch)

	projects, err := load[Project](results)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting projects: %w", err)
	}

	tasks, err := load[Task](results)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting tasks: %w", err)
	}

	err = results.Close()
	if err != nil {
		log.Println(err)
	}

	return projects, tasks, nil
}

func load[T any](results pgx.BatchResults) ([]T, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return pgx.CollectRows[T](rows, pgx.RowToStructByName)
}

func getTasks(dbpool *pgxpool.Pool, ctx context.Context, id string) ([]Task, error) {
	rows, err := dbpool.Query(ctx, tasksQuery, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return pgx.CollectRows[Task](rows, pgx.RowToStructByName)
}

func addTask(dbpool *pgxpool.Pool, ctx context.Context, task NewTask) error {
	_, err := dbpool.Exec(ctx, "INSERT INTO tasks VALUES ($1, false, $2)", task.Name, task.ProjectId)
	return err
}
