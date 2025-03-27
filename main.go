package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "bookstore/handlers"
)

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/books", handlers.GetBooks).Methods("GET")
    router.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
    router.HandleFunc("/books", handlers.CreateBook).Methods("POST")
    router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
    router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")

    router.HandleFunc("/authors", handlers.GetAuthors).Methods("GET")
    router.HandleFunc("/authors", handlers.CreateAuthor).Methods("POST")

   
    router.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
    router.HandleFunc("/categories", handlers.CreateCategory).Methods("POST")

    log.Println("Server running on :8000")
    log.Fatal(http.ListenAndServe(":8000", router))
}
