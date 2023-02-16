
# TodoList
A simple todo list app for practice CRUD functionality using ECO web framework and postgreSQL database

# Run this app
after cloning this repo use this command
`
docker-compose up --build  
`
## Features

- Creating a new task
- Updating a task
- Delete a task
- Get all tasks
- Get a single task

# API endpoints
## GET
#### get all todos
`http://localhost:8080/api/todos`
#### Response

```
{
"TodoList": [
{
"id": 62,
"title": "This is a test updating",
"complete": true
},
{
"id": 63,
"title": "Wake up at 8 am",
"complete": false
},
{
"id": 64,
"title": "Attend meeting at 9.30 am",
"complete": false
}
]
}
```
#### get single todo
`http://localhost:8080/api/todos/{id}`
#### Response

```
{
"id": 63,
"title": "Wake up at 8 am",
"complete": false
}
```

## POST
#### create new todo
`http://localhost:8080/api/todos/create`
- Request Body
```
{
"title": string
}
```
- Response
```
{
"id": integer,
"title": string,
"complete": boolean
}
```

## PATCH
#### update todo
`http://localhost:8080/api/todos/update/{id}`
- Request Body
```
{
"title": string,
"complete": boolean
}
```
- Response
200
```
{
"id": integer,
"title": string,
"complete": boolean
}
```
404 Error
```
id not found
```

## DELETE
#### update todo
`http://localhost:8080/api/todos/delete/{id}`
- Response
200
```
todo deleted
```
404 Error
```
id not found
```

# Swagger Documentation
`
http://localhost:8080/swagger/index.html
`