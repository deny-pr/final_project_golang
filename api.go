package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var host = "127.0.0.1"
var port = "8080"
var connString = "admin:admin123@tcp(127.0.0.1:3306)/todo?charset=utf8&parseTime=True&loc=Local"

// @title To Do API
// @version 1.0
// @description This is a service for final project golang's course
// @termsOfService http://swagger.io/terms/
// @contact.email deny.prasetyo555@gmail.com
// @host 127.0.0.1:8080
// @BasePath /api/v1
func main() {
	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.PathPrefix("/todos").HandlerFunc(GetToDos).Methods(http.MethodGet)
	apiRouter.PathPrefix("/todo").HandlerFunc(GetToDoById).Methods(http.MethodGet)
	apiRouter.PathPrefix("/todo").HandlerFunc(CreateToDo).Methods(http.MethodPost)
	apiRouter.PathPrefix("/todo").HandlerFunc(UpdateToDo).Methods(http.MethodPut)
	apiRouter.PathPrefix("/todo").HandlerFunc(DeleteToDo).Methods(http.MethodDelete)
	fmt.Println("Listening on port: ", port)
	http.ListenAndServe(host+":"+port, router)

}

// GetToDos godoc
// @Summary Get ToDos
// @Description Get ToDos
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} Todos
// @Router /todos [get]
func GetToDos(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}
	defer db.Close()

	var tasks []Task
	rows, err := db.Query("SELECT * from tasks;")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong.")
		return
	}

	defer rows.Close()
	for rows.Next() {
		var eachtask Task
		var id int
		var description sql.NullString
		var taskType sql.NullString
		var time sql.NullTime
		var assigned_to sql.NullString

		rows.Scan(&id, &description, &taskType, &time, &assigned_to)
		eachtask.ID = id
		eachtask.Description = description.String
		eachtask.Type = taskType.String
		eachtask.Time = time.Time
		eachtask.Assigned_To = assigned_to.String
		tasks = append(tasks, eachtask)
	}

	respondWithJSON(w, http.StatusOK, tasks)

}

// GetToDo godoc
// @Summary Get ToDo
// @Description Get ToDo
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} Todos
// @Router /todos/id [get]
func GetToDoById(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}

	defer db.Close()

	id := r.URL.Query().Get("id")
	var description sql.NullString
	var taskType sql.NullString
	var time sql.NullTime
	var assigned_to sql.NullString
	err = db.QueryRow("SELECT description, type, time, assigned_to FROM tasks WHERE id=?", id).Scan(&description, &taskType, &time, &assigned_to)
	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, "No task found with the id="+id)
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	default:
		var eachtask Task
		eachtask.ID, _ = strconv.Atoi(id)
		eachtask.Description = description.String
		eachtask.Type = taskType.String
		eachtask.Time = time.Time
		eachtask.Assigned_To = assigned_to.String
		respondWithJSON(w, http.StatusOK, eachtask)

	}

}

// CreateTodo godoc
// @Summary Create ToDo
// @Description Create ToDo
// @Tags todos
// @Accept json
// @Produce json
// @Success 200
// @Router /todos [post]
func CreateToDo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}

	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var task Task
	err = decoder.Decode(&task)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}
	statement, err := db.Prepare("insert into tasks (description, type, time, assigned_to) values(?,?,?,?)")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}
	defer statement.Close()
	res, err := statement.Exec(task.Description, task.Type, task.Time, task.Assigned_To)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the task.")
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		task.ID = int(id)
		respondWithJSON(w, http.StatusOK, task)
	}

}

// UpdateTodo godoc
// @Summary Update ToDo
// @Description Update ToDo
// @Tags todos
// @Accept json
// @Produce json
// @Success 200
// @Router /todos/id [put]
func UpdateToDo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}

	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var task Task
	err = decoder.Decode(&task)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}
	statement, err := db.Prepare("UPDATE tasks set description=?, type=?, time=?, assigned_to=? where id=?")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}
	defer statement.Close()
	res, err := statement.Exec(task.Description, task.Type, task.Time, task.Assigned_To, task.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was problem updating the task.")
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		task.ID = int(id)
		respondWithJSON(w, http.StatusOK, task)
	}

}

// Delete ToDO godoc
// @Summary Delete ToDo
// @Description Delete ToDo
// @Tags todos
// @Accept json
// @Produce json
// @Success 200
// @Router /todos [delete]
func DeleteToDo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}

	defer db.Close()

	id := r.URL.Query().Get("id")
	var description sql.NullString
	var taskType sql.NullString
	var time sql.NullTime
	var assigned_to sql.NullString
	err = db.QueryRow("SELECT description, type, time, assigned_to from tasks where id=?", id).Scan(&description, &taskType, &time, &assigned_to)
	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, "No tasks found with the id="+id)
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	default:
		res, err := db.Exec("DELETE from tasks where id=?", id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			return
		}
		if count == 1 {
			var eachtask Task
			eachtask.ID, _ = strconv.Atoi(id)
			eachtask.Description = description.String
			eachtask.Type = taskType.String
			eachtask.Time = time.Time
			eachtask.Assigned_To = assigned_to.String
			respondWithJSON(w, http.StatusOK, eachtask)
		}
	}

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
