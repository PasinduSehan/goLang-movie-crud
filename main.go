package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var db *sql.DB

func initDB() {
	// Step 1: Connect to MySQL server without a password (assumes root@localhost no-auth or socket)
	dataSource := "root@tcp(127.0.0.1:3306)/" // No password - insecure, for testing only
	var err error
	db, err = sql.Open("mysql", dataSource)
	if err != nil {
		log.Fatal("Error connecting to MySQL server:", err)
	}

	// Step 2: Auto-create database if it doesn't exist
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS movie_db")
	if err != nil {
		log.Fatal("Error creating database:", err)
	}
	fmt.Println("Database 'movie_db' created or already exists!")

	// Step 3: Switch to the new database
	_, err = db.Exec("USE movie_db")
	if err != nil {
		log.Fatal("Error switching to database:", err)
	}

	// Step 4: Auto-create table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS movies (
			id VARCHAR(10) PRIMARY KEY,
			isbn VARCHAR(10),
			title VARCHAR(100),
			director_first_name VARCHAR(50),
			director_last_name VARCHAR(50)
		)
	`)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
	fmt.Println("Table 'movies' created or already exists!")

	// Step 5: Test the connection to the new DB
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}
	fmt.Println("Database connected successfully!")
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT id, isbn, title, director_first_name, director_last_name FROM movies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		var d Director
		err := rows.Scan(&m.ID, &m.Isbn, &m.Title, &d.FirstName, &d.LastName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		m.Director = &d
		movies = append(movies, m)
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	var m Movie
	var d Director
	err := db.QueryRow("SELECT id, isbn, title, director_first_name, director_last_name FROM movies WHERE id = ?", id).Scan(&m.ID, &m.Isbn, &m.Title, &d.FirstName, &d.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	m.Director = &d
	json.NewEncoder(w).Encode(m)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	movie.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - improve for production
	_, err := db.Exec("INSERT INTO movies (id, isbn, title, director_first_name, director_last_name) VALUES (?, ?, ?, ?, ?)",
		movie.ID, movie.Isbn, movie.Title, movie.Director.FirstName, movie.Director.LastName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	var updatedMovie Movie
	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}
	result, err := db.Exec("UPDATE movies SET isbn = ?, title = ?, director_first_name = ?, director_last_name = ? WHERE id = ?",
		updatedMovie.Isbn, updatedMovie.Title, updatedMovie.Director.FirstName, updatedMovie.Director.LastName, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	updatedMovie.ID = id
	json.NewEncoder(w).Encode(updatedMovie)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	result, err := db.Exec("DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	initDB()

	// Insert mock data into database
	mockMovies := []Movie{
		{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{FirstName: "John", LastName: "Doe"}},
		{ID: "2", Isbn: "454555", Title: "Movie Two", Director: &Director{FirstName: "Steve", LastName: "Smith"}},
	}
	for _, m := range mockMovies {
		_, err := db.Exec("INSERT IGNORE INTO movies (id, isbn, title, director_first_name, director_last_name) VALUES (?, ?, ?, ?, ?)",
			m.ID, m.Isbn, m.Title, m.Director.FirstName, m.Director.LastName)
		if err != nil {
			log.Println("Error inserting mock data:", err)
		} else {
			fmt.Println("Inserted movie:", m.Title)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
