# To-Do API - Rest api сервис для to-do листов

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/GoX7/SnipLink-api/go.yml)](https://github.com/GoX7/SnipLink-api/actions)

## Features

- **Configuration**: Reading settings from the `config.yml` file
- **Logging**: Integrated logging
- **Database**: SQLite for data storage
- **Authentication**: Encrypted cookies
- **Validation**: Input data
- **Routing**: `chi` is used for routing

**ℹ️ Version:** 1.0.0  
**👤 Author:** GoX7  
**📜 License:** MIT (LICENSE file)  

## API Endpoints
**Tasks**  
- GET /tasks - Get a list of all tasks
- GET /tasks/{id} - Get a task by ID
- POST /tasks - Create a new task
- PATCH /tasks/{id} - Update the task
- DELETE /tasks/{id} - Delete a task
  
**Authentication**
- GET /auth/me - Information about the current user
- POST /auth/sign-in - Login
- POST /auth/sign-up - Registration of a new user

## Start project
```bash
git clone <repo-url>
cd to-do
go mod download
go run main.go
```

## Query examples
Creating a task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "New task", "description": "Task description"}'
```
Creating a user
```bash
curl -X POST http://localhost:8080/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{"username": "<username>", "password": "<password>"}'
```

## Project structure

```text
to-do/
├── config/                 #config api
│   └── config.yaml
├── internal/
│   ├── config/             #config reader
│   ├── http/
│   │   ├── controls/       #http service
│   │   │   ├──handlers/
│   │   │   └──interfaces/
│   │   └── cookie/         #cookie
│   ├── logger/             #logger
│   └── sqlite/             #db service
│       ├── db/
│       ├── sql-main.go    
│       ├── user.go         #work with user
│       └── tasks.go        #work with task 
├── pkg/             
│   ├── mw_logger/          #middleware logger
│   └── response/           #json response
└── main.go
```
