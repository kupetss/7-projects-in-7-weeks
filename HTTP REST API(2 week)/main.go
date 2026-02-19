package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Task struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

var (
	tasks = []Task{
		{Id: 1, Text: "By milk", Done: false},
	}
	nextId   = 2
	response map[string]interface{}
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tasks" && !strings.HasPrefix(r.URL.Path, "/tasks/") {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/tasks" || r.URL.Path == "/tasks/" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks)
		} else {
			http.Error(w, "!", http.StatusBadRequest)
		}
	case http.MethodPost:
		var newTask Task
		err := json.NewDecoder(r.Body).Decode(&newTask)
		if err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		if newTask.Text == "" {
			http.Error(w, "!", http.StatusBadRequest)
			return
		}

		newTask.Id = nextId
		nextId++
		tasks = append(tasks, newTask)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(newTask)
	case http.MethodPut:
		path := r.URL.Path
		if !strings.HasSuffix(path, "/done") {
			http.Error(w, "Invalid path. Use /tasks/{id}/done", http.StatusBadRequest)
			return
		}

		idStr := strings.TrimPrefix(path, "/tasks")
		idStr = strings.TrimPrefix(idStr, "/")
		idStr = strings.TrimSuffix(idStr, "/done")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "!", http.StatusBadRequest)
			return
		}

		found := -1
		for i, task := range tasks {
			if task.Id == id {
				found = i
				break
			}
		}
		if found == -1 {
			http.Error(w, "!", http.StatusNotFound)
			return
		}

		tasks[found].Done = true

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks[found])
	case http.MethodDelete:
		if strings.HasSuffix(r.URL.Path, "/done") {
			http.Error(w, "use del", http.StatusBadRequest)
			return
		}
		idStr := strings.TrimPrefix(r.URL.Path, "/tasks")
		idStr = strings.TrimPrefix(idStr, "/")

		id, _ := strconv.Atoi(idStr)
		found := -1
		for i, task := range tasks {
			if task.Id == id {
				found = i
				break
			}
		}
		if found == -1 {
			http.Error(w, "!", http.StatusNotFound)
			return
		}
		tasks = append(tasks[:found], tasks[found+1:]...)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "!", http.StatusMethodNotAllowed)
	}
}

func IsNullOrEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case string:
		return val == ""
	case []interface{}:
		return len(val) == 0
	case map[string]interface{}:
		return len(val) == 0
	default:
		return reflect.ValueOf(val).IsZero()
	}
}

func main() {
	http.HandleFunc("/tasks/", tasksHandler)
	http.ListenAndServe(":8080", nil)
}
