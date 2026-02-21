package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type LogEntry struct {
	Action    string    `json:"action"`
	TaskId    int       `json:"task_id"`
	TaskText  string    `json:"task_text"`
	Timestamp time.Time `json:"timestamp"`
}

var (
	tasks = []Task{
		{Id: 1, Text: "Buy milk", Done: false},
	}
	nextId  = 2
	logFile = "tasks.log.json"
)

func logToFile(entry LogEntry) error {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(data) + "\n")
	return err
}

func readLogs() ([]LogEntry, error) {
	file, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []LogEntry{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var logs []LogEntry
	decoder := json.NewDecoder(file)
	for {
		var entry LogEntry
		if err := decoder.Decode(&entry); err != nil {
			break
		}
		logs = append(logs, entry)
	}
	return logs, nil
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimRight(r.URL.Path, "/")

	if path != "/tasks" && !strings.HasPrefix(path, "/task") {
		http.Error(w, "!", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if path == "/tasks" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)
			return
		}

		if path == "/log" {
			logs, err := readLogs()
			if err != nil {
				http.Error(w, "!", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(logs)
			return
		}

		http.Error(w, "!", http.StatusBadRequest)

	case http.MethodPost:
		if path != "/task" {
			http.Error(w, "!", http.StatusBadRequest)
			return
		}

		var newTask Task
		if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
			http.Error(w, "!", http.StatusBadRequest)
			return
		}

		if newTask.Text == "" {
			http.Error(w, "!", http.StatusBadRequest)
			return
		}

		newTask.Id = nextId
		nextId++
		newTask.Done = false
		tasks = append(tasks, newTask)

		logToFile(LogEntry{
			Action:    "create",
			TaskId:    newTask.Id,
			TaskText:  newTask.Text,
			Timestamp: time.Now(),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTask)

	case http.MethodPatch:
		if !strings.HasPrefix(path, "/task/") {
			http.Error(w, "!", http.StatusBadRequest)
			return
		}

		idStr := strings.TrimPrefix(path, "/task/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "!", http.StatusBadRequest)
			return
		}

		for i, task := range tasks {
			if task.Id == id {
				tasks[i].Done = !tasks[i].Done

				logToFile(LogEntry{
					Action:    "toggle",
					TaskId:    task.Id,
					TaskText:  task.Text,
					Timestamp: time.Now(),
				})

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tasks[i])
				return
			}
		}
		http.Error(w, "!", http.StatusNotFound)

	default:
		http.Error(w, "!", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", tasksHandler)
	http.ListenAndServe(":8080", nil)
}
