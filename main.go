package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getTasks")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createTask")
	body := json.NewDecoder(r.Body)
	var task task
	err := body.Decode(&task)
	if err != nil {
		fmt.Println(err)
	}
	// Check if the task already exists
	for _, t := range tasks {
		if t.ID == task.ID {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Task already exists")
			return
		}
	}
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteTask")

	// get the id of the task to delete
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	for index, task := range tasks {
		if task.ID == id {
			// delete task
			tasks = append(tasks[:index], tasks[index+1:]...)
			break // Delete only one task
		}
	}
	json.NewEncoder(w).Encode(tasks)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: indexRoute")
	fmt.Fprintf(w, "Bienvenido a mi API!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/task", createTask).Methods("POST")
	router.HandleFunc("/task/{id:[0-9]+}", deleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}
