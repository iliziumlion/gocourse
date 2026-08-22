package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/rs/cors"
)

// Модель задачи
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // pending, in_progress, completed
	Priority    string    `json:"priority"` // low, medium, high
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Хранилище (в памяти для примера)
type TaskStore struct {
	sync.RWMutex
	tasks  map[int]*Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks:  make(map[int]*Task),
		nextID: 1,
	}
}

// Task Handler
type TaskHandler struct {
	store *TaskStore
}

func NewTaskHandler(store *TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

// GET /api/tasks - получить все задачи
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	tasks := make([]*Task, 0, len(h.store.tasks))
	for _, task := range h.store.tasks {
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GET /api/tasks/:id - получить задачу по ID
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := getIDFromURL(r.URL.Path)

	h.store.RLock()
	task, exists := h.store.tasks[id]
	h.store.RUnlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// POST /api/tasks - создать задачу
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.store.Lock()
	task.ID = h.store.nextID
	h.store.nextID++
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	if task.Status == "" {
		task.Status = "pending"
	}
	if task.Priority == "" {
		task.Priority = "medium"
	}
	h.store.tasks[task.ID] = &task
	h.store.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// PUT /api/tasks/:id - обновить задачу
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := getIDFromURL(r.URL.Path)

	h.store.RLock()
	_, exists := h.store.tasks[id]
	h.store.RUnlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	h.store.Lock()
	task := h.store.tasks[id]
	if updatedTask.Title != "" {
		task.Title = updatedTask.Title
	}
	if updatedTask.Description != "" {
		task.Description = updatedTask.Description
	}
	if updatedTask.Status != "" {
		task.Status = updatedTask.Status
	}
	if updatedTask.Priority != "" {
		task.Priority = updatedTask.Priority
	}
	task.UpdatedAt = time.Now()
	h.store.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// DELETE /api/tasks/:id - удалить задачу
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := getIDFromURL(r.URL.Path)

	h.store.Lock()
	_, exists := h.store.tasks[id]
	if !exists {
		h.store.Unlock()
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}
	delete(h.store.tasks, id)
	h.store.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// Вспомогательная функция для получения ID из URL
func getIDFromURL(path string) int {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		var id int
		fmt.Sscanf(parts[len(parts)-1], "%d", &id)
		return id
	}
	return 0
}

// Stats Handler
type StatsHandler struct {
	store *TaskStore
}

func NewStatsHandler(store *TaskStore) *StatsHandler {
	return &StatsHandler{store: store}
}

// GET /api/stats - статистика
func (h *StatsHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	h.store.RLock()
	defer h.store.RUnlock()

	stats := map[string]int{
		"total":       len(h.store.tasks),
		"pending":     0,
		"in_progress": 0,
		"completed":   0,
	}

	for _, task := range h.store.tasks {
		switch task.Status {
		case "pending":
			stats["pending"]++
		case "in_progress":
			stats["in_progress"]++
		case "completed":
			stats["completed"]++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	store := NewTaskStore()
	taskHandler := NewTaskHandler(store)
	statsHandler := NewStatsHandler(store)

	// API Routes
	http.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTask(w, r)
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/stats", statsHandler.GetStats)

	// Serve frontend (после сборки)
	fs := http.FileServer(http.Dir("../frontend/dist"))
	http.Handle("/", fs)

	// CORS
	handler := cors.Default().Handler(http.DefaultServeMux)

	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("API available at http://localhost:8080/api")

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
