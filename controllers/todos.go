package controllers

import (
	"crudpractice/database"
	"crudpractice/models"
	"crudpractice/redis"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetallTodo godoc
// @Summary      get all todos
// @Tags         todos
// @Produce      application/json
// @Success      200  {object} models.TodoList
// @Router       /api/todos [get]
func GetallTodos(c echo.Context) error {
	var db *sql.DB = database.DB()
	redisClient := redis.GetRedisInstance()
	result := models.TodoList{}

	// Checking todos already present in redis cache
	isExists, err := redisClient.Exists("todos").Result()
	if err != nil {
		return err
	}
	if isExists == 1 {
		fmt.Println("Getting todos from redis cache")
		values, err := redisClient.HGetAll("todos").Result()
		if err != nil {
			return err
		}
		for _, item := range values {
			todo := &models.Todo{}
			err := json.Unmarshal([]byte(item), todo)
			if err != nil {
				return err
			}
			result.Todos = append(result.Todos, *todo)

		}
		return c.JSON(http.StatusAccepted, result)
	} else {
		// IF not in redis cache
		// Getting data from postgres DB
		fmt.Println("Getting todos from postgres database ")
		sqlStatement := "SELECT id,title,iscomplete FROM todos order by id"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()

		for rows.Next() {
			todo := models.Todo{}
			err2 := rows.Scan(&todo.Id, &todo.Title, &todo.IsComplete)
			if err2 != nil {
				return err2
			}
			result.Todos = append(result.Todos, todo)
			// Storing in Redis
			jsonValue, err := json.Marshal(todo)
			if err != nil {
				return err
			}
			redisClient.HSet("todos", strconv.Itoa(todo.Id), jsonValue)
		}
		
		return c.JSON(http.StatusAccepted, result)
	}

}

// GetSingleTodo godoc
// @Summary      get a single todos
// @Tags         todos
// @Produce      application/json
// @Param        id path int false "enter a id"
// @Success      200  {object} models.Todo
// @Failure      404  string message "id not found"
// @Router       /api/todos/{id} [get]
func GetbyId(c echo.Context) error {
	id := c.Param("id")
	var db *sql.DB = database.DB()
	todo := &models.Todo{}
	// Checking todo already present in redis cache
	redisClient := redis.GetRedisInstance()
	isExists, err := redisClient.HExists("todos", id).Result()
	if err != nil {
		return err
	}
	if isExists {
		fmt.Println("Getting single todo from redis")
		val, err := redisClient.HGet("todos", id).Result()
		if err != nil {
			return err
		}

		err2 := json.Unmarshal([]byte(val), todo)
		if err2 != nil {
			return err2
		}

	} else {
		fmt.Println("Getting single todo from postgres db")
		sqlStatement := "SELECT id,title,iscomplete FROM todos WHERE id=$1"
		rows, err := db.Query(sqlStatement, id)
		if err != nil {
			fmt.Println(err)
		}

		defer rows.Close()
		for rows.Next() {
			err2 := rows.Scan(&todo.Id, &todo.Title, &todo.IsComplete)
			if err2 != nil {
				return err2
			}
		}
		// Storing in Redis
		jsonValue, err := json.Marshal(todo)
		if err != nil {
			return err
		}
		redisClient.HSet("todos", strconv.Itoa(todo.Id), jsonValue)
	}
	if todo.Id == 0 && todo.Title == "" {
		return c.String(http.StatusNotFound, "id not found")
	}
	return c.JSON(http.StatusOK, todo)
}

// CreateTodo godoc
// @Summary      create todo
// @Description  store a new todo in list
// @Tags         todos
// @Param        value body models.CreateTodo false  "enter your todo title"
// @Accept       application/json
// @Produce      application/json
// @Success      201  {object} models.Todo
// @Router       /api/todos/create [post]
func CreateNewTodo(c echo.Context) error {
	var db *sql.DB = database.DB()
	redisClient := redis.GetRedisInstance()
	todo := new(models.CreateTodo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	if todo.Title == "" {
		return c.String(http.StatusBadRequest, "please enter a title")
	}
	sqlStatement := "INSERT INTO todos (title) VALUES ($1) RETURNING title,id,iscomplete"
	rows, err := db.Query(sqlStatement, todo.Title)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	resultTodo := models.Todo{}
	for rows.Next() {
		err2 := rows.Scan(&resultTodo.Title, &resultTodo.Id, &resultTodo.IsComplete)
		if err2 != nil {
			return err2
		}

	}
	// Storing in Redis
	jsonValue, err := json.Marshal(resultTodo)
	if err != nil {
		return err
	}
	redisClient.HSet("todos", strconv.Itoa(resultTodo.Id), jsonValue)
	return c.JSON(http.StatusCreated, resultTodo)
}

// UpdateTodo godoc
// @Summary      update todo
// @Description  update a todo by its id
// @Tags         todos
// @Param        id path int false  "enter id to update"
// @Param        value body models.UpdateTodo false "enter values"
// @Accept       application/json
// @Produce      application/json
// @Success      200  {object} models.Todo
// @Failure      404  string message "id not found"
// @Router       /api/todos/update/{id} [patch]
func UpdateTodo(c echo.Context) error {
	id := c.Param("id")
	var db *sql.DB = database.DB()
	redisClient := redis.GetRedisInstance()
	todo := new(models.UpdateTodo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	sqlStatement := "UPDATE todos SET title=$1,iscomplete=$2 WHERE id=$3 RETURNING title,id,iscomplete"
	rows, err := db.Query(sqlStatement, todo.Title, todo.IsComplete, id)
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
	jsonValue, err := json.Marshal(restodo)
	if err != nil {
		return err
	}
	redisClient.HSet("todos", strconv.Itoa(restodo.Id), jsonValue)

	if restodo.Id == 0 && restodo.Title == "" {
		return c.JSON(http.StatusNotFound, "id not found")
	}
	return c.JSON(http.StatusAccepted, restodo)
}

// DeleteTodo godoc
// @Summary      delete todo
// @Description  delete a todo by its id
// @Tags         todos
// @Param        id path int false "enter id to update"
// @Accept       application/json
// @Success      200  string message "id not found"
// @Failure      404  string message "id not found"
// @Router       /api/todos/delete/{id} [delete]
func DeleteTodo(c echo.Context) error {
	id := c.Param("id")
	var db *sql.DB = database.DB()
	redisClient := redis.GetRedisInstance()
	sqlStatement := "DELETE FROM todos WHERE id=$1"
	_, err := db.Query(sqlStatement, id)
	if err != nil {
		fmt.Println(err)
	} else {
		numDeleted, err := redisClient.HDel("todos", id).Result()
		if numDeleted == 0 {
			return c.JSON(http.StatusNotFound, "id not found")
		}
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, "todo deleted")
	}
	return c.String(http.StatusOK, "ok")
}
