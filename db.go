package main

type Task struct {
	Title string `json:"title"`
	Done bool `json:"done"`
}

type Project struct {
	Name string
	Id string
}

// Test data
var projects = []Project{
	Project{Name: "Personal", Id: "personal"},
	Project{Name: "Work", Id: "work"},
}

var personalTasks = []Task{
	Task{"Find Walter", false},
	Task{"Achieve enlightenment", true},
}

var workTasks = []Task{
	Task{"Steal the moon", false},
}

func getTasks(id string) []Task {
	switch id {
		case "personal":
			return personalTasks
		default:
			return workTasks
	}
}
