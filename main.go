package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Init books var as a slice Book struct
var books []Book

// Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get Single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) //get params
	//Loop through books and find with id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create Book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) //Mock ID - Not safe
	books = append(books, book)
	json.NewEncoder(w).Encode(&book)
}

// Update Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(&book)
			return
		}
	}
	json.NewEncoder(w).Encode(&books)
}

// Delete Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(&books)
}

func main() {

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

	//Init Router
	r := mux.NewRouter()

	//Mock Data
	books = append(books, Book{ID: "1", ISBN: "564322", Title: "Book One", Author: &Author{Firstname: "Mahedi", Lastname: "Hasan"}})
	books = append(books, Book{ID: "2", ISBN: "664323", Title: "Book two", Author: &Author{Firstname: "Petr", Lastname: "Louda"}})

	// Route Handlers / Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PATCH")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	//http.ListenAndServe(":8000", r)

	go func() {

		if err := http.ListenAndServe(":8000", r); err != nil && err != http.ErrServerClosed {

			logger.Fatalf("cannot listen on defined port: %s\n", err)

		}
	}()

	logger.Println("This is an info message.")
	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
}
