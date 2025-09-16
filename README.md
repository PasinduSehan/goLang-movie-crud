# Go Movie CRUD

A simple RESTful API built with Go (Golang) for managing movie tickets using CRUD operations. It uses Gorilla Mux for routing and MySQL for data persistence, with auto-creation of the database and table.

## Features
- **Auto-Database Setup**: Creates `movie_db` and `movies` table on startup.
- **CRUD Endpoints**:
  - GET `/movies`: List all movies.
  - GET `/movies/{id}`: Get a specific movie.
  - POST `/movies`: Create a new movie (body: `{ "isbn": "123", "title": "Movie", "director": { "first_name": "John", "last_name": "Doe" } }`).
  - PUT `/movies/{id}`: Update a movie.
  - DELETE `/movies/{id}`: Delete a movie.
- Mock data inserted on startup for testing.

## Tech Stack
- Go (Golang)
- Gorilla Mux (routing)
- MySQL (database, with go-sql-driver)
- JSON API responses

## Setup and Run
1. Install Go and MySQL.
2. Update MySQL credentials in `initDB()` (no password for testing; use `skip-grant-tables` in `my.ini`).
3. Run `go mod tidy` to install dependencies.
4. Run `go run main.go`.
5. Test with Postman on `http://localhost:8000/movies`.

## Projects Relation
This project builds on my skills in full-stack development, similar to my Bank Transaction Web App (Spring Boot/MySQL) and Hotel Management System (Java/MySQL) from my BSc in Software Engineering.

## Author
Pasindu Weerathunga  
GitHub: [PasinduSehan](https://github.com/PasinduSehan)  
LinkedIn: [pasinduSehan](https://linkedin.com/in/pasinduSehan)  
Email: psehan12@gmail.com

## License
MIT License (or add your preferred).
"# goLang-movie-crud" 
