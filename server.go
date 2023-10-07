package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Todo struct {
	ID    int
	Title string
	Done  bool
}

func getIdFromPath(w http.ResponseWriter, r *http.Request, path string) (int, error) {
	var slug int
	_, err := fmt.Sscanf(r.URL.Path, path, &slug)
	if err != nil {
		return -1, err
	}
	return slug, nil
}

func getTmpl(name string) *template.Template {
	tmpl := template.Must(template.ParseFiles("index.html"))
	switch name {
	case "": // root
		return tmpl
	case "todo":
		return tmpl.Lookup("todo")
	case "edit":
		return tmpl.Lookup("edit")
	default:
		return nil
	}
}

func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	id := 4
	todos := []Todo{
		{1, "Task 1", true},
		{2, "Task 2", true},
		{3, "Task 3", false},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		todos := map[string][]Todo{
			"Todos": todos,
		}
		getTmpl("").Execute(w, todos)
	})

	http.HandleFunc("/add-todo", func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")

		todo := Todo{id, title, false}
		id++
		todos = append(todos, todo)
		getTmpl("todo").Execute(w, todo)
	})

	toggleTodo := func(w http.ResponseWriter, r *http.Request, id int) {
		var todo Todo
		for i := 0; i < len(todos); i++ {
			if todos[i].ID == id {
				todos[i].Done = !todos[i].Done
				todo = todos[i]

				getTmpl("todo").Execute(w, todo)
				return
			}
		}
		http.NotFound(w, r)
	}

	http.HandleFunc("/toggle-todo/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}

		id, err := getIdFromPath(w, r, "/toggle-todo/%d")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		toggleTodo(w, r, id)
	})

	deleteTodo := func(w http.ResponseWriter, r *http.Request, id int) {
		for i := 0; i < len(todos); i++ {
			if todos[i].ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		http.NotFound(w, r)
	}

	updateTitle := func(w http.ResponseWriter, r *http.Request, id int, title string) {
		for i := 0; i < len(todos); i++ {
			if todos[i].ID == id {
				todos[i].Title = title
				getTmpl("todo").Execute(w, todos[i])
				return
			}
		}
		http.NotFound(w, r)
	}

	http.HandleFunc("/todo/", func(w http.ResponseWriter, r *http.Request) {
		id, err := getIdFromPath(w, r, "/todo/%d")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodDelete:
			deleteTodo(w, r, id)
			break
		case http.MethodPut:
			title := r.PostFormValue("title")
			updateTitle(w, r, id, title)
			break
		default:
			methodNotAllowed(w)
			break
		}
	})

	http.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		id, err := getIdFromPath(w, r, "/edit/%d")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var todo Todo
		for i := 0; i < len(todos); i++ {
			if todos[i].ID == id {
				todo = todos[i]
			}
		}
		if todo == (Todo{}) {
			http.NotFound(w, r)
			return
		}

		getTmpl("edit").Execute(w, todo)
	})

	fmt.Println("---------------")
	fmt.Println("Server started at http://localhost:4242")
	log.Fatal(http.ListenAndServe(":4242", nil))
}
