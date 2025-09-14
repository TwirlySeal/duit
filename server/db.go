package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Project struct {
	Id   int
	Name string
}

type PageData struct {
	Projects []Project
	Tasks    []Task
	ActiveProject Project
}

const tasksQuery string = "SELECT id, title, done FROM tasks WHERE project_id=$1 AND done=FALSE"

func homeData(dbPool *pgxpool.Pool, ctx context.Context, id int) (PageData, error) {
	batch := pgx.Batch{}
	batch.Queue("SELECT * FROM projects")
	batch.Queue(tasksQuery, id)
	batch.Queue("SELECT name FROM projects WHERE id=$1", id)

	results := dbPool.SendBatch(ctx, &batch)
	defer func() {
		if err := results.Close(); err != nil {
			log.Println(err)
		}
	}()

	projects, err := load[Project](results)
	if err != nil {
		return PageData{}, fmt.Errorf("Failed to get projects: %w", err)
	}

	tasks, err := load[Task](results)
	if err != nil {
		return PageData{}, fmt.Errorf("Failed to get tasks: %w", err)
	}

	activeProject := Project{Id: id}
	if err := results.QueryRow().Scan(&activeProject.Name); err != nil {
		return PageData{}, fmt.Errorf("Failed to get project name: %w", err)
	}

	return PageData{
		projects,
		tasks,
		activeProject,
	}, nil
}

func load[T any](results pgx.BatchResults) ([]T, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return pgx.CollectRows[T](rows, pgx.RowToStructByName)
}

func getTasks(dbPool *pgxpool.Pool, ctx context.Context, id int) ([]Task, error) {
	rows, err := dbPool.Query(ctx, tasksQuery, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	return pgx.CollectRows[Task](rows, pgx.RowToStructByName)
}

type NewTask struct {
	Title     string `json:"title"`
	ProjectId int    `json:"projectId"`
}

func addTask(dbPool *pgxpool.Pool, ctx context.Context, task NewTask) (int, error) {
	var id int
	err := dbPool.QueryRow(ctx, "INSERT INTO tasks (title, done, project_id) VALUES ($1, false, $2) RETURNING id", task.Title, task.ProjectId).Scan(&id)
	return id, err
}

var missingId = errors.New("'id' is a required field")

type TaskPatch struct {
	Id        *int    `json:"id"`
	Title     *string `json:"title"`
	Done      *bool   `json:"done"`
	ProjectId *string `json:"projectId"`
}

func patchTask(dbPool *pgxpool.Pool, ctx context.Context, patch TaskPatch) error {
	if patch.Id == nil {
		return missingId
	}

	args := []any{*patch.Id}
	var setClauses []string
	addField := func(name string) {
		setClauses = append(setClauses, name+"=$"+strconv.Itoa(len(args)+1))
	}

	if patch.Title != nil {
		addField("title")
		args = append(args, *patch.Title)
	}
	if patch.Done != nil {
		addField("done")
		args = append(args, *patch.Done)
	}
	if patch.ProjectId != nil {
		addField("project_id")
		args = append(args, *patch.ProjectId)
	}

	if len(setClauses) == 0 {
		return nil
	}
	query := "UPDATE tasks SET " + strings.Join(setClauses, ", ") + " WHERE id=$1"
	_, err := dbPool.Exec(ctx, query, args...)
	return err
}
