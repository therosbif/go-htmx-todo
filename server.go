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

func main() {
	fmt.Println("---------------")
	fmt.Println("hello world")

	todos := []Todo{
		{1, "Task 1", true},
		{2, "Task 2", true},
		{3, "Task 3", false},
	}
	getRootTmpl := func() *template.Template {
		return template.Must(template.ParseFiles("index.html"))
	}
	getTodoTmpl := func() *template.Template {
		return getRootTmpl().Lookup("todo")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		todos := map[string][]Todo{
			"Todos": todos,
		}
		getRootTmpl().Execute(w, todos)
	})

	http.HandleFunc("/add-todo", func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")

		todo := Todo{4, title, false}
		todos = append(todos, todo)
		getTodoTmpl().Execute(w, todo)
	})

	http.HandleFunc("/toggle-todo/", func(w http.ResponseWriter, r *http.Request) {
		var slug int
		_, err := fmt.Sscanf(r.URL.Path, "/toggle-todo/%d", &slug)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var todo Todo
		for i := 0; i < len(todos); i++ {
			if todos[i].ID == slug {
				todos[i].Done = !todos[i].Done
				todo = todos[i]
			}
		}
		if todo == (Todo{}) {
			http.NotFound(w, r)
			return
		}

		getTodoTmpl().Execute(w, todo)
	})

	log.Fatal(http.ListenAndServe(":4242", nil))
}
