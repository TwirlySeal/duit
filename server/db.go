package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Date *Date `json:"date"`
	Time *Time `json:"time"`
}

func (t Task) FormatDate() string {
	datePart := t.Date.String()
	if t.Time == nil {
		return datePart
	}

	return datePart + " " + t.Time.String()
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

const tasksQuery string = "SELECT id, title, date, time FROM tasks WHERE project_id=$1 AND done=FALSE"

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
	Title string `json:"title"`
	ProjectId int `json:"projectId"`
	Date *string `json:"date"`
	Time *string `json:"time"`
}

func addTask(dbPool *pgxpool.Pool, ctx context.Context, task NewTask) (int, error) {
	var id int

	var columns []string
	var values []any
	var placeholders []string

	addParam := func(column string, value any) {
		columns = append(columns, column)
		values = append(values, value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(placeholders) + 1))
	}

	addParam("title", task.Title)
	addParam("project_id", task.ProjectId)

	join := func(elems []string) string {
		return strings.Join(elems, ", ")
	}

	if task.Date != nil {
		addParam("date", *task.Date)
	}
	if task.Time != nil {
		addParam("time", *task.Time)
	}

	err := dbPool.QueryRow(
		ctx,
		fmt.Sprintf(
			"INSERT INTO tasks (%s) VALUES (%s) RETURNING id",
			join(columns),
			join(placeholders),
		),
		values...,
	).Scan(&id)

	return id, err
}

var missingId = errors.New("'id' is a required field")

type TaskPatch struct {
	Id *int `json:"id"`
	Title *string `json:"title"`
	Done *bool `json:"done"`
	ProjectId *string `json:"projectId"`
	Date *string `json:"date"`
	Time *string `json:"time"`
}

func patchTask(dbPool *pgxpool.Pool, ctx context.Context, patch TaskPatch) error {
	if patch.Id == nil {
		return missingId
	}

	s := NewSQLMap(", ", *patch.Id)

	if patch.Title != nil {
		s.Param("title", *patch.Title)
	}
	if patch.Done != nil {
		s.Param("done", *patch.Done)
	}
	if patch.ProjectId != nil {
		s.Param("project_id", *patch.ProjectId)
	}
	if patch.Date != nil {
		s.Param("date", *patch.Date)
	}
	if patch.Time != nil {
		s.Param("time", *patch.Time)
	}

	if len(s.clauses) == 0 {
		return nil
	}
	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $1", s)
	_, err := dbPool.Exec(ctx, query, s.args...)
	return err
}
