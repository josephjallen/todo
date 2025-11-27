# README
This project is a basic todo application as outlined in the file [GolangCourse.md](GolangCourse.md); the solution satisfies the requirements outlined in points 1) through to 5).

## Basic Structure

### Packages
- **todostore** - contains the main logic for the todo lists 
- **web** - contains the api handlers and the web pages 
- **actors** - contains the actor pattern implementation for concurrent access to the todostore
- **filestorage** - contains r/w disk access and backup logic 
- **backups** - backup location for updated todo lists 
- **logger** - shared logging

### Main files
- **api.go** - starts the http server to accept api requests - to run: `go run api.go`
- **cli.go** - accepts cli commands to access the todostore, backups in the folder /backups will be made for each change to an existing todo - to run:
    - `go run cli.go -todoList=todod1 -additemname=monday -additemdescription=gotoshop`
    - `go run cli.go -todoList=todod1 -updateitemname=monday -updateitemdescription=gotoshop_updated`
    - `go run cli.go -todoList=todod1 -updateitemname=monday -updateitemstatus=started`
    - `go run cli.go -todoList=todod1 -deleteitemname=monday`
    
### Test files
- **api_test.go** - to run: `go test api_test.go -parallel=2`
- **actors/actor_test.go** - to run: `go test actor_test.go -parallel=2`
- **todostore/todo_test.go** - to run: `go test todo_test.go`

## Web Pages

First start HTTP Server with: `go run api.go`.

- http://127.0.0.1:8080/about
- http://127.0.0.1:8080/list/TODONAME e.g. http://127.0.0.1:8080/list/weeklytodo

## API Requests

First start HTTP Server with: `go run api.go`.

N.B. You can pass in your own traceid (see Get List Curl example below for the header) if required, if you don't pass in a traceid one will be allocated

### Create List
`curl --location --request GET '127.0.0.1:8080/createlist' \
--header 'Content-Type: application/json' \
--data '{
    "TodoListName": "Shopping"
}'`

### Get List
`curl --location --request GET '127.0.0.1:8080/getlist' \
--header 'X-Trace-ID: 123' \
--header 'Content-Type: application/json' \
--data '{
    "TodoListName": "Shopping"
}'`

### Add Item
`curl --location --request GET '127.0.0.1:8080/additem' \
--header 'Content-Type: application/json' \
--data '{
    "TodoListName": "Shopping",
    "ItemName": "Item1",
	"ItemDescription": "Buy Bread"
}'`

### Update Item Description
`curl --location --request GET '127.0.0.1:8080/updateitemdescription' \
--header 'Content-Type: application/json' \
--data '{
    "TodoListName": "Shopping",
    "ItemName": "Item1",
	"ItemDescription": "Buy White Bread"
}'`

### Update Item Status
`curl --location --request GET '127.0.0.1:8080/updateitemstatus' \
--header 'Content-Type: application/json' \
--data '{
    "TodoListName": "Shopping",
    "ItemName": "Item1",
	"ItemStatus": "started"
}'`
