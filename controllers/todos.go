package controllers

import (
	"crudpractice/database"
	"crudpractice/models"
	"database/sql"
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
)

// GetallTodo godoc
// @Summary      get all todos
// @Tags         todos
// @Produce      json
// @Success      200  {object} models.TodoList
// @Router       /api/todos [get]
func GetallTodos(c echo.Context) error {
	var db *sql.DB = database.DB()
	sqlStatement := "SELECT id,title,iscomplete FROM todos"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	result := models.TodoList{}
	for rows.Next() {
		todo := models.Todo{}
		err2 := rows.Scan(&todo.Id, &todo.Title, &todo.IsComplete)
		if err2 != nil {
			return err2
		}
		result.Todos = append(result.Todos, todo)
	}
	return c.JSON(http.StatusCreated, result)
}
// GetSingleTodo godoc
// @Summary      get a single todos
// @Tags         todos
// @Produce      application/json
// @Success      200  {object} 
// @Router       /api/todos [get]
func GetbyId(c echo.Context) error {
	id := c.Param("id")
	var db *sql.DB = database.DB()
	sqlStatement := "SELECT id,title,iscomplete FROM todos WHERE id=$1"
	rows, err := db.Query(sqlStatement, id)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	todo := models.Todo{}
	for rows.Next() {
		err2 := rows.Scan(&todo.Id, &todo.Title, &todo.IsComplete)
		if err2 != nil {
			return err2
		}

	}
	return c.JSON(http.StatusCreated, todo)
}
// CreateTodo godoc
// @Summary      create todo
// @Description  store a new todo in list
// @Tags         todos
// @Param        title body object false  "enter your todo title"
// @Accept       json
// @Produce      json
// @Success      201  {object} models.Todo
// @Router       /api/todos/create [post]
func CreateNewTodo(c echo.Context) error {
	var db *sql.DB = database.DB()
	todo := new(models.Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	if todo.Title==""{
		return c.String(http.StatusBadRequest, "please enter a title")
	}
	sqlStatement := "INSERT INTO todos (title) VALUES ($1) RETURNING title,id,iscomplete"
	rows,err := db.Query(sqlStatement, todo.Title)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	restodo := models.Todo{}
	for rows.Next() {
		err2 := rows.Scan(&restodo.Title, &restodo.Id, &restodo.IsComplete)
		if err2 != nil {
			return err2
		}

	}
	return c.JSON(http.StatusCreated, restodo)
}

func UpdateTodo(c echo.Context) error {
	id := c.Param("id")
	var db *sql.DB = database.DB()
	todo := new(models.Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	sqlStatement := "UPDATE todos SET title=$1,iscomplete=$2 WHERE id=$3"
	_, err := db.Query(sqlStatement, todo.Title,todo.IsComplete,id)
	if err != nil {
		fmt.Println(err)
	} else {
		return c.JSON(http.StatusCreated, "todo updated")
	}
	return c.String(http.StatusOK, "ok")
}

func DeleteTodo(c echo.Context) error {
	id := c.Param("id")
	var db *sql.DB = database.DB()
	sqlStatement := "DELETE FROM todos WHERE id=$1"
	_, err := db.Query(sqlStatement,id)
	if err != nil {
		fmt.Println(err)
	} else {
		return c.JSON(http.StatusCreated, "todo deleted")
	}
	return c.String(http.StatusOK, "ok")
}




