
# Shelf Master

Shelf Master is a modern library management system designed to streamline the tracking of books and users. 
This project aims to improve library efficiency and user experience.



## Features

- User and Book Tracking: Manage user accounts and track borrowed and returned books efficiently.
- Multi-book Borrowing: Allow users to borrow multiple books simultaneously with clear status tracking.
- Search and Filter: Advanced search and filtering options for books, authors, and categories.
- User Roles: Separate admin, librarian, and user functionalities for better management.
- Activity Logs: Maintain detailed logs for borrowing, returning, and overdue books.
- Notifications: Automated reminders for due dates and overdue books.
- Analytics: Generate reports on library usage and inventory trends.


## Prerequisites for Running Shelf Master

- Go Language
- Node.js
- Docker
- Docker Setup

Make sure these dependencies are installed and properly configured to run the project smoothly.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file or copy the .env.example

```sh
cp .env.example .env

```

`PORT`
`POSTGRES_PASSWORD`
`POSTGRES_USER`
`POSTGRES_DB`
`POSTGRES_PORT`
`POSTGRES_HOST`
`CONNECTION_STR`
`ENV`
`SMTP_USERNAME`
`SMTP_PASSWORD`
`SMTP_HOST`
`SMTP_PORT`
`SMTP_EMAIL`
`FRONTEND_URL`

## Run Locally

Clone the project

```bash
  git clone https://github.com/punpundada/Shelf-Master.git
```

Go to the project directory

```bash
  cd Shelf-Master
```
Copy env file

```bash
cp .env.example .env
```

Run on docker 

```bash
  docker compose up --build
```

connect database in terminal

```bash
  docker ps
  docker exec -it <container_name> psql -U <username> -d <database>
```
## Run Backend Saperately
From Parent directory CD into backend directory

```bash
cd backend
```

Start server

```bash
go run cmd/main.go
```
For hot reloading use AIR
```bash
go install github.com/air-verse/air@latest
```

To run project from `air` go to backend dir
```bash
air
```

From Parent directory CD into client directory
```bash
cd client
```

Run React app
```bash
npm run dev
```

## Makefile Info
- To run project through Makefile you must have it installed in your system
- Make file is for backend server

To run backend
```bash
make run
```
This will run `go run cmd/main.go`

To build backend
```bash
make build
```
This will run `go build -o bin/app cmd/main.go`

To Generate SQLC
```bash
make gen
```
This will run `docker run --rm -v ${shell pwd}:/src -w /src sqlc/sqlc generate`
Running a docker instance with sqlc and generating the schema

To run tests
```bash
make test
```
This will run `go test ./...`

# Hi, I'm Prajwal! 👋


## 🚀 About Me
I'm a passionate full-stack developer with expertise in building scalable and efficient web applications. Proficient in modern technologies like React, Node.js, Go, and Docker, I thrive on solving complex problems and delivering impactful solutions.
## 🔗 Links
[![GitHub](https://img.shields.io/badge/GitHub-%23121011.svg?logo=github&logoColor=white)](https://github.com/punpundada)
[![portfolio](https://img.shields.io/badge/my_portfolio-000?style=for-the-badge&logo=ko-fi&logoColor=white)](https://katherineoelsner.com/)
[![linkedin](https://img.shields.io/badge/linkedin-0A66C2?style=for-the-badge&logo=linkedin&logoColor=white)](https://www.linkedin.com/in/prajwal-parashkar-09a00a205/)
[![Bluesky](https://img.shields.io/badge/Bluesky-0285FF?logo=bluesky&logoColor=fff)](https://bsky.app/profile/prajwals.bsky.social)

