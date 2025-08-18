package main

import (
	"context"
	"fmt"
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
	defer results.Close()

	projects, err := load[Project](results)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting projects: %w", err)
	}

	tasks, err := load[Task](results)
	if err != nil {
		return nil, nil, fmt.Errorf("Error getting tasks: %w", err)
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
