package main

import (
	"crudpractice/controllers"
	"crudpractice/database"
	_ "crudpractice/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	database.ConnectDatabase()
}

// @title Todo List App
// @version 1.0
// @description This is a simple todo list for practice CRUD Operation with Eco framework and PostgreSQL DB

// @contact.name Mursalin Sk
// @contact.email mursalin.sk@quantuminventions.com
// @server http://localhost:5000/
func main() {
	e := echo.New()
	e.Static("/static", "public")
	e.File("/", "public/index.html")
	e.File("/task/:id", "public/task.html")
	e.GET("/api/todos",controllers.GetallTodos)
	e.GET("/api/todos/:id",controllers.GetbyId)
	e.POST("/api/todos/create",controllers.CreateNewTodo)
	e.PATCH("/api/todos/update/:id",controllers.UpdateTodo)
	e.DELETE("/api/todos/delete/:id",controllers.DeleteTodo)
	e.GET("/swagger/*",echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}




